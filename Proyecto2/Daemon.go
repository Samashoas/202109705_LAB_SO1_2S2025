package main

import (
    "log"
    "os/exec"
    "time"
    "os"
    "os/signal"
    "syscall"
    "path/filepath"
    "encoding/json"
    "io/ioutil"
    "sort"
    "fmt"
    "database/sql"
    _ "github.com/mattn/go-sqlite3"
)

const (
    LimiteCPU	= 80.0 
    LimiteMemoria = 80.0
    TiemposVerificacion = 20 * time.Second
    DatabasePath = "./monitor.db"
)

// Estructuras para deserialización de la información de /proc
type RAMInfo struct {
    TotalKb int64 `json:"total_kb"`
    UsedKb  int64 `json:"used_kb"`
    FreeKb  int64 `json:"free_kb"`
}

type SysInfo struct {
    RAM RAMInfo `json:"ram"`
}

type ContainerProcess struct {
    PID         int    `json:"pid"`
    Name        string `json:"name"`
    ContainerID string `json:"container_id"`
    CMD         string `json:"cmd"`
    VszKb       int64  `json:"vsz_kb"`
    RssKb       int64  `json:"rss_kb"`
    PctMem      int    `json:"pct_mem"`
    PctCPU      int    `json:"pct_cpu"`
}

type ContainerInfo struct {
    ContainerProcesses []ContainerProcess `json:"container_processes"`
}

func setupDatabase() (*sql.DB, error) {
    // Eliminar base de datos existente si está corrupta
    if _, err := os.Stat(DatabasePath); err == nil {
        log.Printf("Base de datos existente encontrada, verificando integridad...")
    }
    
    db, err := sql.Open("sqlite3", DatabasePath)
    if err != nil {
        return nil, fmt.Errorf("error abriendo base de datos: %v", err)
    }

    // Verificar que la conexión funciona
    if err := db.Ping(); err != nil {
        return nil, fmt.Errorf("error conectando a la base de datos: %v", err)
    }

    // Crear tablas por separado para mejor control de errores
    createSystemMetricsSQL := `
    CREATE TABLE IF NOT EXISTS system_metrics (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
        total_ram_kb INTEGER,
        used_ram_kb INTEGER,
        free_ram_kb INTEGER,
        memory_percent REAL
    );`

    createContainerMetricsSQL := `
    CREATE TABLE IF NOT EXISTS container_metrics (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
        container_id TEXT,
        container_name TEXT,
        pid INTEGER,
        vsz_kb INTEGER,
        rss_kb INTEGER,
        mem_percent INTEGER,
        cpu_percent INTEGER,
        status TEXT
    );`

    // Ejecutar creación de tabla system_metrics
    if _, err := db.Exec(createSystemMetricsSQL); err != nil {
        return nil, fmt.Errorf("error creando tabla system_metrics: %v", err)
    }
    log.Println("✓ Tabla system_metrics creada/verificada")

    // Ejecutar creación de tabla container_metrics
    if _, err := db.Exec(createContainerMetricsSQL); err != nil {
        return nil, fmt.Errorf("error creando tabla container_metrics: %v", err)
    }
    log.Println("✓ Tabla container_metrics creada/verificada")

    // Verificar que las tablas se crearon correctamente
    var systemTableCount, containerTableCount int
    
    err = db.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='system_metrics';").Scan(&systemTableCount)
    if err != nil {
        return nil, fmt.Errorf("error verificando tabla system_metrics: %v", err)
    }
    
    err = db.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='container_metrics';").Scan(&containerTableCount)
    if err != nil {
        return nil, fmt.Errorf("error verificando tabla container_metrics: %v", err)
    }
    
    if systemTableCount == 0 {
        return nil, fmt.Errorf("tabla system_metrics no fue creada")
    }
    
    if containerTableCount == 0 {
        return nil, fmt.Errorf("tabla container_metrics no fue creada")
    }
    
    log.Printf("✓ Base de datos configurada correctamente con %d tablas", systemTableCount + containerTableCount)
    return db, nil
}

func EjecutarScript(ruta string) error{
    cmd := exec.Command("bash", ruta)
    output, err := cmd.CombinedOutput()
    if err != nil {
        log.Printf("Error ejecutando script %s: %v\nSalida: %s", ruta, err, string(output))
        return err	
    }
    log.Printf("Script %s ejecutado correctamente: %s", ruta, string(output))
    return nil
}

func LeerArchivoProcSysinfo() (*SysInfo, error) {
    contenido, err := ioutil.ReadFile("/proc/sysinfo_so1_202109705")
    if err != nil {
        return nil, err
    }

    var sysinfo SysInfo
    if err := json.Unmarshal(contenido, &sysinfo); err != nil {
        return nil, err
    }

    return &sysinfo, nil
}

func LeerArchivoProcContinfo() (*ContainerInfo, error) {
    contenido, err := ioutil.ReadFile("/proc/continfo_so1_202109705")
    if err != nil {
        return nil, err
    }

    var continfo ContainerInfo
    if err := json.Unmarshal(contenido, &continfo); err != nil {
        return nil, err
    }

    return &continfo, nil
}

func GuardarMetricasSistema(db *sql.DB, sysInfo *SysInfo) error {
    pctMemoria := float64(sysInfo.RAM.UsedKb) * 100.0 / float64(sysInfo.RAM.TotalKb)
    
    _, err := db.Exec(`
        INSERT INTO system_metrics 
        (total_ram_kb, used_ram_kb, free_ram_kb, memory_percent) 
        VALUES (?, ?, ?, ?)`,
        sysInfo.RAM.TotalKb, sysInfo.RAM.UsedKb, sysInfo.RAM.FreeKb, pctMemoria)
    
    if err != nil {
        return fmt.Errorf("error al insertar métricas de sistema: %v", err)
    }
    
    log.Printf("Métricas del sistema guardadas en BD: Memoria %.2f%%", pctMemoria)
    return nil
}

func GuardarMetricasContenedores(db *sql.DB, contInfo *ContainerInfo, status string) error {
    for _, proc := range contInfo.ContainerProcesses {
        _, err := db.Exec(`
            INSERT INTO container_metrics 
            (container_id, container_name, pid, vsz_kb, rss_kb, mem_percent, cpu_percent, status) 
            VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
            proc.ContainerID, proc.Name, proc.PID, proc.VszKb, proc.RssKb, 
            proc.PctMem, proc.PctCPU, status)
        
        if err != nil {
            return fmt.Errorf("error al insertar métricas de contenedor: %v", err)
        }
    }
    
    log.Printf("Métricas de %d contenedores guardadas en BD con status: %s", len(contInfo.ContainerProcesses), status)
    return nil
}

func AnalizarYGestionarContenedores(continfo *ContainerInfo, cleanPath string, db *sql.DB) error {
    if len(continfo.ContainerProcesses) == 0 {
        log.Println("No se encontraron procesos de contenedores para analizar")
        return nil
    }
    
    // Crear copias para ordenamiento
    byRAM := make([]ContainerProcess, len(continfo.ContainerProcesses))
    byCPU := make([]ContainerProcess, len(continfo.ContainerProcesses))
    byVSZ := make([]ContainerProcess, len(continfo.ContainerProcesses))
    byRSS := make([]ContainerProcess, len(continfo.ContainerProcesses))
    
    copy(byRAM, continfo.ContainerProcesses)
    copy(byCPU, continfo.ContainerProcesses)
    copy(byVSZ, continfo.ContainerProcesses)
    copy(byRSS, continfo.ContainerProcesses)
    
    // Ordenar por diferentes criterios
    sort.Slice(byRAM, func(i, j int) bool {
        return byRAM[i].PctMem > byRAM[j].PctMem
    })
    
    sort.Slice(byCPU, func(i, j int) bool {
        return byCPU[i].PctCPU > byCPU[j].PctCPU
    })
    
    sort.Slice(byVSZ, func(i, j int) bool {
        return byVSZ[i].VszKb > byVSZ[j].VszKb
    })
    
    sort.Slice(byRSS, func(i, j int) bool {
        return byRSS[i].RssKb > byRSS[j].RssKb
    })
    
    log.Println("Top contenedores por uso de RAM:")
    for i, p := range byRAM[:min(3, len(byRAM))] {
        log.Printf("%d. Container: %s (PID: %d), Mem: %d%%, CPU: %d%%", i+1, p.Name, p.PID, p.PctMem, p.PctCPU)
    }
    
    log.Println("Top contenedores por uso de CPU:")
    for i, p := range byCPU[:min(3, len(byCPU))] {
        log.Printf("%d. Container: %s (PID: %d), CPU: %d%%, Mem: %d%%", i+1, p.Name, p.PID, p.PctCPU, p.PctMem)
    }
    
    // Decisión de eliminar contenedores
    debeEliminar := false
    statusContenedores := "running"
    
    // Verificar si hay contenedores con alto consumo
    for i := 0; i < min(2, len(byRAM)); i++ {
        if byRAM[i].PctMem > 30 {
            log.Printf("Contenedor con alto uso de RAM detectado: %s (%d%%)", byRAM[i].Name, byRAM[i].PctMem)
            debeEliminar = true
            break
        }
    }
    
    for i := 0; i < min(2, len(byCPU)); i++ {
        if byCPU[i].PctCPU > 30 {
            log.Printf("Contenedor con alto uso de CPU detectado: %s (%d%%)", byCPU[i].Name, byCPU[i].PctCPU)
            debeEliminar = true
            break
        }
    }
    
    // Si detectamos contenedores problemáticos, ejecutar limpieza general
    if debeEliminar {
        log.Println("Ejecutando limpieza general de contenedores debido a alto consumo de recursos...")
        statusContenedores = "terminated"
        err := EjecutarScript(cleanPath)
        if err != nil {
            log.Printf("Error ejecutando limpieza: %v", err)
        }
    }
    
    // Guardar métricas de contenedores en la base de datos
    if err := GuardarMetricasContenedores(db, continfo, statusContenedores); err != nil {
        log.Printf("Error guardando métricas de contenedores: %v", err)
    }
    
    return nil
}

func min(a, b int) int {
    if a < b {
        return a
    }
    return b
}

func main(){
    log.Println("Iniciando el daemon de monitoreo de recursos...")
    
    // Configurar base de datos
    db, err := setupDatabase()
    if err != nil {
        log.Fatalf("Error configurando la base de datos: %v", err)
    }
    defer db.Close()
    
    log.Printf("Base de datos SQLite configurada en: %s", DatabasePath)
    
    // Obtener ruta absoluta al directorio donde está Daemon.go
    dirActual, err := os.Getwd()
    if err != nil {
        log.Fatalf("Error obteniendo directorio actual: %v", err)
    }
    
    // Construir rutas absolutas para los scripts
    // Scripts de cargas
    setupPath := filepath.Join(dirActual, "bash", "setup_cronjob.sh")
    loadSysinfoko := filepath.Join(dirActual, "bash", "load_sysinfoko.sh")
    loadContinfoko := filepath.Join(dirActual, "bash", "load_continfoko.sh")
    defaultContainerPath := filepath.Join(dirActual, "bash", "default_container.sh")

    // Scripts de limpieza
    cleanPath := filepath.Join(dirActual, "bash", "cleaning_container.sh")
    removeSysinfoko := filepath.Join(dirActual, "bash", "remove_sysinfoko.sh")
    removeContinfoko := filepath.Join(dirActual, "bash", "remove_continfoko.sh")
    deletePath := filepath.Join(dirActual, "bash", "delete_cronjob.sh")
    deleteDefContainer := filepath.Join(dirActual, "bash", "delete_defcontainer.sh")

    log.Printf("Cargando kernel: %s", loadSysinfoko)
    
    // Verificar que los scripts existen
    if _, err := os.Stat(loadSysinfoko); os.IsNotExist(err) {
        log.Printf("ADVERTENCIA: El script no existe en la ruta: %s", loadSysinfoko)
    } else {
        if err := EjecutarScript(loadSysinfoko); err != nil {
            log.Printf("Error al cargar el módulo del kernel: %v", err)
        } else {
            log.Println("Módulo del kernel cargado exitosamente")
        }
    }

    if _, err := os.Stat(loadContinfoko); os.IsNotExist(err) {
        log.Printf("ADVERTENCIA: El script no existe en la ruta: %s", loadContinfoko)
    } else {
        if err := EjecutarScript(loadContinfoko); err != nil {
            log.Printf("Error al cargar el módulo del kernel: %v", err)
        } else {
            log.Println("Módulo del kernel cargado exitosamente")
        }
    }

    log.Printf("Contenedor por defecto: %s", defaultContainerPath)
    
    // Contenedores default 
    log.Println("Creando contenedores por defecto...")
    if err := EjecutarScript(defaultContainerPath); err != nil{
        log.Printf("Error al crear contenedores por defecto: %v", err)
    }

    if _, err := os.Stat(setupPath); os.IsNotExist(err) {
        log.Printf("ADVERTENCIA: El script no existe en la ruta: %s", setupPath)
    }

    // Ejecutar script de configuración
    if err := EjecutarScript(setupPath); err != nil{
        log.Printf("Error al ejecutar el cronjob: %v", err)
    }

    // Configurar señales para terminar el daemon
    sigs := make(chan os.Signal, 1)
    signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
    done := make(chan bool, 1)
    go func(){
        sig := <-sigs
        log.Printf("Recibida señal %v. Terminando Daemon...", sig)
        done <- true
    }()
    
    // Loop principal
    loop:
        for{
            select{
            case <-done:
                break loop
            default:
                // Leer información del sistema
                sysInfo, err := LeerArchivoProcSysinfo()
                if err != nil {
                    log.Printf("Error leyendo /proc/sysinfo_so1_202109705: %v", err)
                } else {
                    log.Printf("RAM: Total: %d KB, Usado: %d KB, Libre: %d KB",
                        sysInfo.RAM.TotalKb, sysInfo.RAM.UsedKb, sysInfo.RAM.FreeKb)
                    
                    // Guardar métricas del sistema
                    if err := GuardarMetricasSistema(db, sysInfo); err != nil {
                        log.Printf("Error guardando métricas del sistema: %v", err)
                    }
                    
                    // Verificar si se exceden los límites
                    pctMemoriaUsada := float64(sysInfo.RAM.UsedKb) * 100.0 / float64(sysInfo.RAM.TotalKb)
                    if pctMemoriaUsada > LimiteMemoria {
                        log.Printf("Alerta: Uso de memoria excedido (%.2f%%), terminando contenedores", pctMemoriaUsada)
                        _ = EjecutarScript(cleanPath)
                    }
                }
                
                // Leer información de contenedores
                contInfo, err := LeerArchivoProcContinfo()
                if err != nil {
                    log.Printf("Error leyendo /proc/continfo_so1_202109705: %v", err)
                } else {
                    log.Printf("Procesos de contenedores detectados: %d", len(contInfo.ContainerProcesses))
                    
                    // Analizar y gestionar contenedores
                    if err := AnalizarYGestionarContenedores(contInfo, cleanPath, db); err != nil {
                        log.Printf("Error en la gestión de contenedores: %v", err)
                    }
                }
                
                time.Sleep(TiemposVerificacion)
            }
        }
        
    // Limpieza al terminar
    log.Println("Eliminando cronjob...")
    if err := EjecutarScript(deletePath); err != nil{
        log.Printf("Error al eliminar el cronjob: %v", err)
    }

    log.Println("Eliminando contenedores por defecto")
    if err := EjecutarScript(deleteDefContainer); err != nil{
        log.Printf("Error al eliminar contenedores por defecto: %v", err)
    }

    log.Println("Eliminando kernel de informacion del sistema...")
    if err := EjecutarScript(removeSysinfoko); err != nil {
        log.Printf("Error al eliminar el kernel de información del sistema: %v", err)
    }

    log.Println("Eliminando kernel de informacion de los contenedores...")
    if err := EjecutarScript(removeContinfoko); err != nil {
        log.Printf("Error al eliminar el kernel de información de los contenedores: %v", err)
    }

    log.Println("Esperando a que terminen los procesos reciduales del cronjob (10 segundos)")
    time.Sleep(10 * time.Second)

    for i := 0; i < 10; i++ {
        if err := EjecutarScript(cleanPath); err != nil{
            log.Printf("Error al eliminar contenedores: %v", err)
        }
    }
    log.Println("Daemon terminado correctamente")
}
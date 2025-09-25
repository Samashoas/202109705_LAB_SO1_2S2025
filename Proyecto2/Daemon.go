package main

import (
    "log"
    "os/exec"
    "time"
    "math/rand"
    "os"
    "os/signal"
    "syscall"
    "path/filepath"
)

const (
    LimiteCPU	= 80.0 
    LimiteMemoria = 80.0
    TiemposVerificacion = 10 * time.Second
)

func MetricasDelSistema() (UsoCPU int, UsoMemoria int){
    UsoCPU = rand.Intn(100)
    UsoMemoria = rand.Intn(100)
    return
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

func main(){
    log.Println("Iniciando el daemon de monitoreo de recursos...")
    
    // Obtener ruta absoluta al directorio donde está Daemon.go
    dirActual, err := os.Getwd()
    if err != nil {
        log.Fatalf("Error obteniendo directorio actual: %v", err)
    }
    
    // Construir rutas absolutas para los scripts
    // Scripts de cargas
    setupPath := filepath.Join(dirActual, "bash", "setup_cronjob.sh")
    loadSysinfoko := filepath.Join(dirActual, "bash", "load_sysinfoko.sh")
    defaultContainerPath := filepath.Join(dirActual, "bash", "default_container.sh")

    // Scripts de limpieza
    cleanPath := filepath.Join(dirActual, "bash", "cleaning_container.sh")
    removeSysinfoko := filepath.Join(dirActual, "bash", "remove_sysinfoko.sh")
    deletePath := filepath.Join(dirActual, "bash", "delete_cronjob.sh")
    deleteDefContainer := filepath.Join(dirActual, "bash", "delete_defcontainer.sh")

    log.Printf("Cargando kernel: %s", loadSysinfoko)
    
    // Verificar que los scripts existen
    if _, err := os.Stat(loadSysinfoko); os.IsNotExist(err) {
        log.Printf("ADVERTENCIA: El script no existe en la ruta: %s", loadSysinfoko)
    } else {
        // AQUÍ ESTÁ LA LÍNEA AÑADIDA: Carga el módulo del kernel
        if err := EjecutarScript(loadSysinfoko); err != nil {
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
                cpu, memoria := MetricasDelSistema()
                log.Printf("Uso CPU: %d%%, Uso Memoria: %d%%", cpu, memoria)
                if float64(cpu) > LimiteCPU || float64(memoria) > LimiteMemoria{
                    log.Println("Alerta: Uso de recursos excedido, terminando contenedores")
                    _ = EjecutarScript(cleanPath)
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

    log.Println("Esperando a que terminen los procesos reciduales del cronjob (10 segundos)")
    time.Sleep(10 * time.Second)

    for i := 0; i < 10; i++ {
        if err := EjecutarScript(cleanPath); err != nil{
            log.Printf("Error al eliminar contenedores: %v", err)
        }
    }
    log.Println("Daemon terminado correctamente")
}
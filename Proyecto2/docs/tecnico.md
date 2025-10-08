# Manual Técnico - Proyecto 2: Sistema de Monitoreo de Recursos

**Autor:** Juan Pablo Samayoa Ruiz  
**Carnet:** 202109705  
**Curso:** Sistemas Operativos 1  
**Fecha:** Septiembre 2025

---

## 1. Estructura del Módulo

### 1.1 Organización de Archivos y Directorios

```
Proyecto2/
├── kernel/                          # Módulos del kernel
│   ├── system_mod/                  # Módulo de información del sistema
│   │   ├── sysinfo.c               # Código fuente del módulo de sistema
│   │   ├── sysinfo.ko              # Módulo compilado
│   │   └── Makefile                # Script de compilación
│   └── container_mod/               # Módulo de información de contenedores
│       ├── continfo.c              # Código fuente del módulo de contenedores
│       ├── continfo.ko             # Módulo compilado
│       └── Makefile                # Script de compilación
├── bash/                           # Scripts de automatización
│   ├── load_sysinfoko.sh          # Cargar módulo de sistema
│   ├── load_continfoko.sh         # Cargar módulo de contenedores
│   ├── remove_sysinfoko.sh        # Remover módulo de sistema
│   ├── remove_continfoko.sh       # Remover módulo de contenedores
│   ├── setup_cronjob.sh           # Configurar cronjob
│   ├── delete_cronjob.sh          # Eliminar cronjob
│   ├── generate_container.sh      # Generar contenedores automáticamente
│   ├── default_container.sh       # Crear contenedores por defecto
│   ├── cleaning_container.sh      # Limpiar contenedores
│   └── delete_defcontainer.sh     # Eliminar contenedores por defecto
├── grafana/                        # Configuración de Grafana
│   ├── provisioning/
│   │   ├── datasources/
│   │   │   └── sqlite.yml         # Configuración del datasource SQLite
│   │   └── dashboards/
│   │       └── dashboards.yml     # Configuración de dashboards
│   └── dashboards/
│       └── system-monitor.json    # Dashboard principal
├── Daemon.go                       # Daemon principal de monitoreo
├── docker-compose.yml             # Configuración de Grafana
├── monitor.db                      # Base de datos SQLite
├── go.mod                         # Dependencias de Go
└── go.sum                         # Checksums de dependencias
```

### 1.2 Funciones Principales y su Propósito

#### Módulo `sysinfo.c`
- **`sysinfo_show()`**: Genera JSON con información de memoria RAM y procesos del sistema
- **`state_char()`**: Convierte estados de procesos a caracteres legibles
- **`sysinfo_open()`**: Maneja la apertura del archivo `/proc`
- **`sysinfo_init()`**: Inicializa el módulo y crea el archivo `/proc/sysinfo_so1_202109705`
- **`sysinfo_exit()`**: Limpia y remueve el módulo

#### Módulo `continfo.c`
- **`get_memory_kb()`**: Calcula uso de memoria VSZ y RSS en KB
- **`continfo_show()`**: Genera JSON con información de contenedores
- **`continfo_open()`**: Maneja la apertura del archivo `/proc`
- **`continfo_init()`**: Inicializa el módulo y crea el archivo `/proc/continfo_so1_202109705`
- **`continfo_exit()`**: Limpia y remueve el módulo

### 1.3 Dependencias Externas

#### Módulos del Kernel:
- `linux/init.h` - Macros de inicialización del kernel
- `linux/module.h` - Funciones básicas de módulos
- `linux/kernel.h` - Funciones del kernel
- `linux/proc_fs.h` - Sistema de archivos `/proc`
- `linux/seq_file.h` - Interfaz de archivos secuenciales
- `linux/mm.h` - Gestión de memoria
- `linux/sched/signal.h` - Señales de procesos
- `linux/sched.h` - Scheduler del kernel
- `linux/sysinfo.h` - Información del sistema

#### Daemon Go:
- `github.com/mattn/go-sqlite3` - Driver SQLite para Go

---

## 2. Compilación del Módulo

### 2.1 Prerrequisitos

```bash
# Instalar headers del kernel
sudo apt-get install linux-headers-$(uname -r)
sudo apt-get install build-essential
```

### 2.2 Compilación del Módulo de Sistema

```bash
# Navegar al directorio del módulo
cd kernel/system_mod/

# Compilar el módulo
make clean
make

# Verificar compilación exitosa
ls -la sysinfo.ko
```

### 2.3 Compilación del Módulo de Contenedores

```bash
# Navegar al directorio del módulo
cd kernel/container_mod/

# Compilar el módulo
make clean  
make

# Verificar compilación exitosa
ls -la continfo.ko
```

### 2.4 Estructura del Makefile

```makefile
obj-m += sysinfo.o  # o continfo.o para el módulo de contenedores

KDIR := /lib/modules/$(shell uname -r)/build
PWD := $(shell pwd)

all:
	$(MAKE) -C $(KDIR) M=$(PWD) modules

clean:
	$(MAKE) -C $(KDIR) M=$(PWD) clean
```

---

## 3. Carga y Descarga del Módulo

### 3.1 Carga Manual de Módulos

```bash
# Cargar módulo de sistema
sudo insmod kernel/system_mod/sysinfo.ko

# Cargar módulo de contenedores
sudo insmod kernel/container_mod/continfo.ko
```

### 3.2 Carga Automatizada con Scripts

```bash
# Cargar módulo de sistema
bash bash/load_sysinfoko.sh

# Cargar módulo de contenedores
bash bash/load_continfoko.sh
```

### 3.3 Verificación de Carga Correcta

```bash
# Verificar que los módulos están cargados
lsmod | grep sysinfo
lsmod | grep continfo

# Ver mensajes del kernel
dmesg | tail -10

# Verificar archivos proc creados
ls -la /proc/sysinfo_so1_202109705
ls -la /proc/continfo_so1_202109705
```

### 3.4 Descarga de Módulos

```bash
# Descarga manual
sudo rmmod sysinfo
sudo rmmod continfo

# Descarga automatizada
bash bash/remove_sysinfoko.sh
bash bash/remove_continfoko.sh
```

---

## 4. Pruebas y Verificación

### 4.1 Prueba del Módulo de Sistema

```bash
# Leer información del sistema
cat /proc/sysinfo_so1_202109705

# Salida esperada (formato JSON):
{
"ram" : {"total_kb": 8000000, "used_kb": 4000000, "free_kb": 4000000},
"processes" : [
  {"pid": 1, "state": "S"},
  {"pid": 2, "state": "S"}
]
}
```

### 4.2 Prueba del Módulo de Contenedores

```bash
# Leer información de contenedores
cat /proc/continfo_so1_202109705

# Salida esperada (formato JSON):
{
 "containers": [
  {"pid": 1234, "name": "container1", "cmd": "", "vsz_kb": 100000, "rss_kb": 50000, "pct_mem": 25, "pct_cpu": 15}
 ]
}
```

### 4.3 Comandos de Validación

```bash
# Validar formato JSON
cat /proc/sysinfo_so1_202109705 | jq '.'
cat /proc/continfo_so1_202109705 | jq '.'

# Monitoreo continuo
watch -n 2 'cat /proc/sysinfo_so1_202109705'

# Verificar logs del sistema
sudo dmesg -w | grep -E "(sysinfo|continfo)"
```

---

## 5. Decisiones de Diseño y Problemas

### 5.1 Decisiones Clave

#### **Uso de `/proc` filesystem:**
- **Decisión**: Exponer información a través de archivos virtuales en `/proc`
- **Razón**: Interfaz estándar de Unix/Linux para información del kernel
- **Beneficio**: Fácil acceso desde userspace sin syscalls customizadas

#### **Formato JSON para salida:**
- **Decisión**: Estructurar datos en formato JSON
- **Razón**: Facilita parsing desde aplicaciones (Go daemon)
- **Beneficio**: Interoperabilidad con herramientas modernas

#### **Separación de módulos:**
- **Decisión**: Módulos separados para sistema y contenedores
- **Razón**: Modularidad y mantenibilidad
- **Beneficio**: Carga/descarga independiente según necesidades

### 5.2 Problemas Encontrados y Soluciones

#### **Problema 1: Campo `state` no disponible en kernel 6.x**
```c
// Problema original:
if (t->state == TASK_INTERRUPTIBLE) return 'S';

// Solución aplicada:
if (READ_ONCE(t->__state) == TASK_INTERRUPTIBLE) return 'S';
```

#### **Problema 2: Función `seq_print` no existente**
```c
// Problema original:
seq_print(m, "data");

// Solución aplicada:
seq_printf(m, "data");
```

#### **Problema 3: Detección de estados de procesos**
- **Problema**: Muchos procesos retornaban estado '?'
- **Solución**: Ampliar función `state_char()` para más estados:
```c
if (state == TASK_IDLE) return 'I';
if (state == TASK_WAKING) return 'W';
if (state == TASK_PARKED) return 'P';
```

#### **Problema 4: Memoria VSZ/RSS en contenedores**
- **Problema**: Calcular memoria virtual y física correctamente
- **Solución**: Usar `get_mm_rss()` y campos de `task_struct->mm`

---

## 6. Estructura del Daemon GO

### 6.1 Funciones Principales

#### **6.1.1 Gestión de Grafana**

```go
func VerificarDependenciasGrafana() error
```
- **Propósito**: Verifica que Docker y Docker Compose estén instalados
- **Verificaciones**: Ejecuta comandos de versión y valida estructura de directorios
- **Retorna**: Error si alguna dependencia falta

```go
func CrearConfiguracionGrafana() error
```
- **Propósito**: Crea archivos de configuración para Grafana
- **Archivos creados**: 
  - `grafana/provisioning/datasources/sqlite.yml`
  - `grafana/provisioning/dashboards/dashboards.yml`
- **Configuración**: Datasource SQLite y provisioning automático

```go
func IniciarGrafana() error
```
- **Propósito**: Inicia contenedor de Grafana usando Docker Compose
- **Proceso**: 
  1. Ejecuta `docker-compose up -d`
  2. Verifica estado del contenedor
  3. Espera hasta que esté listo
- **Timeout**: 30 segundos máximo

```go
func DetenerGrafana() error
```
- **Propósito**: Detiene y remueve contenedor de Grafana
- **Proceso**: Ejecuta `docker-compose down`

#### **6.1.2 Gestión de Base de Datos**

```go
func setupDatabase() (*sql.DB, error)
```
- **Propósito**: Configura base de datos SQLite y crea tablas
- **Tablas creadas**:
  - `system_metrics`: Métricas del sistema (RAM, CPU)
  - `container_metrics`: Métricas de contenedores
- **Validaciones**: Verifica integridad de tablas

```go
func GuardarMetricasSistema(db *sql.DB, sysInfo *SysInfo) error
```
- **Propósito**: Inserta métricas del sistema en la base de datos
- **Datos guardados**: RAM total, usada, libre, porcentaje de memoria
- **Cálculos**: Porcentaje de memoria = (usada / total) * 100

```go
func GuardarMetricasContenedores(db *sql.DB, contInfo *ContainerInfo, status string) error
```
- **Propósito**: Inserta métricas de contenedores en la base de datos
- **Datos guardados**: PID, nombre, memoria VSZ/RSS, CPU, estado
- **Iteración**: Procesa cada contenedor individualmente

#### **6.1.3 Lectura de Información del Kernel**

```go
func LeerArchivoProcSysinfo() (*SysInfo, error)
```
- **Propósito**: Lee datos JSON del módulo de sistema
- **Archivo**: `/proc/sysinfo_so1_202109705`
- **Parsing**: Deserializa JSON a struct `SysInfo`

```go
func LeerArchivoProcContinfo() (*ContainerInfo, error)
```
- **Propósito**: Lee datos JSON del módulo de contenedores  
- **Archivo**: `/proc/continfo_so1_202109705`
- **Parsing**: Deserializa JSON a struct `ContainerInfo`

#### **6.1.4 Análisis y Gestión de Contenedores**

```go
func AnalizarYGestionarContenedores(continfo *ContainerInfo, cleanPath string, db *sql.DB) error
```
- **Propósito**: Analiza contenedores y toma decisiones de limpieza
- **Análisis realizado**:
  1. Ordenamiento por RAM usage
  2. Ordenamiento por CPU usage  
  3. Ordenamiento por VSZ (memoria virtual)
  4. Ordenamiento por RSS (memoria física)
- **Decisiones**: Si uso > 30% RAM o CPU, ejecuta limpieza
- **Logging**: Muestra top 3 contenedores por cada métrica

#### **6.1.5 Ejecución de Scripts**

```go
func EjecutarScript(ruta string) error
```
- **Propósito**: Ejecuta scripts bash del sistema
- **Proceso**: Usa `exec.Command("bash", ruta)`
- **Logging**: Registra salida y errores de ejecución

### 6.2 Estructuras de Datos

```go
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
```

### 6.3 Constantes de Configuración

```go
const (
    LimiteCPU            = 80.0                 // Límite de CPU (%)
    LimiteMemoria        = 80.0                 // Límite de memoria (%)  
    TiemposVerificacion  = 20 * time.Second     // Intervalo de monitoreo
    DatabasePath         = "./monitor.db"       // Ruta de base de datos
)
```

### 6.4 Flujo Principal del Daemon

1. **Inicialización**:
   - Verifica dependencias de Grafana
   - Crea configuración de Grafana
   - Configura base de datos SQLite
   - Inicia Grafana

2. **Carga de módulos**:
   - Carga módulo de sistema (`load_sysinfoko.sh`)
   - Carga módulo de contenedores (`load_continfoko.sh`)
   - Configura cronjob para generación de contenedores

3. **Loop de monitoreo** (cada 20 segundos):
   - Lee métricas del sistema desde `/proc/sysinfo_so1_202109705`
   - Guarda métricas en base de datos
   - Verifica límites de memoria/CPU del sistema
   - Lee métricas de contenedores desde `/proc/continfo_so1_202109705`
   - Analiza y ordena contenedores por recursos
   - Ejecuta limpieza si es necesario

4. **Limpieza al salir**:
   - Detiene Grafana
   - Remueve cronjob
   - Remueve módulos del kernel
   - Limpia contenedores restantes

---

## 7. Archivos de Automatización

### 7.1 Scripts de Módulos del Kernel

#### **load_sysinfoko.sh**
```bash
#!/bin/bash
echo "Cargando módulo de información del sistema..."
sudo insmod kernel/system_mod/sysinfo.ko
echo "Módulo sysinfo cargado exitosamente"
```

#### **load_continfoko.sh**  
```bash
#!/bin/bash
echo "Cargando módulo de información de contenedores..."
sudo insmod kernel/container_mod/continfo.ko
echo "Módulo continfo cargado exitosamente"
```

#### **remove_sysinfoko.sh**
```bash
#!/bin/bash
echo "Removiendo módulo de información del sistema..."
sudo rmmod sysinfo
echo "Módulo sysinfo removido exitosamente"
```

#### **remove_continfoko.sh**
```bash
#!/bin/bash  
echo "Removiendo módulo de información de contenedores..."
sudo rmmod continfo
echo "Módulo continfo removido exitosamente"
```

### 7.2 Scripts de Gestión de Contenedores

#### **generate_container.sh**
```bash
#!/bin/bash
# Genera 10 contenedores aleatorios con diferentes patrones de consumo

IMAGEN_ALTA_RAM="ubuntu:latest"
IMAGEN_ALTA_CPU="python:3.8-slim"  
IMAGEN_BAJA_RAM="alpine:latest"

for i in {1..10}; do
    imagen_seleccionada=$((RANDOM%3))
    
    if [ $imagen_seleccionada -eq 0 ]; then
        # Contenedor con alto consumo de RAM
        sudo docker run -d --label "proyecto=so1_lab" --name "alto_consumo_ram_$i" \
            $IMAGEN_ALTA_RAM bash -c "dd if=/dev/zero of=/tmp/testfile bs=1M count=100"
    elif [ $imagen_seleccionada -eq 1 ]; then
        # Contenedor con alto consumo de CPU
        sudo docker run -d --label "proyecto=so1_lab" --name "alto_consumo_cpu_$i" \
            $IMAGEN_ALTA_CPU bash -c "while true; do :; done"
    else
        # Contenedor con bajo consumo
        sudo docker run -d --label "proyecto=so1_lab" --name "bajo_consumo_$i" \
            $IMAGEN_BAJA_RAM sleep 300
    fi
done

echo "10 contenedores generados aleatoriamente"
```

#### **cleaning_container.sh**
```bash
#!/bin/bash
echo "Iniciando limpieza de contenedores..."

# Detener contenedores del proyecto
sudo docker stop $(sudo docker ps -q --filter "label=proyecto=so1_lab") 2>/dev/null

# Remover contenedores del proyecto  
sudo docker rm $(sudo docker ps -aq --filter "label=proyecto=so1_lab") 2>/dev/null

echo "Limpieza de contenedores completada"
```

### 7.3 Scripts de Cronjob

#### **setup_cronjob.sh**
```bash
#!/bin/bash
echo "Configurando cronjob para generación automática de contenedores..."

# Agregar cronjob para ejecutar cada minuto
(crontab -l 2>/dev/null; echo "* * * * * $(pwd)/bash/generate_container.sh") | crontab -

echo "Cronjob configurado exitosamente"
```

#### **delete_cronjob.sh**
```bash
#!/bin/bash
echo "Eliminando cronjob..."

# Remover cronjob específico
crontab -l | grep -v "generate_container.sh" | crontab -

echo "Cronjob eliminado exitosamente"
```

### 7.4 Docker Compose Configuration

#### **docker-compose.yml**
```yaml
version: '3.8'

services:
  grafana:
    image: grafana/grafana:latest
    container_name: grafana-monitor
    ports:
      - "3000:3000"
    volumes:
      - grafana-storage:/var/lib/grafana
      - ./monitor.db:/var/lib/grafana/monitor.db:ro
      - ./grafana/provisioning:/etc/grafana/provisioning:ro
    environment:
      - GF_SECURITY_ADMIN_USER=admin
      - GF_SECURITY_ADMIN_PASSWORD=admin
      - GF_INSTALL_PLUGINS=frser-sqlite-datasource
      - GF_LOG_LEVEL=warn
    restart: unless-stopped

volumes:
  grafana-storage:
```

---

## 8. Instalación y Ejecución

### 8.1 Prerrequisitos del Sistema

```bash
# Instalar dependencias del kernel
sudo apt-get update
sudo apt-get install linux-headers-$(uname -r) build-essential

# Instalar Docker
sudo apt-get install docker.io
sudo systemctl start docker
sudo systemctl enable docker

# Instalar Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/download/v2.23.3/docker-compose-$(uname -s)-$(uname -m)" \
    -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

# Instalar Go (si no está instalado)
sudo apt-get install golang-go

# Configurar permisos de Docker
sudo usermod -aG docker $USER
newgrp docker
```

### 8.2 Compilación e Instalación

```bash
# Clonar el proyecto
git clone <repository-url>
cd Proyecto2

# Compilar módulos del kernel
cd kernel/system_mod && make && cd ../..
cd kernel/container_mod && make && cd ../..

# Instalar dependencias de Go
go mod download

# Dar permisos de ejecución a scripts
chmod +x bash/*.sh
```

### 8.3 Ejecución del Sistema Completo

```bash
# Ejecutar el daemon (incluye inicialización automática)
go run Daemon.go

# El sistema automáticamente:
# 1. Verifica dependencias
# 2. Inicia Grafana
# 3. Carga módulos del kernel  
# 4. Configura cronjobs
# 5. Inicia monitoreo continuo

# Acceso a interfaces:
# - Grafana: http://localhost:3000 (admin/admin)
# - Métricas sistema: cat /proc/sysinfo_so1_202109705
# - Métricas contenedores: cat /proc/continfo_so1_202109705
```

### 8.4 Detener el Sistema

```bash
# Presionar Ctrl+C en el terminal del daemon
# El sistema automáticamente:
# 1. Detiene Grafana
# 2. Remueve cronjobs
# 3. Remueve módulos del kernel
# 4. Limpia contenedores
```

---

## 9. Troubleshooting

### 9.1 Errores Comunes de Compilación

**Error**: `No se ha encontrado la orden 'make'`
```bash
sudo apt-get install build-essential
```

**Error**: `linux/module.h: No such file or directory`
```bash
sudo apt-get install linux-headers-$(uname -r)
```

### 9.2 Errores de Módulos del Kernel

**Error**: `Operation not permitted`
```bash
# Usar sudo para cargar módulos
sudo insmod kernel/system_mod/sysinfo.ko
```

**Error**: `Invalid module format`
```bash
# Recompilar para la versión correcta del kernel
make clean && make
```

### 9.3 Errores de Docker

**Error**: `docker: command not found`
```bash
sudo apt-get install docker.io
sudo systemctl start docker
```

**Error**: `permission denied while trying to connect to Docker`
```bash
sudo usermod -aG docker $USER
newgrp docker
```

### 9.4 Errores de Grafana

**Error**: `Grafana no responde en localhost:3000`
```bash
# Verificar estado del contenedor
sudo docker ps | grep grafana
sudo docker logs grafana-monitor
```

**Error**: `SQLite datasource no funciona`
```bash
# Verificar que la base de datos existe
ls -la monitor.db
sqlite3 monitor.db ".tables"
```

---
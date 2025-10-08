# Manual de Usuario - Sistema de Monitoreo de Recursos

**Proyecto:** Sistema de Monitoreo de Contenedores y Recursos del Sistema  
**Autor:** Juan Pablo Samayoa Ruiz  
**Carnet:** 202109705  
**Curso:** Sistemas Operativos 1  
**Fecha:** Septiembre 2025

---

## Índice

1. [Introducción](#1-introducción)
2. [Guía de Instalación](#2-guía-de-instalación)
3. [Cómo Usar el Sistema](#3-cómo-usar-el-sistema)
4. [Dashboard de Grafana](#4-dashboard-de-grafana)
5. [Ejemplos de Uso Práctico](#5-ejemplos-de-uso-práctico)
6. [Arquitectura del Sistema](#6-arquitectura-del-sistema)
7. [Preguntas Frecuentes](#7-preguntas-frecuentes)
8. [Solución de Problemas](#8-solución-de-problemas)

---

## 1. Introducción

### ¿Qué es el Sistema de Monitoreo?

El Sistema de Monitoreo de Recursos es una herramienta completa que te permite:

- **Monitorear recursos del sistema** (RAM, CPU, procesos)
- **Supervisar contenedores Docker** en tiempo real
- **Visualizar datos** a través de dashboards interactivos
- **Automatizar la gestión** de contenedores según el uso de recursos
- **Almacenar historial** de métricas para análisis

### Características Principales

- **Monitoreo en tiempo real** cada 20 segundos
- **Interface web intuitiva** con Grafana
- **Gestión automática de contenedores** cuando hay alto consumo
- **Generación automática de contenedores de prueba**
- **Base de datos integrada** para almacenamiento de métricas
- **Sistema de alertas** por alto uso de recursos

### ¿Quién puede usar este sistema?

- Administradores de sistemas
- Desarrolladores trabajando con contenedores
- Estudiantes de sistemas operativos
- Cualquier usuario que necesite monitorear recursos del sistema

---

## 2. Guía de Instalación

### 2.1 Requisitos del Sistema

#### **Sistema Operativo:**
- Ubuntu 20.04 LTS o superior
- Debian 10 o superior
- Kernel Linux 5.4 o superior

#### **Recursos Mínimos:**
- **RAM:** 4 GB (recomendado 8 GB)
- **Espacio en disco:** 2 GB libres
- **CPU:** 2 núcleos (recomendado 4 núcleos)

#### **Software Necesario:**
- Docker y Docker Compose
- Go 1.19 o superior
- Headers del kernel Linux
- Build tools (make, gcc)

### 2.2 Instalación Paso a Paso

#### **Paso 1: Preparar el Sistema**

```bash
# Actualizar el sistema
sudo apt update && sudo apt upgrade -y

# Instalar herramientas básicas
sudo apt install git curl wget build-essential -y
```

#### **Paso 2: Instalar Docker**

```bash
# Instalar Docker
sudo apt install docker.io -y

# Iniciar y habilitar Docker
sudo systemctl start docker
sudo systemctl enable docker

# Agregar usuario al grupo docker (reiniciar sesión después)
sudo usermod -aG docker $USER
```

#### **Paso 3: Instalar Docker Compose**

```bash
# Descargar Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/download/v2.23.3/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose

# Dar permisos de ejecución
sudo chmod +x /usr/local/bin/docker-compose

# Verificar instalación
docker-compose --version
```

#### **Paso 4: Instalar Go**

```bash
# Instalar Go
sudo apt install golang-go -y

# Verificar instalación
go version
```

#### **Paso 5: Instalar Headers del Kernel**

```bash
# Instalar headers del kernel actual
sudo apt install linux-headers-$(uname -r) -y
```

#### **Paso 6: Descargar el Proyecto**

```bash
# Clonar el repositorio
git clone <URL_DEL_REPOSITORIO>
cd Proyecto2

# O descomprimir si tienes un archivo ZIP
unzip Proyecto2.zip
cd Proyecto2
```

#### **Paso 7: Configurar Permisos**

```bash
# Dar permisos de ejecución a todos los scripts
chmod +x bash/*.sh

# Verificar estructura de archivos
ls -la bash/
```

### 2.3 Verificación de la Instalación

```bash
# Verificar Docker
docker --version
docker ps

# Verificar Docker Compose
docker-compose --version

# Verificar Go
go version

# Verificar headers del kernel
ls /lib/modules/$(uname -r)/build
```

---

## 3. Cómo Usar el Sistema

### 3.1 Iniciar el Sistema

#### **Opción 1: Inicio Simple**

```bash
# Navegar al directorio del proyecto
cd Proyecto2

# Ejecutar el sistema completo
go run Daemon.go
```

#### **Opción 2: Compilar y Ejecutar**

```bash
# Compilar el proyecto
go build -o monitor Daemon.go

# Ejecutar
./monitor
```

### 3.2 Lo que Sucede al Iniciar

Cuando ejecutes el sistema, verás mensajes como estos:

```
Iniciando el daemon de monitoreo de recursos con Grafana...
Verificando dependencias de Grafana...
Docker disponible
Docker Compose disponible
Configuración de datasource creada: grafana/provisioning/datasources/sqlite.yml
Grafana iniciado exitosamente
Grafana disponible en: http://localhost:3000
Base de datos SQLite configurada en: ./monitor.db
Iniciando loop de monitoreo cada 20 segundos...
```

### 3.3 Interfaces Disponibles

Una vez iniciado el sistema, tendrás acceso a:

| Interface | URL/Comando | Descripción |
|-----------|-------------|-------------|
| **Dashboard Grafana** | `http://localhost:3000` | Interface web principal |
| **Métricas del Sistema** | `cat /proc/sysinfo_so1_202109705` | Datos JSON del sistema |
| **Métricas de Contenedores** | `cat /proc/continfo_so1_202109705` | Datos JSON de contenedores |
| **Base de Datos** | `sqlite3 monitor.db` | Base de datos SQLite |

### 3.4 Detener el Sistema

```bash
# Presionar Ctrl+C en el terminal donde está ejecutándose
# El sistema automáticamente:
# Detiene Grafana
# Limpia contenedores
# Remueve módulos del kernel
# Elimina trabajos programados
```

---

## 4. Dashboard de Grafana

### 4.1 Acceder a Grafana

1. **Abrir navegador web** y ir a: `http://localhost:3000`

2. **Iniciar sesión** con estas credenciales:
   - **Usuario:** `admin`
   - **Contraseña:** `admin`

3. **Cambiar contraseña** (opcional): Grafana te pedirá cambiar la contraseña por seguridad

### 4.2 Configurar Fuente de Datos

Si es la primera vez que usas el sistema:

1. **Ir a Configuration** (icono de configuración) → **Data Sources**
2. **Hacer clic en "Add data source"**
3. **Seleccionar "SQLite"**
4. **Configurar:**
   - **Name:** `System Monitor SQLite`
   - **Path:** `/var/lib/grafana/monitor.db`
5. **Hacer clic en "Save & Test"**

### 4.3 Crear Dashboards Personalizados

#### **Dashboard de Memoria del Sistema:**

1. **Ir a "+"** → **Dashboard** → **Add new panel**
2. **En Query, escribir:**
   ```sql
   SELECT 
     timestamp,
     memory_percent
   FROM system_metrics 
   ORDER BY timestamp DESC 
   LIMIT 100
   ```
3. **Configurar Visualization:** Time series
4. **Título:** "Uso de Memoria del Sistema (%)"

#### **Dashboard de Contenedores Activos:**

1. **Crear nuevo panel**
2. **Query:**
   ```sql
   SELECT 
     timestamp,
     COUNT(*) as container_count
   FROM container_metrics 
   WHERE status = 'running'
   GROUP BY timestamp
   ORDER BY timestamp DESC
   LIMIT 50
   ```
3. **Visualization:** Stat
4. **Título:** "Contenedores Activos"

#### **Dashboard de Top Contenedores por RAM:**

1. **Crear nuevo panel**
2. **Query:**
   ```sql
   SELECT 
     container_name,
     AVG(mem_percent) as avg_memory
   FROM container_metrics 
   WHERE timestamp > datetime('now', '-1 hour')
   GROUP BY container_name
   ORDER BY avg_memory DESC
   LIMIT 10
   ```
3. **Visualization:** Bar chart
4. **Título:** "Top Contenedores por Uso de RAM"

### 4.4 Configurar Alertas

1. **En cualquier panel, ir a Alert tab**
2. **Configurar condición, ejemplo:**
   - **Condition:** `IS ABOVE 80` (para uso de memoria > 80%)
   - **Evaluation:** Every `1m` for `2m`
3. **Configurar notificaciones** (email, Slack, etc.)

### 4.5 Variables y Filtros Dinámicos

Para crear dashboards más interactivos:

1. **Ir a Dashboard settings** → **Variables**
2. **Crear variable para filtrar por tiempo:**
   - **Name:** `time_range`
   - **Type:** Interval
   - **Values:** `5m,15m,1h,6h,24h`

3. **Usar en queries:**
   ```sql
   SELECT * FROM system_metrics 
   WHERE timestamp > datetime('now', '-$time_range')
   ```

---

## 5. Ejemplos de Uso Práctico

### 5.1 Monitoreo de un Servidor Web

**Escenario:** Tienes contenedores ejecutando aplicaciones web y necesitas monitorear su rendimiento.

**Pasos:**
1. **Iniciar el sistema de monitoreo**
2. **Crear contenedores web de prueba:**
   ```bash
   # El sistema automáticamente genera contenedores cada minuto
   # Puedes ver cuáles se están ejecutando:
   docker ps
   ```
3. **Observar en Grafana:**
   - Dashboard de memoria del sistema
   - Contenedores activos
   - Uso de CPU por contenedor

**Beneficios:**
- Detectar picos de memoria antes de que afecten el sistema
- Identificar contenedores problemáticos automáticamente
- Historial de rendimiento para análisis

### 5.2 Depuración de Problemas de Memoria

**Escenario:** El sistema se está quedando sin memoria y necesitas identificar la causa.

**Pasos:**
1. **Abrir Grafana** y ir al dashboard principal
2. **Observar el gráfico de memoria del sistema**
3. **Identificar contenedores con alto uso:**
   ```bash
   # Ver datos en tiempo real
   cat /proc/continfo_so1_202109705 | jq '.containers[] | select(.pct_mem > 20)'
   ```
4. **El sistema automáticamente limpiará contenedores problemáticos**

**Resultado:**
- Identificación rápida de consumidores de memoria
- Limpieza automática cuando se superan límites
- Logs detallados del proceso de limpieza

### 5.3 Análisis de Tendencias

**Escenario:** Quieres analizar patrones de uso durante una semana.

**Pasos:**
1. **Dejar el sistema ejecutándose por varios días**
2. **En Grafana, configurar tiempo de análisis:**
   - Time range: Last 7 days
3. **Crear queries para tendencias:**
   ```sql
   SELECT 
     DATE(timestamp) as day,
     AVG(memory_percent) as avg_memory,
     MAX(memory_percent) as peak_memory
   FROM system_metrics 
   WHERE timestamp > datetime('now', '-7 days')
   GROUP BY DATE(timestamp)
   ORDER BY day
   ```

**Insights obtenidos:**
- Patrones diarios de uso de memoria
- Identificación de horarios pico
- Planificación de capacidad

### 5.4 Automatización de Respuestas

**Escenario:** Configurar respuestas automáticas a problemas de recursos.

**Configuración actual del sistema:**
- Si memoria del sistema > 80% → Ejecuta limpieza general
- Si contenedor individual > 30% RAM → Marca para limpieza
- Si contenedor individual > 30% CPU → Marca para limpieza

**Personalización:**
Puedes modificar estos límites en `Daemon.go`:
```go
const (
    LimiteCPU      = 80.0  // Cambiar según necesidades
    LimiteMemoria  = 80.0  // Cambiar según necesidades
)
```

---

## 6. Arquitectura del Sistema

### 6.1 Visión General del Sistema

```
┌─────────────────────────────────────────────────────────────────┐
│                        USUARIO FINAL                            │
└─────────────────────┬───────────────────────────────────────────┘
                      │
                      ▼
┌─────────────────────────────────────────────────────────────────┐
│                     GRAFANA DASHBOARD                           │
│                  http://localhost:3000                          │
│              (Visualización y Análisis)                         │
└─────────────────────┬───────────────────────────────────────────┘
                      │
                      ▼
┌─────────────────────────────────────────────────────────────────┐
│                   BASE DE DATOS SQLITE                          │
│                     monitor.db                                  │
│            (Almacenamiento de Métricas)                         │
└─────────────────────┬───────────────────────────────────────────┘
                      │
                      ▼
┌─────────────────────────────────────────────────────────────────┐
│                    DAEMON GO                                    │
│                 (Daemon.go)                                     │
│              Orquestador Principal                              │
└─────┬─────────────────────────────────────────────────┬─────────┘
      │                                                 │
      ▼                                                 ▼
┌─────────────────┐                           ┌─────────────────┐
│ MÓDULOS KERNEL  │                           │ CONTENEDORES    │
│                 │                           │                 │
│ /proc/sysinfo   │◄──────────┬──────────────►│ Docker Engine   │
│ /proc/continfo  │           │               │                 │
│                 │           │               │ auto_containers │
└─────────────────┘           │               └─────────────────┘
                              │
                              ▼
                    ┌─────────────────┐
                    │ SCRIPTS BASH    │
                    │                 │
                    │ • Carga módulos │
                    │ • Gestiona cron │
                    │ • Limpia sistema│
                    └─────────────────┘
```

### 6.2 Flujo de Datos

```
1. RECOLECCIÓN DE DATOS
   ├── Módulo Kernel (sysinfo.ko) → /proc/sysinfo_so1_202109705
   ├── Módulo Kernel (continfo.ko) → /proc/continfo_so1_202109705
   └── Docker Engine → Lista de contenedores activos

2. PROCESAMIENTO
   ├── Daemon Go lee archivos /proc cada 20 segundos
   ├── Convierte JSON a estructuras Go
   ├── Analiza y ordena datos por consumo
   └── Toma decisiones automáticas de limpieza

3. ALMACENAMIENTO
   ├── Inserta métricas en SQLite (monitor.db)
   ├── Tabla system_metrics (memoria, CPU)
   └── Tabla container_metrics (contenedores individuales)

4. VISUALIZACIÓN
   ├── Grafana lee datos de SQLite
   ├── Genera gráficos en tiempo real
   └── Dashboards interactivos para el usuario
```

### 6.3 Componentes del Sistema

#### **Frontend (Interface de Usuario):**
- **Grafana Dashboard:** Interface web principal
- **Navegador Web:** Chrome, Firefox, Safari, Edge

#### **Backend (Lógica de Negocio):**
- **Daemon Go:** Orquestador principal del sistema
- **Base de Datos SQLite:** Almacenamiento persistente
- **Scripts Bash:** Automatización de tareas del sistema

#### **Kernel Layer (Recolección de Datos):**
- **Módulo sysinfo.ko:** Información del sistema operativo
- **Módulo continfo.ko:** Información de contenedores
- **Sistema /proc:** Interface entre kernel y userspace

#### **Infraestructura:**
- **Docker Engine:** Gestión de contenedores
- **Cron Jobs:** Generación automática de contenedores
- **Sistema de Archivos:** Almacenamiento de configuración

### 6.4 Interacciones Entre Componentes

1. **Daemon ↔ Módulos Kernel:**
   - Daemon lee `/proc/sysinfo_so1_202109705` y `/proc/continfo_so1_202109705`
   - Módulos exponen datos del kernel en formato JSON

2. **Daemon ↔ Base de Datos:**
   - Daemon inserta métricas cada 20 segundos
   - SQLite almacena historial para análisis

3. **Daemon ↔ Docker:**
   - Daemon ejecuta comandos Docker para limpieza
   - Scripts bash gestionan contenedores automáticamente

4. **Grafana ↔ Base de Datos:**
   - Grafana consulta SQLite para generar gráficos
   - Actualizaciones en tiempo real cada 10 segundos

5. **Usuario ↔ Sistema:**
   - Control través de Grafana web interface
   - Monitoreo manual con comandos de consola

---

## 7. Preguntas Frecuentes

## 7. Preguntas Frecuentes

### **¿Qué hago si Grafana no carga?**
```bash
# Verificar que Docker esté ejecutándose
sudo systemctl status docker

# Ver logs de Grafana
sudo docker logs grafana-monitor

# Reiniciar Grafana
sudo docker-compose restart grafana
```

### **¿Cómo cambio los límites de memoria/CPU?**
Edita el archivo `Daemon.go` y cambia estas constantes:
```go
const (
    LimiteCPU      = 80.0  // Tu nuevo límite de CPU
    LimiteMemoria  = 80.0  // Tu nuevo límite de memoria
)
```

### **¿Puedo usar el sistema sin Grafana?**
Sí, puedes comentar las partes de Grafana en `Daemon.go` y usar solo:
```bash
# Ver métricas directamente
cat /proc/sysinfo_so1_202109705
cat /proc/continfo_so1_202109705

# Consultar base de datos
sqlite3 monitor.db "SELECT * FROM system_metrics LIMIT 10;"
```

### **¿Cómo genero más contenedores de prueba?**
```bash
# Ejecutar script manualmente
bash bash/generate_container.sh

# O modificar el cronjob para ejecutar más frecuentemente
crontab -e
# Cambiar de * * * * * a */30 * * * * (cada 30 segundos)
```

### **¿El sistema afecta el rendimiento?**
- **Impacto mínimo:** < 1% de CPU en promedio
- **Memoria:** ~50MB para el daemon + Grafana
- **Disco:** ~100MB incluyendo base de datos
- **Red:** Solo local (puerto 3000)

### **¿Cómo hacer backup de los datos?**
```bash
# Backup de la base de datos
cp monitor.db backup_$(date +%Y%m%d).db

# Backup de configuración de Grafana
sudo docker exec grafana-monitor tar -czf - /var/lib/grafana > grafana_backup.tar.gz
```

### **¿Puedo monitorear múltiples servidores?**
El sistema actual monitorea un servidor. Para múltiples servidores:
- Instalar en cada servidor
- Configurar Grafana central para múltiples datasources
- Usar herramientas como Prometheus para agregación

---

## 8. Solución de Problemas

### 8.1 Problemas de Instalación

#### **Error: "docker: command not found"**
```bash
# Solución:
sudo apt update
sudo apt install docker.io
sudo systemctl start docker
sudo usermod -aG docker $USER
# Reiniciar sesión
```

#### **Error: "permission denied while trying to connect to Docker"**
```bash
# Solución:
sudo usermod -aG docker $USER
newgrp docker
# O reiniciar sesión completamente
```

#### **Error: "go: command not found"**
```bash
# Solución:
sudo apt install golang-go
# Verificar instalación:
go version
```

### 8.2 Problemas de Ejecución

#### **Error: "No se ha encontrado la orden 'make'"**
```bash
# Solución:
sudo apt install build-essential
```

#### **Error: "linux/module.h: No such file or directory"**
```bash
# Solución:
sudo apt install linux-headers-$(uname -r)
```

#### **Error: "insmod: ERROR: could not insert module"**
```bash
# Solución:
# Recompilar módulos para el kernel actual
cd kernel/system_mod && make clean && make
cd ../container_mod && make clean && make
```

### 8.3 Problemas de Grafana

#### **Grafana no responde en localhost:3000**
```bash
# Diagnóstico:
sudo docker ps | grep grafana
sudo docker logs grafana-monitor

# Soluciones:
# 1. Reiniciar Grafana
sudo docker-compose restart grafana

# 2. Verificar puertos
sudo netstat -tulpn | grep :3000

# 3. Usar puerto alternativo
# Editar docker-compose.yml: "3001:3000"
```

#### **Error: "SQLite datasource not working"**
```bash
# Verificar que la base de datos existe:
ls -la monitor.db

# Verificar contenido:
sqlite3 monitor.db ".tables"
sqlite3 monitor.db "SELECT COUNT(*) FROM system_metrics;"

# Si no existe, ejecutar daemon para crearla:
go run Daemon.go
```

### 8.4 Problemas de Rendimiento

#### **Sistema lento o con alta carga**
```bash
# Verificar procesos del sistema:
top
htop

# Verificar contenedores problemáticos:
docker stats

# Reducir frecuencia de monitoreo:
# Editar Daemon.go, cambiar:
TiemposVerificacion = 60 * time.Second  // cada minuto instead de 20 segundos
```

#### **Base de datos muy grande**
```bash
# Ver tamaño actual:
du -h monitor.db

# Limpiar datos antiguos:
sqlite3 monitor.db "DELETE FROM system_metrics WHERE timestamp < datetime('now', '-7 days');"
sqlite3 monitor.db "DELETE FROM container_metrics WHERE timestamp < datetime('now', '-7 days');"
sqlite3 monitor.db "VACUUM;"
```

### 8.5 Logs y Depuración

#### **Habilitar logs detallados**
```bash
# Ver logs del daemon en tiempo real:
go run Daemon.go | tee daemon.log

# Ver logs del kernel:
sudo dmesg | grep -E "(sysinfo|continfo)"

# Ver logs de Docker:
sudo docker logs grafana-monitor

# Ver logs del sistema:
sudo journalctl -f | grep -E "(docker|grafana)"
```

#### **Verificar estado completo del sistema**
```bash
# Script de diagnóstico rápido:
echo "=== ESTADO DEL SISTEMA ==="
echo "Docker: $(docker --version)"
echo "Go: $(go version)"
echo "Kernel: $(uname -r)"
echo ""
echo "=== MÓDULOS CARGADOS ==="
lsmod | grep -E "(sysinfo|continfo)"
echo ""
echo "=== CONTENEDORES ==="
docker ps
echo ""
echo "=== GRAFANA ==="
curl -s http://localhost:3000/api/health | jq '.' 2>/dev/null || echo "Grafana no responde"
echo ""
echo "=== BASE DE DATOS ==="
sqlite3 monitor.db "SELECT COUNT(*) as system_records FROM system_metrics;"
sqlite3 monitor.db "SELECT COUNT(*) as container_records FROM container_metrics;"
```

---

## Soporte

Si tienes problemas que no están cubiertos en este manual:

1. **Verificar logs** usando los comandos de depuración
2. **Reiniciar el sistema** completamente:
   ```bash
   # Detener daemon (Ctrl+C)
   # Limpiar Docker
   sudo docker-compose down
   sudo docker system prune -f
   # Reiniciar
   go run Daemon.go
   ```
3. **Consultar documentación técnica** en `tecnico.md`
4. **Contactar soporte técnico** con los logs del error

---

**¡Gracias por usar el Sistema de Monitoreo de Recursos!**

*Este sistema fue desarrollado como proyecto académico para el curso de Sistemas Operativos 1, Universidad de San Carlos de Guatemala.*
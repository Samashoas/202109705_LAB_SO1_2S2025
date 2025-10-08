#!/bin/bash
# filepath: /home/juan-pablo/Escritorio/Uni/SOPES1/202109705_LAB_SO1_2S2025/Proyecto2/bash/remove_sysinfoko.sh

MODULE_NAME="sysinfo"  # Este es el nombre real del módulo cargado
PROC_FILE="sysinfo_so1_202109705"  # Este es el nombre del archivo en /proc

echo "Verificando si el módulo $MODULE_NAME está cargado..."

# Verificar si el módulo está cargado
if lsmod | grep -q "$MODULE_NAME"; then
    echo "El módulo $MODULE_NAME está cargado."
    echo "Descargando el módulo $MODULE_NAME..."
    sudo rmmod "$MODULE_NAME"
    
    # Verificar si se descargó correctamente
    if ! lsmod | grep -q "$MODULE_NAME"; then
        echo "Módulo $MODULE_NAME descargado correctamente."
        
        # Verificar que el archivo en /proc se eliminó
        if [ ! -f "/proc/$PROC_FILE" ]; then
            echo "Archivo /proc/$PROC_FILE eliminado correctamente."
        else
            echo "ADVERTENCIA: El archivo /proc/$PROC_FILE aún existe a pesar de haber descargado el módulo."
        fi
    else
        echo "Error: No se pudo descargar el módulo $MODULE_NAME."
    fi
else
    echo "El módulo $MODULE_NAME no está cargado según lsmod."
    
    # Comprobar si el archivo en /proc existe a pesar de que el módulo no aparece en lsmod
    if [ -f "/proc/$PROC_FILE" ]; then
        echo "ADVERTENCIA: El archivo /proc/$PROC_FILE existe aunque el módulo no aparece en lsmod."
        echo "Intentando forzar la descarga..."
        sudo rmmod "$MODULE_NAME" 2>/dev/null || echo "No se pudo descargar el módulo."
    fi
fi
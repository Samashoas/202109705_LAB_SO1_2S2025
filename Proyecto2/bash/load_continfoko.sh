#!/bin/bash
# filepath: /home/juan-pablo/Escritorio/Uni/SOPES1/202109705_LAB_SO1_2S2025/Proyecto2/bash/load_continfoko.sh

MODULE_NAME="continfo"
PROC_FILE="continfo_so1_202109705"
MODULE_PATH="/home/juan-pablo/Escritorio/Uni/SOPES1/202109705_LAB_SO1_2S2025/Proyecto2/kernel/container_mod/continfo.ko"

echo "Verificando si el módulo $MODULE_NAME ya está cargado..."

# Verificar si el módulo ya está cargado
if lsmod | grep -q "$MODULE_NAME"; then
    echo "El módulo $MODULE_NAME ya está cargado."
else
    echo "El módulo $MODULE_NAME no está cargado."
    
    # Intentar descargarlo por si acaso
    echo "Intentando descargar el módulo por si acaso..."
    sudo rmmod $MODULE_NAME 2>/dev/null || echo "No se pudo descargar '$MODULE_NAME'"
    
    # Cargar el módulo
    if [ -f "$MODULE_PATH" ]; then
        echo "Cargando el módulo $MODULE_NAME..."
        sudo insmod "$MODULE_PATH"
        
        if lsmod | grep -q "$MODULE_NAME"; then
            echo "Módulo $MODULE_NAME cargado correctamente."
            
            if [ -f "/proc/$PROC_FILE" ]; then
                echo "Archivo /proc/$PROC_FILE creado correctamente:"
                cat /proc/$PROC_FILE
            else
                echo "ADVERTENCIA: El archivo /proc/$PROC_FILE no se creó."
            fi
        else
            echo "Error: No se pudo cargar el módulo $MODULE_NAME."
        fi
    else
        echo "Error: El archivo del módulo no existe en $MODULE_PATH"
    fi
fi
#!/bin/bash

# Definir las imágenes
IMAGEN_ALTA_RAM="ubuntu:latest"
IMAGEN_ALTA_CPU="python:3.8-slim"
IMAGEN_BAJA_RAM="alpine:latest"

# Definir nombres específicos para contenedores por defecto
BAJO_NOMBRES=("default_containerBajo1" "default_containerBajo2" "default_containerBajo3")
ALTO_NOMBRES=("default_containerAltoCPU" "default_containerAltoRAM")

# Verificar contenedores de bajo consumo por nombre
for nombre in "${BAJO_NOMBRES[@]}"; do
    if ! docker ps -q --filter "name=$nombre" | grep -q .; then
        echo "Creando contenedor bajo consumo: $nombre"
        docker run -d --label "proyecto=so1_lab" --name "$nombre" $IMAGEN_BAJA_RAM sleep 300
    else
        echo "El contenedor $nombre ya existe"
    fi
done

# Verificar contenedor de alto consumo CPU
if ! docker ps -q --filter "name=default_containerAltoCPU" | grep -q .; then
    echo "Creando contenedor alto consumo CPU: default_containerAltoCPU"
    docker run -d --label "proyecto=so1_lab" --name "default_containerAltoCPU" $IMAGEN_ALTA_CPU bash -c "while true; do :; done"
else
    echo "El contenedor default_containerAltoCPU ya existe"
fi

# Verificar contenedor de alto consumo RAM
if ! docker ps -q --filter "name=default_containerAltoRAM" | grep -q .; then
    echo "Creando contenedor alto consumo RAM: default_containerAltoRAM"
    docker run -d --label "proyecto=so1_lab" --name "default_containerAltoRAM" $IMAGEN_ALTA_RAM bash -c "dd if=/dev/zero of=/tmp/testfile bs=1M count=100"
else
    echo "El contenedor default_containerAltoRAM ya existe"
fi

echo "Verificación de contenedores por defecto completada."
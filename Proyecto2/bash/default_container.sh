#!/bin/bash
# filepath: /home/juan-pablo/Escritorio/Uni/SOPES1/202109705_LAB_SO1_2S2025/Proyecto2/bash/default_container.sh

# Definir las imágenes (las mismas que en generate_container.sh)
IMAGEN_ALTA_RAM="ubuntu:latest"
IMAGEN_ALTA_CPU="python:3.8-slim"
IMAGEN_BAJA_RAM="alpine:latest"

# Definir nombres específicos para contenedores por defecto
BAJO_NOMBRES=("default_containerBajo1" "default_containerBajo2" "default_containerBajo3")
ALTO_NOMBRES=("default_containerAltoCPU" "default_containerAltoRAM")

echo "Creando contenedores por defecto..."

# Eliminar contenedores por defecto existentes (para evitar conflictos)
for nombre in "${BAJO_NOMBRES[@]}" "${ALTO_NOMBRES[@]}"; do
    docker rm -f "$nombre" 2>/dev/null
done

# Crear contenedores de bajo consumo
for nombre in "${BAJO_NOMBRES[@]}"; do
    echo "Creando contenedor bajo consumo: $nombre"
    docker run -d --label "proyecto=def_so1_lab" --name "$nombre" $IMAGEN_BAJA_RAM \
        sh -c "while true; do echo 'Contenedor activo: $nombre'; sleep 60; done"
done

# Crear contenedor de alto consumo CPU
echo "Creando contenedor alto consumo CPU: default_containerAltoCPU"
docker run -d --label "proyecto=def_so1_lab" --name "default_containerAltoCPU" $IMAGEN_ALTA_CPU \
    bash -c "while true; do :; done"

# Crear contenedor de alto consumo RAM
echo "Creando contenedor alto consumo RAM: default_containerAltoRAM"
docker run -d --label "proyecto=def_so1_lab" --name "default_containerAltoRAM" $IMAGEN_ALTA_RAM \
    bash -c "dd if=/dev/zero of=/tmp/testfile bs=1M count=100 && tail -f /dev/null"

echo "Verificando que los contenedores estén en ejecución..."
docker ps --filter "label=proyecto=def_so1_lab" | grep "default_container"

echo "Contenedores por defecto creados correctamente."
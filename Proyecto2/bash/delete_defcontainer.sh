#!/bin/bash
# filepath: /home/juan-pablo/Escritorio/Uni/SOPES1/202109705_LAB_SO1_2S2025/Proyecto2/bash/delete_default_containers.sh

echo "Eliminando únicamente contenedores por defecto..."

# Lista de nombres de contenedores por defecto
DEFAULT_CONTAINERS=("default_containerBajo1" "default_containerBajo2" "default_containerBajo3" "default_containerAltoCPU" "default_containerAltoRAM")

for container in "${DEFAULT_CONTAINERS[@]}"; do
    if docker ps -a -q --filter "name=$container" | grep -q .; then
        echo "Deteniendo y eliminando $container..."
        docker stop "$container" 2>/dev/null
        docker rm "$container" 2>/dev/null
        echo "$container eliminado."
    else
        echo "$container no existe, nada que eliminar."
    fi
done

echo "Todos los contenedores por defecto han sido eliminados."
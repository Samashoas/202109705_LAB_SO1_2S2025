#!/bin/bash

IMAGEN_ALTA_RAM="ubuntu:latest"
IMAGEN_ALTA_CPU="python:3.8-slim"
IMAGEN_BAJA_RAM="alpine:latest"

for i in {1..10}; do
    imagen_seleccionada=$((RANDOM%3))

    if [ $imagen_seleccionada -eq 0 ]; then
        echo "Creando contenedor RAM: alto_consumo_ram_$i"
        sudo docker run -d --label "proyecto=so1_lab" --name "alto_consumo_ram_$i" $IMAGEN_ALTA_RAM bash -c "dd if=/dev/zero of=/tmp/testfile bs=1M count=100"
    elif [ $imagen_seleccionada -eq 1 ]; then
        echo "Creando contenedor CPU: alto_consumo_cpu_$i"
        sudo docker run -d --label "proyecto=so1_lab" --name "alto_consumo_cpu_$i" $IMAGEN_ALTA_CPU bash -c "while true; do :; done"
    else
        echo "Creando contenedor básico: bajo_consumo_$i"
        sudo docker run -d --label "proyecto=so1_lab" --name "bajo_consumo_$i" $IMAGEN_BAJA_RAM sleep 300
    fi
done

echo "10 contenedores generados aleatoriamente"
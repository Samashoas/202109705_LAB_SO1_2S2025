#!/bin/bash

echo "Eliminando contenedores del proyecto SO1..."

# Parar y eliminar solo los contenedores con la etiqueta del proyecto
sudo docker stop $(sudo docker ps -q --filter "label=proyecto=so1_lab") 2>/dev/null
sudo docker rm $(sudo docker ps -aq --filter "label=proyecto=so1_lab") 2>/dev/null

echo "Contenedores del proyecto eliminados. Grafana y otros contenedores permanecen intactos."
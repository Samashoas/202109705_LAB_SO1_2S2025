IMAGEN_ALTA_RAM="ubuntu:latest"
IMAGEN_ALTA_CPU="python:3.8-slim"
IMAGEN_BAJA_RAM="alpine:latest"

for i in {1..10} do
    imagen_seleccionada=$((RANFOM%3))

    if [$imagen_seleccionada -eq 0]; then
        docker run -d --rm --name "alto consumo_ram_$i" $IMAGEN_ALTA_RAM bashh -c "dd if=/dev/zero of=/tmp/testfile bs=1M count=1000">/dev/null 2>&1
    elif [$imagen_seleccionada -eq 1]; then
        docker run -d --rm --name "alto consumo_cpu_$i" $IMAGEN_ALTA_CPU bash -c "while true; do :; done">/dev/null 2>&1
    else
        docker run -d --rm --name "bajo_consumo_$i" $IMAGEN_BAJA_RAM sleep 300>/dev/null 2>&1
    fi
done

echo "10 contenedores generados aleatoriamente"

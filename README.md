# 202109705_LAB_SO1_2S2025

## 1. Guia de instalación

### 1.1 Instalación de hipervisor KVM

**Verificaciones previas**:

1. Actualizar dependecias del sitema operativo en el que esté realizando la virtualización (en este caso Ubuntu) con el comando:

```bash
sudo apt update
```

2. Verificar que el CPU permita la virtualización
3. Habilitar VMX/SVM desde la BIOS de la computadora en el apartado de ADVANCE
4. Comprobar si la virtualización ha sido activada exitosamente ingresando el siguiente comando:

```bash
egrep  -c '(vmx|vms)' /proc/cpuinfo
```

Si el comando devuelve un número diferente de 0 significa que la virtualización es posible.

5. De igual forma es recomendable instalar un checker de CPU con el siguiente comando:

```bash
sudo apt install cpu-checker -y
```

Esta herramienta ayuda a verificar si se te permite realizar maquinas virtuales utilizando maquinas virtuales a base del kernel (KVM)

6. Ingresar el siguienete comando para verificar si se tiene soporte para KVM y si este se encuentra habilitado

```bash
kvm-ok

la respuesta debe ser similar a la siguiente:

INFO: /dev/kvm exists
KVM acceleration can be used
```

**Instalación Qemu-KVM**

7. Proceder a la sintalación de Qemu-KVM

```bash
sudo apt install -y qemu-kvm virt-manager libvirt-daemon-system virtinst libvirt-clients bridge-utils
```

Una vez realizado esto el sistema estará listo para la virtualización

8. Habilitar e iniciar los servicios de libvirt

```bash
echo Asegura que el servicio se incia cada vez que se bootea la computadora
sudo sytemctl enable --now libvirtd

echo inicia el servicio de libvirtd
sudo sytemctl start libvirtd

echo verificar el estado del servicio (debe indicar activo y corriendo)
sudo systemctl status libvirtd
```

9. Proceder a la creación de las Maquinas virtuales (Realizar con la propia interfaz grafica que proporciona UBUNTU)




## API1-ENDPONTS_TEST

curl http://192.168.122.207:8081/api1/202109705/llamar-api2
curl http://192.168.122.207:8081/api1/202109705/llamar-api3


## API2-ENDPOINTS_TEST

curl http://192.168.122.207:8082/api2/202109705/llamar-api1
curl http://192.168.122.207:8082/api2/202109705/llamar-api3

## API3-ENDPOINTS_TEST

curl http://192.168.122.114:8083/api3/202109705/llamar-api1
curl http://192.168.122.114:8083/api3/202109705/llamar-api2

## Comandos para levantar contenedores en VM1, VM2 y VM3

* sudo nerdctl start <name>
* sudo docker start <name>

## Comandos utiles para verificar espacio en memoria y liberar espacio

* df -h
* sudo rm -rf /opt/zot/data/_upload 2>/dev/null || true

- curl http://192.168.122.158:5000/v2/_catalog
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

9. Crear grupos de usuario para KVM y libvirt (De esta forma se evita escribir el comando sudo cada vez que se tenga que realizar una operación con estas herramientas)

```bash
sudo usermod -aG kvm $USER
sudo usermod -aG libvirt $USER
```

Una vez realizados todos estos pasos ya se puede proceder con la creación y configuración de las maquinar virtuales

## 1.2 Creación Maquinas virtuales 

1. Descargar la imagen de ubuntu server desde la pagina oficila de ubuntu o haciendo click en el siguiente enlace:

[Descargar Ubuntu Server](https://ubuntu.com/download/server)

2. Ingresar a la herramienta visual de Gestor de Maquinas virtuales y crear una nueva maquina virtual y seleccionar la imagen ISO descargada anteriormente para poder usar la maquina virtual con Ubuntu server

3. Asignar recursos a las maquinar virtuales, a continuación se tiene una tabla con los recursos asignados para la elaboración de este proyecto (Nota: En este caso los recurso pueden variar dependiendo de la disposición que tenga la maquina HOST)

<div align="center">

| VM | Tipo de red | RAM | Nucleos del CPU | Almacenamiento |
|:--:|:-----------:|:---:|:---------------:|:--------------:|
|VM1 | RED NAT (default)| 2 GB | 2 | 7 GB |
|VM2 | RED NAT (default)| 2 GB | 2 | 7 GB |
|VM3 | RED NAT (default)| 2 GB | 2 | 10 GB |

</div>

Una vez realizado esto ya se puede proceder a la creación de la maquina virtual

4. Inciar la maquina virtual y configurar UBUNTU server, a continuación se deja una tabla con los datos ingresados para la configuración de UBUNTU server

<div align="center">

| VM | Usuario | Contraseña |
|:--:|:-----------:|:---:|
|VM1 | sopes1vm1|SOPES1VM1|
|VM2 | sopes1vm1|SOPES1VM2|
|VM3 | sopes1vm1|SOPES1VM3|

</div>

5. Reinciar las VM e ingresar a UBUNTU server para descargar las herramientas necesarias

### 1.2.1 Configuración VM1 y VM2

1. Actualizar los paquetes del sistema
```bash
sudo apt update && sudo apt upgrade -y
```

2. Instalar Containerd y configurar el servicio para que se mantenga activo todo el tiempo

```bash
sudo apt install containerd -y
sudo systemctl enable containerd 
sudo systemctl start containerd

echo verificar que el servicio esté corriendo

sudo systemctl status containerd
```

3. Instalar Golang

```bash
sudo apt install containerd -y
sudo systemctl enable containerd 
sudo systemctl start containerd

echo verificar que el servicio esté corriendo

sudo systemctl status containerd
```
4. Instalar nerdctl para poder crear las imagenes de las API's realizadas en las maquinas virtuales

```bash
wget https://github.com/containerd/nerdctl/releases/download/v1.7.7/nerdctl-full-1.7.7-linux-amd64.tar.gz

echo Instalar en /usr/local
sudo tar Cxzvf /usr/local nerdctl-full-1.7.7-linux-amd64.tar.gz

echo Iniciar Buildkit en modo demonio
sudo systemctl enable --now buildkit || true
```

5. Habilitar e iniciar los servicios de Buildkit para que siempre se encuentren disponibles


```bash
sudo systemctl enable --now buildkit || sudo systemctl enable --now buildkitd
sudo systemctl status buildkit || sudo systemctl status buildkitd
```

6. Instalar git e iniciar sesión en github para tener un control y un historial del codigo realizado (Recomendado utilziar estrategias de branching)

```bash
git config --global user.name "Tu Nombre"
git config --global user.email "tu_correo@example.com"

echo autenticación con github (generación de una llave SSH)
ssh-keygen -t ed25519 -C "tu_correo@example.com"

echo copiar la llave publica y colocarla en la pagina de github
cat ~/.ssh/id_ed25519.pub

echo clonar repositorio de 
cd ~
git clone git@github.com:TU_USUARIO/REPOSITORIO.git
cd REPOSITORIO
```

7. (EXTRA) Instalar OPEN SSH para programar con mayor comodidad

```bash
sudo apt install -y openssh-server
sudo systemctl enable --now ssh
sudo systemctl status ssh
hostname -I
```
La IP indicada se debe de colocar en VSCODE, se tiene que tener la extensión Remote-SSH instalada previamente


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
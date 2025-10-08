#!/bin/bash
# filepath: /home/juan-pablo/Escritorio/Uni/SOPES1/202109705_LAB_SO1_2S2025/Proyecto2/bash/remove_cronjob.sh

CRONJOB="* * * * * /bin/bash /home/juan-pablo/Escritorio/Uni/SOPES1/202109705_LAB_SO1_2S2025/Proyecto2/bash/generate_container.sh"
crontab -l | grep -v -F "$CRONJOB" | crontab -
echo "Cronjob eliminado."
#!/bin/bash
# filepath: /home/juan-pablo/Escritorio/Uni/SOPES1/202109705_LAB_SO1_2S2025/Proyecto2/bash/setup_cronjob.sh

CRONJOB="* * * * * /bin/bash /home/juan-pablo/Escritorio/Uni/SOPES1/202109705_LAB_SO1_2S2025/Proyecto2/bash/generate_container.sh >> /tmp/cron_generate_container.log 2>&1"
crontab -l 2>/dev/null | grep -F "$CRONJOB" >/dev/null 2>&1

if [ $? -eq 0 ]; then
    echo "El cronjob ya existe."
else
    (crontab -l 2>/dev/null; echo "$CRONJOB") | crontab -
    echo "Cronjob agregado."
fi
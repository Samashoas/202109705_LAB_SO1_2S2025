from locust import HttpUser, task, between
import json
import random

# Lista de municipios
municipalities = ['mixco', 'guatemala', 'amatitlan', 'chinautla']

# Lista de climas
weathers = ['sunny', 'cloudy', 'rainy', 'foggy']

class WeatherTweetUser(HttpUser):
    # Definir la frecuencia de las peticiones (entre 1 y 3 segundos)
    wait_time = between(1, 3)

    @task
    def send_tweet(self):
        # Generar los datos de un tweet aleatorio
        data = {
            "municipality": random.choice(municipalities),
            "temperature": random.randint(15, 35),  # temperatura entre 15 y 35 grados
            "humidity": random.randint(30, 90),     # humedad entre 30% y 90%
            "weather": random.choice(weathers)
        }

        # Realizar la petición POST a la API REST de Rust
        self.client.post("/tweets", json=data)


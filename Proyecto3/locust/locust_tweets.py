from locust import HttpUser, task, between
import random
import json

class WeatherTweetUser(HttpUser):
    wait_time = between(1, 3)
    
    municipalities = [
        "Guatemala", "Mixco", "Villa Nueva", "Petapa", "San Juan Sacatepéquez",
        "Villa Canales", "Fraijanes", "Santa Catarina Pinula", "San José Pinula",
        "Amatitlán", "Chinautla", "San Pedro Ayampuc"
    ]
    
    weather_conditions = ["Sunny", "Cloudy", "Rainy", "Stormy", "Foggy"]
    
    @task
    def send_weather_tweet(self):
        tweet_data = {
            "municipality": random.choice(self.municipalities),
            "temperature": random.randint(15, 35),
            "humidity": random.randint(40, 90),
            "weather": random.choice(self.weather_conditions)
        }
        
        response = self.client.post(
            "/weather",
            json=tweet_data,
            headers={"Content-Type": "application/json"}
        )
        
        if response.status_code == 200:
            print(f"Sent: {tweet_data}")
        else:
            print(f"Error {response.status_code}: {response.text}")
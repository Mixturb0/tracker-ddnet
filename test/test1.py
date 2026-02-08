import requests
from bs4 import BeautifulSoup


url = "https://ru.ddnet.org/players/?json2=H20"
response = requests.get(url)
data = response.json()
total_hours = sum(day.get("hours_played") for day in data["activity"])
print (total_hours)

    
import requests
from bs4 import BeautifulSoup

url = "https://ru.ddnet.org/ranks"

headers ={
    "User-Agent":
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:147.0) Gecko/20100101 Firefox/147.0"
}

TopLadder = []

response = requests.get(url, headers=headers)
print(response.status_code)  # Обычно 200 для успеха, 404 для "не найдено" и т.д.

soup = BeautifulSoup(response.text, "lxml")

data = soup.find_all("tr")

for i, row in enumerate(data, 1):     # начинаем счёт с 1
    if i > 30:
        break
    
    rank = row.find("td", class_="rankglobal")
    
    name = row.find("a")
    
    url_player = "https://ru.ddnet.org" + row.find("a").get("href")
    
    points = row.find("td", class_="points")
    
    TopLadder.append ({
        "rank": rank.text,
        "name": name.text,
        "points": points.text,
        "url": url_player,
    })
 
print(TopLadder)
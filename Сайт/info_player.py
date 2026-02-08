import requests
from bs4 import BeautifulSoup

HEADERS = {
    "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:147.0) Gecko/20100101 Firefox/147.0"
}

def format_time(ms: int) -> str:
    """Миллисекунды → mm:ss"""
    if ms <= 0:
        return "—"
    total_sec = ms // 1000
    minutes = total_sec // 60
    seconds = total_sec % 60
    return f"{minutes:02d}:{seconds:02d}"

def get_rank(nick: str) -> str:
    url = f"https://ru.ddnet.org/players/{nick}"
    try:
        response = requests.get(url, headers=HEADERS, timeout=10)
        response.raise_for_status()
        soup = BeautifulSoup(response.text, "lxml")
        rank_string = soup.find("p", class_="pers-result")
        if rank_string:
            rank = rank_string.text.split('.')[0].strip()
            return rank
        return "—"
    except Exception:
        return "—"

def get_player_data(nick: str) -> dict:
    json_url = f"https://ru.ddnet.org/players/?json2={nick}"
    try:
        response = requests.get(json_url, headers=HEADERS, timeout=10)
        response.raise_for_status()
        data = response.json()
    except Exception as e:
        return {"error": f"Не удалось загрузить данные: {str(e)}"}

    # Ранг — отдельный запрос (HTML)
    rank = get_rank(nick)

    # Часы из активности
    hours = sum(day.get("hours_played", 0) for day in data.get("activity", []))

    # Общее количество финишей
    finishes = 0
    for info in data.get("types", {}).values():
        for map_data in info.get("maps", {}).values():
            finishes += map_data.get("finishes", 0)

    # Поинты (из JSON)
    points_data = data.get("points", {})
    points = points_data.get("points", 0)
    total_points = points_data.get("total", 0)  # или "of_total" — проверь в JSON

    # Карты по сложностям
    maps = {
        "Novice": [],
        "Moderate": [],
        "Brutal": [],
        "Insane": []
    }
    types_data = data.get("types", {})
    for diff, info in types_data.items():
        if diff not in maps:
            continue
        for map_name, map_data in info.get("maps", {}).items():
            count_finishes = map_data.get("finishes", 0)
            time_ms = map_data.get("time", 0)
            if count_finishes > 0:  # показываем только с финишами > 0 (можно убрать if)
                one_map = {
                    "name": map_name,
                    "finishes": count_finishes,
                    "time": format_time(time_ms)
                }
                maps[diff].append(one_map)
    # Сортировка по имени
    for diff in maps:
        maps[diff].sort(key=lambda x: x["name"])

    return {
        "nickname": nick,  # или data.get("player", nick)
        "rank": rank,
        "points": f"{int(points):,} / {int(total_points):,}" if total_points else "— / —",
        "hours": hours if hours > 0 else "—",
        "finishes": finishes if finishes > 0 else "—",
        "maps": maps  # ← реальные данные по картам
    }
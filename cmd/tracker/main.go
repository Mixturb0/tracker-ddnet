package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"trackerDDnet/internal/tracker"

	"github.com/PuerkitoBio/goquery"
)

func ServersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	defer func() {
		if err := recover(); err != nil {
			fmt.Println("❌ Паника в серверах:", err)
			http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
		}
	}()

	fmt.Println("⏳ Запрашиваем статус серверов...")

	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://ddnet.org/status/", nil)
	if err != nil {
		return
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/120.0.0.0 Safari/537.36")

	res, err := client.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return
	}

	var servers []tracker.ServerInfo

	// Ищем все блоки div, id которых начинается с "server-"
	doc.Find("div[id^='server-']").Each(func(i int, s *goquery.Selection) {
		idAttr, exists := s.Attr("id")
		if !exists {
			return
		}

		// Вырезаем номер из "server-X"
		idStr := strings.TrimPrefix(idAttr, "server-")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return
		}

		// Фильтруем по нужным тебе диапазонам
		region := ""
		if id >= 0 && id <= 133 {
			region = "EUR"
		} else if id >= 134 && id <= 290 {
			region = "RUS"
		} else {
			return // Остальные сервера (например, 291+) пропускаем
		}

		// Достаем название из <h2>
		title := strings.TrimSpace(s.Find("h2").Text())
		if title == "" {
			return
		}

		var players []tracker.PlayerInfo

		// Перебираем строки таблицы внутри этого сервера
		s.Find("table tr").Each(func(j int, tr *goquery.Selection) {
			tds := tr.Find("td")
			if tds.Length() >= 2 {
				timeStr := strings.TrimSpace(tds.Eq(0).Text())
				nameStr := strings.TrimSpace(tds.Eq(1).Text())

				// Отделяем игроков от зрителей (у игроков время содержит двоеточие "02:14")
				if strings.Contains(timeStr, ":") && nameStr != "" {
					players = append(players, tracker.PlayerInfo{
						Time: timeStr,
						Name: nameStr,
					})
				}
			}
		})

		// Добавляем только если на сервере есть игроки (по желанию можешь убрать это условие)
		if len(players) > 0 {
			servers = append(servers, tracker.ServerInfo{
				ID:      id,
				Region:  region,
				Title:   title,
				Players: players,
			})
		}
	})

	fmt.Printf("✅ Найдено серверов с игроками: %d\n", len(servers))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(servers)
}

// Вспомогательная функция: вытаскивает только числа из любого текста
func extractNumber(s string) int {
	var result string
	for _, r := range s {
		if r >= '0' && r <= '9' {
			result += string(r)
		}
	}
	if result == "" {
		return 0
	}
	num, _ := strconv.Atoi(result)
	return num
}

func LadderHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://ddnet.org/ranks/", nil)
	if err != nil {
		return
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/120.0.0.0 Safari/537.36")

	res, err := client.Do(req)
	if err != nil {
		http.Error(w, "Ошибка сети", http.StatusBadGateway)
		return
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return
	}

	var allPlayers []tracker.LadderPlayer

	// Ищем таблицу
	rows := doc.Find("table.tight tr")

	rows.Each(func(i int, s *goquery.Selection) {
		// Останавливаемся, когда собрали ровно Топ-20
		if len(allPlayers) >= 20 {
			return
		}

		tds := s.Find("td")
		if tds.Length() < 3 {
			return
		}

		// Ранг — всегда в первой колонке
		rank := extractNumber(tds.Eq(0).Text())

		// Очки — ищем по классу, который ты показал на скрине
		pointsText := s.Find("td.points").Text()
		if pointsText == "" {
			pointsText = tds.Eq(1).Text() // запасной вариант
		}
		points := extractNumber(pointsText)

		// Никнейм — ВСЕГДА в самой последней колонке
		name := strings.TrimSpace(tds.Last().Text())

		// Если данные валидны, сохраняем игрока
		if rank > 0 && name != "" {
			allPlayers = append(allPlayers, tracker.LadderPlayer{
				Rank:   rank,
				Name:   name,
				Points: points,
			})
		}
	})

	fmt.Printf("✅ Собрано игроков: %d\n", len(allPlayers))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(allPlayers)
}

func PlayerHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	nickname := r.URL.Query().Get("name")
	if nickname == "" {
		http.Error(w, "name required", http.StatusBadRequest)
		return
	}

	stats, err := tracker.PlayerPars(nickname)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	response := stats.ToPlayerGive()

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	http.HandleFunc("/player", PlayerHandler)
	http.HandleFunc("/ladder", LadderHandler)
	http.HandleFunc("/servers", ServersHandler)

	err := http.ListenAndServe(":9091", nil)
	if err != nil {
		fmt.Println("Ошибка при запуске сервера:", err)
	}
}

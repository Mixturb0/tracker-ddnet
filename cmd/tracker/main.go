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

	err := http.ListenAndServe(":9091", nil)
	if err != nil {
		fmt.Println("Ошибка при запуске сервера:", err)
	}
}

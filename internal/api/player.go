package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"trackerDDnet/internal/tracker"
)

var stats tracker.PlayerStats

func Player(nickname string) string {
	fullUrl := fmt.Sprintf("https://ru.ddnet.org/players?json2=%s", url.QueryEscape(nickname))

	info, err := http.Get(fullUrl)
	if err != nil {
		fmt.Println("Ошибка при получения JSON:", err)
		return err.Error()
	}

	defer info.Body.Close()

	body, err := io.ReadAll(info.Body)
	if err != nil {
		fmt.Println("Ошибка при получения body:", err)
		return err.Error()
	}

	err = json.Unmarshal(body, &stats)
	if err != nil {
		fmt.Println("Ошибка при декодировании JSON:", err)
		return err.Error()
	}

	return stats.Player

}

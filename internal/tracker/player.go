package tracker

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func PlayerPars(nickname string) (*PlayerStats, error) {
	fullUrl := fmt.Sprintf("https://ru.ddnet.org/players?json2=%s", url.QueryEscape(nickname))

	info, err := http.Get(fullUrl)
	if err != nil {
		fmt.Println("Ошибка при получения JSON:", err)
		return nil, err
	}

	defer info.Body.Close()

	body, err := io.ReadAll(info.Body)
	if err != nil {
		fmt.Println("Ошибка при получения body:", err)
		return nil, err
	}

	var stats PlayerStats
	err = json.Unmarshal(body, &stats)
	if err != nil {
		fmt.Println("Ошибка при декодировании JSON:", err)
		return nil, err
	}
	return &stats, err
}

func (p *PlayerStats) ToPlayerGive() PlayerGive {
	return PlayerGive{
		Player:     p.Player,
		Points:     p.Points,
		Finish:     p.TotalFinishes(),
		TotalHours: p.TotalHourse(),
	}
}

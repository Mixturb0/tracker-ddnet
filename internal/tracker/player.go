package tracker

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func PlayerPars(nickname string) error {
	fullUrl := fmt.Sprintf("https://ru.ddnet.org/players?json2=%s", url.QueryEscape(nickname))

	info, err := http.Get(fullUrl)
	if err != nil {
		fmt.Println("Ошибка при получения JSON:", err)
		return err
	}

	defer info.Body.Close()

	body, err := io.ReadAll(info.Body)
	if err != nil {
		fmt.Println("Ошибка при получения body:", err)
		return err
	}

	err = json.Unmarshal(body, PlayerGive{})
	if err != nil {
		fmt.Println("Ошибка при декодировании JSON:", err)
		return err
	}
	return nil
}
func (p *PlayerStats) PlayerGet() PlayerGive {
	return PlayerGive{
		Player: p.Player,
		Points: p.Points,
		Rank:   p.Rank,
		Finish: p.TotalFinishes(),
	}
}

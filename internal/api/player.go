package api

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func Player(nickname string) string {
	baseUrl := "https://ru.ddnet.org/players/?json2="
	fullUrl, err := url.JoinPath(baseUrl, nickname)
	if err != nil {
		return err.Error()
	}
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

	bodyString := string(body)

	return bodyString

}

package main

import (
	"fmt"
	"net/http"
	"trackerDDnet/internal/api"
)

func PlayerHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	nickname := r.URL.Query().Get("name")
	JsonResponse := api.Player(nickname)
	fmt.Fprint(w, JsonResponse)

}

func main() {
	http.HandleFunc("/player", PlayerHandler)

	err := http.ListenAndServe(":9091", nil)
	if err != nil {
		fmt.Println("Ошибка при запуске сервера:", err)
	}
}

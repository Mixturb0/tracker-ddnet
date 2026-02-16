package main

import (
	"fmt"
	"net/http"
	"trackerDDnet/internal/tracker"
)

func PlayerHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	nickname := r.URL.Query().Get("name")
	JsonResponse := tracker.PlayerPars(nickname)
	if _, err := fmt.Fprint(w, JsonResponse); err != nil {
		fmt.Println("error writing response", err.Error())
	}

}

func main() {
	http.HandleFunc("/player", PlayerHandler)

	err := http.ListenAndServe(":9091", nil)
	if err != nil {
		fmt.Println("Ошибка при запуске сервера:", err)
	}
}

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"trackerDDnet/internal/tracker"
)

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

	err := http.ListenAndServe(":9091", nil)
	if err != nil {
		fmt.Println("Ошибка при запуске сервера:", err)
	}
}

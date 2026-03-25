package handler

import (
	"encoding/json"

	"net/http"
)

func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed!", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	//записываем ответ
	json.NewEncoder(w).Encode("")

}

package handlers

import (
	"encoding/json"
	"go-basic-rest-api/models"
	"net/http"
	"fmt"

	"github.com/jinzhu/gorm"
)

type HandlerRoute func(*gorm.DB, http.ResponseWriter, *http.Request)

func respondJSON(w http.ResponseWriter, status int, payload interface{models.Serialize}) {
	response, err := payload.Serialize()

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))

		return
	}

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.WriteHeader(status)
	w.Write([]byte(response))
}

func ProcessJSON(w http.ResponseWriter, status int, response []byte) {
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.WriteHeader(status)
	w.Write([]byte(response))
}

func respondError(w http.ResponseWriter, code int, message string) {
	response, err := json.Marshal(map[string]string{"error": message})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))

		return
	}

	ProcessJSON(w, code, response)
}

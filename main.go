package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func GetMessages(w http.ResponseWriter, r *http.Request) {
	var messages []Message
	if result := DB.Find(&messages); result.Error != nil {
		http.Error(w, "Не удалось получить сообщения", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}

func CreateMessage(w http.ResponseWriter, r *http.Request) {
	var requestBody Message
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Невозможно декодировать JSON", http.StatusBadRequest)
		return
	}

	if result := DB.Create(&requestBody); result.Error != nil {
		http.Error(w, "Не удалось сохранить сообщение", http.StatusInternalServerError)
		return
	}

	// Устанавливаем заголовок ответа как JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	// Формируем ответ
	response := map[string]interface{}{
		"message": "Сообщение успешно создано",
		"data":    requestBody,
	}

	// Кодируем ответ в JSON
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Не удалось закодировать ответ в JSON", http.StatusInternalServerError)
		return
	}
}

func main() {
	// Вызываем метод InitDB() из файла db.go
	InitDB()

	// Автоматическая миграция модели Message 
	DB.AutoMigrate(&Message{})

	router := mux.NewRouter()
	router.HandleFunc("/api/messages", CreateMessage).Methods("POST")
	router.HandleFunc("/api/messages", GetMessages).Methods("GET")
	fmt.Println("Сервер запущен")
	http.ListenAndServe(":8080", router)
}

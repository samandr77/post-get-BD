package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Функция для обновления сообщения
func UpdateMessage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var message Message
	if err := DB.First(&message, id).Error; err != nil {
		http.Error(w, "Message not found", http.StatusNotFound)
		return
	}

	var updatedData struct {
		Text string `json:"text"`
	}

	if err := json.NewDecoder(r.Body).Decode(&updatedData); err != nil {
		http.Error(w, "Failed to decode JSON", http.StatusBadRequest)
		return
	}

	message.Text = updatedData.Text

	if err := DB.Save(&message).Error; err != nil {
		http.Error(w, "Failed to update message", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Функция для получения всех сообщений
func GetMessages(w http.ResponseWriter, r *http.Request) {
	var messages []Message
	if result := DB.Find(&messages); result.Error != nil {
		http.Error(w, "Не удалось получить сообщения", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}

// Функция для создания нового сообщения
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	response := map[string]interface{}{
		"message": "Сообщение успешно создано",
		"data":    requestBody,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Не удалось закодировать ответ в JSON", http.StatusInternalServerError)
		return
	}
}

func main() {
	InitDB()

	// Автоматическая миграция модели Message
	DB.AutoMigrate(&Message{})

	router := mux.NewRouter()
	router.HandleFunc("/api/messages", CreateMessage).Methods("POST")
	router.HandleFunc("/api/messages", GetMessages).Methods("GET")
	router.HandleFunc("/api/messages/{id:[0-9]+}", UpdateMessage).Methods("PATCH")


	fmt.Println("Сервер запущен на порту 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

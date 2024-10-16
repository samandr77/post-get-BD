package main 

import (
	"gorm.io/gorm"
)

type Message struct {
    ID      uint   `gorm:"primaryKey"`
    Content string `json:"content"`

	gorm.Model
	Text string `json:"text"` // Наш сервер будет ожидать JSON с полем text
}


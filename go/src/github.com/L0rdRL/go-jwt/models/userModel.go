package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Login         string
	Name          string
	Lastname      string
	Surname       string
	Position      string
	Password      string
	Role          string
	Verified      bool
	AllowedToDocs bool
	Documents     []Document // Добавляем связь: пользователь может иметь несколько документов
}

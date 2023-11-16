package initializers

import (
	"github.com/models"
)

func SyncDataBase() {
	DB.AutoMigrate(&models.User{})
}

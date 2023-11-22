package initializers

import "notes_application/models"

func InitialMigrations(){
	DB.AutoMigrate(&models.User{}, &models.Note{})
}
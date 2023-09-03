package main

import (
	"github.com/notalent2code/task-5-vix-btpns-rihlan/database"
	"github.com/notalent2code/task-5-vix-btpns-rihlan/models"
	"github.com/notalent2code/task-5-vix-btpns-rihlan/router"
)

func main() {
	db := database.SetupDB()
	db.AutoMigrate(&models.User{})

	r := router.SetupRoutes(db)
	r.Run()
}

package main

import (
	"github.com/tantowish/task-5-pbi-btpns-tantowishahhanif/database"
	"github.com/tantowish/task-5-pbi-btpns-tantowishahhanif/models"
	"github.com/tantowish/task-5-pbi-btpns-tantowishahhanif/router"
)

func main() {
	db := database.SetupDB()
	db.AutoMigrate(&models.User{})

	r := router.SetupRoutes(db)
	r.Run()
}

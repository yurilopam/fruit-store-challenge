package main

import (
	"github.com/yurilopam/fruit-store-challenge/internal/database"
	"github.com/yurilopam/fruit-store-challenge/internal/routes"
)

func main() {
	database.Connect()
	routes.HandleRequests()
}

package routes

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/yurilopam/fruit-store-challenge/internal/controllers"
)

func HandleRequests() {
	r := gin.Default()

	r.GET("/fruits", controllers.GetFruits)
	r.GET("/fruits/:id", controllers.GetFruit)
	r.POST("/fruits", controllers.CreateFruit)
	r.PUT("/fruits/:id", controllers.UpdateFruit)
	r.DELETE("/fruits/:id", controllers.DeleteFruit)

	r.GET("/users", controllers.GetUsers)
	r.GET("/users/:id", controllers.GetUser)
	r.POST("/users", controllers.CreateUser)
	r.PUT("/users/:id", controllers.UpdateUser)
	r.DELETE("/users/:id", controllers.DeleteUser)

	if err := r.Run(); err != nil {
		log.Fatal(err)
	}
}

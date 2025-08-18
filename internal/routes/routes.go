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

	if err := r.Run(); err != nil {
		log.Fatal(err)
	}
}

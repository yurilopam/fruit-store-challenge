package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yurilopam/fruit-store-challenge/internal/database"
	"github.com/yurilopam/fruit-store-challenge/internal/models"
)

func GetFruits(c *gin.Context) {
	var fruits []models.Fruit
	database.DB.Find(&fruits)
	c.JSON(http.StatusOK, fruits)
}

func GetFruit(c *gin.Context) {
	id := c.Param("id")
	var fruit models.Fruit

	if err := database.DB.First(&fruit, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Fruit not found"})
		return
	}

	c.JSON(http.StatusOK, fruit)
}

func CreateFruit(c *gin.Context) {
	var fruit models.Fruit
	if err := c.ShouldBindJSON(&fruit); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	database.DB.Create(&fruit)
	c.JSON(http.StatusCreated, fruit)
}

func UpdateFruit(c *gin.Context) {
	id := c.Param("id")
	var fruit models.Fruit

	if err := database.DB.First(&fruit, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Fruit not found"})
		return
	}

	var input models.Fruit
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	database.DB.Model(&fruit).Updates(input)
	c.JSON(http.StatusOK, fruit)
}

func DeleteFruit(c *gin.Context) {
	id := c.Param("id")
	var fruit models.Fruit

	if err := database.DB.First(&fruit, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Fruit not found"})
		return
	}

	database.DB.Delete(&fruit)
	c.JSON(http.StatusOK, gin.H{"message": "Fruit deleted"})
}

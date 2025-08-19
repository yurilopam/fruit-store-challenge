package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yurilopam/fruit-store-challenge/internal/app"
	"github.com/yurilopam/fruit-store-challenge/internal/models"
	"github.com/yurilopam/fruit-store-challenge/internal/util"
)

type FruitReq struct {
	Name     string  `json:"name" binding:"required"`
	Price    float64 `json:"price" binding:"required"`
	Quantity int     `json:"quantity" binding:"required"`
}

type FruitsRoutes struct{ App *app.App }

func NewFruitsRoutes(a *app.App) *FruitsRoutes { return &FruitsRoutes{App: a} }

func (f *FruitsRoutes) Create(c *gin.Context) {
	var r FruitReq
	if err := c.ShouldBindJSON(&r); err != nil {
		util.BadRequest(c, err.Error())
		return
	}
	fr := models.Fruit{Name: r.Name, Price: r.Price, Quantity: r.Quantity}
	if err := f.App.DB.Create(&fr).Error; err != nil {
		util.BadRequest(c, err.Error())
		return
	}
	f.App.InvalidateFruits(c.Request.Context())
	util.JSON(c, http.StatusCreated, fr)
}

func (f *FruitsRoutes) List(c *gin.Context) {
	if cached, ok := f.App.GetCachedFruits(c.Request.Context()); ok {
		util.JSON(c, http.StatusOK, gin.H{"cached": true, "data": cached})
		return
	}
	var fruits []models.Fruit
	if err := f.App.DB.Find(&fruits).Error; err != nil {
		util.BadRequest(c, err.Error())
		return
	}
	f.App.CacheFruits(c.Request.Context(), fruits)
	util.JSON(c, http.StatusOK, gin.H{"cached": false, "data": fruits})
}

func (f *FruitsRoutes) Get(c *gin.Context) {
	id := c.Param("id")
	var fr models.Fruit
	if err := f.App.DB.Where("id = ?", id).First(&fr).Error; err != nil {
		util.NotFound(c)
		return
	}
	util.JSON(c, http.StatusOK, fr)
}

func (f *FruitsRoutes) Update(c *gin.Context) {
	id := c.Param("id")
	var r FruitReq
	if err := c.ShouldBindJSON(&r); err != nil {
		util.BadRequest(c, err.Error())
		return
	}
	var fr models.Fruit
	if err := f.App.DB.Where("id = ?", id).First(&fr).Error; err != nil {
		util.NotFound(c)
		return
	}
	fr.Name, fr.Price, fr.Quantity = r.Name, r.Price, r.Quantity
	if err := f.App.DB.Save(&fr).Error; err != nil {
		util.BadRequest(c, err.Error())
		return
	}
	f.App.InvalidateFruits(c.Request.Context())
	util.JSON(c, http.StatusOK, fr)
}

func (f *FruitsRoutes) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := f.App.DB.Delete(&models.Fruit{}, "id = ?", id).Error; err != nil {
		util.BadRequest(c, err.Error())
		return
	}
	f.App.InvalidateFruits(c.Request.Context())
	c.Status(http.StatusNoContent)
}

package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yurilopam/fruit-store-challenge/internal/models"
	"github.com/yurilopam/fruit-store-challenge/internal/services"
	"github.com/yurilopam/fruit-store-challenge/internal/util"
	"gorm.io/gorm"
)

type LoginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AuthRoutes struct {
	DB  *gorm.DB
	JWT *services.JWTService
}

func NewAuthRoutes(db *gorm.DB, jwt *services.JWTService) *AuthRoutes {
	return &AuthRoutes{DB: db, JWT: jwt}
}

func (a *AuthRoutes) Login(c *gin.Context) {
	var r LoginReq
	if err := c.ShouldBindJSON(&r); err != nil {
		util.BadRequest(c, err.Error())
		return
	}
	var u models.User
	if err := a.DB.Where("username = ?", r.Username).First(&u).Error; err != nil {
		util.Unauthorized(c)
		return
	}
	if !services.CheckPassword(u.Password, r.Password) {
		util.Unauthorized(c)
		return
	}
	tok, _ := a.JWT.Generate(u.ID.String(), u.Role)
	util.JSON(c, http.StatusOK, gin.H{"token": tok})
}

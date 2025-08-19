package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yurilopam/fruit-store-challenge/internal/app"
	"github.com/yurilopam/fruit-store-challenge/internal/models"
	"github.com/yurilopam/fruit-store-challenge/internal/services"
	"github.com/yurilopam/fruit-store-challenge/internal/util"
)

type UserReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Role     string `json:"role" binding:"required,oneof=admin user"`
}

type UsersRoutes struct {
	App *app.App
}

func NewUsersRoutes(a *app.App) *UsersRoutes { return &UsersRoutes{App: a} }

func (u *UsersRoutes) Create(c *gin.Context) {
	var r UserReq
	if err := c.ShouldBindJSON(&r); err != nil {
		util.BadRequest(c, err.Error())
		return
	}
	hash, _ := services.HashPassword(r.Password)
	user := models.User{Username: r.Username, Password: hash, Role: r.Role}
	if err := u.App.DB.Create(&user).Error; err != nil {
		util.BadRequest(c, err.Error())
		return
	}
	// Pub event in Kafka
	if u.App.Kafka != nil {
		ctx, cancel := app.Ctx()
		defer cancel()
		_ = u.App.Kafka.PublishUser(ctx, services.UserEvent{ID: user.ID.String(), Username: user.Username, Role: user.Role})
	}
	util.JSON(c, http.StatusCreated, gin.H{"id": user.ID, "username": user.Username, "role": user.Role})
}

func (u *UsersRoutes) List(c *gin.Context) {
	var users []models.User
	if err := u.App.DB.Select("id, username, role, created_at, updated_at").Find(&users).Error; err != nil {
		util.BadRequest(c, err.Error())
		return
	}
	util.JSON(c, http.StatusOK, users)
}

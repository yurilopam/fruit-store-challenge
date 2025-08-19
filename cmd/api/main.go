package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/yurilopam/fruit-store-challenge/internal/app"
	"github.com/yurilopam/fruit-store-challenge/internal/config"
	"github.com/yurilopam/fruit-store-challenge/internal/database"
	"github.com/yurilopam/fruit-store-challenge/internal/middleware"
	"github.com/yurilopam/fruit-store-challenge/internal/models"
	"github.com/yurilopam/fruit-store-challenge/internal/routes"
	"github.com/yurilopam/fruit-store-challenge/internal/services"
)

func main() {
	_ = godotenv.Load()
	cfg := config.MustLoad()

	// DB
	dsn := database.DSN(cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)
	db, err := database.NewPostgres(dsn)
	if err != nil {
		log.Fatal(err)
	}
	if err := db.AutoMigrate(&models.User{}, &models.Fruit{}); err != nil {
		log.Fatal(err)
	}

	// Seed admin (if not exists)
	var count int64
	db.Model(&models.User{}).Where("username = ?", cfg.SeedAdminUsername).Count(&count)
	if count == 0 {
		h, _ := services.HashPassword(cfg.SeedAdminPassword)
		db.Create(&models.User{Username: cfg.SeedAdminUsername, Password: h, Role: "admin"})
		log.Println("Admin user seeded")
	}

	// Redis
	rds := services.NewRedis(cfg.RedisAddr, cfg.RedisDB, cfg.RedisPassword, services.WithRedisTTLSeconds(cfg.RedisTTLSeconds))

	// JWT
	jwtSvc := services.NewJWT(cfg.JWTSecret)

	// Kafka producer
	producer := services.NewKafkaProducer(cfg.KafkaBrokers, cfg.KafkaTopicUsers)

	// App (DI via Functional Options)
	application := app.New(
		app.WithDB(db),
		app.WithRedis(rds),
		app.WithJWT(jwtSvc),
		app.WithKafka(producer))

	r := gin.Default()

	// Healthcheck
	r.GET("/health", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"status": "ok"}) })

	// Auth
	authRoutes := routes.NewAuthRoutes(db, jwtSvc)
	r.POST("/auth/login", authRoutes.Login)

	// Protected
	authMw := middleware.AuthRequired(middleware.AuthConfig{Secret: jwtSvc.Secret()})

	// Users (admin only)
	usersRoutes := routes.NewUsersRoutes(application)
	users := r.Group("/users", authMw, middleware.RequireRole("admin"))
	{
		users.POST("", usersRoutes.Create)
		users.GET("", usersRoutes.List)
	}

	// Fruits
	fruitRoutes := routes.NewFruitsRoutes(application)
	fr := r.Group("/fruits", authMw)
	{
		fr.GET("", fruitRoutes.List)
		fr.GET(":id", fruitRoutes.Get)
	}
	// Admin-only modifications
	frAdmin := r.Group("/fruits", authMw, middleware.RequireRole("admin"))
	{
		frAdmin.POST("", fruitRoutes.Create)
		frAdmin.PUT(":id", fruitRoutes.Update)
		frAdmin.DELETE(":id", fruitRoutes.Delete)
	}

	log.Println("API listening on :" + cfg.Port)
	r.Run(":" + cfg.Port)
}

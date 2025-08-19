package app

import (
	"context"
	"encoding/json"
	"time"

	"github.com/yurilopam/fruit-store-challenge/internal/models"
	"github.com/yurilopam/fruit-store-challenge/internal/services"
	"gorm.io/gorm"
)

type App struct {
	DB    *gorm.DB
	Redis *services.RedisService
	JWT   *services.JWTService
	Kafka *services.KafkaProducer
}

type Option func(*App)

func WithDB(db *gorm.DB) Option                  { return func(a *App) { a.DB = db } }
func WithRedis(r *services.RedisService) Option  { return func(a *App) { a.Redis = r } }
func WithJWT(j *services.JWTService) Option      { return func(a *App) { a.JWT = j } }
func WithKafka(k *services.KafkaProducer) Option { return func(a *App) { a.Kafka = k } }

func New(opts ...Option) *App {
	app := &App{}
	for _, opt := range opts {
		opt(app)
	}
	return app
}

func (a *App) CacheFruits(ctx context.Context, fruits []models.Fruit) {
	if a.Redis == nil {
		return
	}
	b, _ := json.Marshal(fruits)
	_ = a.Redis.Set(ctx, "fruits:all", string(b))
}

func (a *App) InvalidateFruits(ctx context.Context) {
	if a.Redis != nil {
		_ = a.Redis.Del(ctx, "fruits:all")
	}
}

func (a *App) GetCachedFruits(ctx context.Context) ([]models.Fruit, bool) {
	if a.Redis == nil {
		return nil, false
	}
	s, err := a.Redis.Get(ctx, "fruits:all")
	if err != nil {
		return nil, false
	}
	var out []models.Fruit
	if err := json.Unmarshal([]byte(s), &out); err != nil {
		return nil, false
	}
	return out, true
}

func Ctx() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 5*time.Second)
}

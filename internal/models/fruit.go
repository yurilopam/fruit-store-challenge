package models

import "gorm.io/gorm"

type Fruit struct {
	gorm.Model

	Name          string  `json:"name"`
	Price         float64 `json:"price"`
	StockQuantity int     `json:"stock_quantity"`
}

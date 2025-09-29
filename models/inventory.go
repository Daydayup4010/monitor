package models

import (
	"time"
	"uu/config"
)

type Inventory struct {
	Id        int64     `json:"id" gorm:"type:int;primaryKey"`
	Name      string    `json:"name"  gorm:"type:varchar(255)"`
	BuffPrice float64   `json:"buff_price" gorm:"type:double"`
	UPrice    float64   `json:"u_price" gorm:"type:double"`
	Category  string    `json:"category" gorm:"type:varchar(20)"`
	ImageUrl  string    `json:"image_url" gorm:"type:varchar(255)"`
	UpdatedAt time.Time `json:"updated_at" gorm:"type:datetime"`
}

// todo
func UpdateInventory() {
	config.DB.Model(&UItemsInfo{}).Select("u.name, u.image_url, u.cache")
}

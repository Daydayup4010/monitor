package models

import (
	"time"
	"uu/config"
)

type SkinItem struct {
	ID          string    `json:"id" gorm:"type:varchar(20);primaryKey"`
	Name        string    `json:"name"  gorm:"type:varchar(255)"`
	BuffPrice   float64   `json:"buff_price" gorm:"type:double"`
	YoupinPrice float64   `json:"youpin_price" gorm:"type:double"`
	BuffUrl     string    `json:"buff_url" gorm:"type:varchar(255)"`
	YoupinUrl   string    `json:"youpin_url" gorm:"type:varchar(255)"`
	Category    string    `json:"category" gorm:"type:varchar(20)"`
	ImageUrl    string    `json:"image_url" gorm:"type:varchar(255)"`
	LastUpdated time.Time `json:"last_updated" gorm:"type:datetime"`
	PriceDiff   float64   `json:"price_diff" gorm:"type:double"`
	ProfitRate  float64   `json:"profit_rate" gorm:"type:double"`
}

func GetSkinItems(pageSize, pageNum int) ([]SkinItem, int64) {
	var skins []SkinItem
	var total int64
	config.DB.Limit(pageSize).Offset((pageNum - 1) * pageSize).Find(&skins).Count(&total)
	return skins, total
}

package models

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"uu/config"
)

type BuffItem struct {
	MarketHashName string    `json:"market_hash_name" gorm:"type:varchar(255);uniqueIndex"`
	Name           string    `json:"name" gorm:"type:varchar(255)"`
	ID             int64     `json:"id" gorm:"primaryKey"`
	SellMinPrice   string    `json:"sell_min_price" gorm:"type:decimal(10,2)"`
	SellNum        int64     `json:"sell_num" gorm:"type:int"`
	GoodsInfo      GoodsInfo `json:"goods_info" gorm:"foreignKey:BuffItemID;references:ID"`
}

type GoodsInfo struct {
	gorm.Model
	IconURL         string `json:"icon_url" gorm:"type:varchar(255)"`
	BuffItemID      int64  // 外键，指向 BuffItem 的 ID
	SteamPriceCNY   string `json:"steam_price_cny" gorm:"type:decimal(10,2)"`
	OriginalIconURL string `json:"original_icon_url" gorm:"type:text"`
}

type BuffInventory struct {
	MarketHashName string `json:"market_hash_name" gorm:"type:varchar(255);uniqueIndex"`
	Name           string `json:"name" gorm:"type:varchar(255)"`
	ID             int64  `json:"goods_id" gorm:"primaryKey"`
	SellMinPrice   string `json:"sell_min_price" gorm:"type:decimal(10,2)"`
}

func BatchAddBuffItem(buff []*BuffItem) {
	err := config.DB.Clauses(clause.OnConflict{UpdateAll: true}).CreateInBatches(buff, 80).Error
	if err != nil {
		config.Log.Errorf("batch insert buff item error: %s", err)
	}
}

func BatchAddBuffInventory(buff []*BuffInventory) {
	err := config.DB.Clauses(clause.OnConflict{UpdateAll: true}).CreateInBatches(buff, 50).Error
	if err != nil {
		config.Log.Errorf("insert buff inventory error: %s", err)
	}
}

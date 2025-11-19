package models

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"uu/config"
)

type C5 struct {
	Id             string  `json:"platformItemId" gorm:"primaryKey"`
	MarketHashName string  `json:"marketHashName" gorm:"type:varchar(255);uniqueIndex;not null"`
	SellPrice      float64 `json:"sellPrice"`
	SellCount      int64   `json:"sellCount"`
	BiddingPrice   float64 `json:"biddingPrice"`
	BiddingCount   int64   `json:"biddingCount"`
	UpdateTime     int64   `json:"updateTime"`
	TurnOver       int64   `json:"turn_over"`
	Link           string  `json:"link"`
}

func GetC5Goods(hashName string) *C5 {
	var c5 C5
	config.DB.Where("market_hash_name = ?", hashName).Find(&c5)
	return &c5
}

func BatchUpdateC5Goods(c5 []*C5) {
	err := config.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.OnConflict{UpdateAll: true}).CreateInBatches(c5, 100).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		config.Log.Errorf("Update C5 Goods fail: %v", err)
		return
	}
	config.Log.Info("Update C5 Goods Success")
}

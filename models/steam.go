package models

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"uu/config"
)

type Steam struct {
	gorm.Model
	Id             uint    `json:"id" gorm:"primaryKey"`
	MarketHashName string  `json:"marketHashName" gorm:"type:varchar(255);uniqueIndex;not null"`
	SellPrice      float64 `json:"sellPrice"`
	SellCount      int64   `json:"sellCount"`
	BiddingPrice   float64 `json:"biddingPrice"`
	BiddingCount   int64   `json:"biddingCount"`
	UpdateTime     int64   `json:"updateTime"`
	TurnOver       int64   `json:"turn_over"`
	Link           string  `json:"link"`
}

func BatchGetSteamGoods(hashNames []string) map[string]*Steam {
	var steams []Steam
	result := make(map[string]*Steam)
	config.DB.Where("market_hash_name in ?", hashNames).Find(&steams)
	for i := range steams {
		result[steams[i].MarketHashName] = &steams[i]
	}
	return result
}

func BatchUpdateSteamGoods(steam []*Steam) {
	err := config.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "market_hash_name"}},
			DoUpdates: clause.AssignmentColumns([]string{"sell_price", "sell_count", "bidding_price", "bidding_count", "update_time", "turn_over", "link"}),
		}).CreateInBatches(steam, 100).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		config.Log.Errorf("Update steam Goods fail: %v", err)
		return
	}
	config.Log.Info("Update steam Goods Success")
}

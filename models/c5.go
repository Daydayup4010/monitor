package models

import (
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"uu/config"
)

type C5 struct {
	Id             string  `json:"platformItemId" gorm:"primaryKey"`
	MarketHashName string  `json:"marketHashName" gorm:"type:varchar(255);uniqueIndex;not null"`
	SellPrice      float64 `json:"sellPrice" gorm:"index"`
	SellCount      int64   `json:"sellCount" gorm:"index"`
	BiddingPrice   float64 `json:"biddingPrice" gorm:"index"`
	BiddingCount   int64   `json:"biddingCount"`
	BeforeTime     int64   `json:"beforeTime"`
	BeforeCount    int64   `json:"beforeCount"`
	UpdateTime     int64   `json:"updateTime"`
	TurnOver       int64   `json:"turn_over"`
	Link           string  `json:"link"`
}

func BatchGetC5Goods(hashNames []string) map[string]*C5 {
	var c5s []C5
	result := make(map[string]*C5)
	config.DB.Where("market_hash_name in ?", hashNames).Find(&c5s)
	for i := range c5s {
		result[c5s[i].MarketHashName] = &c5s[i]
	}
	return result
}

func BatchUpdateC5Goods(c5 []*C5) {
	maxRetries := 3
	var err error
	for i := 0; i < maxRetries; i++ {
		err = config.DB.Transaction(func(tx *gorm.DB) error {
			if err := tx.Clauses(clause.OnConflict{UpdateAll: true}).CreateInBatches(c5, 50).Error; err != nil {
				return err
			}
			return nil
		})
		if err == nil {
			return
		}
		if strings.Contains(err.Error(), "Deadlock") || strings.Contains(err.Error(), "SAVEPOINT") {
			config.Log.Warnf("Update C5 Goods deadlock, retrying (%d/%d)...", i+1, maxRetries)
			time.Sleep(time.Millisecond * time.Duration(100*(i+1)))
			continue
		}
		break
	}
	if err != nil {
		config.Log.Errorf("Update C5 Goods fail: %v", err)
	}
}

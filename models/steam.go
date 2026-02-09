package models

import (
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"uu/config"
)

type Steam struct {
	MarketHashName string  `json:"marketHashName" gorm:"type:varchar(255);primaryKey"`
	Id             string  `json:"platformItemId" gorm:"type:varchar(50)"`
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
	maxRetries := 3
	var err error
	for i := 0; i < maxRetries; i++ {
		err = config.DB.Transaction(func(tx *gorm.DB) error {
			if err := tx.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "market_hash_name"}},
				DoUpdates: clause.AssignmentColumns([]string{"id", "sell_price", "sell_count", "bidding_price", "bidding_count", "update_time", "turn_over", "link"}),
			}).CreateInBatches(steam, 50).Error; err != nil {
				return err
			}
			return nil
		})
		if err == nil {
			return
		}
		if strings.Contains(err.Error(), "Deadlock") || strings.Contains(err.Error(), "SAVEPOINT") {
			config.Log.Warnf("Update Steam Goods deadlock, retrying (%d/%d)...", i+1, maxRetries)
			time.Sleep(time.Millisecond * time.Duration(100*(i+1)))
			continue
		}
		break
	}
	if err != nil {
		config.Log.Errorf("Update Steam Goods fail: %v", err)
	}
}

// GetSteamsWithoutItemNameId 获取所有没有 item_nameid 的商品
func GetSteamsWithoutItemNameId() ([]Steam, error) {
	var steams []Steam
	err := config.DB.Where("(id = '' OR id IS NULL) AND sell_price < 500").Find(&steams).Error
	return steams, err
}

// BatchUpdateSteamItemNameIds 批量更新 item_nameid
func BatchUpdateSteamItemNameIds(updates map[string]string) (int, error) {
	if len(updates) == 0 {
		return 0, nil
	}

	// 使用事务批量更新
	tx := config.DB.Begin()
	count := 0
	for marketHashName, itemNameId := range updates {
		if err := tx.Model(&Steam{}).Where("market_hash_name = ?", marketHashName).Update("id", itemNameId).Error; err != nil {
			tx.Rollback()
			return count, err
		}
		count++
	}
	tx.Commit()
	return count, nil
}

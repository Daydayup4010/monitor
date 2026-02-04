package models

import (
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"uu/config"
)

type BuffItem struct {
	HashName string `json:"market_hash_name" gorm:"type:varchar(255);uniqueIndex"`
	Name     string `json:"name" gorm:"type:varchar(255)"`
	ID       int64  `json:"id" gorm:"primaryKey"`
	Price    string `json:"sell_min_price" gorm:"type:decimal(10,2)"`
	Count    int64  `json:"sell_num" gorm:"type:int"`
	//GoodsInfo      GoodsInfo `json:"goods_info" gorm:"foreignKey:BuffItemID;references:ID"`
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

// -------------------------------------------------------v2------------------------------------------------------------
// data from steamDT

type Buff struct {
	Id             string  `json:"platformItemId" gorm:"primaryKey"`
	MarketHashName string  `json:"marketHashName" gorm:"type:varchar(255);uniqueIndex;not null"`
	SellPrice      float64 `json:"sellPrice" gorm:"index"`
	SellCount      int64   `json:"sellCount" gorm:"index"`
	BiddingPrice   float64 `json:"biddingPrice" gorm:"index"`
	BiddingCount   int64   `json:"biddingCount"`
	UpdateTime     int64   `json:"updateTime"`
	BeforeTime     int64   `json:"beforeTime"`
	BeforeCount    int64   `json:"beforeCount"`
	TurnOver       int64   `json:"turn_over"`
	Link           string  `json:"link"`
}

func BatchGetBuffGoods(hashNames []string) map[string]*Buff {
	var buffs []Buff
	result := make(map[string]*Buff)
	config.DB.Where("market_hash_name in ?", hashNames).Find(&buffs)
	for i := range buffs {
		result[buffs[i].MarketHashName] = &buffs[i]
	}
	return result
}

func BatchUpdateBuffGoods(buff []*Buff) {
	maxRetries := 3
	var err error
	for i := 0; i < maxRetries; i++ {
		err = config.DB.Transaction(func(tx *gorm.DB) error {
			if err := tx.Clauses(clause.OnConflict{UpdateAll: true}).CreateInBatches(buff, 50).Error; err != nil {
				return err
			}
			return nil
		})
		if err == nil {
			return
		}
		if strings.Contains(err.Error(), "Deadlock") || strings.Contains(err.Error(), "SAVEPOINT") {
			config.Log.Warnf("Update Buff Goods deadlock, retrying (%d/%d)...", i+1, maxRetries)
			time.Sleep(time.Millisecond * time.Duration(100*(i+1)))
			continue
		}
		break
	}
	if err != nil {
		config.Log.Errorf("Update Buff Goods fail: %v", err)
	}
}

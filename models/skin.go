package models

import (
	"gorm.io/gorm/clause"
	"time"
	"uu/config"
)

type SkinItem struct {
	Id          int64     `json:"id" gorm:"type:int;primaryKey"`
	Name        string    `json:"name"  gorm:"type:varchar(255)"`
	BuffPrice   float64   `json:"buff_price" gorm:"type:double"`
	UPrice      float64   `json:"u_price" gorm:"type:double"`
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

func UpdateSkinItems(minDiff, minSellPrice float64, minSellNum int) {
	var skins []SkinItem
	err := config.DB.Model(&UItem{}).Select(`select u.id as id, select u.commodity_name as name, u.icon_url as image_url, u.type_name as category, uitem.price as u_price, buff_item.sell_min_price as buff_price, (uitem.price - buff_item.sell_min_price) as price_diff, ROUND((u.price - b.sell_min_price)/b.sell_min_price,1)as profit_rate`).
		Joins("join buff_item ON uitem.commodity_hash_name = buff_item.market_hash_name").
		Where("u.price - b.sell_min_price > ? and b.sell_num > ? and b.sell_min_price < ?", minDiff, minSellNum, minSellPrice).Scan(&skins).Error
	if err != nil {
		config.Log.Errorf("get price diff data fail: %s", err)
	}
	err = config.DB.Clauses(clause.OnConflict{UpdateAll: true}).CreateInBatches(skins, 100).Error
	if err != nil {
		config.Log.Errorf("update skin item fail: %s", err)
	}
}

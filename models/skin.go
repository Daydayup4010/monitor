package models

import (
	"context"
	"gorm.io/gorm/clause"
	"time"
	"uu/config"
)

type SkinItem struct {
	Id         int64     `json:"id" gorm:"type:int;primaryKey"`
	Name       string    `json:"name"  gorm:"type:varchar(255)"`
	BuffPrice  float64   `json:"buff_price" gorm:"type:double"`
	UPrice     float64   `json:"u_price" gorm:"type:double"`
	Category   string    `json:"category" gorm:"type:varchar(20)"`
	ImageUrl   string    `json:"image_url" gorm:"type:varchar(255)"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"type:datetime"`
	PriceDiff  float64   `json:"price_diff" gorm:"type:double"`
	ProfitRate float64   `json:"profit_rate" gorm:"type:double"`
}

func GetSkinItems(pageSize, pageNum int) ([]SkinItem, int64) {
	var skins []SkinItem
	var total int64
	err := config.DB.Limit(pageSize).Offset((pageNum - 1) * pageSize).Find(&skins).Count(&total).Error
	if err != nil {
		config.Log.Errorf("Get skins error : %s", err)
	}
	return skins, total
}

func UpdateSkinItems() {
	var c = context.Background()
	var settings Settings
	err := settings.GetSettings(c)
	if err != nil {
		config.Log.Errorf("Get settings error: %s", err)
	}

	if err = config.DB.Exec("delete from skin_item").Error; err != nil {
		config.Log.Errorf("delete skin table fail: %s", err)
	}
	var skins []SkinItem
	err = config.DB.Model(&UItem{}).Select("uitem.id as id, uitem.commodity_name as name, uitem.icon_url as image_url, uitem.type_name as category, uitem.price as u_price, buff_item.sell_min_price as buff_price, (uitem.price - buff_item.sell_min_price) as price_diff, ROUND((uitem.price - buff_item.sell_min_price)/buff_item.sell_min_price,2) as profit_rate").
		Joins("join buff_item ON uitem.commodity_hash_name = buff_item.market_hash_name").
		Where("(uitem.price - buff_item.sell_min_price) > ? and buff_item.sell_num > ? and buff_item.sell_min_price < ? and buff_item.sell_min_price > ?", settings.MinDiff, settings.MinSellNum, settings.MaxSellPrice, settings.MinSellPrice).Scan(&skins).Error
	if err != nil {
		config.Log.Errorf("get price diff data fail: %s", err)
	}
	err = config.DB.Clauses(clause.OnConflict{UpdateAll: true}).CreateInBatches(skins, 100).Error
	if err != nil {
		config.Log.Errorf("update skin item fail: %s", err)
	}
}

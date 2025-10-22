package models

import (
	"fmt"
	"gorm.io/gorm/clause"
	"time"
	"uu/config"
)

type SkinItem struct {
	Id         int64     `json:"id" gorm:"type:int;primaryKey"`
	UserId     string    `json:"user_id" gorm:"type:char(36);index"`
	Name       string    `json:"name"  gorm:"type:varchar(255)"`
	BuffPrice  float64   `json:"buff_price" gorm:"type:double"`
	UPrice     float64   `json:"u_price" gorm:"type:double"`
	Category   string    `json:"category" gorm:"type:varchar(20)"`
	ImageUrl   string    `json:"image_url" gorm:"type:varchar(255)"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"type:datetime"`
	PriceDiff  float64   `json:"price_diff" gorm:"type:double"`
	ProfitRate float64   `json:"profit_rate" gorm:"type:double"`
}

type Goods struct {
	Id         int64   `json:"id"`
	UserId     string  `json:"user_id"`
	Name       string  `json:"name"`
	BuffPrice  float64 `json:"buff_price"`
	UPrice     float64 `json:"u_price"`
	Category   string  `json:"category"`
	ImageUrl   string  `json:"image_url"`
	PriceDiff  float64 `json:"price_diff"`
	ProfitRate float64 `json:"profit_rate"`
}

func GetSkinItems(pageSize, pageNum int, isDesc bool, sortField, category string) ([]SkinItem, int64) {
	var skins []SkinItem
	var total int64
	var err error
	validFields := map[string]bool{
		"buff_price":  true,
		"u_price":     true,
		"price_diff":  true,
		"profit_rate": true,
	}
	if !validFields[sortField] {
		sortField = "id" // 默认排序字段
	}

	order := sortField
	if isDesc {
		order += " DESC"
	}

	if category != "" {
		config.DB.Model(&SkinItem{}).Where("category = ?", category).Count(&total)
		err = config.DB.Order(order).Limit(pageSize).Offset((pageNum-1)*pageSize).Where("category = ?", category).Find(&skins).Error
	} else {
		config.DB.Model(&SkinItem{}).Count(&total)
		err = config.DB.Order(order).Limit(pageSize).Offset((pageNum - 1) * pageSize).Find(&skins).Error
	}

	if err != nil {
		config.Log.Errorf("Get skins error : %s", err)
	}
	return skins, total
}

func UpdateSkinItems(id string) {
	settings, err := GetUserSetting(id)
	if err != nil {
		config.Log.Errorf("Get user: %s settings error: %s", id, err)
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

func GetGoods(userId string, pageSize, pageNum int, isDesc bool, sortField, category string) (*[]Goods, int64, error) {
	var goods []Goods
	var total int64
	validFields := map[string]bool{
		"buff_price":  true,
		"u_price":     true,
		"price_diff":  true,
		"profit_rate": true,
	}

	if !validFields[sortField] {
		sortField = "id" // 默认排序字段
	}

	order := sortField
	if isDesc {
		order += " DESC"
	}

	settings, err := GetUserSetting(userId)
	if err != nil {
		config.Log.Errorf("Get user: %s settings error: %s", userId, err)
		return &goods, 0, fmt.Errorf("get user settings error")
	}

	err = config.DB.Model(&UItem{}).
		Joins("JOIN buff_item ON uitem.commodity_hash_name = buff_item.market_hash_name").
		Where("(uitem.price - buff_item.sell_min_price) > ?", settings.MinDiff).
		Where("buff_item.sell_num > ?", settings.MinSellNum).
		Where("buff_item.sell_min_price < ?", settings.MaxSellPrice).
		Where("buff_item.sell_min_price > ?", settings.MinSellPrice).
		Count(&total).Error
	if err != nil {
		config.Log.Errorf("get goods total fail: %v", err)
		return &goods, 0, fmt.Errorf("get goods total fail")
	}

	err = config.DB.Model(&UItem{}).Select("uitem.id as id, uitem.commodity_name as name, uitem.icon_url as image_url, uitem.type_name as category, uitem.price as u_price, buff_item.sell_min_price as buff_price, (uitem.price - buff_item.sell_min_price) as price_diff, ROUND((uitem.price - buff_item.sell_min_price)/buff_item.sell_min_price,2) as profit_rate").
		Joins("join buff_item ON uitem.commodity_hash_name = buff_item.market_hash_name").
		Where("(uitem.price - buff_item.sell_min_price) > ? and buff_item.sell_num > ? and buff_item.sell_min_price < ? and buff_item.sell_min_price > ?", settings.MinDiff, settings.MinSellNum, settings.MaxSellPrice, settings.MinSellPrice).Order(order).Limit(pageSize).Offset((pageNum - 1) * pageSize).Scan(&goods).Error
	if err != nil {
		config.Log.Errorf("get price diff data fail: %s", err)
		return &goods, 0, fmt.Errorf("get price diff data fail")
	}
	return &goods, total, nil
}

package models

import (
	"gorm.io/gorm/clause"
	"uu/config"
)

type UItem struct {
	CommodityHashName  string `json:"commodityHashName" gorm:"type:varchar(255)"`
	CommodityName      string `json:"commodityName" gorm:"type:varchar(255)"`
	IconUrl            string `json:"iconUrl" gorm:"type:varchar(255)"`
	Id                 int64  `json:"id" gorm:"type:int;primaryKey"`
	LongLeaseUnitPrice string `json:"longLeaseUnitPrice" gorm:"type:varchar(20)"`
	OnLeaseCount       int64  `json:"onLeaseCount" gorm:"type:int"`
	OnSaleCount        int64  `json:"onSaleCount" gorm:"type:int"`
	Price              string `json:"price" gorm:"type:varchar(20)"`
	RarityColor        string `json:"rarityColor" gorm:"type:varchar(20)"`
	Rent               string `json:"rent" gorm:"type:varchar(20)"`
	SortId             int64  `json:"sortId" gorm:"type:int"`
	SteamPrice         string `json:"steamPrice" gorm:"type:varchar(20)"`
	TypeName           string `json:"typeName" gorm:"type:varchar(255)"`
}

func BatchAddUUItem(uu []*UItem) {
	err := config.DB.Clauses(clause.OnConflict{UpdateAll: true}).CreateInBatches(uu, 80).Error
	if err != nil {
		config.Log.Errorf("insert uu item fail: %s", err)
	}
}

func GetUUItems() []UItem {
	var uu []UItem
	config.DB.Find(&uu)
	return uu
}

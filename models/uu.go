package models

import (
	"gorm.io/gorm/clause"
	"uu/config"
)

type UItem struct {
	HashName string `json:"commodityHashName" gorm:"type:varchar(255);uniqueIndex"`
	Name     string `json:"commodityName" gorm:"type:varchar(255)"`
	IconUrl  string `json:"iconUrl" gorm:"type:varchar(255)"`
	Id       int64  `json:"id" gorm:"type:int;primaryKey"`
	//LongLeaseUnitPrice string `json:"longLeaseUnitPrice" gorm:"type:decimal(10,2)"`
	//OnLeaseCount       int64  `json:"onLeaseCount" gorm:"type:int"`
	Count int64  `json:"onSaleCount" gorm:"type:int"`
	Price string `json:"price" gorm:"type:decimal(10,2)"`
	//RarityColor        string `json:"rarityColor" gorm:"type:varchar(20)"`
	//Rent               string `json:"rent" gorm:"type:decimal(10,2)"`
	//SortId             int64  `json:"sortId" gorm:"type:int"`
	SteamPrice string `json:"steamPrice" gorm:"type:decimal(10,2)"`
	TypeName   string `json:"typeName" gorm:"type:varchar(255)"`
}

type UItemsInfo struct {
	MarketHashName      string  `json:"marketHashName" gorm:"type:varchar(255);uniqueIndex"`
	Name                string  `json:"Name" gorm:"type:varchar(255)"`
	ImageUrl            string  `json:"imageUrl" gorm:"type:varchar(255)"`
	Id                  int64   `gorm:"type:int;primaryKey"`
	CacheExpirationDesc string  `json:"cacheExpirationDesc" gorm:"type:varchar(20)"`
	AssetMergeCount     int64   `json:"assetMergeCount" gorm:"type:int"`
	Price               float64 `json:"price" gorm:"type:decimal(10,2)"`
}

func BatchAddUUItem(uu []*UItem) {
	err := config.DB.Clauses(clause.OnConflict{UpdateAll: true}).CreateInBatches(uu, 100).Error
	if err != nil {
		config.Log.Errorf("batch insert uu item fail: %s", err)
	}
}

func BatchAddUUInventory(uu []*UItemsInfo) {
	err := config.DB.Clauses(clause.OnConflict{UpdateAll: true}).CreateInBatches(uu, 100).Error
	if err != nil {
		config.Log.Errorf("barch insert uu inventory fail: %s", err)
	}
}

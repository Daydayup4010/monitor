package models

import (
	"gorm.io/gorm"
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

type UBaseInfo struct {
	Id       int    `json:"item_id" gorm:"primaryKey"`
	HashName string `json:"hash_name" gorm:"type:varchar(255);uniqueIndex;not null"`
	IconUrl  string `json:"icon_url" gorm:"index"`
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

// -------------------------------------------------------v2------------------------------------------------------------
// data from steamDT

type U struct {
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

func BatchGetUUGoods(hashNames []string) map[string]*U {
	var uList []U
	err := config.DB.Where("market_hash_name in ?", hashNames).Find(&uList).Error
	if err != nil {
		config.Log.Errorf("Batch get uu goods error: %v", err)
	}
	result := make(map[string]*U)
	for i := range uList {
		result[uList[i].MarketHashName] = &uList[i]
	}
	return result
}

func BatchUpdateUUGoods(uu []*U) {
	err := config.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.OnConflict{UpdateAll: true}).CreateInBatches(uu, 100).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		config.Log.Errorf("Update UU Goods fail: %v", err)
		return
	}
	//config.Log.Info("Update UU Goods Success")
}

func BatchQueryHashIcon() ([]UBaseInfo, error) {
	var Infos []UBaseInfo
	err := config.DB.Select("hash_name, icon_url").Find(&Infos).Error
	if err != nil {
		config.Log.Errorf("Get uu icon url error: %v", err)
	}
	return Infos, err
}

func BatchUpdateIcon(infos []*UBaseInfo) {
	err := config.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.OnConflict{UpdateAll: true}).CreateInBatches(infos, 200).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		config.Log.Errorf("Update UU Goods Icon fail: %v", err)
		return
	}
}

func QueryAllUUHashName() []string {
	var hashNames []string
	err := config.DB.Model(&U{}).Pluck("market_hash_name", &hashNames).Error
	if err != nil {
		config.Log.Errorf("Get all uu hash name error: %v", err)
	}
	return hashNames
}

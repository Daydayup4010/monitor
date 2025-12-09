package models

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"math"
	"strconv"
	"time"
	"uu/config"
	"uu/utils"
)

type Goods struct {
	Id               int64   `json:"id"`
	MarketHashName   string  `json:"market_hash_name"`
	UserId           string  `json:"user_id"`
	Name             string  `json:"name"`
	SourcePrice      float64 `json:"source_price"`
	TargetPrice      float64 `json:"target_price"`
	SourceUpdateTime int64   `json:"source_update_time"`
	TargetUpdateTime int64   `json:"target_update_time"`
	BiddingPrice     float64 `json:"bidding_price"`
	BiddingCount     int64   `json:"bidding_count"`
	//Category    string  `json:"category"`
	ImageUrl     string      `json:"image_url"`
	PriceDiff    float64     `json:"price_diff"`
	ProfitRate   float64     `json:"profit_rate"`
	SellCount    int64       `json:"sell_count"`
	TurnOver     int64       `json:"turn_over"`
	PlatformList []*Platform `json:"platform_list" gorm:"-"`
}

type BaseGoods struct {
	ID             uint   `gorm:"primaryKey"`
	Name           string `json:"name" gorm:"type:varchar(255);not null;index"`
	MarketHashName string `json:"marketHashName" gorm:"type:varchar(255);uniqueIndex;not null"`
	IconUrl        string `json:"icon_url"`
	//PlatformList   []*Platform `json:"platformList" gorm:"-"` // 数据库忽略该字段
}

type Platform struct {
	Id           string  `json:"platformItemId"`
	Name         string  `json:"platformName" gorm:"-"`
	SellPrice    float64 `json:"sellPrice"`
	SellCount    int64   `json:"sellCount"`
	BiddingPrice float64 `json:"biddingPrice"`
	BiddingCount int64   `json:"biddingCount"`
	PriceDiff    float64 `json:"price_diff"`
	UpdateTime   int64   `json:"updateTime"`
	Link         string  `json:"link"`
}

func UpdateBaseGoods(base []*BaseGoods) {
	UInfos, err := BatchQueryHashIcon()
	iconMap := make(map[string]string)
	for _, info := range UInfos {
		iconMap[info.HashName] = info.IconUrl
	}

	if err != nil {
		return
	}
	for i := range base {
		if icon, exist := iconMap[base[i].MarketHashName]; exist {
			base[i].IconUrl = icon
		}
	}
	err = config.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "market_hash_name"}},
			DoUpdates: clause.AssignmentColumns([]string{"name", "icon_url"}),
		}).CreateInBatches(base, 100).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		config.Log.Errorf("Update Base Goods fail: %v", err)
		return
	}
	config.Log.Info("Update Base Goods Success")
}

func UpdateBaseGoodsIcon() {
	UInfos, err := BatchQueryHashIcon()
	if err != nil {
		config.Log.Errorf("Get icon info failed: %v", err)
		return
	}

	iconMap := make(map[string]string)
	for _, info := range UInfos {
		if info.IconUrl != "" {
			iconMap[info.HashName] = info.IconUrl
		}
	}

	if len(iconMap) == 0 {
		config.Log.Info("No icon data to update")
		return
	}

	// 只查询 icon_url 为空的 BaseGoods
	var baseGoodsList []BaseGoods
	err = config.DB.Where("icon_url = '' OR icon_url IS NULL").Find(&baseGoodsList).Error
	if err != nil {
		config.Log.Errorf("Query BaseGoods failed: %v", err)
		return
	}

	if len(baseGoodsList) == 0 {
		config.Log.Info("No BaseGoods need to update icon")
		return
	}

	// 筛选出能匹配到 icon 的记录
	var toUpdate []BaseGoods
	for i := range baseGoodsList {
		if icon, exist := iconMap[baseGoodsList[i].MarketHashName]; exist {
			baseGoodsList[i].IconUrl = icon
			toUpdate = append(toUpdate, baseGoodsList[i])
		}
	}

	if len(toUpdate) == 0 {
		config.Log.Info("No matching icon found for BaseGoods")
		return
	}

	// 批量更新到数据库
	batchSize := 500
	for i := 0; i < len(toUpdate); i += batchSize {
		end := i + batchSize
		if end > len(toUpdate) {
			end = len(toUpdate)
		}
		batch := toUpdate[i:end]

		err = config.DB.Transaction(func(tx *gorm.DB) error {
			for _, goods := range batch {
				if err := tx.Model(&BaseGoods{}).
					Where("market_hash_name = ?", goods.MarketHashName).
					Update("icon_url", goods.IconUrl).Error; err != nil {
					return err
				}
			}
			return nil
		})
		if err != nil {
			config.Log.Errorf("Batch update BaseGoods icon failed: %v", err)
			continue
		}
	}

	config.Log.Infof("Update BaseGoods Icon Success, updated %d records", len(toUpdate))
}

func GetHashNames() ([]string, error) {
	var hashNames []string
	err := config.DB.Model(&BaseGoods{}).Pluck("market_hash_name", &hashNames).Error
	return hashNames, err
}

// GetPlatformListBatch 批量获取多个商品的平台列表
func GetPlatformListBatch(marketHashNames []string) map[string][]*Platform {
	result := make(map[string][]*Platform)

	if len(marketHashNames) == 0 {
		return result
	}

	// 批量查询 UU 平台
	var uList []U
	config.DB.Where("market_hash_name IN ?", marketHashNames).Find(&uList)
	for _, u := range uList {
		result[u.MarketHashName] = append(result[u.MarketHashName], &Platform{
			Id:           u.Id,
			SellPrice:    u.SellPrice,
			SellCount:    u.SellCount,
			BiddingPrice: u.BiddingPrice,
			BiddingCount: u.BiddingCount,
			UpdateTime:   u.UpdateTime,
			Link:         u.Link,
			Name:         "悠悠",
		})
	}

	// 批量查询 Buff 平台
	var buffList []Buff
	config.DB.Where("market_hash_name IN ?", marketHashNames).Find(&buffList)
	for _, buff := range buffList {
		result[buff.MarketHashName] = append(result[buff.MarketHashName], &Platform{
			Id:           buff.Id,
			SellPrice:    buff.SellPrice,
			SellCount:    buff.SellCount,
			BiddingPrice: buff.BiddingPrice,
			BiddingCount: buff.BiddingCount,
			UpdateTime:   buff.UpdateTime,
			Link:         buff.Link,
			Name:         "BUFF",
		})
	}

	// 批量查询 C5 平台
	var c5List []C5
	config.DB.Where("market_hash_name IN ?", marketHashNames).Find(&c5List)
	for _, c5 := range c5List {
		result[c5.MarketHashName] = append(result[c5.MarketHashName], &Platform{
			Id:           c5.Id,
			SellPrice:    c5.SellPrice,
			SellCount:    c5.SellCount,
			BiddingPrice: c5.BiddingPrice,
			BiddingCount: c5.BiddingCount,
			UpdateTime:   c5.UpdateTime,
			Link:         c5.Link,
			Name:         "C5GAME",
		})
	}

	// 批量查询 Steam 平台
	var steamList []Steam
	config.DB.Where("market_hash_name IN ?", marketHashNames).Find(&steamList)
	for _, steam := range steamList {
		result[steam.MarketHashName] = append(result[steam.MarketHashName], &Platform{
			Id:           fmt.Sprintf("%d", steam.Id),
			SellPrice:    steam.SellPrice,
			SellCount:    steam.SellCount,
			BiddingPrice: steam.BiddingPrice,
			BiddingCount: steam.BiddingCount,
			UpdateTime:   steam.UpdateTime,
			Link:         steam.Link,
			Name:         "Steam",
		})
	}

	return result
}

func GetGoods(userId string, pageSize, pageNum int, isDesc bool, sortField, search, source, target string) (*[]Goods, int64, int) {
	var goods []Goods
	var total int64
	validFields := map[string]bool{
		"price_diff":  true,
		"profit_rate": true,
	}
	var targetMap map[string]interface{}
	tableMap := map[string]string{
		"uu":    "u",
		"buff":  "buff",
		"c5":    "c5",
		"steam": "steam",
	}

	sourceTable := tableMap[source]

	if !validFields[sortField] {
		sortField = "id" // 默认排序字段
	}

	order := sortField
	if isDesc {
		order += " DESC"
	}

	settings, code := GetUserSetting(userId)

	switch target {
	case "uu":
		targetMap = map[string]interface{}{
			"model": &U{},
			"table": tableMap[target],
		}
	case "buff":
		targetMap = map[string]interface{}{
			"model": &Buff{},
			"table": tableMap[target],
		}
	case "steam":
		targetMap = map[string]interface{}{
			"model": &Steam{},
			"table": tableMap[target],
		}
	case "c5":
		targetMap = map[string]interface{}{
			"model": &C5{},
			"table": tableMap[target],
		}
	default:
		targetMap = map[string]interface{}{
			"model": &U{},
			"table": tableMap[target],
		}
	}

	query1 := config.DB.Model(targetMap["model"]).
		Select(fmt.Sprintf("%s.id as id, %s.sell_count as sell_count, %s.turn_over as turn_over, %s.bidding_count as bidding_count, %s.bidding_price as bidding_price, base_goods.market_hash_name as market_hash_name, base_goods.name as name, base_goods.icon_url as image_url, %s.sell_price as target_price, %s.sell_price as source_price, (%s.sell_price - %s.sell_price) as price_diff, ROUND((%s.sell_price - %s.sell_price)/%s.sell_price,4) as profit_rate, %s.update_time as target_update_time, %s.update_time as source_update_time", targetMap["table"], targetMap["table"], targetMap["table"], targetMap["table"], targetMap["table"], targetMap["table"], sourceTable, targetMap["table"], sourceTable, targetMap["table"], sourceTable, sourceTable, targetMap["table"], sourceTable)).
		Joins(fmt.Sprintf("join %s ON %s.market_hash_name = %s.market_hash_name", sourceTable, targetMap["table"], sourceTable)).
		Joins(fmt.Sprintf("join base_goods ON %s.market_hash_name = base_goods.market_hash_name", targetMap["table"])).
		Where(fmt.Sprintf("(%s.sell_price - %s.sell_price) > ? and %s.sell_count > ? and %s.sell_price < ? and %s.sell_price > ?", targetMap["table"], sourceTable, targetMap["table"], sourceTable, sourceTable),
			settings.MinDiff, settings.MinSellNum, settings.MaxSellPrice, settings.MinSellPrice)
	query2 := config.DB.Model(targetMap["model"]).
		Joins(fmt.Sprintf("join %s ON %s.market_hash_name = %s.market_hash_name", sourceTable, targetMap["table"], sourceTable)).
		Joins(fmt.Sprintf("join base_goods ON %s.market_hash_name = base_goods.market_hash_name", targetMap["table"])).
		Where(fmt.Sprintf("(%s.sell_price - %s.sell_price) > ? and %s.sell_count > ? and %s.sell_price < ? and %s.sell_price > ?", targetMap["table"], sourceTable, targetMap["table"], sourceTable, sourceTable),
			settings.MinDiff, settings.MinSellNum, settings.MaxSellPrice, settings.MinSellPrice)
	//if category != "" {
	//	query1 = query1.Where(fmt.Sprintf("%s.type_name = ?", targetMap["table"]), category)
	//	query2 = query2.Where(fmt.Sprintf("%s.type_name = ?", targetMap["table"]), category)
	//}

	if search != "" {
		query1 = query1.Where("name LIKE ?", "%"+search+"%")
		query2 = query2.Where("base_goods.name LIKE ?", "%"+search+"%")
	}

	err := query2.Count(&total).Error
	if err != nil {
		config.Log.Errorf("get goods total fail: %v", err)
		return &goods, 0, utils.ErrCodeGetGoodsTotal
	}
	err = query1.
		Order(order).
		Limit(pageSize).
		Offset((pageNum - 1) * pageSize).
		Scan(&goods).
		Error
	if err != nil {
		config.Log.Errorf("get price diff data fail: %s", err)
		return &goods, 0, utils.ErrCodeGetGoods
	}
	hashNames := make([]string, 0, len(goods))
	for i := range goods {
		hashNames = append(hashNames, goods[i].MarketHashName)
	}
	platformList := GetPlatformListBatch(hashNames)
	for i := range goods {
		goods[i].PlatformList = platformList[goods[i].MarketHashName]
		for _, p := range goods[i].PlatformList {
			p.PriceDiff = goods[i].TargetPrice - p.SellPrice
			p.PriceDiff = math.Abs(math.Round(p.PriceDiff*100) / 100)
		}
	}
	return &goods, total, code
}

func GetCategory() (*[]string, int) {
	var category []string
	err := config.DB.Model(UItem{}).Select("DISTINCT(type_name)").Scan(&category).Error
	if err != nil {
		config.Log.Errorf("Get category error: %v", err)
		return &category, utils.ErrCodeGetGoodsCategory
	}
	return &category, utils.SUCCESS
}

func SetLastIndex(index int) {
	ctx := context.Background()
	err := config.RDB.Set(ctx, "hash_name_index", index, time.Minute*5).Err()
	if err != nil {
		config.Log.Errorf("Set hash name index error: %v", err)
	}
}

func GetLastIndex() int {
	ctx := context.Background()
	index, err := config.RDB.Get(ctx, "hash_name_index").Result()
	if err != nil {
		config.Log.Errorf("Get hash name index error: %v", err)
		return 0
	}
	newIndex, _ := strconv.Atoi(index)

	return newIndex
}

// PublicHomeData GetPublicHomeData 获取公开首页数据（不需要登录）
// 包含：饰品榜单前10，挂刀搬砖利润率前10（使用管理员settings）
type PublicHomeData struct {
	RankingList []PriceIncreaseItem `json:"rankingList"` // 饰品涨幅榜前10
	BrickMoving []Goods             `json:"brickMoving"` // 搬砖数据前10
}

func GetPublicHomeData() (*PublicHomeData, error) {
	result := &PublicHomeData{}

	// 1. 获取饰品涨幅榜前10
	rankingList, err := GetPriceIncrease("YOUPIN", "", true, 10)
	if err != nil {
		config.Log.Errorf("GetPublicHomeData get ranking error: %v", err)
		rankingList = []PriceIncreaseItem{}
	}
	result.RankingList = rankingList

	// 2. 获取搬砖数据前10（使用管理员settings）
	adminSettings, _ := GetAdminSetting()

	var goods []Goods
	// 默认 uu -> steam 的搬砖数据
	sourceTable := "u"
	targetTable := "steam"

	query := config.DB.Model(&Buff{}).
		Select(fmt.Sprintf("%s.id as id, %s.sell_count as sell_count, %s.turn_over as turn_over, %s.bidding_count as bidding_count, %s.bidding_price as bidding_price, base_goods.market_hash_name as market_hash_name, base_goods.name as name, base_goods.icon_url as image_url, %s.sell_price as target_price, %s.sell_price as source_price, (%s.sell_price - %s.sell_price) as price_diff, ROUND((%s.sell_price - %s.sell_price)/%s.sell_price,4) as profit_rate, %s.update_time as target_update_time, %s.update_time as source_update_time",
			targetTable, targetTable, targetTable, targetTable, targetTable, targetTable, sourceTable, targetTable, sourceTable, targetTable, sourceTable, sourceTable, targetTable, sourceTable)).
		Joins(fmt.Sprintf("join %s ON %s.market_hash_name = %s.market_hash_name", sourceTable, targetTable, sourceTable)).
		Joins(fmt.Sprintf("join base_goods ON %s.market_hash_name = base_goods.market_hash_name", targetTable)).
		Where(fmt.Sprintf("(%s.sell_price - %s.sell_price) > ? and %s.sell_count > ? and %s.sell_price < ? and %s.sell_price > ?",
			targetTable, sourceTable, targetTable, sourceTable, sourceTable),
			adminSettings.MinDiff, adminSettings.MinSellNum, adminSettings.MaxSellPrice, adminSettings.MinSellPrice).
		Order("profit_rate DESC").
		Limit(10)

	err = query.Scan(&goods).Error
	if err != nil {
		config.Log.Errorf("GetPublicHomeData get brick moving error: %v", err)
		goods = []Goods{}
	}

	// 获取平台列表
	if len(goods) > 0 {
		hashNames := make([]string, 0, len(goods))
		for i := range goods {
			hashNames = append(hashNames, goods[i].MarketHashName)
		}
		platformList := GetPlatformListBatch(hashNames)
		for i := range goods {
			goods[i].PlatformList = platformList[goods[i].MarketHashName]
			for _, p := range goods[i].PlatformList {
				p.PriceDiff = goods[i].TargetPrice - p.SellPrice
				p.PriceDiff = math.Abs(math.Round(p.PriceDiff*100) / 100)
			}
		}
	}
	result.BrickMoving = goods

	return result, nil
}

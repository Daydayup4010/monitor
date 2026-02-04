package models

import (
	"context"
	"fmt"
	"strings"

	"math"
	"strconv"
	"time"
	"uu/config"
	"uu/utils"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Goods struct {
	Id               int64       `json:"id"`
	MarketHashName   string      `json:"market_hash_name"`
	UserId           string      `json:"user_id"`
	Name             string      `json:"name"`
	SourcePrice      float64     `json:"source_price"`
	TargetPrice      float64     `json:"target_price"`
	SourceUpdateTime int64       `json:"source_update_time"`
	TargetUpdateTime int64       `json:"target_update_time"`
	BiddingPrice     float64     `json:"bidding_price"`
	BiddingCount     int64       `json:"bidding_count"`
	TypeName         string      `json:"type_name"`
	ImageUrl         string      `json:"image_url"`
	PriceDiff        float64     `json:"price_diff"`
	ProfitRate       float64     `json:"profit_rate"`
	SellCount        int64       `json:"sell_count"`
	TurnOver         int64       `json:"turn_over"`
	PlatformList     []*Platform `json:"platform_list" gorm:"-"`
}

type BaseGoods struct {
	ID             uint   `gorm:"primaryKey"`
	Name           string `json:"name" gorm:"type:varchar(255);not null;index"`
	MarketHashName string `json:"marketHashName" gorm:"type:varchar(255);uniqueIndex;not null"`
	IconUrl        string `json:"icon_url"`
	//PlatformList   []*Platform `json:"platformList" gorm:"-"` // 数据库忽略该字段
}

// SearchResult 搜索结果
type SearchResult struct {
	Name           string `json:"name" gorm:"column:name"`
	MarketHashName string `json:"marketHashName" gorm:"column:market_hash_name"`
	IconUrl        string `json:"iconUrl" gorm:"column:icon_url"`
}

// SearchGoodsByKeyword 根据关键词搜索商品（支持模糊匹配）
// 将关键词拆分为字符，按顺序匹配，例如 "蝴蝶刀蓝钢" -> "%蝴%蝶%刀%蓝%钢%"
func SearchGoodsByKeyword(keyword string, limit int) ([]SearchResult, error) {
	var results []SearchResult

	// 将关键词转为字符顺序匹配模式
	// 例如: "蝴蝶刀蓝钢" -> "%蝴%蝶%刀%蓝%钢%"
	chars := strings.Split(keyword, "")
	pattern := "%" + strings.Join(chars, "%") + "%"

	err := config.DB.Model(&BaseGoods{}).
		Select("name, market_hash_name, icon_url").
		Where("name LIKE ?", pattern).
		Limit(limit).
		Find(&results).Error
	return results, err
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

// GetGoods 获取搬砖数据
// buyType: sell(在售价购买) / bidding(求购价购买)
// sellType: sell(在售价出售) / bidding(求购价出售)
func GetGoods(userId string, pageSize, pageNum int, isDesc bool, sortField, search, source, target, category, buyType, sellType string) (*[]Goods, int64, int) {
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
		sortField = "profit_rate" // 默认按利润率排序
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

	targetTable := targetMap["table"].(string)

	// 根据 buyType 和 sellType 确定使用的价格字段
	// buyType: 购买方案 - sell(在售价购买) / bidding(求购价购买)
	// sellType: 出售方案 - sell(在售价出售) / bidding(求购价出售)
	var sourcePriceField, targetPriceField string
	if buyType == "bidding" {
		sourcePriceField = "bidding_price" // 求购价购买：用来源平台的求购价作为买入价
	} else {
		sourcePriceField = "sell_price" // 在售价购买：用来源平台的在售价作为买入价
	}
	if sellType == "bidding" {
		targetPriceField = "bidding_price" // 求购价出售：用目标平台的求购价作为卖出价
	} else {
		targetPriceField = "sell_price" // 在售价出售：用目标平台的在售价作为卖出价
	}

	// 构建 SELECT 语句
	// source_price: 买入价（来源平台）
	// target_price: 卖出价（目标平台）
	// price_diff: 差价 = 卖出价 - 买入价
	// profit_rate: 利润率 = 差价 / 买入价
	selectSQL := fmt.Sprintf(`
		%s.id as id,
		%s.sell_count as sell_count,
		%s.turn_over as turn_over,
		%s.bidding_count as bidding_count,
		%s.bidding_price as bidding_price,
		base_goods.market_hash_name as market_hash_name,
		base_goods.name as name,
		base_goods.icon_url as image_url,
		u_base_info.type_name as type_name,
		%s.%s as target_price,
		%s.%s as source_price,
		(%s.%s - %s.%s) as price_diff,
		ROUND((%s.%s - %s.%s) / %s.%s, 4) as profit_rate,
		%s.update_time as target_update_time,
		%s.update_time as source_update_time`,
		targetTable,
		targetTable,
		targetTable,
		targetTable,
		targetTable,
		targetTable, targetPriceField,
		sourceTable, sourcePriceField,
		targetTable, targetPriceField, sourceTable, sourcePriceField,
		targetTable, targetPriceField, sourceTable, sourcePriceField, sourceTable, sourcePriceField,
		targetTable,
		sourceTable)

	// 构建 WHERE 条件
	// 差价 > 最小差价
	// 目标平台在售数量 > 最小在售数量
	// 来源平台买入价 < 最大价格 且 > 最小价格
	// 同时要求使用的价格字段 > 0（避免除零错误和无效数据）
	whereSQL := fmt.Sprintf(`
		(%s.%s - %s.%s) > ? AND
		%s.sell_count > ? AND
		%s.%s < ? AND
		%s.%s > ? AND
		%s.%s > 0 AND
		%s.%s > 0`,
		targetTable, targetPriceField, sourceTable, sourcePriceField,
		targetTable,
		sourceTable, sourcePriceField,
		sourceTable, sourcePriceField,
		sourceTable, sourcePriceField,
		targetTable, targetPriceField)

	query1 := config.DB.Model(targetMap["model"]).
		Select(selectSQL).
		Joins(fmt.Sprintf("join %s ON %s.market_hash_name = %s.market_hash_name", sourceTable, targetTable, sourceTable)).
		Joins(fmt.Sprintf("join base_goods ON %s.market_hash_name = base_goods.market_hash_name", targetTable)).
		Joins(fmt.Sprintf("left join u_base_info ON %s.market_hash_name = u_base_info.hash_name", targetTable)).
		Where(whereSQL, settings.MinDiff, settings.MinSellNum, settings.MaxSellPrice, settings.MinSellPrice)

	query2 := config.DB.Model(targetMap["model"]).
		Joins(fmt.Sprintf("join %s ON %s.market_hash_name = %s.market_hash_name", sourceTable, targetTable, sourceTable)).
		Joins(fmt.Sprintf("join base_goods ON %s.market_hash_name = base_goods.market_hash_name", targetTable)).
		Joins(fmt.Sprintf("left join u_base_info ON %s.market_hash_name = u_base_info.hash_name", targetTable)).
		Where(whereSQL, settings.MinDiff, settings.MinSellNum, settings.MaxSellPrice, settings.MinSellPrice)

	if category != "" {
		// 支持多个类别，逗号分隔
		categories := strings.Split(category, ",")
		query1 = query1.Where("u_base_info.type_name IN ?", categories)
		query2 = query2.Where("u_base_info.type_name IN ?", categories)
	}

	if search != "" {
		query1 = query1.Where("base_goods.name LIKE ?", "%"+search+"%")
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

// BigItemBidding 大件求购数据结构
type BigItemBidding struct {
	Id             int64       `json:"id"`
	MarketHashName string      `json:"market_hash_name"`
	Name           string      `json:"name"`
	ImageUrl       string      `json:"image_url"`
	TypeName       string      `json:"type_name"`
	SellPrice      float64     `json:"sell_price"`
	SellCount      int64       `json:"sell_count"`
	BiddingPrice   float64     `json:"bidding_price"`
	BiddingCount   int64       `json:"bidding_count"`
	PriceDiff      float64     `json:"price_diff"`  // 价差 = sell_price - bidding_price
	ProfitRate     float64     `json:"profit_rate"` // 利润率 = (sell_price - bidding_price) / bidding_price
	UpdateTime     int64       `json:"update_time"`
	PlatformList   []*Platform `json:"platform_list" gorm:"-"`
}

func GetBigItemBidding(pageSize, pageNum int, isDesc bool, sortField, search, platform, category string) (*[]BigItemBidding, int64, int) {
	var items []BigItemBidding
	var total int64

	validFields := map[string]bool{
		"price_diff":    true,
		"profit_rate":   true,
		"sell_price":    true,
		"bidding_price": true,
	}

	tableMap := map[string]string{
		"uu":   "u",
		"buff": "buff",
		"c5":   "c5",
	}

	platformTable, ok := tableMap[platform]
	if !ok {
		platformTable = "u" // 默认悠悠
	}

	if !validFields[sortField] {
		sortField = "profit_rate" // 默认按利润率排序
	}

	order := sortField
	if isDesc {
		order += " DESC"
	}

	// 构建查询
	selectFields := fmt.Sprintf(`
		%s.id as id,
		%s.market_hash_name as market_hash_name,
		base_goods.name as name,
		base_goods.icon_url as image_url,
		u_base_info.type_name as type_name,
		%s.sell_price as sell_price,
		%s.sell_count as sell_count,
		%s.bidding_price as bidding_price,
		%s.bidding_count as bidding_count,
		(%s.sell_price - %s.bidding_price) as price_diff,
		ROUND((%s.sell_price - %s.bidding_price) / %s.bidding_price, 4) as profit_rate,
		%s.update_time as update_time
	`, platformTable, platformTable, platformTable, platformTable, platformTable, platformTable,
		platformTable, platformTable, platformTable, platformTable, platformTable, platformTable)

	query1 := config.DB.Table(platformTable).
		Select(selectFields).
		Joins(fmt.Sprintf("JOIN base_goods ON %s.market_hash_name = base_goods.market_hash_name", platformTable)).
		Joins(fmt.Sprintf("LEFT JOIN u_base_info ON %s.market_hash_name = u_base_info.hash_name", platformTable)).
		Where(fmt.Sprintf("%s.bidding_price > 0 AND %s.sell_price > 0 AND %s.sell_count > 10", platformTable, platformTable, platformTable)).
		Where(fmt.Sprintf("%s.sell_price > %s.bidding_price", platformTable, platformTable))

	query2 := config.DB.Table(platformTable).
		Joins(fmt.Sprintf("JOIN base_goods ON %s.market_hash_name = base_goods.market_hash_name", platformTable)).
		Joins(fmt.Sprintf("LEFT JOIN u_base_info ON %s.market_hash_name = u_base_info.hash_name", platformTable)).
		Where(fmt.Sprintf("%s.bidding_price > 0 AND %s.sell_price > 0 AND %s.sell_count > 10", platformTable, platformTable, platformTable)).
		Where(fmt.Sprintf("%s.sell_price > %s.bidding_price", platformTable, platformTable))

	// 类别筛选（默认手套和刀具）
	if category == "" || category == "all" {
		query1 = query1.Where("u_base_info.type_name IN ?", []string{"手套", "匕首"})
		query2 = query2.Where("u_base_info.type_name IN ?", []string{"手套", "匕首"})
	} else {
		query1 = query1.Where("u_base_info.type_name = ?", category)
		query2 = query2.Where("u_base_info.type_name = ?", category)
	}

	// 搜索
	if search != "" {
		query1 = query1.Where("base_goods.name LIKE ?", "%"+search+"%")
		query2 = query2.Where("base_goods.name LIKE ?", "%"+search+"%")
	}

	// 计算总数
	err := query2.Count(&total).Error
	if err != nil {
		config.Log.Errorf("Get big item bidding total fail: %v", err)
		return &items, 0, utils.ErrCodeGetGoods
	}

	// 查询数据
	err = query1.
		Order(order).
		Limit(pageSize).
		Offset((pageNum - 1) * pageSize).
		Scan(&items).
		Error
	if err != nil {
		config.Log.Errorf("Get big item bidding data fail: %v", err)
		return &items, 0, utils.ErrCodeGetGoods
	}

	// 获取各平台数据
	hashNames := make([]string, 0, len(items))
	for i := range items {
		hashNames = append(hashNames, items[i].MarketHashName)
	}
	platformList := GetPlatformListBatch(hashNames)
	for i := range items {
		items[i].PlatformList = platformList[items[i].MarketHashName]
	}

	return &items, total, utils.SUCCESS
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

// RelatedWearItem 关联磨损项
type RelatedWearItem struct {
	HashName    string  `json:"hash_name"`
	Wear        string  `json:"wear"`
	WearShort   string  `json:"wear_short"`
	IconUrl     string  `json:"icon_url"`
	Price       float64 `json:"price"`
	QualityName string  `json:"quality_name"`
}

// RelatedWearsResponse 关联磨损响应
type RelatedWearsResponse struct {
	CurrentWear    string                       `json:"current_wear"`
	CurrentQuality string                       `json:"current_quality"`
	Qualities      []string                     `json:"qualities"`
	Wears          map[string][]RelatedWearItem `json:"wears"`
}

// 磨损等级映射
var wearMap = map[string]struct {
	Short string
	Order int
}{
	"Factory New":    {"FN", 1},
	"Minimal Wear":   {"MW", 2},
	"Field-Tested":   {"FT", 3},
	"Well-Worn":      {"WW", 4},
	"Battle-Scarred": {"BS", 5},
}

// 磨损等级中文映射
var wearCNMap = map[string]string{
	"Factory New":    "崭新出厂",
	"Minimal Wear":   "略有磨损",
	"Field-Tested":   "久经沙场",
	"Well-Worn":      "破损不堪",
	"Battle-Scarred": "战痕累累",
}

// extractBaseName 从 hash_name 中提取基础名称（去掉磨损等级）
func extractBaseName(hashName string) string {
	// 磨损等级列表
	wears := []string{"(Factory New)", "(Minimal Wear)", "(Field-Tested)", "(Well-Worn)", "(Battle-Scarred)"}
	result := hashName
	for _, wear := range wears {
		if strings.HasSuffix(result, wear) {
			result = strings.TrimSuffix(result, wear)
			break
		}
	}
	// 去掉末尾空格
	result = strings.TrimSpace(result)
	return result
}

// extractWear 从 hash_name 中提取磨损等级
func extractWear(hashName string) string {
	wears := []string{"Factory New", "Minimal Wear", "Field-Tested", "Well-Worn", "Battle-Scarred"}
	for _, wear := range wears {
		suffix := "(" + wear + ")"
		if strings.HasSuffix(hashName, suffix) {
			return wear
		}
	}
	return ""
}

// extractKnifeBaseName 从刀具名称中提取基础名称（去掉 ★ 和 StatTrak™ 前缀）
// 例如：★ StatTrak™ Butterfly Knife -> Butterfly Knife
// 例如：★ Butterfly Knife -> Butterfly Knife
func extractKnifeBaseName(hashName string) string {
	result := hashName
	// 去掉 ★ 前缀
	result = strings.TrimPrefix(result, "★ ")
	result = strings.TrimPrefix(result, "★")
	// 去掉 StatTrak™ 前缀
	result = strings.TrimPrefix(result, "StatTrak™ ")
	result = strings.TrimPrefix(result, "StatTrak™")
	result = strings.TrimSpace(result)
	return result
}

// inferQualityFromHashName 从 hash_name 推断 quality_name
// 如果数据库中 quality_name 为空，则根据 hash_name 推断
func inferQualityFromHashName(hashName string) string {
	if strings.Contains(hashName, "★ StatTrak™") {
		return "★ StatTrak™"
	}
	if strings.Contains(hashName, "StatTrak™") {
		return "StatTrak™"
	}
	if strings.HasPrefix(hashName, "Souvenir ") {
		return "纪念品"
	}
	if strings.HasPrefix(hashName, "★") {
		return "★"
	}
	return "普通"
}

// 贴纸变体类型列表
var stickerVariants = []string{"(Holo)", "(Gold)", "(Foil)", "(Embroidered)", "(Champion)", "(Lenticular)"}

// isSticker 判断是否是贴纸类饰品（包括 Sticker 和 Sticker Slab）
func isSticker(hashName string) bool {
	return strings.HasPrefix(hashName, "Sticker |") || strings.HasPrefix(hashName, "Sticker Slab |")
}

// extractStickerBaseName 提取贴纸基础名称（去掉变体类型和 Slab 前缀）
// 例如：Sticker | The Mongolz (Holo) | Budapest 2025 -> The Mongolz | Budapest 2025
// 例如：Sticker Slab | The Mongolz (Gold) | Budapest 2025 -> The Mongolz | Budapest 2025
func extractStickerBaseName(hashName string) string {
	result := hashName
	// 去掉变体类型
	for _, variant := range stickerVariants {
		result = strings.Replace(result, " "+variant, "", 1)
	}
	// 去掉 Sticker Slab | 或 Sticker | 前缀，得到纯名称
	result = strings.TrimPrefix(result, "Sticker Slab | ")
	result = strings.TrimPrefix(result, "Sticker | ")
	return result
}

// extractStickerVariant 提取贴纸变体类型
// 返回 "Holo", "Gold", "普通" 等
func extractStickerVariant(hashName string) string {
	for _, variant := range stickerVariants {
		if strings.Contains(hashName, variant) {
			// 去掉括号，返回变体名称
			return strings.Trim(variant, "()")
		}
	}
	return "普通"
}

// GetRelatedWears 获取同款饰品的不同磨损和品质版本
func GetRelatedWears(hashName string) (*RelatedWearsResponse, error) {
	baseName := extractBaseName(hashName)
	currentWear := extractWear(hashName)

	if baseName == "" {
		return nil, fmt.Errorf("invalid hash_name")
	}

	var baseInfos []UBaseInfo

	// 判断饰品类型
	if isSticker(hashName) {
		// 贴纸类饰品：按变体类型分组（普通、Holo、Gold 等）
		// 同时查询 Sticker 和 Sticker Slab 版本
		stickerBaseName := extractStickerBaseName(hashName) // 例如：The Mongolz | Budapest 2025
		config.Log.Infof("GetRelatedWears (sticker): hashName=%s, stickerBaseName=%s", hashName, stickerBaseName)

		var conditions []string
		var args []interface{}

		// Sticker 和 Sticker Slab 两种前缀
		prefixes := []string{"Sticker | ", "Sticker Slab | "}

		for _, prefix := range prefixes {
			// 普通版本（精确匹配）
			conditions = append(conditions, "hash_name = ?")
			args = append(args, prefix+stickerBaseName)

			// 各种变体版本
			for _, variant := range stickerVariants {
				// 在名称中插入变体标识
				// 例如：Sticker | The Mongolz (Holo) | Budapest 2025
				parts := strings.SplitN(stickerBaseName, " |", 2)
				if len(parts) == 2 {
					variantName := prefix + parts[0] + " " + variant + " |" + parts[1]
					conditions = append(conditions, "hash_name = ?")
					args = append(args, variantName)
				}
			}
		}

		query := strings.Join(conditions, " OR ")
		err := config.DB.Where(query, args...).Find(&baseInfos).Error
		if err != nil {
			config.Log.Errorf("GetRelatedWears sticker query error: %v", err)
			return nil, err
		}
	} else if currentWear != "" {
		// 有磨损等级：需要查询所有品质版本
		config.Log.Infof("GetRelatedWears: hashName=%s, baseName=%s", hashName, baseName)

		// 判断是否是刀具（以 ★ 开头）
		if strings.HasPrefix(baseName, "★ ") {
			// 刀具格式：★ Knife | Skin 或 ★ StatTrak™ Knife | Skin
			// 提取纯刀具名称（去掉 ★ 和 StatTrak™）
			pureKnifeName := baseName
			pureKnifeName = strings.TrimPrefix(pureKnifeName, "★ StatTrak™ ")
			pureKnifeName = strings.TrimPrefix(pureKnifeName, "★ ")
			pureKnifeName = strings.TrimSpace(pureKnifeName)

			config.Log.Infof("GetRelatedWears (knife): pureKnifeName=%s", pureKnifeName)

			// 查询普通★版本和★ StatTrak™版本
			err := config.DB.Where("hash_name LIKE ? OR hash_name LIKE ?",
				"★ "+pureKnifeName+" (%",
				"★ StatTrak™ "+pureKnifeName+" (%").Find(&baseInfos).Error
			if err != nil {
				config.Log.Errorf("GetRelatedWears knife query error: %v", err)
				return nil, err
			}
		} else {
			// 普通武器：去掉 StatTrak™ 和 Souvenir 前缀
			pureBaseName := baseName
			pureBaseName = strings.TrimPrefix(pureBaseName, "StatTrak™ ")
			pureBaseName = strings.TrimPrefix(pureBaseName, "Souvenir ")
			pureBaseName = strings.TrimSpace(pureBaseName)

			config.Log.Infof("GetRelatedWears (weapon): pureBaseName=%s", pureBaseName)

			// 查询普通版本、StatTrak™ 版本和 Souvenir 版本
			err := config.DB.Where("hash_name LIKE ? OR hash_name LIKE ? OR hash_name LIKE ?",
				pureBaseName+" (%",
				"StatTrak™ "+pureBaseName+" (%",
				"Souvenir "+pureBaseName+" (%").Find(&baseInfos).Error
			if err != nil {
				config.Log.Errorf("GetRelatedWears query error: %v", err)
				return nil, err
			}
		}
	} else {
		// 无磨损等级（如刀具）：查询同名但不同品质的版本
		// 例如：★ Butterfly Knife 和 ★ StatTrak™ Butterfly Knife
		config.Log.Infof("GetRelatedWears (no wear): hashName=%s", hashName)

		// 提取刀具基础名称（去掉 ★ 和 StatTrak™ 前缀）
		knifeBaseName := extractKnifeBaseName(hashName)
		if knifeBaseName != "" {
			// 查询所有包含这个刀具名称的饰品
			err := config.DB.Where("hash_name LIKE ?", "%"+knifeBaseName).Find(&baseInfos).Error
			if err != nil {
				config.Log.Errorf("GetRelatedWears knife query error: %v", err)
				return nil, err
			}
		} else {
			// 如果无法提取，直接查询完全匹配的
			err := config.DB.Where("hash_name = ?", hashName).Find(&baseInfos).Error
			if err != nil {
				return nil, err
			}
		}
	}
	config.Log.Infof("GetRelatedWears: found %d items", len(baseInfos))

	// 获取当前饰品的品质
	var currentQuality string
	for _, info := range baseInfos {
		if info.HashName == hashName {
			currentQuality = info.QualityName
			if currentQuality == "" {
				currentQuality = inferQualityFromHashName(info.HashName)
			}
			break
		}
	}

	// 获取价格数据（从悠悠平台）
	hashNames := make([]string, len(baseInfos))
	for i, info := range baseInfos {
		hashNames[i] = info.HashName
	}
	priceMap := make(map[string]float64)
	if len(hashNames) > 0 {
		var uList []U
		config.DB.Where("market_hash_name IN ?", hashNames).Find(&uList)
		for _, u := range uList {
			priceMap[u.MarketHashName] = u.SellPrice
		}
	}

	// 按品质分组
	qualityMap := make(map[string][]RelatedWearItem)
	qualitiesSet := make(map[string]bool)

	// 判断是否是贴纸
	isStickerItem := len(baseInfos) > 0 && isSticker(baseInfos[0].HashName)

	for _, info := range baseInfos {
		var wear, wearShort string
		var wearOrder int

		if isStickerItem {
			// 贴纸：用变体类型代替磨损等级
			variant := extractStickerVariant(info.HashName)
			wear = variant
			wearShort = variant
			// 变体排序：普通 < 闪亮/刺绣 < 全息 < 金色
			variantOrder := map[string]int{
				"普通": 1, "Foil": 2, "Embroidered": 2, "Holo": 3, "Gold": 4, "Champion": 5, "Lenticular": 6,
			}
			wearOrder = variantOrder[variant]
			if wearOrder == 0 {
				wearOrder = 99
			}
		} else {
			wear = extractWear(info.HashName)
			if wear != "" {
				wearInfo := wearMap[wear]
				wearShort = wearInfo.Short
				wearOrder = wearInfo.Order
			} else {
				// 无磨损的饰品（如刀具），使用特殊标记
				wear = "NO_WEAR"
				wearShort = "-"
				wearOrder = 0
			}
		}

		// 处理空的 quality_name，从 hash_name 推断
		qualityName := info.QualityName
		if qualityName == "" {
			qualityName = inferQualityFromHashName(info.HashName)
		}

		item := RelatedWearItem{
			HashName:    info.HashName,
			Wear:        wear,
			WearShort:   wearShort,
			IconUrl:     info.IconUrl,
			Price:       priceMap[info.HashName],
			QualityName: qualityName,
		}

		qualityMap[qualityName] = append(qualityMap[qualityName], item)
		qualitiesSet[qualityName] = true
		_ = wearOrder // 避免未使用警告
	}

	// 对每个品质的磨损/变体列表按顺序排序
	// 排序：普通 < 闪亮/刺绣 < 全息 < 金色
	variantOrder := map[string]int{
		"普通": 1, "Foil": 2, "Embroidered": 2, "Holo": 3, "Gold": 4, "Champion": 5, "Lenticular": 6,
	}

	for quality := range qualityMap {
		items := qualityMap[quality]
		// 冒泡排序（简单实现）
		for i := 0; i < len(items); i++ {
			for j := i + 1; j < len(items); j++ {
				var orderI, orderJ int

				if isStickerItem {
					// 贴纸用变体顺序
					orderI = variantOrder[items[i].Wear]
					orderJ = variantOrder[items[j].Wear]
					if orderI == 0 {
						orderI = 99
					}
					if orderJ == 0 {
						orderJ = 99
					}
				} else {
					// 普通饰品用磨损顺序
					if items[i].Wear != "NO_WEAR" {
						orderI = wearMap[items[i].Wear].Order
					}
					if items[j].Wear != "NO_WEAR" {
						orderJ = wearMap[items[j].Wear].Order
					}
				}

				if orderI > orderJ {
					items[i], items[j] = items[j], items[i]
				}
			}
		}
		qualityMap[quality] = items
	}

	// 获取品质列表
	qualities := make([]string, 0, len(qualitiesSet))
	for q := range qualitiesSet {
		qualities = append(qualities, q)
	}

	return &RelatedWearsResponse{
		CurrentWear:    currentWear,
		CurrentQuality: currentQuality,
		Qualities:      qualities,
		Wears:          qualityMap,
	}, nil
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
	// 默认 uu -> buff 的搬砖数据
	sourceTable := "u"
	targetTable := "buff"

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

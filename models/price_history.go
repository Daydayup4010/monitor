package models

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"uu/config"
)

// getLocalToday 获取本地时区的今天 00:00:00
func getLocalToday() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
}

type PriceHistory struct {
	ID             uint      `gorm:"primaryKey;autoIncrement"`
	MarketHashName string    `gorm:"type:varchar(255);index:idx_hash_platform_date,priority:1;index:idx_query,priority:3;not null"`
	Platform       string    `gorm:"type:varchar(20);index:idx_hash_platform_date,priority:2;index:idx_platform_date,priority:1;index:idx_query,priority:1;not null"`
	SellPrice      float64   `json:"sellPrice" gorm:"index:idx_query,priority:4"`
	SellCount      int64     `json:"sellCount"`
	RecordDate     time.Time `gorm:"type:date;index:idx_hash_platform_date,priority:3;index:idx_date;index:idx_platform_date,priority:2;index:idx_query,priority:2;"` // 记录日期
}

// PriceIncreaseItem 带价格趋势和各时间段涨幅的数据
type PriceIncreaseItem struct {
	MarketHashName  string   `json:"marketHashName" gorm:"column:market_hash_name"`
	Name            string   `json:"name" gorm:"column:name"`
	IconUrl         string   `json:"iconUrl" gorm:"column:icon_url"`
	Platform        string   `json:"platform" gorm:"column:platform"`
	TodayPrice      float64  `json:"todayPrice" gorm:"column:today_price"`
	YesterdayPrice  float64  `json:"yesterdayPrice" gorm:"column:yesterday_price"`
	Price3DaysAgo   *float64 `json:"price3DaysAgo" gorm:"column:price_3_days_ago"`
	Price7DaysAgo   *float64 `json:"price7DaysAgo" gorm:"column:price_7_days_ago"`
	Price15DaysAgo  *float64 `json:"price15DaysAgo" gorm:"column:price_15_days_ago"`
	Price30DaysAgo  *float64 `json:"price30DaysAgo" gorm:"column:price_30_days_ago"`
	PriceChange     float64  `json:"priceChange" gorm:"column:price_change"`
	IncreaseRate1D  float64  `json:"increaseRate1D" gorm:"column:increase_rate_1d"`
	IncreaseRate3D  *float64 `json:"increaseRate3D" gorm:"column:increase_rate_3d"`
	IncreaseRate7D  *float64 `json:"increaseRate7D" gorm:"column:increase_rate_7d"`
	IncreaseRate15D *float64 `json:"increaseRate15D" gorm:"column:increase_rate_15d"`
	IncreaseRate30D *float64 `json:"increaseRate30D" gorm:"column:increase_rate_30d"`
	// 在售数相关字段
	TodaySellCount     int64    `json:"todaySellCount" gorm:"column:today_sell_count"`
	YesterdaySellCount int64    `json:"yesterdaySellCount" gorm:"column:yesterday_sell_count"`
	SellCount3DaysAgo  *int64   `json:"sellCount3DaysAgo" gorm:"column:sell_count_3_days_ago"`
	SellCount7DaysAgo  *int64   `json:"sellCount7DaysAgo" gorm:"column:sell_count_7_days_ago"`
	SellCount15DaysAgo *int64   `json:"sellCount15DaysAgo" gorm:"column:sell_count_15_days_ago"`
	SellCount30DaysAgo *int64   `json:"sellCount30DaysAgo" gorm:"column:sell_count_30_days_ago"`
	SellCountChange    int64    `json:"sellCountChange" gorm:"column:sell_count_change"`
	SellCountRate1D    float64  `json:"sellCountRate1D" gorm:"column:sell_count_rate_1d"`
	SellCountRate3D    *float64 `json:"sellCountRate3D" gorm:"column:sell_count_rate_3d"`
	SellCountRate7D    *float64 `json:"sellCountRate7D" gorm:"column:sell_count_rate_7d"`
	SellCountRate15D   *float64 `json:"sellCountRate15D" gorm:"column:sell_count_rate_15d"`
	SellCountRate30D   *float64 `json:"sellCountRate30D" gorm:"column:sell_count_rate_30d"`
}

// TableName 自定义表名
func (PriceHistory) TableName() string {
	return "price_history"
}

// PriceHistoryItem 返回给前端的历史数据项
type PriceHistoryItem struct {
	Date      string  `json:"date"`
	SellPrice float64 `json:"sellPrice"`
	SellCount int64   `json:"sellCount"`
}

// PriceHistoryResponse 历史数据响应
type PriceHistoryResponse struct {
	MarketHashName string                        `json:"marketHashName"`
	Platforms      map[string][]PriceHistoryItem `json:"platforms"` // key: platform name
}

// PriceChangeItem 价格变化项
type PriceChangeItem struct {
	Label      string  `json:"label"`      // 今日、本周、本月
	PriceDiff  float64 `json:"priceDiff"`  // 价格差
	ChangeRate float64 `json:"changeRate"` // 涨跌幅百分比
	IsUp       bool    `json:"isUp"`       // 是否上涨
}

// GoodsDetailResponse 商品详情响应
type GoodsDetailResponse struct {
	MarketHashName string                        `json:"marketHashName"`
	Name           string                        `json:"name"`
	IconUrl        string                        `json:"iconUrl"`
	RarityName     string                        `json:"rarityName"`
	QualityName    string                        `json:"qualityName"`
	PriceHistory   map[string][]PriceHistoryItem `json:"priceHistory"` // 所有平台的历史数据，key: 平台名
	PlatformList   []*GoodsPlatformInfo          `json:"platformList"` // 各平台当前在售信息
	PriceChange    []PriceChangeItem             `json:"priceChange"`  // 悠悠平台的涨幅信息（今日、本周、本月）
}

// GoodsPlatformInfo 各平台在售信息
type GoodsPlatformInfo struct {
	Platform     string  `json:"platform"`
	PlatformName string  `json:"platformName"`
	SellPrice    float64 `json:"sellPrice"`
	SellCount    int64   `json:"sellCount"`
	BiddingPrice float64 `json:"biddingPrice"`
	BiddingCount int64   `json:"biddingCount"`
	UpdateTime   int64   `json:"updateTime"`
	Link         string  `json:"link"`
}

// BatchCreatePriceHistory 批量创建历史记录
func BatchCreatePriceHistory(histories []*PriceHistory) error {
	if len(histories) == 0 {
		return nil
	}
	return config.DB.CreateInBatches(histories, 500).Error
}

// GetPriceHistoryByHashName 获取指定商品所有平台最近N天的历史记录
func GetPriceHistoryByHashName(marketHashName string, days int) (*PriceHistoryResponse, error) {
	var histories []PriceHistory
	startDate := getLocalToday().AddDate(0, 0, -days)
	err := config.DB.Where("market_hash_name = ? AND record_date >= ?", marketHashName, startDate).
		Order("platform ASC, record_date ASC").
		Find(&histories).Error

	if err != nil {
		return nil, err
	}

	response := &PriceHistoryResponse{
		MarketHashName: marketHashName,
		Platforms:      make(map[string][]PriceHistoryItem),
	}

	for _, h := range histories {
		item := PriceHistoryItem{
			Date:      h.RecordDate.Format("2006-01-02"),
			SellPrice: h.SellPrice,
			SellCount: h.SellCount,
		}
		response.Platforms[h.Platform] = append(response.Platforms[h.Platform], item)
	}

	return response, nil
}

// GetPriceHistoryByPlatform 获取指定商品指定平台最近N天的历史记录
func GetPriceHistoryByPlatform(marketHashName, platform string, days int) ([]PriceHistoryItem, error) {
	var histories []PriceHistory
	startDate := getLocalToday().AddDate(0, 0, -days)

	err := config.DB.Where("market_hash_name = ? AND platform = ? AND record_date >= ?",
		marketHashName, platform, startDate).
		Order("record_date ASC").
		Find(&histories).Error

	if err != nil {
		return nil, err
	}

	result := make([]PriceHistoryItem, 0, len(histories))
	for _, h := range histories {
		result = append(result, PriceHistoryItem{
			Date:      h.RecordDate.Format("2006-01-02"),
			SellPrice: h.SellPrice,
			SellCount: h.SellCount,
		})
	}

	return result, nil
}

// GetGoodsDetail 获取商品详情（包含基础信息、所有平台历史数据、各平台在售信息）
func GetGoodsDetail(marketHashName string, days int) (*GoodsDetailResponse, error) {
	// 1. 获取商品基础信息
	var baseGoods BaseGoods
	err := config.DB.Where("market_hash_name = ?", marketHashName).First(&baseGoods).Error
	if err != nil {
		return nil, err
	}

	var uBase UBaseInfo
	config.DB.Where("hash_name = ?", marketHashName).First(&uBase)

	// 2. 获取所有平台的历史数据
	historyResponse, err := GetPriceHistoryByHashName(marketHashName, days)
	if err != nil {
		return nil, err
	}

	// 3. 获取各平台当前在售信息
	platformList := GetAllPlatformInfo(marketHashName)

	// 4. 计算悠悠平台的涨幅信息
	priceChange := calculatePriceChange(marketHashName)

	return &GoodsDetailResponse{
		MarketHashName: marketHashName,
		Name:           baseGoods.Name,
		IconUrl:        baseGoods.IconUrl,
		RarityName:     uBase.RarityName,
		QualityName:    uBase.QualityName,
		PriceHistory:   historyResponse.Platforms,
		PlatformList:   platformList,
		PriceChange:    priceChange,
	}, nil
}

// calculatePriceChange 计算悠悠平台的涨幅信息（今日、本周、本月）
func calculatePriceChange(marketHashName string) []PriceChangeItem {
	today := getLocalToday()
	result := []PriceChangeItem{}

	// 获取当前价格（最新一天的价格）
	var currentHistory PriceHistory
	err := config.DB.Where("market_hash_name = ? AND platform = ?", marketHashName, "YOUPIN").
		Order("record_date DESC").First(&currentHistory).Error
	if err != nil {
		return result
	}
	currentPrice := currentHistory.SellPrice

	// 计算今日涨幅（与昨天对比）
	yesterday := today.AddDate(0, 0, -1)
	result = append(result, calcPriceChangeItem("今日", marketHashName, currentPrice, yesterday))

	// 计算本周涨幅（与7天前对比）
	weekAgo := today.AddDate(0, 0, -7)
	result = append(result, calcPriceChangeItem("本周", marketHashName, currentPrice, weekAgo))

	// 计算本月涨幅（与30天前对比）
	monthAgo := today.AddDate(0, 0, -30)
	result = append(result, calcPriceChangeItem("本月", marketHashName, currentPrice, monthAgo))

	return result
}

// calcPriceChangeItem 计算单个涨幅项
func calcPriceChangeItem(label, marketHashName string, currentPrice float64, compareDate time.Time) PriceChangeItem {
	var oldHistory PriceHistory
	err := config.DB.Where("market_hash_name = ? AND platform = ? AND record_date <= ?",
		marketHashName, "YOUPIN", compareDate).
		Order("record_date DESC").First(&oldHistory).Error

	if err != nil || oldHistory.SellPrice == 0 {
		return PriceChangeItem{
			Label:      label,
			PriceDiff:  0,
			ChangeRate: 0,
			IsUp:       false,
		}
	}

	priceDiff := currentPrice - oldHistory.SellPrice
	changeRate := (priceDiff / oldHistory.SellPrice) * 100

	return PriceChangeItem{
		Label:      label,
		PriceDiff:  priceDiff,
		ChangeRate: changeRate,
		IsUp:       priceDiff >= 0,
	}
}

// GetAllPlatformInfo 获取指定商品在所有平台的当前在售信息
func GetAllPlatformInfo(marketHashName string) []*GoodsPlatformInfo {
	result := make([]*GoodsPlatformInfo, 0, 4)

	// 悠悠有品
	var u U
	if err := config.DB.Where("market_hash_name = ?", marketHashName).First(&u).Error; err == nil {
		result = append(result, &GoodsPlatformInfo{
			Platform:     "YOUPIN",
			PlatformName: "悠悠有品",
			SellPrice:    u.SellPrice,
			SellCount:    u.SellCount,
			BiddingPrice: u.BiddingPrice,
			BiddingCount: u.BiddingCount,
			UpdateTime:   u.UpdateTime,
			Link:         u.Link,
		})
	}

	// BUFF
	var buff Buff
	if err := config.DB.Where("market_hash_name = ?", marketHashName).First(&buff).Error; err == nil {
		result = append(result, &GoodsPlatformInfo{
			Platform:     "BUFF",
			PlatformName: "BUFF",
			SellPrice:    buff.SellPrice,
			SellCount:    buff.SellCount,
			BiddingPrice: buff.BiddingPrice,
			BiddingCount: buff.BiddingCount,
			UpdateTime:   buff.UpdateTime,
			Link:         buff.Link,
		})
	}

	// C5GAME
	var c5 C5
	if err := config.DB.Where("market_hash_name = ?", marketHashName).First(&c5).Error; err == nil {
		result = append(result, &GoodsPlatformInfo{
			Platform:     "C5",
			PlatformName: "C5GAME",
			SellPrice:    c5.SellPrice,
			SellCount:    c5.SellCount,
			BiddingPrice: c5.BiddingPrice,
			BiddingCount: c5.BiddingCount,
			UpdateTime:   c5.UpdateTime,
			Link:         c5.Link,
		})
	}

	// Steam
	var steam Steam
	if err := config.DB.Where("market_hash_name = ?", marketHashName).First(&steam).Error; err == nil {
		result = append(result, &GoodsPlatformInfo{
			Platform:     "STEAM",
			PlatformName: "Steam",
			SellPrice:    steam.SellPrice,
			SellCount:    steam.SellCount,
			BiddingPrice: steam.BiddingPrice,
			BiddingCount: steam.BiddingCount,
			UpdateTime:   steam.UpdateTime,
			Link:         steam.Link,
		})
	}

	return result
}

// CleanOldHistory 清理超过指定天数的历史数据
func CleanOldHistory(days int) {
	cutoffDate := getLocalToday().AddDate(0, 0, -days)
	result := config.DB.Where("record_date < ?", cutoffDate).Delete(&PriceHistory{})
	if result.Error != nil {
		config.Log.Errorf("Clean old history error: %v", result.Error)
	}
}

// CheckTodayRecordExists 检查今天是否已经记录过
func CheckTodayRecordExists() bool {
	var count int64
	today := time.Now().Format("2006-01-02")
	config.DB.Model(&PriceHistory{}).
		Where("DATE(record_date) = ?", today).
		Limit(1).
		Count(&count)
	return count > 0
}

// GetTurnOverFromHistory 从历史数据计算成交量（最近一天的变化）
func GetTurnOverFromHistory(marketHashName, platform string) int64 {
	var histories []PriceHistory
	// 获取最近两天的数据
	startDate := getLocalToday().AddDate(0, 0, -2)

	err := config.DB.Where("market_hash_name = ? AND platform = ? AND record_date >= ?",
		marketHashName, platform, startDate).
		Order("record_date DESC").
		Limit(2).
		Find(&histories).Error

	if err != nil || len(histories) < 2 {
		return 0
	}

	// 计算最近一天的成交量变化
	turnOver := histories[0].SellCount - histories[1].SellCount
	if turnOver < 0 {
		turnOver = -turnOver
	}
	return turnOver
}

// getCacheExpiration 获取缓存过期时间（到第二天凌晨2点）
func getCacheExpiration() time.Duration {
	now := time.Now()
	// 第二天凌晨3点（在每日记录任务 2:00 之后）
	next := time.Date(now.Year(), now.Month(), now.Day(), 2, 0, 0, 0, now.Location())
	if next.Before(now) {
		next = next.Add(24 * time.Hour)
	}
	return next.Sub(now)
}

// GetPriceIncrease 获取涨幅排行（带多日价格趋势和各时间段涨幅）
// platform: 平台名称 (YOUPIN, BUFF, C5, STEAM)
// sortBy: 排序字段 (1d, 3d, 7d, 15d, 30d)
// isDesc: 是否降序
// limit: 返回数量
func GetPriceIncrease(platform, sortBy string, isDesc bool, limit int) ([]PriceIncreaseItem, error) {
	var results []PriceIncreaseItem
	ctx := context.Background()

	// 生成缓存 key
	cacheKey := fmt.Sprintf("price_increase:%s:%s:%v:%d", platform, sortBy, isDesc, limit)

	// 1. 先查缓存
	cached, err := config.RDB.Get(ctx, cacheKey).Result()
	if err == nil && cached != "" {
		if err := json.Unmarshal([]byte(cached), &results); err == nil {
			return results, nil
		}
	}

	// 2. 缓存未命中，查询数据库
	today := getLocalToday()
	yesterday := today.AddDate(0, 0, -1)
	day3 := today.AddDate(0, 0, -3)
	day7 := today.AddDate(0, 0, -7)
	day15 := today.AddDate(0, 0, -15)
	day30 := today.AddDate(0, 0, -30)

	// 排序方向
	sortOrder := "ASC"
	if isDesc {
		sortOrder = "DESC"
	}

	// 根据 sortBy 确定排序的基准日期
	var sortDate time.Time
	switch sortBy {
	case "3d":
		sortDate = day3
	case "7d":
		sortDate = day7
	case "15d":
		sortDate = day15
	case "30d":
		sortDate = day30
	default:
		sortDate = yesterday
	}

	// 优化：使用子查询先筛选出 Top N，再 JOIN 获取其他数据
	sql := `
		SELECT 
			top_items.market_hash_name,
			bg.name,
			bg.icon_url,
			top_items.platform,
			top_items.today_price,
			top_items.yesterday_price,
			d3.sell_price AS price_3_days_ago,
			d7.sell_price AS price_7_days_ago,
			d15.sell_price AS price_15_days_ago,
			d30.sell_price AS price_30_days_ago,
			top_items.price_change,
			top_items.increase_rate_1d,
			ROUND((top_items.today_price - d3.sell_price) / NULLIF(d3.sell_price, 0) * 100, 2) AS increase_rate_3d,
			ROUND((top_items.today_price - d7.sell_price) / NULLIF(d7.sell_price, 0) * 100, 2) AS increase_rate_7d,
			ROUND((top_items.today_price - d15.sell_price) / NULLIF(d15.sell_price, 0) * 100, 2) AS increase_rate_15d,
			ROUND((top_items.today_price - d30.sell_price) / NULLIF(d30.sell_price, 0) * 100, 2) AS increase_rate_30d,
			top_items.today_sell_count,
			top_items.yesterday_sell_count,
			d3.sell_count AS sell_count_3_days_ago,
			d7.sell_count AS sell_count_7_days_ago,
			d15.sell_count AS sell_count_15_days_ago,
			d30.sell_count AS sell_count_30_days_ago,
			top_items.sell_count_change,
			top_items.sell_count_rate_1d,
			ROUND((top_items.today_sell_count - d3.sell_count) / NULLIF(d3.sell_count, 0) * 100, 2) AS sell_count_rate_3d,
			ROUND((top_items.today_sell_count - d7.sell_count) / NULLIF(d7.sell_count, 0) * 100, 2) AS sell_count_rate_7d,
			ROUND((top_items.today_sell_count - d15.sell_count) / NULLIF(d15.sell_count, 0) * 100, 2) AS sell_count_rate_15d,
			ROUND((top_items.today_sell_count - d30.sell_count) / NULLIF(d30.sell_count, 0) * 100, 2) AS sell_count_rate_30d
		FROM (
			SELECT 
				t.market_hash_name,
				t.platform,
				t.sell_price AS today_price,
				s.sell_price AS yesterday_price,
				(t.sell_price - s.sell_price) AS price_change,
				ROUND((t.sell_price - s.sell_price) / s.sell_price * 100, 2) AS increase_rate_1d,
				t.sell_count AS today_sell_count,
				s.sell_count AS yesterday_sell_count,
				(t.sell_count - s.sell_count) AS sell_count_change,
				ROUND((t.sell_count - s.sell_count) / NULLIF(s.sell_count, 0) * 100, 2) AS sell_count_rate_1d
			FROM price_history t
			INNER JOIN price_history s 
				ON t.market_hash_name = s.market_hash_name 
				AND t.platform = s.platform
				AND s.record_date = ?
			WHERE t.platform = ?
				AND t.record_date = ?
				AND s.sell_price > 0
				AND t.sell_price > 1
				AND t.sell_count > 100
			ORDER BY increase_rate_1d ` + sortOrder + `
			LIMIT ?
		) AS top_items
		LEFT JOIN price_history d3 
			ON top_items.market_hash_name = d3.market_hash_name 
			AND top_items.platform = d3.platform
			AND d3.record_date = ?
		LEFT JOIN price_history d7 
			ON top_items.market_hash_name = d7.market_hash_name 
			AND top_items.platform = d7.platform
			AND d7.record_date = ?
		LEFT JOIN price_history d15 
			ON top_items.market_hash_name = d15.market_hash_name 
			AND top_items.platform = d15.platform
			AND d15.record_date = ?
		LEFT JOIN price_history d30 
			ON top_items.market_hash_name = d30.market_hash_name 
			AND top_items.platform = d30.platform
			AND d30.record_date = ?
		INNER JOIN base_goods bg 
			ON top_items.market_hash_name = bg.market_hash_name
		ORDER BY top_items.increase_rate_1d ` + sortOrder + `
	`

	err = config.DB.Raw(sql,
		sortDate, // s.record_date (排序基准日期)
		platform, // t.platform
		today,    // t.record_date
		limit,    // LIMIT
		day3,     // d3.record_date
		day7,     // d7.record_date
		day15,    // d15.record_date
		day30,    // d30.record_date
	).Scan(&results).Error

	if err != nil {
		config.Log.Errorf("Get price increase error: %v", err)
		return nil, err
	}

	// 3. 存入缓存
	if len(results) > 0 {
		if data, err := json.Marshal(results); err == nil {
			expiration := getCacheExpiration()
			config.RDB.Set(ctx, cacheKey, data, expiration)
		}
	}

	return results, nil
}

// ClearPriceIncreaseCache 清除涨幅缓存（在每日数据更新后调用）
func ClearPriceIncreaseCache() {
	ctx := context.Background()
	// 使用通配符删除所有涨幅缓存
	keys, err := config.RDB.Keys(ctx, "price_increase:*").Result()
	if err != nil {
		config.Log.Errorf("Get price increase cache keys error: %v", err)
		return
	}
	if len(keys) > 0 {
		config.RDB.Del(ctx, keys...)
		config.Log.Infof("Cleared %d price increase cache keys", len(keys))
	}
}

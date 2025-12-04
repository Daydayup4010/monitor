package models

import (
	"time"
	"uu/config"
)

// PriceHistory 价格历史记录表
type PriceHistory struct {
	ID             uint      `gorm:"primaryKey;autoIncrement"`
	MarketHashName string    `gorm:"type:varchar(255);index:idx_hash_platform_date,priority:1;not null"`
	Platform       string    `gorm:"type:varchar(20);index:idx_hash_platform_date,priority:2;not null"`
	SellPrice      float64   `json:"sellPrice"`
	SellCount      int64     `json:"sellCount"`
	RecordDate     time.Time `gorm:"type:date;index:idx_hash_platform_date,priority:3;index:idx_date"` // 记录日期
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
	startDate := time.Now().AddDate(0, 0, -days)

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
	startDate := time.Now().AddDate(0, 0, -days)

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

// CleanOldHistory 清理超过指定天数的历史数据
func CleanOldHistory(days int) {
	cutoffDate := time.Now().AddDate(0, 0, -days)
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
	startDate := time.Now().AddDate(0, 0, -2)

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

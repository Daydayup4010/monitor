package models

import (
	"gorm.io/gorm"
	"uu/config"
)

//type Settings struct {
//	MinSellNum   int     `json:"min_sell_num"`
//	MinDiff      float64 `json:"min_diff"`
//	MaxSellPrice float64 `json:"max_sell_price"`
//	MinSellPrice float64 `json:"min_sell_price"`
//}
//
//func (s *Settings) UpdateSetting(c context.Context) error {
//	err := config.RDB.HSet(c, "settings", "min_sell_num", s.MinSellNum, "min_diff", s.MinDiff, "max_sell_price", s.MaxSellPrice, "min_sell_price", s.MinSellPrice).Err()
//	return err
//}
//
//func (s *Settings) GetSettings(c context.Context) error {
//	// 使用HGetAll一次性获取所有字段
//	result, err := config.RDB.HGetAll(c, "settings").Result()
//	if err != nil {
//		return err
//	}
//
//	if minSellNumStr, ok := result["min_sell_num"]; ok && minSellNumStr != "" {
//		if num, err := strconv.Atoi(minSellNumStr); err == nil {
//			s.MinSellNum = num
//		}
//	}
//
//	if minDiffStr, ok := result["min_diff"]; ok && minDiffStr != "" {
//		if num, err := strconv.ParseFloat(minDiffStr, 64); err == nil {
//			s.MinDiff = num
//		}
//	}
//
//	if minSellPriceStr, ok := result["min_sell_price"]; ok && minSellPriceStr != "" {
//		if num, err := strconv.ParseFloat(minSellPriceStr, 64); err == nil {
//			s.MinSellPrice = num
//		}
//	}
//
//	if maxSellPriceStr, ok := result["max_sell_price"]; ok && maxSellPriceStr != "" {
//		if num, err := strconv.ParseFloat(maxSellPriceStr, 64); err == nil {
//			s.MaxSellPrice = num
//		}
//	}
//
//	return nil
//}

type Settings struct {
	gorm.Model
	UserId       string  `json:"user_id" gorm:"type:char(36);uniqueIndex:idx_user_setting"`
	MinSellNum   int     `json:"min_sell_num" gorm:"default:200"`
	MinDiff      float64 `json:"min_diff" gorm:"default:1"`
	MaxSellPrice float64 `json:"max_sell_price" gorm:"default:10000"`
	MinSellPrice float64 `json:"min_sell_price" gorm:"default:0"`
}

type SettingsResponse struct {
	UserId       string  `json:"user_id"`
	MinSellNum   int     `json:"min_sell_num"`
	MinDiff      float64 `json:"min_diff"`
	MaxSellPrice float64 `json:"max_sell_price"`
	MinSellPrice float64 `json:"min_sell_price"`
}

func CreateDefaultSetting(id string) error {
	var setting Settings
	setting.UserId = id
	err := config.DB.Create(&setting).Error
	return err
}

func GetUserSetting(id string) (*SettingsResponse, error) {
	var setting SettingsResponse
	err := config.DB.Model(&Settings{}).Where("user_id = ?", id).First(&setting).Error
	return &setting, err
}

func UpdateSetting(id string, setting Settings) error {
	err := config.DB.Where("user_id = ?", id).Updates(&setting).Error
	return err
}

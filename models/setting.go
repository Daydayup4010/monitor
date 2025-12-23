package models

import (
	"gorm.io/gorm"
	"uu/config"
	"uu/utils"
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
	MinSellNum   int     `json:"min_sell_num" gorm:"default:100"`
	MinDiff      float64 `json:"min_diff" gorm:"default:0.8"`
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

func CreateDefaultSetting(id string) int {
	var setting Settings
	setting.UserId = id
	err := config.DB.Create(&setting).Error
	if err != nil {
		config.Log.Errorf("create default setting error: %v", err)
		return utils.ErrCodeCreateDefaultSetting
	}
	return utils.SUCCESS
}

func GetUserSetting(id string) (*SettingsResponse, int) {
	var setting SettingsResponse
	var code int
	err := config.DB.Model(&Settings{}).Where("user_id = ?", id).First(&setting).Error
	if err != nil {
		config.Log.Errorf("get settings user: %s err: %s", id, err)
		code = utils.ErrCodeGetSettings
	} else {
		code = utils.SUCCESS
	}
	return &setting, code
}

func UpdateSetting(id string, setting Settings) int {
	err := config.DB.Where("user_id = ?", id).Updates(&setting).Error
	if err != nil {
		config.Log.Errorf("update setting err: %s", err)
		return utils.ErrCodeUpdateSetting
	}
	return utils.SUCCESS
}

// GetAdminSetting 获取管理员的设置（用于公开首页）
func GetAdminSetting() (*SettingsResponse, int) {
	var setting SettingsResponse
	// 查找管理员用户的setting
	err := config.DB.Model(&Settings{}).
		Joins("JOIN user ON settings.user_id = user.id").
		Where("user.role = ?", RoleAdmin).
		First(&setting).Error
	if err != nil {
		config.Log.Errorf("get admin settings err: %s", err)
		// 返回默认设置
		return &SettingsResponse{
			MinSellNum:   200,
			MinDiff:      1,
			MaxSellPrice: 10000,
			MinSellPrice: 0,
		}, utils.SUCCESS
	}
	return &setting, utils.SUCCESS
}

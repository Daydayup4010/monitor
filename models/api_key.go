package models

import (
	"encoding/json"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"os"
	"time"
	"uu/config"
)

type APIKey struct {
	ID        uint       `gorm:"primaryKey"`
	Key       string     `gorm:"type:varchar(255);uniqueIndex;not null"`
	LastUsed  *time.Time `gorm:"index"`
	FailCount int        `gorm:"default:0"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type APIKeyConfig struct {
	Keys []string `json:"key"`
}

func LoadAPIKeys(filename string) ([]string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		config.Log.Errorf("Read key json file error:%v", err)
		return nil, err
	}

	var keyConfig APIKeyConfig
	err = json.Unmarshal(data, &keyConfig)
	if err != nil {
		config.Log.Errorf("Unmarshal key json fail: %v", err)
		return nil, err
	}

	return keyConfig.Keys, nil
}

func InitKeys() {
	Keys, err := LoadAPIKeys("key.json")
	if err != nil {
		return
	}
	var ApiKeys []APIKey
	for _, key := range Keys {
		ApiKeys = append(ApiKeys, APIKey{Key: key})
	}
	err = config.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "key"}},
			DoUpdates: clause.AssignmentColumns([]string{"updated_at"}),
		}).Create(&ApiKeys).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		config.Log.Errorf("Init Api key fail: %v", err)
		return
	}
	config.Log.Info("Init Api key success")
}

func GetActivateKey() []APIKey {
	var apiKey []APIKey
	err := config.DB.Where("last_useD is NULL or last_useD < ?", time.Now().Add(-60*time.Second)).Find(&apiKey).Error
	if err != nil {
		config.Log.Errorf("Get activate key error: %v", err)
	}
	return apiKey
}

func UpdateLastUsed(key *APIKey) {
	err := config.DB.Model(&key).Update("last_used", time.Now()).Error
	if err != nil {
		config.Log.Errorf("Update %v last_used fail: %v", key, err)
	}
}

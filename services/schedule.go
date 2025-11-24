package services

import (
	"sync"
	"time"
	"uu/config"
)

func StartBuffFullUpdateScheduler() {
	go UpdateBuffFullData()
	ticker := time.NewTicker(30 * time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		go UpdateBuffFullData()
	}
}

func StartVerifyToken() {
	ticker := time.NewTicker(2 * time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		go func() {
			VerifyBuffToken()
			VerifyUUToken()
		}()
	}
}

func StartUUFullUpdateScheduler() {
	go UpdateUUFullData()
	ticker := time.NewTicker(8 * time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		go UpdateUUFullData()
	}
}

func UpdateBaseGoodsScheduler() {
	go UpdateBaseGoodsToDb()
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()
	for range ticker.C {
		go UpdateBaseGoodsToDb()
	}
}

func UpdateAllGoodsScheduler() {
	ticker := time.NewTicker(100 * time.Second)
	defer ticker.Stop()
	var running bool
	var mutex sync.Mutex

	for range ticker.C {
		mutex.Lock()
		if running {
			config.Log.Info("上一个任务还在执行")
			mutex.Unlock()
			continue
		}
		running = true
		mutex.Unlock()
		go func() {
			defer func() {
				mutex.Lock()
				running = false
				config.Log.Info("任务执行完，释放锁")
				mutex.Unlock()
			}()
			UpdateAllPlatformData()
		}()
	}
}

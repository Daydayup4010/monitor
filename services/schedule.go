package services

import (
	"sync"
	"time"
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
			mutex.Unlock()
			continue
		}
		running = true
		mutex.Unlock()
		go func() {
			defer func() {
				mutex.Lock()
				running = false
				mutex.Unlock()
			}()
			UpdateAllPlatformData()
		}()
	}
}

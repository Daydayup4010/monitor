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
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()
	for range ticker.C {
		SafeGo(UpdateBaseGoodsToDb)
	}
}

func UpdateIconScheduler() {
	SafeGo(UpdateUUGoods)
	ticker := time.NewTicker(13 * time.Hour)
	defer ticker.Stop()
	for range ticker.C {
		SafeGo(UpdateUUGoods)
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
			SafeGo(UpdateAllPlatformData)
		}()
	}
}

// StartDailyPriceHistoryScheduler 每天凌晨2点记录一次价格历史
func StartDailyPriceHistoryScheduler() {
	// 计算距离下一个凌晨2点的时间
	now := time.Now()
	next := time.Date(now.Year(), now.Month(), now.Day(), 2, 0, 0, 0, now.Location())
	if next.Before(now) {
		next = next.Add(24 * time.Hour)
	}

	// 先等待到第一次执行时间
	waitDuration := next.Sub(now)
	time.Sleep(waitDuration)

	// 执行第一次记录
	SafeGo(RecordDailyPriceHistory)

	// 然后每24小时执行一次
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()
	for range ticker.C {
		SafeGo(RecordDailyPriceHistory)
	}
}

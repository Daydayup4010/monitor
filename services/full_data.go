package services

import (
	"context"
	"strconv"
	"strings"
	"sync"
	"time"
	"uu/config"
	"uu/models"

	"golang.org/x/time/rate"
)

var UUMaxPageSize = 100
var BuffMaxPageSize = "60"
var taskBuff sync.Mutex
var taskUU sync.Mutex
var uuLimiter = rate.NewLimiter(3, 1) // 速率: 3 tokens/s, 突发容量: 1
var buffLimiter = rate.NewLimiter(rate.Every(2000*time.Millisecond), 1)
var RequestDelay = time.Second * 10
var newLimiter = rate.NewLimiter(rate.Every(310*time.Second), 1)

func UpdateAllUUItems() {
	_, total, _ := GetUUItems(20, 1)
	totalPages := total/UUMaxPageSize + 1
	for page := 1; page <= totalPages+1; page++ {
		if err := uuLimiter.Wait(context.Background()); err != nil {
			config.Log.Errorf("wait limiter error: %v", err)
			return
		}

		items, _, err := GetUUItems(UUMaxPageSize, page)
		if err != nil {
			if isRateLimitError(err) {
				// 动态退避：遇到429时增加延迟
				handleRateLimitError()
				page-- // 重试当前页
				continue
			}
			config.Log.Errorf("request uu %d page fail: %v", page, err)
			continue
		}
		models.BatchAddUUItem(items)
		//config.Log.Infof("Full Update uu item pageName: %d, success", page)
	}
}

// 判断是否为速率限制错误
func isRateLimitError(err error) bool {
	return strings.Contains(err.Error(), "429") ||
		strings.Contains(err.Error(), "too many requests")
}

func handleRateLimitError() {
	config.Log.Warnf("limits sleep %v", RequestDelay)
	time.Sleep(RequestDelay)
}

func UpdateAllBuffItems() {
	_, total, _ := GetBuffItems("20", "1")
	size, _ := strconv.Atoi(BuffMaxPageSize)
	pageNum := total/size + 1
	for i := 1; i <= pageNum+1; i++ {
		if err := buffLimiter.Wait(context.Background()); err != nil {
			config.Log.Errorf("wait buff limiter error: %v", err)
			return
		}
		items, _, err := GetBuffItems(BuffMaxPageSize, strconv.Itoa(i))
		if err != nil {
			if isRateLimitError(err) {
				// 动态退避：遇到429时增加延迟
				handleRateLimitError()
				i-- // 重试当前页
				continue
			}
			config.Log.Errorf("request buff %d page fail: %v", i, err)
			continue
		}
		config.Log.Infof("Full Update buff item pageName: %d, success", i)
		models.BatchAddBuffItem(items)
	}
}

func UpdateUUGoods() {
	hashNames := models.QueryAllUUHashName()
	n := len(hashNames) / 200
	remainder := len(hashNames) % 200
	if remainder > 0 {
		n++
	}
	for i := 0; i < n; i++ {
		if err := newLimiter.Wait(context.Background()); err != nil {
			config.Log.Errorf("wait limiter error: %v", err)
			return
		}
		start := i * 200
		end := start + 200
		goods := GetUUGoods(hashNames[start:end])
		models.BatchUpdateIcon(goods)
	}
	models.UpdateBaseGoodsIcon()
}

//func UpdateFullData() int {
//	if !task.TryLock() {
//		config.Log.Info("full update running")
//		return utils.ErrCodeFullUpdateRunning
//	}
//	defer task.Unlock()
//	config.Log.Info("Start full update")
//	wg.Add(2)
//	go UpdateAllUUItems()
//	go UpdateAllBuffItems()
//	wg.Wait()
//	config.Log.Info("Full update completed")
//	return utils.SUCCESS
//}

func UpdateInventory() {
	uu := GetUUInventory()
	models.BatchAddUUInventory(uu)
	buff := GetBuffInventory()
	models.BatchAddBuffInventory(buff)
}

func UpdateUUFullData() {
	if !taskUU.TryLock() {
		config.Log.Info("uu full update running")
	}
	defer taskUU.Unlock()
	config.Log.Info("Start uu full update")
	UpdateAllUUItems()
	config.Log.Info("uu full update completed")
}

func UpdateBuffFullData() {
	if !taskBuff.TryLock() {
		config.Log.Info("buff full update running")
	}
	defer taskBuff.Unlock()
	config.Log.Info("Start buff full update")
	UpdateAllBuffItems()
	config.Log.Info("buff full update completed")
}

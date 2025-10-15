package services

import (
	"context"
	"golang.org/x/time/rate"
	"strconv"
	"strings"
	"sync"
	"time"
	"uu/config"
	"uu/models"
)

var UUMaxPageSize = 100
var BuffMaxPageSize = "80"
var RequestDelay = time.Second * 8
var wg sync.WaitGroup
var task sync.Mutex
var uuLimiter = rate.NewLimiter(1, 1) // 速率: 3 tokens/s, 突发容量: 1
var buffLimiter = rate.NewLimiter(rate.Every(2500*time.Millisecond), 1)

func UpdateAllUUItems() {
	defer wg.Done()
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

// 处理速率限制错误（动态退避）
func handleRateLimitError() {
	config.Log.Warnf("limits sleep %v", RequestDelay)
	time.Sleep(RequestDelay)
}

func UpdateAllBuffItems() {
	defer wg.Done()
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
		//config.Log.Infof("Full Update buff item pageName: %d, success", i)
		models.BatchAddBuffItem(items)
	}
}

func UpdateFullData() {
	if !task.TryLock() {
		config.Log.Info("update full data running")
		return
	}
	defer task.Unlock()
	config.Log.Info("Start full update")
	wg.Add(2)
	go UpdateAllUUItems()
	go UpdateAllBuffItems()
	wg.Wait()
	models.UpdateSkinItems()
	config.Log.Info("Full update completed")
}

func UpdateInventory() {
	uu := GetUUInventory()
	models.BatchAddUUInventory(uu)
	buff := GetBuffInventory()
	models.BatchAddBuffInventory(buff)
}

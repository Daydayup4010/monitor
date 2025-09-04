package services

import (
	"strconv"
	"sync"
	"time"
	"uu/config"
	"uu/models"
)

var UUMaxPageSize = 100
var BuffMaxPageSize = "80"
var RequestDelay = time.Second * 2
var wg sync.WaitGroup

func UpdateAllUUItems() {
	defer wg.Done()
	_, total, _ := GetUUItems(20, 1)
	pageNum := total/UUMaxPageSize + 1
	for i := 1; i <= pageNum+1; i++ {
		items, _, err := GetUUItems(UUMaxPageSize, i)
		models.BatchAddUUItem(items)
		if err == nil {
			config.Log.Infof("Full Update uu item pageName: %d, success", i)
		}
		time.Sleep(RequestDelay)
	}
}

func UpdateAllBuffItems() {
	//defer wg.Done()
	_, total, _, _ := GetBuffItems("20", "1")
	size, _ := strconv.Atoi(BuffMaxPageSize)
	pageNum := total/size + 1
	config.Log.Infof("buff page_num: %d", pageNum)
	for i := 1; i <= pageNum+1; i++ {
		items, _, err, code := GetBuffItems(BuffMaxPageSize, strconv.Itoa(i))
		if err == nil {
			config.Log.Infof("Full Update buff item pageName: %d, success", i)
		} else {
			if code == 429 {
				config.Log.Infof("Full Update buff item pageName: %d, retry", i)
				time.Sleep(time.Second * 10)
				items, _, err, code = GetBuffItems(BuffMaxPageSize, strconv.Itoa(i))
			}
		}
		models.BatchAddBuffItem(items)
		time.Sleep(RequestDelay)
	}
}

func UpdateFullData() chan bool {
	var c chan bool
	go func() {
		wg.Add(2)
		go UpdateAllUUItems()
		go UpdateAllBuffItems()
		wg.Wait()
		c <- true
	}()
	return c
}

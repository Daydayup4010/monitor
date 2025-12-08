package services

import (
	"fmt"
	"math"
	"time"
	"uu/config"
	"uu/models"
	"uu/utils"
)

var steamClient = utils.CreateClient("https://open.steamdt.com")

type BaseGoodsResponse struct {
	Data      []*models.BaseGoods `json:"data"`
	ErrorCode int64               `json:"errorCode"`
}

type BatchPriceResponse struct {
	Success   bool         `json:"success"`
	Data      []*PriceList `json:"data"`
	ErrorCode int64        `json:"errorCode"`
	ErrorMsg  string       `json:"errorMsg"`
}

type PriceList struct {
	MarketHashName string      `json:"marketHashName"`
	DataList       []*Platform `json:"dataList"`
}

type Platform struct {
	Platform       string  `json:"platform"`
	PlatformItemId string  `json:"platformItemId"`
	SellPrice      float64 `json:"sellPrice"`
	SellCount      int64   `json:"sellCount"`
	BiddingPrice   float64 `json:"biddingPrice"`
	BiddingCount   int64   `json:"biddingCount"`
	UpdateTime     int64   `json:"updateTime"`
}

func getHeader() map[string]string {
	return map[string]string{
		"Authorization": "Bearer " + config.CONFIG.SteamDt.Key,
		"Content-Type":  "application/json",
	}
}

func getHeaderFormatKey(key string) map[string]string {
	return map[string]string{
		"Authorization": "Bearer " + key,
		"Content-Type":  "application/json",
	}
}

func GetBaseGoods() ([]*models.BaseGoods, error) {
	header := getHeader()
	var base BaseGoodsResponse
	var opts = utils.RequestOptions{
		Headers: header,
		Result:  &base,
	}
	rep, err := steamClient.DoRequest("GET", "open/cs2/v1/base", opts)
	if err != nil || rep.StatusCode() != 200 {
		config.Log.Errorf("Request steamDT base api error:%v, code: %v", err, rep.StatusCode())
		return base.Data, err
	}
	if base.ErrorCode == 4005 {
		config.Log.Warningf("Request api %s limit", "open/cs2/v1/base")
		return base.Data, fmt.Errorf("request api %s limit", "open/cs2/v1/base")
	}
	return base.Data, err
}

func UpdateBaseGoodsToDb() {
	goods, err := GetBaseGoods()
	if err != nil {
		return
	}
	models.UpdateBaseGoods(goods)
}

func BatchGetPrice() ([]*PriceList, error) {
	var allPrice []*PriceList
	hashNames, err := models.GetHashNames()
	if err != nil {
		config.Log.Errorf("Get hash name error: %v", err)
		return allPrice, err
	}
	index := models.GetLastIndex()
	n := len(hashNames) / 100
	remainder := len(hashNames) % 100
	if remainder > 0 {
		n++
	}

	keys := models.GetActivateKey()
	if len(keys) == 0 {
		return allPrice, fmt.Errorf("no activate key")
	}

	for i := index; i < n; i++ {
		var rep BatchPriceResponse
		start := i * 100
		end := start + 100
		if end > len(hashNames) {
			end = len(hashNames)
		}
		key := keys[0]
		if len(keys) > 1 {
			keys = keys[1:]
			if end >= len(hashNames) {
				i = 0
				continue
			}
		} else {
			models.UpdateLastUsed(&key)
			keys = models.GetActivateKey()
			if len(keys) == 0 {
				models.SetLastIndex(i)
				return allPrice, fmt.Errorf("no activate key")
			}
		}
		hashList := hashNames[start:end]
		header := getHeaderFormatKey(key.Key)
		opts := utils.RequestOptions{
			Headers: header,
			Body: map[string][]string{
				"marketHashNames": hashList,
			},
			Result: &rep,
		}
		res, err := steamClient.DoRequest("POST", "open/cs2/v1/price/batch", opts)
		if err != nil || res.StatusCode() != 200 {
			config.Log.Errorf("Request open/cs2/v1/price/batch error: %v", err)
		}
		if rep.ErrorCode == 4005 {
			config.Log.Warningf("Request api %s limit, key: %s", "open/cs2/v1/price/batch", key)
			config.Log.Info(rep.ErrorMsg)
		} else {
			models.UpdateLastUsed(&key)
		}
		allPrice = append(allPrice, rep.Data...)
	}

	return allPrice, err
}

func UpdateAllPlatformData() {
	var uList []*models.U
	var buffList []*models.Buff
	var c5List []*models.C5
	var steamList []*models.Steam
	all, err := BatchGetPrice()
	if err != nil && len(all) == 0 {
		return
	}

	hashNames, _ := models.GetHashNames()
	uMap := models.BatchGetUUGoods(hashNames)
	buffMap := models.BatchGetBuffGoods(hashNames)
	c5Map := models.BatchGetC5Goods(hashNames)
	steamMap := models.BatchGetSteamGoods(hashNames)

	for _, data := range all {
		dataList := data.DataList
		for i, _ := range dataList {
			switch dataList[i].Platform {
			case "YOUPIN":
				u := uMap[data.MarketHashName]
				if u == nil {
					u = &models.U{}
				}
				if dataList[i].UpdateTime-u.BeforeTime >= 43200 {
					turnOver := int64(math.Abs(float64(dataList[i].SellCount - u.BeforeCount)))
					u.BeforeTime = dataList[i].UpdateTime
					u.BeforeCount = dataList[i].SellCount
					u.TurnOver = turnOver
				}
				u.Id = dataList[i].PlatformItemId
				u.MarketHashName = data.MarketHashName
				u.SellPrice = dataList[i].SellPrice
				u.SellCount = dataList[i].SellCount
				u.BiddingPrice = dataList[i].BiddingPrice
				u.BiddingCount = dataList[i].BiddingCount
				u.UpdateTime = dataList[i].UpdateTime
				u.Link = fmt.Sprintf("https://www.youpin898.com/market/goods-list?listType=10&templateId=%s&gameId=730", dataList[i].PlatformItemId)
				uList = append(uList, u)
			case "BUFF":
				buff := buffMap[data.MarketHashName]
				if buff == nil {
					buff = &models.Buff{}
				}
				if dataList[i].UpdateTime-buff.BeforeTime >= 43200 {
					turnOver := int64(math.Abs(float64(dataList[i].SellCount - buff.BeforeCount)))
					buff.BeforeTime = dataList[i].UpdateTime
					buff.BeforeCount = dataList[i].SellCount
					buff.TurnOver = turnOver
				}
				buff.Id = dataList[i].PlatformItemId
				buff.MarketHashName = data.MarketHashName
				buff.SellPrice = dataList[i].SellPrice
				buff.SellCount = dataList[i].SellCount
				buff.BiddingPrice = dataList[i].BiddingPrice
				buff.BiddingCount = dataList[i].BiddingCount
				buff.UpdateTime = dataList[i].UpdateTime
				buff.Link = fmt.Sprintf("https://buff.163.com/goods/%s?from=market#tab=selling", dataList[i].PlatformItemId)
				buffList = append(buffList, buff)
			case "C5":
				c5 := c5Map[data.MarketHashName]
				if c5 == nil {
					c5 = &models.C5{}
				}
				if dataList[i].UpdateTime-c5.BeforeTime >= 43200 {
					turnOver := int64(math.Abs(float64(dataList[i].SellCount - c5.BeforeCount)))
					c5.BeforeTime = dataList[i].UpdateTime
					c5.BeforeCount = dataList[i].SellCount
					c5.TurnOver = turnOver
				}
				c5.Id = dataList[i].PlatformItemId
				c5.MarketHashName = data.MarketHashName
				c5.SellPrice = dataList[i].SellPrice
				c5.SellCount = dataList[i].SellCount
				c5.BiddingPrice = dataList[i].BiddingPrice
				c5.BiddingCount = dataList[i].BiddingCount
				c5.UpdateTime = dataList[i].UpdateTime
				c5.Link = fmt.Sprintf("https://www.c5game.com/csgo/%s/%s/sell", dataList[i].PlatformItemId, data.MarketHashName)
				c5List = append(c5List, c5)
			case "STEAM":
				steam := steamMap[data.MarketHashName]
				if steam == nil {
					steam = &models.Steam{}
				}
				if dataList[i].UpdateTime-steam.BeforeTime >= 43200 {
					turnOver := int64(math.Abs(float64(dataList[i].SellCount - steam.BeforeCount)))
					steam.BeforeCount = dataList[i].SellCount
					steam.BeforeTime = dataList[i].UpdateTime
					steam.TurnOver = turnOver
				}
				steam.MarketHashName = data.MarketHashName
				steam.SellPrice = dataList[i].SellPrice
				steam.SellCount = dataList[i].SellCount
				steam.BiddingPrice = dataList[i].BiddingPrice
				steam.BiddingCount = dataList[i].BiddingCount
				steam.UpdateTime = dataList[i].UpdateTime
				steam.Link = fmt.Sprintf("https://steamcommunity.com/market/listings/730/%s", data.MarketHashName)
				steamList = append(steamList, steam)
			}
		}
	}
	models.BatchUpdateUUGoods(uList)
	models.BatchUpdateBuffGoods(buffList)
	models.BatchUpdateC5Goods(c5List)
	models.BatchUpdateSteamGoods(steamList)
}

// RecordDailyPriceHistory 每天记录一次价格历史
func RecordDailyPriceHistory() {
	// 检查今天是否已记录
	if models.CheckTodayRecordExists() {
		config.Log.Info("Today's price history already recorded, skip.")
		return
	}

	// 获取本地时区的今天 00:00:00
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	var histories []*models.PriceHistory

	hashNames, err := models.GetHashNames()
	if err != nil {
		config.Log.Errorf("Get hash names error: %v", err)
		return
	}

	// 获取各平台当前数据
	uMap := models.BatchGetUUGoods(hashNames)
	buffMap := models.BatchGetBuffGoods(hashNames)
	c5Map := models.BatchGetC5Goods(hashNames)
	steamMap := models.BatchGetSteamGoods(hashNames)

	for _, u := range uMap {
		if u.SellPrice > 0 { // 只记录有价格的数据
			histories = append(histories, &models.PriceHistory{
				MarketHashName: u.MarketHashName,
				Platform:       "YOUPIN",
				SellPrice:      u.SellPrice,
				SellCount:      u.SellCount,
				RecordDate:     today,
			})
		}
	}

	for _, buff := range buffMap {
		if buff.SellPrice > 0 {
			histories = append(histories, &models.PriceHistory{
				MarketHashName: buff.MarketHashName,
				Platform:       "BUFF",
				SellPrice:      buff.SellPrice,
				SellCount:      buff.SellCount,
				RecordDate:     today,
			})
		}
	}

	for _, c5 := range c5Map {
		if c5.SellPrice > 0 {
			histories = append(histories, &models.PriceHistory{
				MarketHashName: c5.MarketHashName,
				Platform:       "C5",
				SellPrice:      c5.SellPrice,
				SellCount:      c5.SellCount,
				RecordDate:     today,
			})
		}
	}

	for _, steam := range steamMap {
		if steam.SellPrice > 0 {
			histories = append(histories, &models.PriceHistory{
				MarketHashName: steam.MarketHashName,
				Platform:       "STEAM",
				SellPrice:      steam.SellPrice,
				SellCount:      steam.SellCount,
				RecordDate:     today,
			})
		}
	}

	// 批量保存
	if err := models.BatchCreatePriceHistory(histories); err != nil {
		config.Log.Errorf("Record daily price history error: %v", err)
		return
	}
	config.Log.Infof("Recorded %d price history entries for today", len(histories))

	// 清理超过一年的旧数据
	models.CleanOldHistory(366)

	// 清除涨幅缓存，让下次查询获取最新数据
	models.ClearPriceIncreaseCache()
}

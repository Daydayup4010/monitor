package services

import (
	"fmt"
	"math"
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

	n := len(hashNames) / 100
	remainder := len(hashNames) % 10
	if remainder > 0 {
		n++
	}

	keys := models.GetActivateKey()
	if len(keys) == 0 {
		config.Log.Warnf("no activate key")
		return allPrice, fmt.Errorf("no activate key")
	}

	for i := 0; i < n; i++ {
		var rep BatchPriceResponse
		start := i * 100
		end := start + 100
		if end > len(hashNames) {
			end = len(hashNames)
		}
		hashList := hashNames[start:end]
		key := keys[0]
		if len(keys) != 1 {
			keys = keys[1:]
		} else {
			keys = models.GetActivateKey()
			if len(keys) == 0 {
				config.Log.Warnf("no activate key")
				return allPrice, fmt.Errorf("no activate key")
			}
		}

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
			config.Log.Errorf("Request open/cs2/v1/price/batch error:%v, code: %v", err, res.StatusCode())
		}
		if rep.ErrorCode == 4005 {
			config.Log.Warningf("Request api %s limit", "open/cs2/v1/price/batch")
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
	for _, data := range all {
		dataList := data.DataList
		for i, _ := range dataList {
			switch dataList[i].Platform {
			case "YOUPIN":
				u := models.GetUUGoods(data.MarketHashName)
				if dataList[i].UpdateTime-u.UpdateTime >= 86400 {
					turnOver := int64(math.Abs(float64(dataList[i].SellCount - u.SellCount)))
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
				buff := models.GetBuffGoods(data.MarketHashName)
				if dataList[i].UpdateTime-buff.UpdateTime >= 86400 {
					turnOver := int64(math.Abs(float64(dataList[i].SellCount - buff.SellCount)))
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
				c5 := models.GetC5Goods(data.MarketHashName)
				if dataList[i].UpdateTime-c5.UpdateTime >= 86400 {
					turnOver := int64(math.Abs(float64(dataList[i].SellCount - c5.SellCount)))
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
				steam := models.GetSteamGoods(data.MarketHashName)
				if dataList[i].UpdateTime-steam.UpdateTime >= 86400 {
					turnOver := int64(math.Abs(float64(dataList[i].SellCount - steam.SellCount)))
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

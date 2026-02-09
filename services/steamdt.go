package services

import (
	"fmt"
	"math"
	neturl "net/url"
	"regexp"
	"strconv"
	"strings"
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
				steam.Id = dataList[i].PlatformItemId
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

// Steam 社区市场客户端
var steamCommunityClient = utils.CreateClient("https://steamcommunity-a.akamaihd.net")

// FetchSteamItemNameId 从 Steam 商品详情页获取 item_nameid
// marketHashName: 商品的 market_hash_name，如 "AK-47 | Redline (Field-Tested)"
func FetchSteamItemNameId(marketHashName string) (string, error) {
	// 直接构建路径，不依赖数据库的 link 字段
	// URL 编码 marketHashName
	encodedName := neturl.PathEscape(marketHashName)
	url := "/market/listings/730/" + encodedName

	var result string
	opts := utils.RequestOptions{
		Result: &result,
		Headers: map[string]string{
			"User-Agent":      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
			"Accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8",
			"Accept-Language": "en-US,en;q=0.5",
		},
	}

	resp, err := steamCommunityClient.DoRequest("GET", url, opts)
	if err != nil {
		return "", err
	}
	if resp.StatusCode() != 200 {
		return "", fmt.Errorf("steam response status: %d", resp.StatusCode())
	}

	html := resp.String()

	// 正则提取 item_nameid
	// 方法1: ItemActivityTicker.Start( 730, 'item_nameid' )
	re1 := regexp.MustCompile(`ItemActivityTicker\.Start\(\s*\d+\s*,\s*'(\d+)'`)
	if matches := re1.FindStringSubmatch(html); len(matches) > 1 {
		return matches[1], nil
	}

	// 方法2: Market_LoadOrderSpread( item_nameid )
	re2 := regexp.MustCompile(`Market_LoadOrderSpread\(\s*(\d+)`)
	if matches := re2.FindStringSubmatch(html); len(matches) > 1 {
		return matches[1], nil
	}

	return "", fmt.Errorf("item_nameid not found in page")
}

// UpdateSteamItemNameIds 批量更新 Steam 表中的 item_nameid
// limit: 每次执行的最大数量，0 表示不限制（建议设置 500-1000，分多次执行）
func UpdateSteamItemNameIds() {
	config.Log.Info("Starting to update Steam item_nameid...")

	// 获取所有没有 item_nameid 的商品
	steams, err := models.GetSteamsWithoutItemNameId()
	if err != nil {
		config.Log.Errorf("Failed to get steams without item_nameid: %v", err)
		return
	}

	if len(steams) == 0 {
		config.Log.Info("All steam items already have item_nameid")
		return
	}

	total := len(steams)

	config.Log.Infof("Found %d items without item_nameid, processing %d this time", total, len(steams))

	// 收集待更新的数据，批量写入
	pendingUpdates := make(map[string]string) // marketHashName -> itemNameId
	batchSize := 20
	successCount := 0
	failCount := 0
	consecutive429 := 0 // 连续 429 错误计数

	for i, item := range steams {
		// 尝试获取，遇到 429 会重试
		var itemNameId string
		var fetchErr error
		for retry := 0; retry < 3; retry++ {
			itemNameId, fetchErr = FetchSteamItemNameId(item.MarketHashName)
			if fetchErr == nil {
				consecutive429 = 0 // 成功，重置计数
				break
			}
			// 检查是否是 429 错误
			if strings.Contains(fetchErr.Error(), "429") {
				consecutive429++
				waitTime := time.Duration(30*(retry+1)) * time.Second // 30s, 60s, 90s
				config.Log.Warnf("[%d/%d] Rate limited (429), waiting %v before retry %d/3...", i+1, len(steams), waitTime, retry+1)
				time.Sleep(waitTime)
			} else {
				break // 非 429 错误不重试
			}
		}

		if fetchErr != nil {
			config.Log.Warnf("[%d/%d] Failed to get item_nameid for %s: %v", i+1, len(steams), item.MarketHashName, fetchErr)
			failCount++
			// 如果连续多次 429，暂停更长时间
			if consecutive429 >= 5 {
				config.Log.Warnf("Too many rate limits, waiting 5 minutes...")
				time.Sleep(5 * time.Minute)
				consecutive429 = 0
			}
		} else {
			pendingUpdates[item.MarketHashName] = itemNameId
			config.Log.Infof("[%d/%d] Fetched %s -> %s", i+1, len(steams), item.MarketHashName, itemNameId)
		}

		// 每 batchSize 条批量写入一次
		if len(pendingUpdates) >= batchSize {
			count, err := models.BatchUpdateSteamItemNameIds(pendingUpdates)
			if err != nil {
				config.Log.Warnf("Batch update failed: %v", err)
				failCount += len(pendingUpdates)
			} else {
				successCount += count
			}
			pendingUpdates = make(map[string]string)
		}

		time.Sleep(3000 * time.Millisecond)
	}

	// 处理剩余的数据
	if len(pendingUpdates) > 0 {
		count, err := models.BatchUpdateSteamItemNameIds(pendingUpdates)
		if err != nil {
			config.Log.Warnf("Batch update failed: %v", err)
			failCount += len(pendingUpdates)
		} else {
			successCount += count
		}
	}

	config.Log.Infof("Update steam item_nameid complete. Success: %d, Failed: %d, Remaining: %d", successCount, failCount, total-len(steams))
}

// SteamOrderHistogram Steam 市场订单数据响应结构
type SteamOrderHistogram struct {
	Success          int    `json:"success"`
	LowestSellOrder  string `json:"lowest_sell_order"`  // 最低售价（分）
	HighestBuyOrder  string `json:"highest_buy_order"`  // 最高求购价（分）
	SellOrderSummary string `json:"sell_order_summary"` // 在售数量的 HTML
	BuyOrderSummary  string `json:"buy_order_summary"`  // 求购数量的 HTML
}

// SteamOrderData 解析后的订单数据
type SteamOrderData struct {
	SellPrice    float64 // 最低售价（元）
	SellCount    int64   // 在售数量
	BiddingPrice float64 // 最高求购价（元）
	BiddingCount int64   // 求购数量
}

// FetchSteamOrderData 获取 Steam 市场订单数据（求购价、售价、数量）
func FetchSteamOrderData(itemNameId string) (*SteamOrderData, error) {
	if itemNameId == "" {
		return nil, fmt.Errorf("item_nameid is empty")
	}

	// currency=23 是人民币
	url := fmt.Sprintf("market/itemordershistogram?country=CN&language=schinese&currency=23&item_nameid=%s", itemNameId)

	var result SteamOrderHistogram
	opts := utils.RequestOptions{
		Result: &result,
	}

	resp, err := steamCommunityClient.DoRequest("GET", url, opts)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("steam response status: %d", resp.StatusCode())
	}

	if result.Success != 1 {
		return nil, fmt.Errorf("steam api returned success=%d", result.Success)
	}

	data := &SteamOrderData{}

	// 解析最低售价（分 -> 元）
	if result.LowestSellOrder != "" {
		if price, err := strconv.ParseInt(result.LowestSellOrder, 10, 64); err == nil {
			data.SellPrice = float64(price) / 100.0
		}
	}

	// 解析最高求购价（分 -> 元）
	if result.HighestBuyOrder != "" {
		if price, err := strconv.ParseInt(result.HighestBuyOrder, 10, 64); err == nil {
			data.BiddingPrice = float64(price) / 100.0
		}
	}

	// 从 sell_order_summary 提取在售数量
	// 格式: <span class="market_commodity_orders_header_promote">164</span> 个出售中...
	sellCountRe := regexp.MustCompile(`<span[^>]*>(\d+)</span>\s*个出售中`)
	if matches := sellCountRe.FindStringSubmatch(result.SellOrderSummary); len(matches) > 1 {
		if count, err := strconv.ParseInt(matches[1], 10, 64); err == nil {
			data.SellCount = count
		}
	}

	// 从 buy_order_summary 提取求购数量
	// 格式: <span class="market_commodity_orders_header_promote">6872</span> 人请求...
	buyCountRe := regexp.MustCompile(`<span[^>]*>(\d+)</span>\s*人请求`)
	if matches := buyCountRe.FindStringSubmatch(result.BuyOrderSummary); len(matches) > 1 {
		if count, err := strconv.ParseInt(matches[1], 10, 64); err == nil {
			data.BiddingCount = count
		}
	}

	return data, nil
}

// UpdateSteamPricesFromMarket 从 Steam 市场更新价格数据（需要先有 item_nameid）
func UpdateSteamPricesFromMarket() {
	config.Log.Info("Starting to update Steam prices from market...")

	// 获取有 item_nameid 且 sell_price < 500 的商品
	var steams []models.Steam
	config.DB.Where("id != '' AND id IS NOT NULL AND sell_price < 500").Find(&steams)

	if len(steams) == 0 {
		config.Log.Info("No steam items with item_nameid and sell_price < 500 found")
		return
	}

	config.Log.Infof("Found %d steam items with item_nameid and sell_price < 500", len(steams))

	successCount := 0
	failCount := 0

	for i, item := range steams {
		orderData, err := FetchSteamOrderData(item.Id)
		if err != nil {
			config.Log.Warnf("[%d/%d] Failed to get order data for %s: %v", i+1, len(steams), item.MarketHashName, err)
			failCount++
		} else {
			// 更新数据库
			updates := map[string]interface{}{
				"sell_price":    orderData.SellPrice,
				"sell_count":    orderData.SellCount,
				"bidding_price": orderData.BiddingPrice,
				"bidding_count": orderData.BiddingCount,
				"update_time":   time.Now().Unix(),
			}
			if err := config.DB.Model(&models.Steam{}).Where("market_hash_name = ?", item.MarketHashName).Updates(updates).Error; err != nil {
				config.Log.Warnf("[%d/%d] Failed to update %s: %v", i+1, len(steams), item.MarketHashName, err)
				failCount++
			} else {
				config.Log.Infof("[%d/%d] Updated %s: sell=%.2f(%d), bid=%.2f(%d)",
					i+1, len(steams), item.MarketHashName,
					orderData.SellPrice, orderData.SellCount,
					orderData.BiddingPrice, orderData.BiddingCount)
				successCount++
			}
		}
		time.Sleep(800 * time.Millisecond)
	}

	config.Log.Infof("Update steam prices complete. Success: %d, Failed: %d", successCount, failCount)
}

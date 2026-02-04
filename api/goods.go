package api

import (
	"net/http"
	"strconv"
	"uu/models"
	"uu/utils"

	"github.com/gin-gonic/gin"
)

func GetGoods(c *gin.Context) {
	pageSize, _ := strconv.Atoi(c.Query("page_size"))
	pageNum, _ := strconv.Atoi(c.Query("page_num"))
	if pageNum == 0 || pageSize == 0 {
		pageSize = 25
		pageNum = 1
	}
	userId := getUserIdFromContext(c)
	sort := c.Query("sort")
	desc, _ := strconv.ParseBool(c.Query("desc"))
	category := c.Query("category")
	search := c.Query("search")
	source := c.Query("source")
	target := c.Query("target")

	// 新增：购买方案和出售方案
	// buy_type: sell(在售价购买) / bidding(求购价购买)，默认 sell
	// sell_type: sell(在售价出售) / bidding(求购价出售)，默认 sell
	buyType := c.Query("buy_type")
	sellType := c.Query("sell_type")
	if buyType != "sell" && buyType != "bidding" {
		buyType = "sell"
	}
	if sellType != "sell" && sellType != "bidding" {
		sellType = "sell"
	}

	s, total, code := models.GetGoods(userId, pageSize, pageNum, desc, sort, search, source, target, category, buyType, sellType)
	c.JSON(http.StatusOK, gin.H{
		"code":  code,
		"data":  s,
		"total": total,
		"msg":   utils.ErrorMessage(code),
	})
}

func GetGoodsCategory(c *gin.Context) {
	category, code := models.GetCategory()
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"data": category,
		"msg":  utils.ErrorMessage(code),
	})
}

// GetPriceHistory 获取商品价格历史
func GetPriceHistory(c *gin.Context) {
	marketHashName := c.Query("market_hash_name")
	platform := c.Query("platform")
	daysStr := c.Query("days")

	if marketHashName == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": utils.ErrCodeInvalidParams,
			"msg":  "market_hash_name is required",
		})
		return
	}

	days := 30 // 默认30天
	if daysStr != "" {
		if d, err := strconv.Atoi(daysStr); err == nil && d > 0 && d <= 365 {
			days = d
		}
	}

	// 如果指定了平台，只返回该平台的数据
	if platform != "" {
		history, err := models.GetPriceHistoryByPlatform(marketHashName, platform, days)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": utils.ErrCodeGetGoods,
				"msg":  "Failed to get price history",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code": utils.SUCCESS,
			"data": gin.H{
				"marketHashName": marketHashName,
				"platform":       platform,
				"history":        history,
			},
			"msg": utils.ErrorMessage(utils.SUCCESS),
		})
		return
	}

	// 返回所有平台的数据
	history, err := models.GetPriceHistoryByHashName(marketHashName, days)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": utils.ErrCodeGetGoods,
			"msg":  "Failed to get price history",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": utils.SUCCESS,
		"data": history,
		"msg":  utils.ErrorMessage(utils.SUCCESS),
	})
}

// GetGoodsDetail 获取商品详情（包含基础信息、所有平台历史数据、各平台在售信息）
func GetGoodsDetail(c *gin.Context) {
	marketHashName := c.Query("market_hash_name")
	daysStr := c.Query("days")

	if marketHashName == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": utils.ErrCodeInvalidParams,
			"msg":  "market_hash_name is required",
		})
		return
	}

	// 默认30天，最多365天
	days := 30
	if daysStr != "" {
		if d, err := strconv.Atoi(daysStr); err == nil && d > 0 && d <= 365 {
			days = d
		}
	}

	detail, err := models.GetGoodsDetail(marketHashName, days)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": utils.ErrCodeGetGoods,
			"msg":  "Failed to get goods detail",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": utils.SUCCESS,
		"data": detail,
		"msg":  utils.ErrorMessage(utils.SUCCESS),
	})
}

func GetPriceIncreaseByU(c *gin.Context) {
	limit, _ := strconv.Atoi(c.Query("limit"))
	isDesc, _ := strconv.ParseBool(c.Query("is_desc"))
	if limit > 500 {
		c.JSON(http.StatusOK, gin.H{
			"code": utils.ErrCodeGetGoods,
			"msg":  "Over limit",
		})
		return
	}
	increase, err := models.GetPriceIncrease("YOUPIN", "", isDesc, limit)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": utils.ErrCodeGetGoods,
			"msg":  "Failed to get price increase data",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": utils.SUCCESS,
		"data": increase,
		"msg":  utils.ErrorMessage(utils.SUCCESS),
	})
}

// GetPublicHomeData 获取公开首页数据（不需要登录）
func GetPublicHomeData(c *gin.Context) {
	data, err := models.GetPublicHomeData()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": utils.ErrCodeGetGoods,
			"msg":  "Failed to get public home data",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": utils.SUCCESS,
		"data": data,
		"msg":  utils.ErrorMessage(utils.SUCCESS),
	})
}

// GetRelatedWears 获取同款饰品的不同磨损和品质版本
func GetRelatedWears(c *gin.Context) {
	marketHashName := c.Query("market_hash_name")
	if marketHashName == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": utils.ErrCodeInvalidParams,
			"msg":  "market_hash_name is required",
		})
		return
	}

	result, err := models.GetRelatedWears(marketHashName)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": utils.ErrCodeGetGoods,
			"msg":  "Failed to get related wears",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": utils.SUCCESS,
		"data": result,
		"msg":  utils.ErrorMessage(utils.SUCCESS),
	})
}

// SearchGoods 搜索商品（根据名称模糊匹配）
func SearchGoods(c *gin.Context) {
	keyword := c.Query("keyword")
	if keyword == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": utils.SUCCESS,
			"data": []interface{}{},
			"msg":  utils.ErrorMessage(utils.SUCCESS),
		})
		return
	}

	limitStr := c.Query("limit")
	limit := 20 // 默认返回20条
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 50 {
			limit = l
		}
	}

	results, err := models.SearchGoodsByKeyword(keyword, limit)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": utils.ErrCodeGetGoods,
			"msg":  "Failed to search goods",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": utils.SUCCESS,
		"data": results,
		"msg":  utils.ErrorMessage(utils.SUCCESS),
	})
}

// GetBigItemBidding 获取大件求购数据（手套和刀具的求购价差）
func GetBigItemBidding(c *gin.Context) {
	pageSize, _ := strconv.Atoi(c.Query("page_size"))
	pageNum, _ := strconv.Atoi(c.Query("page_num"))
	if pageNum == 0 || pageSize == 0 {
		pageSize = 25
		pageNum = 1
	}

	sort := c.Query("sort")
	desc, _ := strconv.ParseBool(c.Query("desc"))
	search := c.Query("search")
	platform := c.Query("platform")
	category := c.Query("category")

	// 默认平台为悠悠
	if platform == "" {
		platform = "uu"
	}

	data, total, code := models.GetBigItemBidding(pageSize, pageNum, desc, sort, search, platform, category)
	c.JSON(http.StatusOK, gin.H{
		"code":  code,
		"data":  data,
		"total": total,
		"msg":   utils.ErrorMessage(code),
	})
}

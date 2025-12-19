package api

import (
	"fmt"
	"net/http"
	"strings"
	"time"
	"uu/config"
	"uu/models"
	"uu/services"
	"uu/utils"

	"github.com/gin-gonic/gin"
)

// 获取VIP价格信息
func GetVipPrice(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code": utils.SUCCESS,
		"msg":  utils.ErrorMessage(utils.SUCCESS),
		"data": gin.H{
			"plans": models.VipPlans,
		},
	})
}

// 创建支付订单
func CreatePayOrder(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": utils.ErrCodePermissionDenied,
			"msg":  utils.ErrorMessage(utils.ErrCodePermissionDenied),
		})
		return
	}

	// 解析请求参数
	var req struct {
		Months int `json:"months"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		req.Months = 1 // 默认1个月
	}

	// 获取套餐信息
	plan, ok := models.GetVipPlan(req.Months)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": utils.InvalidParameter,
			"msg":  "无效的套餐",
		})
		return
	}

	// 生成商户订单号：VIP + 时间戳 + 随机数
	outTradeNo := fmt.Sprintf("VIP%s%d", time.Now().Format("20060102150405"), time.Now().UnixNano()%10000)

	// 创建订单记录
	order, err := models.CreatePaymentOrder(userID, outTradeNo, plan.Price, plan.Months)
	if err != nil {
		config.Log.Errorf("CreatePayOrder: create order failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": utils.ERROR,
			"msg":  "创建订单失败",
		})
		return
	}

	// 商品描述
	var body string
	if plan.Months == 1 {
		body = "CS Goods VIP会员月卡"
	} else {
		body = fmt.Sprintf("CS Goods VIP会员%d个月", plan.Months)
	}

	// 调用YunGouOS Native支付
	payResp, err := services.CreateNativePay(
		outTradeNo,
		plan.Price,
		body,
		fmt.Sprintf("%s|%d", userID, plan.Months), // 附加用户ID和月数，回调时使用
	)
	if err != nil {
		config.Log.Errorf("CreatePayOrder: call native pay failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": utils.ERROR,
			"msg":  "创建支付失败",
		})
		return
	}

	if payResp.Code != 0 {
		config.Log.Errorf("CreatePayOrder: native pay return error: code=%d, msg=%s", payResp.Code, payResp.Msg)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": utils.ERROR,
			"msg":  fmt.Sprintf("支付创建失败: %s", payResp.Msg),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": utils.SUCCESS,
		"msg":  utils.ErrorMessage(utils.SUCCESS),
		"data": gin.H{
			"order_no":   order.OutTradeNo,
			"qrcode_img": payResp.Data, // base64二维码图片
			"amount":     plan.Price,
			"created_at": order.CreatedAt,
		},
	})
}

// 支付回调通知
func PayNotify(c *gin.Context) {
	paymentConfig := config.CONFIG.Payment
	if paymentConfig == nil {
		config.Log.Error("PayNotify: payment config not found")
		c.String(http.StatusOK, "FAIL")
		return
	}

	// 获取回调参数
	code := c.PostForm("code")
	orderNo := c.PostForm("orderNo")       // YunGouOS订单号
	outTradeNo := c.PostForm("outTradeNo") // 商户订单号
	payNo := c.PostForm("payNo")           // 微信支付单号
	money := c.PostForm("money")           // 支付金额
	mchId := c.PostForm("mchId")           // 商户号
	payChannel := c.PostForm("payChannel") // 支付渠道
	payTime := c.PostForm("time")          // 支付时间
	attach := c.PostForm("attach")         // 附加数据（用户ID）
	sign := c.PostForm("sign")             // 签名

	config.Log.Infof("PayNotify: code=%s, orderNo=%s, outTradeNo=%s, payNo=%s, money=%s, mchId=%s, attach=%s",
		code, orderNo, outTradeNo, payNo, money, mchId, attach)

	// 构建验签参数（根据文档，sign和time不参与签名）
	params := map[string]string{
		"code":       code,
		"orderNo":    orderNo,
		"outTradeNo": outTradeNo,
		"payNo":      payNo,
		"money":      money,
		"mchId":      mchId,
		"payChannel": payChannel,
	}
	if attach != "" {
		params["attach"] = attach
	}

	// 验证签名
	if !services.VerifySign(params, sign, paymentConfig.ApiKey) {
		config.Log.Error("PayNotify: sign verify failed")
		c.String(http.StatusOK, "FAIL")
		return
	}

	// 检查支付状态
	if code != "1" {
		config.Log.Warnf("PayNotify: payment not success, code=%s", code)
		c.String(http.StatusOK, "FAIL")
		return
	}

	// 查询订单
	order, err := models.GetPaymentOrderByOutTradeNo(outTradeNo)
	if err != nil {
		config.Log.Errorf("PayNotify: order not found, outTradeNo=%s, err=%v", outTradeNo, err)
		c.String(http.StatusOK, "FAIL")
		return
	}

	// 防止重复处理
	if order.Status == models.PaymentStatusPaid {
		config.Log.Infof("PayNotify: order already paid, outTradeNo=%s", outTradeNo)
		c.String(http.StatusOK, "SUCCESS")
		return
	}

	// 解析支付时间
	payTimeValue, err := time.Parse("2006-01-02 15:04:05", payTime)
	if err != nil {
		config.Log.Warnf("PayNotify: parse pay time failed, use now: %v", err)
		payTimeValue = time.Now()
	}

	// 更新订单状态
	if err := models.UpdatePaymentOrderPaid(outTradeNo, orderNo, payTimeValue); err != nil {
		config.Log.Errorf("PayNotify: update order failed, err=%v", err)
		c.String(http.StatusOK, "FAIL")
		return
	}

	// 解析附加数据：userID|months
	var userID string
	var months int
	if attach != "" {
		parts := strings.Split(attach, "|")
		userID = parts[0]
		if len(parts) > 1 {
			fmt.Sscanf(parts[1], "%d", &months)
		}
	}
	if userID == "" {
		userID = order.UserID
	}
	if months == 0 {
		// 从订单中获取月数
		months = order.Months
		if months == 0 {
			months = 1 // 默认1个月
		}
	}

	// 更新用户VIP状态（按月份计算）
	if err := models.UpdateUserVipAfterPayment(userID, months); err != nil {
		config.Log.Errorf("PayNotify: update user vip failed, userID=%s, months=%d, err=%v", userID, months, err)
		// 这里不返回FAIL，因为订单已经更新成功，VIP状态可以后续补偿
	}

	config.Log.Infof("PayNotify: payment success, outTradeNo=%s, userID=%s", outTradeNo, userID)
	c.String(http.StatusOK, "SUCCESS")
}

// 查询订单状态
func QueryPayOrder(c *gin.Context) {
	orderNo := c.Query("order_no")
	userID := c.GetString("userID")

	if orderNo == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": utils.InvalidParameter,
			"msg":  "订单号不能为空",
		})
		return
	}

	order, err := models.GetPaymentOrderByOutTradeNo(orderNo)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code": utils.ERROR,
			"msg":  "订单不存在",
		})
		return
	}

	// 检查订单是否属于当前用户
	if order.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{
			"code": utils.ErrCodePermissionDenied,
			"msg":  "无权查看该订单",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": utils.SUCCESS,
		"msg":  utils.ErrorMessage(utils.SUCCESS),
		"data": gin.H{
			"order_no":   order.OutTradeNo,
			"amount":     order.Amount,
			"status":     order.Status,
			"pay_time":   order.PayTime,
			"created_at": order.CreatedAt,
		},
	})
}

// 获取用户订单列表
func GetUserOrders(c *gin.Context) {
	userID := c.GetString("userID")
	pageSize := 10
	pageNum := 1

	if ps := c.Query("page_size"); ps != "" {
		fmt.Sscanf(ps, "%d", &pageSize)
	}
	if pn := c.Query("page_num"); pn != "" {
		fmt.Sscanf(pn, "%d", &pageNum)
	}

	orders, total, err := models.GetUserPaymentOrders(userID, pageSize, pageNum)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": utils.ERROR,
			"msg":  "查询订单失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":      utils.SUCCESS,
		"msg":       utils.ErrorMessage(utils.SUCCESS),
		"data":      orders,
		"total":     total,
		"page_size": pageSize,
		"page_num":  pageNum,
	})
}

// 获取所有订单（管理员）
func GetAllOrders(c *gin.Context) {
	pageSize := 20
	pageNum := 1
	status := -1 // -1 表示全部

	if ps := c.Query("page_size"); ps != "" {
		fmt.Sscanf(ps, "%d", &pageSize)
	}
	if pn := c.Query("page_num"); pn != "" {
		fmt.Sscanf(pn, "%d", &pageNum)
	}
	if s := c.Query("status"); s != "" {
		fmt.Sscanf(s, "%d", &status)
	}
	keyword := c.Query("keyword")

	orders, total, err := models.GetAllPaymentOrders(pageSize, pageNum, status, keyword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": utils.ERROR,
			"msg":  "查询订单失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":      utils.SUCCESS,
		"msg":       utils.ErrorMessage(utils.SUCCESS),
		"data":      orders,
		"total":     total,
		"page_size": pageSize,
		"page_num":  pageNum,
	})
}

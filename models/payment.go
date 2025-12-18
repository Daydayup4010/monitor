package models

import (
	"time"
	"uu/config"

	"github.com/google/uuid"
)

// 支付订单状态
const (
	PaymentStatusPending  = 0 // 待支付
	PaymentStatusPaid     = 1 // 已支付
	PaymentStatusCanceled = 2 // 已取消
	PaymentStatusRefunded = 3 // 已退款
)

// VIP套餐价格配置
type VipPlan struct {
	Months int     `json:"months"`
	Price  float64 `json:"price"`
	Days   int     `json:"days"`
}

// VIP套餐列表
var VipPlans = map[int]VipPlan{
	1:  {Months: 1, Price: 19.9, Days: 30},
	3:  {Months: 3, Price: 49.9, Days: 90},
	6:  {Months: 6, Price: 89.9, Days: 180},
	12: {Months: 12, Price: 169.9, Days: 365},
}

// 获取VIP套餐
func GetVipPlan(months int) (VipPlan, bool) {
	plan, ok := VipPlans[months]
	return plan, ok
}

// 支付订单
type PaymentOrder struct {
	ID         string     `json:"id" gorm:"type:char(36);primaryKey"`
	UserID     string     `json:"user_id" gorm:"type:char(36);index"`
	OutTradeNo string     `json:"out_trade_no" gorm:"type:varchar(64);uniqueIndex"` // 商户订单号
	Amount     float64    `json:"amount"`                                           // 支付金额
	Months     int        `json:"months" gorm:"default:1"`                          // 购买月数
	Status     int        `json:"status" gorm:"default:0"`                          // 订单状态
	PayTime    *time.Time `json:"pay_time"`                                         // 支付时间
	YunOrderNo string     `json:"yun_order_no" gorm:"type:varchar(64)"`             // YunGouOS订单号
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

// 创建支付订单
func CreatePaymentOrder(userID, outTradeNo string, amount float64, months int) (*PaymentOrder, error) {
	order := &PaymentOrder{
		ID:         uuid.New().String(),
		UserID:     userID,
		OutTradeNo: outTradeNo,
		Amount:     amount,
		Months:     months,
		Status:     PaymentStatusPending,
	}
	if err := config.DB.Create(order).Error; err != nil {
		return nil, err
	}
	return order, nil
}

// 根据商户订单号查询订单
func GetPaymentOrderByOutTradeNo(outTradeNo string) (*PaymentOrder, error) {
	var order PaymentOrder
	if err := config.DB.Where("out_trade_no = ?", outTradeNo).First(&order).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

// 更新订单为已支付
func UpdatePaymentOrderPaid(outTradeNo, yunOrderNo string) error {
	now := time.Now()
	return config.DB.Model(&PaymentOrder{}).
		Where("out_trade_no = ? AND status = ?", outTradeNo, PaymentStatusPending).
		Updates(map[string]interface{}{
			"status":       PaymentStatusPaid,
			"pay_time":     now,
			"yun_order_no": yunOrderNo,
		}).Error
}

// 更新用户VIP状态
func UpdateUserVipAfterPayment(userID string, days int) error {
	var user User
	if err := config.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		return err
	}

	// 如果当前VIP未过期，在原基础上累加；否则从现在开始计算
	baseTime := time.Now()
	if user.VipExpiry.After(baseTime) {
		baseTime = user.VipExpiry
	}
	newExpiry := baseTime.AddDate(0, 0, days)

	return config.DB.Model(&user).Updates(map[string]interface{}{
		"role":       RoleVip,
		"vip_expiry": newExpiry,
	}).Error
}

// 获取用户支付订单列表
func GetUserPaymentOrders(userID string, pageSize, pageNum int) ([]PaymentOrder, int64, error) {
	var orders []PaymentOrder
	var total int64

	db := config.DB.Model(&PaymentOrder{}).Where("user_id = ?", userID)
	db.Count(&total)

	offset := (pageNum - 1) * pageSize
	if err := db.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&orders).Error; err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

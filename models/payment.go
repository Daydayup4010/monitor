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
}

// VIP套餐列表
var VipPlans = map[int]VipPlan{
	1:  {Months: 1, Price: 0.1},
	3:  {Months: 3, Price: 0.2},
	6:  {Months: 6, Price: 0.3},
	12: {Months: 12, Price: 0.4},
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
func UpdatePaymentOrderPaid(outTradeNo, yunOrderNo string, payTime time.Time) error {
	return config.DB.Model(&PaymentOrder{}).
		Where("out_trade_no = ? AND status = ?", outTradeNo, PaymentStatusPending).
		Updates(map[string]interface{}{
			"status":       PaymentStatusPaid,
			"pay_time":     payTime,
			"yun_order_no": yunOrderNo,
		}).Error
}

// 更新用户VIP状态（按月份计算）
func UpdateUserVipAfterPayment(userID string, months int) error {
	var user User
	if err := config.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		return err
	}

	// 如果当前VIP未过期，在原基础上累加；否则从现在开始计算
	baseTime := time.Now()
	if user.VipExpiry.After(baseTime) {
		baseTime = user.VipExpiry
	}
	newExpiry := baseTime.AddDate(0, months, 0) // 按月份计算

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

// 获取用户已支付的VIP开通记录
func GetUserPaidOrders(userID string, pageSize, pageNum int) ([]PaymentOrder, int64, error) {
	var orders []PaymentOrder
	var total int64

	db := config.DB.Model(&PaymentOrder{}).Where("user_id = ? AND status = ?", userID, PaymentStatusPaid)
	db.Count(&total)

	offset := (pageNum - 1) * pageSize
	if err := db.Order("pay_time DESC").Offset(offset).Limit(pageSize).Find(&orders).Error; err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

// 订单详情（包含用户信息）
type PaymentOrderDetail struct {
	PaymentOrder
	Username string `json:"username"`
	Email    string `json:"email"`
}

// 获取所有支付订单列表（管理员用）
func GetAllPaymentOrders(pageSize, pageNum int, status int, keyword string, startTime, endTime string) ([]PaymentOrderDetail, int64, error) {
	var orders []PaymentOrderDetail
	var total int64

	query := config.DB.Model(&PaymentOrder{}).
		Select("payment_order.*, user.email").
		Joins("LEFT JOIN user ON payment_order.user_id = user.id")

	// 状态筛选：-1表示全部
	if status >= 0 {
		query = query.Where("payment_order.status = ?", status)
	}

	// 关键词搜索（订单号或邮箱）
	if keyword != "" {
		query = query.Where(
			"payment_order.out_trade_no LIKE ? OR user.email LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%",
		)
	}

	// 时间范围筛选
	if startTime != "" {
		query = query.Where("payment_order.created_at >= ?", startTime)
	}
	if endTime != "" {
		query = query.Where("payment_order.created_at <= ?", endTime+" 23:59:59")
	}

	query.Count(&total)

	offset := (pageNum - 1) * pageSize
	if err := query.Order("payment_order.created_at DESC").Offset(offset).Limit(pageSize).Find(&orders).Error; err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

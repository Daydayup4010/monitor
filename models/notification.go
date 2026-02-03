package models

import (
	"time"
	"uu/config"
	"uu/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Notification 通知表（管理员发布的全局通知）
type Notification struct {
	ID        uuid.UUID      `json:"id" gorm:"type:char(36);primaryKey"`
	Title     string         `json:"title" gorm:"type:varchar(255);not null"`
	Content   string         `json:"content" gorm:"type:text;not null"`
	ImageURL  string         `json:"image_url" gorm:"type:varchar(500)"` // 图片URL（可选）
	CreatedAt time.Time      `json:"created_at" gorm:"type:datetime;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"type:datetime;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"type:datetime;index"`
}

// NotificationRead 通知已读记录表
type NotificationRead struct {
	ID             uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID         uuid.UUID `json:"user_id" gorm:"type:char(36);index:idx_user_notification,unique"`
	NotificationID uuid.UUID `json:"notification_id" gorm:"type:char(36);index:idx_user_notification,unique"`
	ReadAt         time.Time `json:"read_at" gorm:"type:datetime;default:CURRENT_TIMESTAMP"`
}

// NotificationResponse 通知响应结构
type NotificationResponse struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	ImageURL  string    `json:"image_url,omitempty"`
	IsRead    bool      `json:"is_read"`
	CreatedAt time.Time `json:"created_at"`
}

// CreateNotification 创建通知（管理员）
func CreateNotification(title, content, imageURL string) (*Notification, int) {
	notification := &Notification{
		ID:       uuid.New(),
		Title:    title,
		Content:  content,
		ImageURL: imageURL,
	}
	err := config.DB.Create(notification).Error
	if err != nil {
		config.Log.Errorf("Create notification error: %v", err)
		return nil, utils.ErrCodeCreateNotification
	}
	return notification, utils.SUCCESS
}

// GetNotifications 获取通知列表（带用户已读状态）
func GetNotifications(userID uuid.UUID, pageSize, pageNum int) ([]NotificationResponse, int64, int) {
	var notifications []Notification
	var total int64

	// 获取通知总数
	config.DB.Model(&Notification{}).Count(&total)

	// 获取通知列表
	err := config.DB.Order("created_at DESC").
		Limit(pageSize).
		Offset((pageNum - 1) * pageSize).
		Find(&notifications).Error

	if err != nil {
		config.Log.Errorf("Get notifications error: %v", err)
		return nil, 0, utils.ErrCodeGetNotifications
	}

	// 获取用户已读的通知ID列表
	var readNotificationIDs []uuid.UUID
	config.DB.Model(&NotificationRead{}).
		Where("user_id = ?", userID).
		Pluck("notification_id", &readNotificationIDs)

	// 转换为map方便查询
	readMap := make(map[string]bool)
	for _, id := range readNotificationIDs {
		readMap[id.String()] = true
	}

	// 构建响应
	var result []NotificationResponse
	for _, n := range notifications {
		result = append(result, NotificationResponse{
			ID:        n.ID.String(),
			Title:     n.Title,
			Content:   n.Content,
			ImageURL:  n.ImageURL,
			IsRead:    readMap[n.ID.String()],
			CreatedAt: n.CreatedAt,
		})
	}

	return result, total, utils.SUCCESS
}

// GetUnreadCount 获取用户未读通知数量
func GetUnreadCount(userID uuid.UUID) (int64, int) {
	var totalCount int64
	var readCount int64

	// 获取通知总数
	config.DB.Model(&Notification{}).Count(&totalCount)

	// 获取用户已读数量
	config.DB.Model(&NotificationRead{}).
		Where("user_id = ?", userID).
		Count(&readCount)

	unreadCount := totalCount - readCount
	if unreadCount < 0 {
		unreadCount = 0
	}

	return unreadCount, utils.SUCCESS
}

// MarkAsRead 标记通知为已读
func MarkAsRead(userID uuid.UUID, notificationID uuid.UUID) int {
	// 检查通知是否存在
	var notification Notification
	if err := config.DB.First(&notification, "id = ?", notificationID).Error; err != nil {
		return utils.ErrCodeNotificationNotFound
	}

	// 检查是否已读
	var existingRead NotificationRead
	err := config.DB.Where("user_id = ? AND notification_id = ?", userID, notificationID).First(&existingRead).Error
	if err == nil {
		// 已经已读
		return utils.SUCCESS
	}

	// 创建已读记录
	read := &NotificationRead{
		UserID:         userID,
		NotificationID: notificationID,
	}
	if err := config.DB.Create(read).Error; err != nil {
		config.Log.Errorf("Mark notification as read error: %v", err)
		return utils.ErrCodeMarkNotificationRead
	}

	return utils.SUCCESS
}

// MarkAllAsRead 标记所有通知为已读
func MarkAllAsRead(userID uuid.UUID) int {
	// 获取所有通知ID
	var notificationIDs []uuid.UUID
	config.DB.Model(&Notification{}).Pluck("id", &notificationIDs)

	// 获取用户已读的通知ID
	var readNotificationIDs []uuid.UUID
	config.DB.Model(&NotificationRead{}).
		Where("user_id = ?", userID).
		Pluck("notification_id", &readNotificationIDs)

	// 转换为map
	readMap := make(map[string]bool)
	for _, id := range readNotificationIDs {
		readMap[id.String()] = true
	}

	// 批量创建未读的已读记录
	var newReads []NotificationRead
	for _, id := range notificationIDs {
		if !readMap[id.String()] {
			newReads = append(newReads, NotificationRead{
				UserID:         userID,
				NotificationID: id,
			})
		}
	}

	if len(newReads) > 0 {
		if err := config.DB.Create(&newReads).Error; err != nil {
			config.Log.Errorf("Mark all notifications as read error: %v", err)
			return utils.ErrCodeMarkNotificationRead
		}
	}

	return utils.SUCCESS
}

// DeleteNotification 删除通知（管理员）
func DeleteNotification(notificationID uuid.UUID) int {
	result := config.DB.Delete(&Notification{}, "id = ?", notificationID)
	if result.Error != nil {
		config.Log.Errorf("Delete notification error: %v", result.Error)
		return utils.ErrCodeDeleteNotification
	}
	if result.RowsAffected == 0 {
		return utils.ErrCodeNotificationNotFound
	}

	// 同时删除已读记录
	config.DB.Delete(&NotificationRead{}, "notification_id = ?", notificationID)

	return utils.SUCCESS
}

// GetAllNotifications 获取所有通知（管理员）
func GetAllNotifications(pageSize, pageNum int) ([]NotificationResponse, int64, int) {
	var notifications []Notification
	var total int64

	config.DB.Model(&Notification{}).Count(&total)

	err := config.DB.Order("created_at DESC").
		Limit(pageSize).
		Offset((pageNum - 1) * pageSize).
		Find(&notifications).Error

	if err != nil {
		config.Log.Errorf("Get all notifications error: %v", err)
		return nil, 0, utils.ErrCodeGetNotifications
	}

	var result []NotificationResponse
	for _, n := range notifications {
		result = append(result, NotificationResponse{
			ID:        n.ID.String(),
			Title:     n.Title,
			Content:   n.Content,
			ImageURL:  n.ImageURL,
			IsRead:    false, // 管理员视角不需要已读状态
			CreatedAt: n.CreatedAt,
		})
	}

	return result, total, utils.SUCCESS
}

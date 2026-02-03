package api

import (
	"net/http"
	"strconv"
	"uu/models"
	"uu/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CreateNotificationRequest 创建通知请求
type CreateNotificationRequest struct {
	Title    string `json:"title" binding:"required,min=1,max=255"`
	Content  string `json:"content" binding:"required,min=1"`
	ImageURL string `json:"image_url"` // 图片URL（可选）
}

// CreateNotification 创建通知（管理员）
func CreateNotification(c *gin.Context) {
	var req CreateNotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": utils.ErrCodeInvalidParams,
			"msg":  utils.ErrorMessage(utils.ErrCodeInvalidParams),
		})
		return
	}

	notification, code := models.CreateNotification(req.Title, req.Content, req.ImageURL)
	if code != utils.SUCCESS {
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"msg":  utils.ErrorMessage(code),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": utils.SUCCESS,
		"msg":  "创建成功",
		"data": notification,
	})
}

// GetNotifications 获取通知列表（用户）
func GetNotifications(c *gin.Context) {
	userIDVal, _ := c.Get("userID")
	userID, ok := userIDVal.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"code": utils.ErrCodeInvalidToken,
			"msg":  utils.ErrorMessage(utils.ErrCodeInvalidToken),
		})
		return
	}

	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	pageNum, _ := strconv.Atoi(c.DefaultQuery("page_num", "1"))

	if pageSize <= 0 {
		pageSize = 20
	}
	if pageNum <= 0 {
		pageNum = 1
	}

	notifications, total, code := models.GetNotifications(userID, pageSize, pageNum)
	if code != utils.SUCCESS {
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"msg":  utils.ErrorMessage(code),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":      utils.SUCCESS,
		"data":      notifications,
		"total":     total,
		"page_size": pageSize,
		"page_num":  pageNum,
	})
}

// GetUnreadCount 获取未读通知数量
func GetUnreadCount(c *gin.Context) {
	userIDVal, _ := c.Get("userID")
	userID, ok := userIDVal.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"code": utils.ErrCodeInvalidToken,
			"msg":  utils.ErrorMessage(utils.ErrCodeInvalidToken),
		})
		return
	}

	count, code := models.GetUnreadCount(userID)
	if code != utils.SUCCESS {
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"msg":  utils.ErrorMessage(code),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": utils.SUCCESS,
		"data": gin.H{
			"unread_count": count,
		},
	})
}

// MarkAsReadRequest 标记已读请求
type MarkAsReadRequest struct {
	NotificationID string `json:"notification_id" binding:"required"`
}

// MarkAsRead 标记通知为已读
func MarkAsRead(c *gin.Context) {
	userIDVal, _ := c.Get("userID")
	userID, ok := userIDVal.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"code": utils.ErrCodeInvalidToken,
			"msg":  utils.ErrorMessage(utils.ErrCodeInvalidToken),
		})
		return
	}

	var req MarkAsReadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": utils.ErrCodeInvalidParams,
			"msg":  utils.ErrorMessage(utils.ErrCodeInvalidParams),
		})
		return
	}

	notificationID, err := uuid.Parse(req.NotificationID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": utils.ErrCodeInvalidParams,
			"msg":  "无效的通知ID",
		})
		return
	}

	code := models.MarkAsRead(userID, notificationID)
	if code != utils.SUCCESS {
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"msg":  utils.ErrorMessage(code),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": utils.SUCCESS,
		"msg":  "标记成功",
	})
}

// MarkAllAsRead 标记所有通知为已读
func MarkAllAsRead(c *gin.Context) {
	userIDVal, _ := c.Get("userID")
	userID, ok := userIDVal.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"code": utils.ErrCodeInvalidToken,
			"msg":  utils.ErrorMessage(utils.ErrCodeInvalidToken),
		})
		return
	}

	code := models.MarkAllAsRead(userID)
	if code != utils.SUCCESS {
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"msg":  utils.ErrorMessage(code),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": utils.SUCCESS,
		"msg":  "已全部标记为已读",
	})
}

// DeleteNotification 删除通知（管理员）
func DeleteNotification(c *gin.Context) {
	notificationIDStr := c.Query("notification_id")
	if notificationIDStr == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": utils.ErrCodeInvalidParams,
			"msg":  "缺少通知ID",
		})
		return
	}

	notificationID, err := uuid.Parse(notificationIDStr)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": utils.ErrCodeInvalidParams,
			"msg":  "无效的通知ID",
		})
		return
	}

	code := models.DeleteNotification(notificationID)
	if code != utils.SUCCESS {
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"msg":  utils.ErrorMessage(code),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": utils.SUCCESS,
		"msg":  "删除成功",
	})
}

// GetAllNotifications 获取所有通知（管理员）
func GetAllNotifications(c *gin.Context) {
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	pageNum, _ := strconv.Atoi(c.DefaultQuery("page_num", "1"))

	if pageSize <= 0 {
		pageSize = 20
	}
	if pageNum <= 0 {
		pageNum = 1
	}

	notifications, total, code := models.GetAllNotifications(pageSize, pageNum)
	if code != utils.SUCCESS {
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"msg":  utils.ErrorMessage(code),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":      utils.SUCCESS,
		"data":      notifications,
		"total":     total,
		"page_size": pageSize,
		"page_num":  pageNum,
	})
}

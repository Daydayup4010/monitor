package api

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"uu/config"
	"uu/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// 允许的图片类型
var allowedImageTypes = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".gif":  true,
	".webp": true,
}

// 最大文件大小 5MB
const maxFileSize = 5 * 1024 * 1024

// UploadImage 上传图片（管理员）
func UploadImage(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": utils.ErrCodeInvalidParams,
			"msg":  "请选择要上传的图片",
		})
		return
	}

	// 检查文件大小
	if file.Size > maxFileSize {
		c.JSON(http.StatusOK, gin.H{
			"code": utils.ErrCodeInvalidParams,
			"msg":  "图片大小不能超过5MB",
		})
		return
	}

	// 检查文件类型
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !allowedImageTypes[ext] {
		c.JSON(http.StatusOK, gin.H{
			"code": utils.ErrCodeInvalidParams,
			"msg":  "只支持 jpg、png、gif、webp 格式的图片",
		})
		return
	}

	// 生成唯一文件名
	filename := uuid.New().String() + ext

	// 确保上传目录存在
	uploadDir := "./uploads/notifications"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		config.Log.Errorf("Create upload dir error: %v", err)
		c.JSON(http.StatusOK, gin.H{
			"code": utils.ERROR,
			"msg":  "创建上传目录失败",
		})
		return
	}

	// 保存文件
	dst := filepath.Join(uploadDir, filename)
	if err := c.SaveUploadedFile(file, dst); err != nil {
		config.Log.Errorf("Save uploaded file error: %v", err)
		c.JSON(http.StatusOK, gin.H{
			"code": utils.ERROR,
			"msg":  "保存图片失败",
		})
		return
	}

	// 返回图片URL
	imageURL := "/uploads/notifications/" + filename

	c.JSON(http.StatusOK, gin.H{
		"code": utils.SUCCESS,
		"msg":  "上传成功",
		"data": gin.H{
			"url": imageURL,
		},
	})
}

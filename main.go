package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func requireAuth(requiredPerm Permission) gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := getAPIKeyFromHeader(c)
		if apiKey == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "API key is required"})
			c.Abort()
			return
		}

		if !validateAPIKey(apiKey, requiredPerm) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func main() {
	// 加载配置
	if err := loadConfig(); err != nil {
		fmt.Printf("Warning: Failed to load config: %v\n", err)
		return
	}

	// 创建必要的目录
	os.MkdirAll("templates/live", 0755)
	os.MkdirAll("templates/file", 0755)

	// 初始化 Git 仓库
	gitRepo := NewGitRepo(".")
	if err := gitRepo.Init(); err != nil {
		fmt.Printf("Warning: Failed to initialize git repo: %v\n", err)
	}

	metadataStore := NewTemplateMetadata("templates/metadata.json")
	metadataStore.Load()

	r := gin.Default()

	// CORS 中间件
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-API-Key")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// 列出模板文件 - 需要读权限
	r.GET("/api/templates/list", requireAuth(PermissionRead), func(c *gin.Context) {
		templateType := c.Query("type")

		metadataStore.RLock()
		defer metadataStore.RUnlock()

		result := make([]TemplateInfo, 0, len(metadataStore.Templates))
		for _, tpl := range metadataStore.Templates {
			if templateType == "" || tpl.Type == templateType {
				result = append(result, tpl)
			}
		}

		// 按创建时间倒序排序
		sort.SliceStable(result, func(i, j int) bool {
			if result[i].CreateTime != result[j].CreateTime {
				return result[i].CreateTime > result[j].CreateTime
			}
			return result[i].DisplayName < result[j].DisplayName
		})

		c.JSON(http.StatusOK, result)
	})

	// 下载模板 - 需要读权限
	r.GET("/api/templates/:type/:name", requireAuth(PermissionRead), func(c *gin.Context) {
		templateType := c.Param("type")
		fileName := c.Param("name")

		if templateType != "live" && templateType != "file" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid template type"})
			return
		}

		filePath := filepath.Join("templates", templateType, fileName)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Template not found"})
			return
		}

		c.File(filePath)
	})

	// 上传模板 - 需要写权限
	r.POST("/api/templates/upload/:type", requireAuth(PermissionWrite), func(c *gin.Context) {
		templateType := c.Param("type")
		if templateType != "live" && templateType != "file" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid template type"})
			return
		}

		displayName := c.PostForm("displayName")
		if displayName == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Display name is required"})
			return
		}

		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
			return
		}

		ext := filepath.Ext(file.Filename)
		if ext != ".zip" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Only .zip files are allowed"})
			return
		}

		// 生成唯一文件名
		filename := fmt.Sprintf("%s_%s%s", displayName, uuid.New().String()[:8], ext)
		targetPath := filepath.Join("templates", templateType, filename)

		// 检查是否已存在同名模板
		metadataStore.Lock()
		for _, tpl := range metadataStore.Templates {
			if tpl.DisplayName == displayName && tpl.Type == templateType {
				metadataStore.Unlock()
				c.JSON(http.StatusBadRequest, gin.H{"error": "Template with this name already exists"})
				return
			}
		}

		if err := c.SaveUploadedFile(file, targetPath); err != nil {
			metadataStore.Unlock()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
			return
		}

		// 保存元数据，添加创建时间
		metadataStore.Templates[filename] = TemplateInfo{
			DisplayName: displayName,
			FileName:    filename,
			Type:        templateType,
			CreateTime:  time.Now().Unix(),
		}
		if err := metadataStore.Save(); err != nil {
			metadataStore.Unlock()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save metadata"})
			return
		}

		// Git 提交
		if err := gitRepo.AddTemplate(templateType, filename); err != nil {
			fmt.Printf("Warning: Failed to commit template addition: %v\n", err)
		}

		metadataStore.Unlock()

		c.JSON(http.StatusOK, gin.H{
			"message":     "Template uploaded successfully",
			"displayName": displayName,
			"fileName":    filename,
			"type":        templateType,
		})
	})

	// 删除模板 - 需要管理员权限
	r.DELETE("/api/templates/:type/:name", requireAuth(PermissionAdmin), func(c *gin.Context) {
		templateType := c.Param("type")
		fileName := c.Param("name")

		if templateType != "live" && templateType != "file" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid template type"})
			return
		}

		// 检查文件是否存在
		filePath := filepath.Join("templates", templateType, fileName)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Template not found"})
			return
		}

		metadataStore.Lock()
		defer metadataStore.Unlock()

		// 检查模板是否存在于元数据中
		if _, exists := metadataStore.Templates[fileName]; !exists {
			c.JSON(http.StatusNotFound, gin.H{"error": "Template metadata not found"})
			return
		}

		// 删除文件
		if err := os.Remove(filePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete template file"})
			return
		}

		// 删除元数据
		delete(metadataStore.Templates, fileName)
		if err := metadataStore.Save(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update metadata"})
			return
		}

		// Git 提交
		if err := gitRepo.DeleteTemplate(templateType, fileName); err != nil {
			fmt.Printf("Warning: Failed to commit template deletion: %v\n", err)
		}

		c.JSON(http.StatusOK, gin.H{
			"message":  "Template deleted successfully",
			"fileName": fileName,
			"type":     templateType,
		})
	})

	r.Run(fmt.Sprintf(":%d", config.Port))
}

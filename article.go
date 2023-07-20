package main

import (
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// Article 数据模型
type Article struct {
	ID          uint      `gorm:"primary_key" json:"id"`
	Title       string    `json:"title" binding:"required"`
	Author      string    `json:"author" binding:"required"`
	Content     string    `json:"content" binding:"required"`
	Cover       string    `form:"cover"`
	Summary     string    `json:"summary"`
	PublishDate time.Time `json:"publish_date"`
	Views       int       `json:"views" gorm:"default:0"`
	Likes       int       `json:"likes" gorm:"default:0"`
}

// 静云谷文库主页
func getArticlesHandler(c *gin.Context) {
	const articlesPerPage = 6
	pageStr := c.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的页码"})
		return
	}

	var articles []Article
	searchQuery := c.Query("searchQuery")
	var totalArticlesCount int64

	if searchQuery != "" {
		// 使用 Like 条件进行模糊搜索
		if err := db.Where("title LIKE ?", "%"+searchQuery+"%").Find(&articles).Count(&totalArticlesCount).Error; err != nil {
			c.HTML(http.StatusInternalServerError, "error.html", gin.H{
				"message": "数据库查询失败",
			})
			return
		}
	} else {
		// 没有搜索关键词，返回所有文章
		if err := db.Find(&articles).Count(&totalArticlesCount).Error; err != nil {
			c.HTML(http.StatusInternalServerError, "error.html", gin.H{
				"message": "数据库查询失败",
			})
			return
		}
	}

	totalPages := int(math.Ceil(float64(totalArticlesCount) / float64(articlesPerPage)))
	if page > totalPages {
		c.JSON(http.StatusNotFound, gin.H{"error": "页码超出范围"})
		return
	}

	offset := (page - 1) * articlesPerPage
	if searchQuery != "" {
		// 使用 Like 条件进行模糊搜索
		if err := db.Where("title LIKE ?", "%"+searchQuery+"%").Offset(offset).Limit(articlesPerPage).Find(&articles).Error; err != nil {
			c.HTML(http.StatusInternalServerError, "error.html", gin.H{
				"message": "数据库查询失败",
			})
			return
		}
	} else {
		// 没有搜索关键词，返回所有文章
		if err := db.Offset(offset).Limit(articlesPerPage).Find(&articles).Error; err != nil {
			c.HTML(http.StatusInternalServerError, "error.html", gin.H{
				"message": "数据库查询失败",
			})
			return
		}
	}

	c.HTML(http.StatusOK, "articles.html", gin.H{
		"articles":    articles,
		"searchQuery": searchQuery,
		"currentPage": page,
		"totalPages":  totalPages,
	})
}

// 文章详情页面
func getArticlePageHandler(c *gin.Context) {
	id := c.Param("id")
	var article Article
	if err := db.First(&article, id).Error; err != nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{
			"message": "文章不存在",
		})
		return
	}

	// 更新文章的访问量
	article.Views++
	db.Save(&article)
	if err := db.Save(&article).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": "访问量增加失败"})
		return
	}

	c.HTML(http.StatusOK, "article_views.html", gin.H{
		"article": article,
	})
}

// 点赞文章
func likeArticleHandler(c *gin.Context) {
	id := c.Param("id")
	var article Article
	if err := db.First(&article, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
		return
	}

	// 更新文章的点赞数
	article.Likes++
	if err := db.Save(&article).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "点赞失败"})
		return
	}

	// 返回点赞成功后的新点赞数
	c.JSON(http.StatusOK, gin.H{"likes": article.Likes})
}

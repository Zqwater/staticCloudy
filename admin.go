package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// containsString 检查字符串是否包含子字符串
func containsString(s, substr string) bool {
	return strings.Contains(s, substr)
}

// 管理员登录页面
func adminLoginHandler(c *gin.Context) {
	session := sessions.Default(c)
	admin := session.Get("admin")
	if admin != nil {
		c.Redirect(http.StatusFound, "/admin/dashboard")
		return
	}
	c.HTML(http.StatusOK, "admin_login.html", nil)
}

// 管理员登录表单提交
func adminLoginSubmitHandler(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	// 输出接收到的用户名和密码
	fmt.Println("Received username:", username)
	fmt.Println("Received password:", password)
	if username == "admin" && password == "admin123" {
		session := sessions.Default(c)
		session.Set("admin", username)
		session.Save()
		c.JSON(http.StatusOK, gin.H{"success": true})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": false, "message": "用户名或密码不正确"})
}

// 管理员登录后的仪表盘页面
func adminDashboardHandler(c *gin.Context) {
	session := sessions.Default(c)
	admin := session.Get("admin")
	if admin == nil {
		c.Redirect(http.StatusFound, "/admin")
		return
	}

	if adminUsername, ok := admin.(string); ok {
		c.HTML(http.StatusOK, "admin_dashboard.html", gin.H{"username": adminUsername})
	} else {
		c.Redirect(http.StatusFound, "/admin")
	}
}

// 管理员文章管理页面
func adminArticlesHandler(c *gin.Context) {
	session := sessions.Default(c)
	admin := session.Get("admin")
	if adminUsername, ok := admin.(string); ok {
		var articles []Article
		if err := db.Find(&articles).Error; err != nil {
			c.HTML(http.StatusInternalServerError, "admin_articles.html", gin.H{"error": "无法获取文章列表"})
			return
		}

		searchKeyword := c.Query("search")
		if searchKeyword != "" {
			var filteredArticles []Article
			for _, article := range articles {
				if containsString(strings.ToLower(article.Title), strings.ToLower(searchKeyword)) ||
					containsString(strings.ToLower(article.Author), strings.ToLower(searchKeyword)) {
					filteredArticles = append(filteredArticles, article)
				}
			}
			articles = filteredArticles
		}

		c.HTML(http.StatusOK, "admin_articles.html", gin.H{
			"username": adminUsername,
			"articles": articles,
			"contains": containsString, // 注册自定义函数到模板
		})
	} else {
		c.Redirect(http.StatusFound, "/admin")
	}
}

// 管理员创建新文章页面
func adminNewArticleHandler(c *gin.Context) {
	session := sessions.Default(c)
	admin := session.Get("admin")
	if adminUsername, ok := admin.(string); ok {
		c.HTML(http.StatusOK, "admin_new_article.html", gin.H{"username": adminUsername})
	} else {
		c.Redirect(http.StatusFound, "/admin")
	}
}

// 管理员创建新文章
func adminCreateArticleHandler(c *gin.Context) {
	session := sessions.Default(c)
	admin := session.Get("admin")
	if adminUsername, ok := admin.(string); ok {
		var article Article

		// 解析表单数据
		article.Title = c.PostForm("Title")
		article.Author = c.PostForm("Author")
		article.Content = c.PostForm("Content")
		article.Summary = c.PostForm("Summary")

		// 处理封面图片上传
		coverFile, err := c.FormFile("Cover")
		if err != nil {
			c.HTML(http.StatusBadRequest, "admin_new_article.html", gin.H{"username": adminUsername, "error": "无法获取封面图片"})
			return
		}

		// 保存封面图片到指定路径
		coverPath := fmt.Sprintf("static/image/upload/%d_%s", time.Now().Unix(), coverFile.Filename)
		if err := c.SaveUploadedFile(coverFile, coverPath); err != nil {
			c.HTML(http.StatusInternalServerError, "admin_new_article.html", gin.H{"username": adminUsername, "error": "无法保存封面图片"})
			return
		}

		article.Cover = coverPath
		article.PublishDate, _ = time.Parse("2006-01-02T15:04", c.PostForm("PublishDate"))
		article.Views = 0
		article.Likes = 0

		// 进行数据库操作
		if err := db.Create(&article).Error; err != nil {
			c.HTML(http.StatusInternalServerError, "admin_new_article.html", gin.H{"username": adminUsername, "error": "无法创建文章"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "文章创建成功"})
	} else {
		c.Redirect(http.StatusFound, "/admin")
	}
}

// 管理员编辑文章页面
func adminEditArticleHandler(c *gin.Context) {
	session := sessions.Default(c)
	admin := session.Get("admin")
	if adminUsername, ok := admin.(string); ok {
		id := c.Param("id")
		var article Article
		if err := db.First(&article, id).Error; err != nil {
			c.HTML(http.StatusNotFound, "admin_edit_article.html", gin.H{"username": adminUsername, "error": "文章不存在"})
			return
		}
		c.HTML(http.StatusOK, "admin_edit_article.html", gin.H{"username": adminUsername, "article": article})
	} else {
		c.Redirect(http.StatusFound, "/admin")
	}
}

func adminUpdateArticleHandler(c *gin.Context) {
	session := sessions.Default(c)
	admin := session.Get("admin")
	if adminUsername, ok := admin.(string); ok {
		id := c.Param("id")
		var article Article
		var coverPath string
		if err := db.First(&article, id).Error; err != nil {
			c.HTML(http.StatusNotFound, "admin_edit_article.html", gin.H{"username": adminUsername, "error": "文章不存在"})
			return
		}
		// 获取要更新的字段值
		article.Title = c.PostForm("Title")
		article.Author = c.PostForm("Author")
		article.Content = c.PostForm("Content")
		article.Summary = c.PostForm("Summary")

		//处理封面图片上传
		coverFile, err := c.FormFile("Cover")

		if err != nil && coverFile == nil { //用户未上传新的图片
			var articleOld Article
			//根据id查找Cover路径
			if err := db.Select("cover").First(&articleOld, id).Error; err != nil {
				c.HTML(http.StatusBadRequest, "admin_edit_article.html", gin.H{"error": "无法获取封面图片"})
				return
			}
			coverPath = articleOld.Cover
		}

		if err != nil && coverFile != nil {
			c.HTML(http.StatusBadRequest, "admin_edit_article.html", gin.H{"username": adminUsername, "error": "无法获取封面图片"})
			return
		}

		// 如果有新图片上传，保存封面图片到指定路径
		if coverFile != nil {
			coverPath = fmt.Sprintf("static/image/upload/%d_%s", time.Now().Unix(), coverFile.Filename)
			if err := c.SaveUploadedFile(coverFile, coverPath); err != nil {
				c.HTML(http.StatusInternalServerError, "admin_new_article.html", gin.H{"username": adminUsername, "error": "无法保存封面图片"})
				return
			}
		}

		article.Cover = coverPath
		article.PublishDate, _ = time.Parse("2006-01-02T15:04", c.PostForm("PublishDate"))

		//保存新的表单
		if err := db.Save(&article).Error; err != nil {
			c.HTML(http.StatusInternalServerError, "admin_edit_article.html", gin.H{"username": adminUsername, "error": "无法更新文章"})
			return
		}

		c.Redirect(http.StatusFound, "/admin/articles")
	} else {
		c.Redirect(http.StatusFound, "/admin")
	}
}

// 管理员删除文章
func adminDeleteArticleHandler(c *gin.Context) {
	session := sessions.Default(c)
	admin := session.Get("admin")
	if admin != nil {
		id := c.Param("id")
		var article Article
		if err := db.Where("id = ?", id).First(&article).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "文章不存在"})
			return
		}
		if err := db.Delete(&article).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "无法删除文章"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"success": true, "message": "文章删除成功"})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "未登录"})
	}
}

// 注销登录
func adminLogoutHandler(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
	c.JSON(http.StatusOK, gin.H{"success": true})
}

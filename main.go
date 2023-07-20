// main.go

package main

import (
	"fmt"
	"html/template"
	"net/http"
	"sort"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var db *gorm.DB

// 定义分页函数plus 和 minus 函数
var templateFuncs = template.FuncMap{
	"plus":  func(a, b int) int { return a + b },
	"minus": func(a, b int) int { return a - b },
}

func main() {
	router := gin.Default()

	//处理跨域问题
	router.Use(cors.Default())

	// 设置session中间件
	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("session", store))

	// 注册 html/template 作为模板引擎，并定义 plus 和 minus 函数
	html := template.New("").Funcs(templateFuncs)
	html, errs := html.ParseGlob("templates/*.html")
	if errs != nil {
		panic("无法解析模板：" + errs.Error())
	}
	router.SetHTMLTemplate(html)

	router.Static("/static", "./static")

	// 配置上传文件保存目录
	router.MaxMultipartMemory = 15 << 20 // 8MB

	// 连接 MySQL 数据库
	var err error
	db, err = gorm.Open("mysql", "root:295295@tcp(localhost:3306)/staticCloudy?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		panic("无法连接到数据库")
	}
	defer db.Close()

	// 自动迁移数据库表结构
	db.AutoMigrate(&Article{})

	router.GET("/", homeHandler)
	router.GET("/articles", getArticlesHandler)
	router.GET("/articles/:id/view", getArticlePageHandler)
	router.POST("/articles/:id/like", likeArticleHandler)

	// 添加文章管理路由
	adminGroup := router.Group("/admin")
	adminGroup.GET("", adminLoginHandler)
	adminGroup.POST("/login", adminLoginSubmitHandler)
	adminGroup.GET("/dashboard", adminDashboardHandler)
	adminGroup.GET("/articles", adminArticlesHandler)
	adminGroup.GET("/articles/new", adminNewArticleHandler)
	adminGroup.POST("/articles", adminCreateArticleHandler)
	adminGroup.GET("/articles/:id/edit", adminEditArticleHandler)
	adminGroup.POST("/articles/:id", adminUpdateArticleHandler)
	adminGroup.DELETE("/articles/:id", adminDeleteArticleHandler)
	adminGroup.POST("/logout", adminLogoutHandler)

	fmt.Println("服务器已启动,访问地址:http://0.0.0.0:8080")
	router.Run(":8080")
}

func homeHandler(c *gin.Context) {
	// 假设您的数据库表名为 articles，使用 your-database-driver 进行数据库查询
	var articles []Article
	if err := db.Find(&articles).Error; err != nil {
		// 处理错误
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"message": "数据库查询失败",
		})
		return
	}

	// 根据发布时间排序
	sort.SliceStable(articles, func(i, j int) bool {
		return articles[i].PublishDate.After(articles[j].PublishDate)
	})

	// 只保留时间最新的两篇文章
	var latestArticles []Article
	if len(articles) > 2 {
		latestArticles = articles[:2]
	} else {
		latestArticles = articles
	}

	// 根据访问量排序
	sort.SliceStable(articles, func(i, j int) bool {
		return articles[i].Views > articles[j].Views
	})

	// 只保留访问量最高的两篇文章
	var popularArticles []Article
	if len(articles) > 2 {
		popularArticles = articles[:2]
	} else {
		popularArticles = articles
	}

	// 根据喜欢数排序
	sort.SliceStable(articles, func(i, j int) bool {
		return articles[i].Likes > articles[j].Likes
	})

	// 只保留喜欢数最大的两篇文章
	var topLikedArticles []Article
	if len(articles) > 2 {
		topLikedArticles = articles[:2]
	} else {
		topLikedArticles = articles
	}

	// 渲染模板并传递数据给前端
	c.HTML(http.StatusOK, "home.html", gin.H{
		"latestArticles":   latestArticles,
		"popularArticles":  popularArticles,
		"topLikedArticles": topLikedArticles,
	})
}

package main

import (
	"github.com/gin-gonic/gin"
)

func router() *gin.Engine {
	r := gin.Default()

	// 获取用户身份
	r.Use(authen())

	r.LoadHTMLGlob("templates/*")

	oauthApiRoute(r)
	userApiRoute(r, "/api")
	postApiRoute(r, "/api")
	replyApiRoute(r, "/api")
	taxApiRoute(r, "/api")

	r.GET("/test", func(c *gin.Context) {
		exist, user, err := getCurrentUser()
		if exist || err != nil {
			c.String(500, "exist or error")
			return
		}
		c.String(200, "hello "+user.Name)
	})
	return r
}

func getCurrentUser() (bool, *User, error) {
	return false, &User{Name: "guotie"}, nil
}

func oauthApiRoute(r *gin.Engine) {
	r.GET("/oauth/login", oauthLogin)
	r.GET("/oauth/callback/github", githubLogin)
}

func userApiRoute(r *gin.Engine, prefix string) {
	g := r.Group(prefix + "/user")
	{
		g.GET("/current", userCurrent)
		g.GET("/exist", userNameExist)

		g.GET("/info/:username", getUserInfo)
		g.POST("/info/:username", modifyUserInfo)

		g.GET("/t/:username/topics", getUserTopics)
		g.GET("/t/:username/reply", getUserReplies)

		g.GET("/a/:username/favor", getUserFavor)
		g.GET("/a/:username/up", getUserUp)
		g.GET("/a/:username/down", getUserDown)
		g.GET("/a/:username/pay", getUserPay)
	}
}

func postApiRoute(r *gin.Engine, prefix string) {
	r.GET(prefix+"/topics", getTopics)
	g := r.Group(prefix + "/topics")
	{
		g.GET("/", getTopics)
		g.GET("/digest", getTopicsByDigest)
		g.GET("/latest", getTopicsByLatest)
		g.GET("/rocket", getTopicsByRocket)
		g.GET("/controversy", getTopicsByControversy)
	}
	gt := r.Group(prefix + "/topic")
	{
		gt.GET("/:id", getTopic)
		gt.POST("/:id", createNewTopic)
		gt.PUT("/:id", modifyTopic)
		gt.DELETE("/:id", deleteTopic)
	}
}

func replyApiRoute(r *gin.Engine, prefix string) {
	r.GET(prefix+"/reply/:id", getReply)
}

func taxApiRoute(r *gin.Engine, prefix string) {
	r.GET(prefix+"/cat/topics", getTaxonomyTopics)
	g := r.Group(prefix + "/tax")
	{
		g.GET("/list", getCategoryList)
		g.POST("/create/:id", needLogin(), createCategory)
	}
}

package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"strings"
)

/*
用户：
GET /api/user/current
*/
func userCurrent(c *gin.Context) {
	v, ok := c.Get(CurrentUserKey)
	if !ok {
		c.JSON(200, RespResult{ErrCodeVisitor, ErrMsgVisitor, nil})
		return
	}
	c.JSON(200, RespResult{0, "ok", v.(*User)})
}

// 用户名是否存在：
// GET /api/user/exist?name=xxx
func userNameExist(c *gin.Context) {
	exist, err := nameHasExist(c.Query("name"))
	if err != nil {
		c.JSON(200, RespResult{ErrCodeDBQuery, err.Error(), nil})
		return
	}
	c.JSON(200, RespResult{0, "ok", exist})
}

/*
用户资料：
GET /api/user/info/{username}
*/
func getUserInfo(c *gin.Context) {
	name := strings.TrimSpace(c.Param("username"))
	if name == "" {
		c.JSON(200, RespResult{ErrCodeParam, "param username is empty", nil})
		return
	}
	u := getUserByName(name)
	if u == nil {
		c.JSON(200, RespResult{ErrCodeNotFound, "Not found User by username: " + name, nil})
		return
	}
	// 不显示email
	u.Email = ""
	c.JSON(200, RespResult{0, "ok", u})
}

/*
修改用户资料：
POST /api/user/info/{username}
*/
func modifyUserInfo(c *gin.Context) {
	var u User

	ct := c.Request.Header.Get("Content-Type")
	// content-type必须包含"application/json"
	// "application/json"
	// "application/json; charset=utf-8"
	if strings.Contains(ct, "application/json") == false {
		c.JSON(200, RespResult{ErrCodeContentType, "Post Content-Type invalid: " + ct, nil})
		return
	}

	defer c.Request.Body.Close()
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(200, RespResult{ErrCodeReadRequest, err.Error(), nil})
		return
	}
	if err = json.Unmarshal(body, &u); err != nil {
		c.JSON(200, RespResult{ErrCodeUnmarshalUser, err.Error(), nil})
		return
	}
	// 写入数据库中
	if err = updateUserById(&u); err != nil {
		c.JSON(200, RespResult{ErrCodeUpdateUser, err.Error(), nil})
		return
	}
	c.JSON(200, RespResult{0, "ok", nil})
	return
}

// getUserTopics
// getUserReplies
func getUserPostByType(c *gin.Context, typ string) {
	var (
		u     *User
		start int
		count int
	)

	name := strings.TrimSpace(c.Param("username"))
	if name == "" {
		c.JSON(200, RespResult{ErrCodeParam, "param username is empty", nil})
		return
	}
	if u = getUserByName(name); u == nil {
		c.JSON(200, RespResult{ErrCodeNotFound, "Not found User by username: " + name, nil})
		return
	}

	start = getQueryIntDefault(c, "start", 0)
	count = getQueryIntDefault(c, "count", 10)

	posts, err := getPostsByUser(u, start, count, map[string]string{"object_type": typ})
	if err != nil {
		c.JSON(200, RespResult{ErrCodeGetUserPosts, err.Error(), "object type: " + typ})
		return
	}
	c.JSON(200, RespResult{0, "ok", posts})
	return
}

/*
获取用户文章
GET /api/user/t/{username}/topics?start=x&count=y
*/
func getUserTopics(c *gin.Context) {
	getUserPostByType(c, "post")
}

/*
获取用户评论
GET /api/user/t/{username}/reply?start=x&count=y
*/
func getUserReplies(c *gin.Context) {
	getUserPostByType(c, "comment")
}

/*
获取用户收藏
GET /api/user/a/{username}/favor?start=x&count=y
获取用户点赞
GET /api/user/a/{username}/up?start=x&count=y
获取用户反对
GET /api/user/a/{username}/down?start=x&count=y
获得用户打赏：
GET /api/user/a/{username}/pay?start=x&count=y
*/
func getUserFavor(c *gin.Context) {
	getUserPostByAction(c, "favor")
}
func getUserUp(c *gin.Context) {
	getUserPostByAction(c, "up")

}
func getUserDown(c *gin.Context) {
	getUserPostByAction(c, "down")

}
func getUserPay(c *gin.Context) {
	getUserPostByAction(c, "pay")

}

func getUserPostByAction(c *gin.Context, action string) {
	var (
		u     *User
		start int
		count int
	)

	name := strings.TrimSpace(c.Param("username"))
	if name == "" {
		c.JSON(200, RespResult{ErrCodeParam, "param username is empty", nil})
		return
	}
	if u = getUserByName(name); u == nil {
		c.JSON(200, RespResult{ErrCodeNotFound, "Not found User by username: " + name, nil})
		return
	}

	start = getQueryIntDefault(c, "start", 0)
	count = getQueryIntDefault(c, "count", 10)

	posts, err := getPostsByUserAction(u, action, start, count)
	if err != nil {
		c.JSON(200, RespResult{ErrCodeGetPostsByAction, err.Error(), "action type: " + action})
		return
	}
	c.JSON(200, RespResult{0, "ok", posts})
	return
}

/*
** 获取用户关注列表：

GET /api/user/{username}/follow?start=x&count=y

** 获取block列表

GET /api/user/{username}/block?start=x&count=y

用户动作：
block用户/unblock 用户
POST/DELETE /api/block/{uid1, uid2,uid3...}

follow用户
POST/DELETE /api/block/{uid1,uid2,uid3...}
*/

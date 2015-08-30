package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/smtc/glog"
	"io/ioutil"
	"net/http"
)

/*
# 文章：

获取默认文章列表：
GET /api/topics?start=xx&count=xx
*/
func getTopics(c *gin.Context) {
	getTopicByPoint(c, "default")
}

/*
获取优质文章列表：
GET /api/topics/digest
*/
func getTopicsByDigest(c *gin.Context) {
	getTopicByPoint(c, "digest")
}

/*
获取最新创建文章列表：
GET /api/topics/latest
*/
func getTopicsByLatest(c *gin.Context) {
	getTopicByPoint(c, "latest")
}

/*
获取上升最快
GET /api/topics/rocket
*/
func getTopicsByRocket(c *gin.Context) {
	getTopicByPoint(c, "rocket")
}

/*
获取争议
GET /api/topics/controversy
*/
func getTopicsByControversy(c *gin.Context) {
	getTopicByPoint(c, "controversy")
}

func getTopicByPoint(c *gin.Context, point string) {
	start := getQueryIntDefault(c, "start", 0)
	count := getQueryIntDefault(c, "count", 20)

	posts, ret, err := getPostsByPoint(point, start, count)
	if ret != 0 {
		glog.Error("getPostsByPoint failed: point=%s ret=%d err=%s\n", point, ret, err.Error())
		c.JSON(200, RespResult{ret, "get post by " + point + " failed: " + err.Error(), nil})
		return
	}
	c.JSON(200, RespResult{0, "ok", posts})
}

/*
GET /api/topic/{id}
*/
func getTopic(c *gin.Context) {
	post, err := getPostById(c.Param("id"), true)
	if err != nil {
		c.JSON(200, RespResult{ErrCodeNotFound, err.Error(), nil})
		return
	}
	c.JSON(200, RespResult{0, "ok", post})
	return
}

/*
创建文章：
POST /api/topic/{id}

1 post的body是一个json结构体，post从这个结构体unmarshal得到
2 post的类别信息也在这个结构体中，字段设置如下：
	taxonomy   string
	name       string
	或者：
	tax_id   string
*/
func createNewTopic(c *gin.Context) {
	if c.Request.Method != "POST" {
		c.JSON(200, RespResult{ErrCodeInvalidMethod, "invalid method, should be post", nil})
		return
	}

	user := getCurrent(c)
	if user == nil {
		c.JSON(200, RespResult{ErrCodeNeedLogin, "need login", nil})
		return
	}

	post, tax, err := readPostFromRequest(c)
	if err != nil {
		return
	}
	id := c.Param("id")
	if post.Id != id {
		c.JSON(200, RespResult{ErrCodeParam, "post id not equal with path param id", nil})
		return
	}

	if post.AuthorId != "" {
		if user.Id != post.AuthorId {
			c.JSON(200, RespResult{ErrCodeNeeAuthen, "you have no permission", nil})
			return
		}
	} else {
		post.AuthorId = user.Id
		post.AuthorName = user.Name
	}

	if err := createPost(post, user); err != nil {
		c.JSON(200, RespResult{ErrCodeDBQuery, err.Error(), nil})
		return
	} else {
		// 保存post与taxnomy的关系
		err = setPostTaxonmoy(post, tax)
		if err != nil {
			glog.Error("set post-taxonomy failed: post id %s\n", post.Id)
		}
	}
	c.JSON(200, RespResult{0, "ok", nil})
}

// 从post header中解出post和post的taxonomy
func readPostFromRequest(c *gin.Context) (*Post, []*Taxonomy, error) {
	var (
		r       *http.Request = c.Request
		post    Post
		taxes   []*Taxonomy
		taxInfo CategoryInfo
	)

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		c.JSON(200, RespResult{ErrCodeReadRequest, err.Error(), nil})
		return &post, nil, err
	}

	err = json.Unmarshal(body, &post)
	if err != nil {
		c.JSON(200, RespResult{ErrCodeUnmarshalPost, "unmarshal post request body failed: " + err.Error(), nil})
		return &post, nil, err
	}

	if err = json.Unmarshal(body, &taxInfo); err == nil {
		//glog.Info("category info: %v %s\n", taxInfo, string(body))
		taxes = getTaxFromInfos(taxInfo.Infoes)
	} else {
		glog.Error("unmarshal categroy info failed: %s %s\n", err.Error(), string(body))
		taxes = []*Taxonomy{&UnCategory}
	}

	return &post, taxes, nil
}

/*
修改文章：
PUT /api/topic/{id}
*/
func modifyTopic(c *gin.Context) {
	if c.Request.Method != "PUT" {
		c.JSON(200, RespResult{ErrCodeInvalidMethod, "invalid method, should be put", nil})
		return
	}

	post, err := getPostById(c.Param("id"), false)
	if err != nil {
		c.JSON(200, RespResult{ErrCodeNotFound, err.Error(), nil})
		return
	}

	user := getCurrent(c)
	if user == nil {
		c.JSON(200, RespResult{ErrCodeNeedLogin, "need login", nil})
		return
	}

	// Todo: 管理员帐号授权
	if post.AuthorId != user.Id {
		c.JSON(200, RespResult{ErrCodeNeeAuthen, "you have no permission", nil})
		return
	}

	mpost, _, err := readPostFromRequest(c)
	if err != nil {
		return
	}
	post.modifyTitle(mpost.Title)
	post.modifyContent(mpost.Content)
	if err = post.flush(); err != nil {
		c.JSON(200, RespResult{ErrCodeDBQuery, err.Error(), nil})
		return
	}
	c.JSON(200, RespResult{0, "ok", nil})
}

/*
删除文章

DELETE /api/topic/{id}
*/
func deleteTopic(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(200, RespResult{ErrCodeParam, "param id is empty", nil})
		return
	}
	// todo: 鉴权
	user := getCurrent(c)
	if user == nil {
		c.JSON(200, RespResult{ErrCodeNeePerm, "need permission", nil})
		return
	}

	post, err := getPostById(id, true)
	if err != nil {
		c.JSON(200, RespResult{ErrCodeDBQuery, err.Error(), nil})
		return
	}
	if user.Id == "administrator" {
		user.parseUserCap()
	}
	if post.AuthorId != user.Id && user.capability["delete_others_posts"] == false {
		c.JSON(200, RespResult{ErrCodeNeePerm, "need permission", nil})
		return
	}

	err = db.Delete(&Post{}).Where("id=?", id).Error
	if err != nil {
		c.JSON(200, RespResult{ErrCodeDBQuery, "delete post " + id + " failed: " + err.Error(), nil})
		return
	}
	c.JSON(200, RespResult{0, "ok", nil})
}

/*
todo: 设置文章的类别，label等属性
*/

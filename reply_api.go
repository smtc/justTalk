package main

import (
	"fmt"
	"io/ioutil"

	"github.com/gin-gonic/gin"
)

/*
# 回复：

获取回复：
GET /api/reply/:id
*/
func getReply(c *gin.Context) {
	reply, err := getReplyById(c.Param("id"))
	if err != nil {
		c.JSON(200, RespResult{ErrCodeNotFound, err.Error(), nil})
		return
	}
	c.JSON(200, RespResult{0, "ok", reply})
	return
}

/*
id is the post's id
POST /api/reply/t/:id
*/
func createReplyAPI(c *gin.Context) {
	user := getCurrent(c)
	if user == nil {
		c.JSON(200, RespResult{ErrCodeNeedLogin, "login to reply", nil})
		return
	}

	pid := c.Param("id")
	if pid == "" {
		c.JSON(200, RespResult{ErrCodeParam, "param id is empty.", nil})
		return
	}
	post, err := getPostById(pid, false)
	if err != nil || post == nil {
		c.JSON(200, RespResult{ErrCodeGetPost, fmt.Sprintf("get post by id %s failed: %v", pid, err), nil})
		return
	}
	// read http post header
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(200, RespResult{ErrCodeReadRequest, err.Error(), nil})
		return
	}
	// 读取reply
	reply, err := readPostFromRequest(body)
	if err != nil {
		c.JSON(200, RespResult{ErrCodeReadReqPost, err.Error(), nil})
		return
	}
	setReplyDefaultValue(reply, post)
	if err = createReply(reply, post, user); err != nil {
		c.JSON(200, RespResult{ErrCodeCreateReply, err.Error(), nil})
		return
	}
	c.JSON(200, RespResult{0, "ok", ""})
	return
}

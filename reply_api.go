package main

import "github.com/gin-gonic/gin"

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

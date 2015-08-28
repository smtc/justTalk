package main

import (
	"github.com/gin-gonic/gin"
)

const (
	XAuthToken     = "X-Auth-Token"
	CurrentUserKey = "current_user"
)

// 身份中间件
func authen() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			tok string
			uid string
			cu  *User
		)

		if tok = c.Request.Header.Get(XAuthToken); tok == "" {
			return
		}
		// 从redis中获取用户的uid
		if uid = getFromCacheString(tok); uid == "" {
			return
		}
		if cu = getUserById(uid); cu == nil {
			return
		}
		c.Set(CurrentUserKey, cu)
	}
}

// 需要登陆
func needLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, ok := c.Get(CurrentUserKey)
		if !ok {
			// Todo: 直接重定向到登陆页面
			c.JSON(200, RespResult{ErrCodeNeedLogin, ErrMsgNeedLogin, nil})
			c.Abort()
			return
		}
	}
}

func getCurrent(c *gin.Context) *User {
	v, ok := c.Get(CurrentUserKey)
	if !ok {
		return nil
	}

	return v.(*User)
}

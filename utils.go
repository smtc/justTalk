package main

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	//"github.com/smtc/glog"
)

const (
	ErrCodeNeedLogin        = 100
	ErrMsgNeedLogin         = "need login"
	ErrCodeVisitor          = 101
	ErrMsgVisitor           = "current user is visitor"
	ErrCodeDBQuery          = 102
	ErrMsgDBQuery           = "db query failed"
	ErrCodeParam            = 103
	ErrCodeNotFound         = 104
	ErrCodeContentType      = 105
	ErrCodeReadRequest      = 106
	ErrCodeUnmarshalUser    = 107
	ErrCodeUpdateUser       = 108
	ErrCodeGetUserPosts     = 109
	ErrCodeGetUserReplies   = 110
	ErrCodeGetPostsByAction = 111
	ErrCodeInvalidMethod    = 112
	ErrCodeUnmarshalPost    = 113
	ErrCodeNeeAuthen        = 114
	ErrCodeNeePerm          = 115
	ErrCodeUnmarshaLTax     = 116
	ErrCodeInvalidUUID      = 117
	ErrCodeGetPost          = 118
	ErrCodeReadReqPost      = 119
	ErrCodeCreateReply      = 120
)

var (
	// xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
	// 54c7dae5-dc4b-4121-bf40-bcbd1f2958d8
	idRe = regexp.MustCompile("[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}")
)

// 返回的数据格式
// code:    0:成功；其他值失败
// message: 对code的解释
// data:    返回的数据，失败时，可能为空
type RespResult struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func getQueryInt(c *gin.Context, key string) int {
	v := strings.TrimSpace(c.Query(key))
	if v == "" {
		return 0
	}
	i, err := strconv.Atoi(v)
	if err != nil {
		return 0
	}
	return int(i)
}

// 获取query参数，如果没有或参数非法，返回默认值
func getQueryIntDefault(c *gin.Context, key string, d int) int {
	v := strings.TrimSpace(c.Query(key))
	if v == "" {
		return d
	}
	i, err := strconv.Atoi(v)
	if err != nil {
		return d
	}
	return int(i)
}

// 校验uuid是否合法
func validUUID(id string) bool {
	return idRe.MatchString(id)
}

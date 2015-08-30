package main

import (
	"errors"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/smtc/glog"
)

func getReplyById(id string) (post *Post, err error) {
	if id == "" {
		err = errors.New("id is empty")
		return
	}
	err = db.Where("id=?", id).First(post).Error

	return
}

// 查找一个主题的回复
func getReplyByPid(pid string, start, count int, options ...interface{}) (posts []*Post, err error) {
	q := db.Where("object_type=?", "reply").
		Where("post_parent=?", pid)

	if len(options) > 0 {
		opt := options[0].(map[string]interface{})
		for key, val := range opt {
			switch key {
			case "points":
				sval := strings.ToLower(val.(string))
				if sval != "desc" && sval != "asc" {
					glog.Warn("option for points invalid: %s, should be asc or desc.\n", sval)
					sval = "desc"
				}
				q = q.Order("points " + sval)
			}
		}
	} else {
		q = q.Order("created_at desc")
	}

	err = q.Offset(start).
		Limit(count).
		Find(&posts).Error
	if err != nil && err != gorm.RecordNotFound {
		return
	}

	return
}

// 查找用户的回复
func getReplyByUser(user *User, start, count int) ([]*Post, error) {
	return getPostsByUser(user, start, count, map[string]interface{}{
		"object_type": "reply",
	})
}

// 创建回复
func createReply(reply *Post, post *Post, user *User) (err error) {
	return
}

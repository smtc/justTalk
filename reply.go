package main

import (
	"errors"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/smtc/glog"
)

// 设置reply的默认值
func setReplyDefaultValue(reply, post *Post) {
	reply.ReplyCount = 0
	reply.MenuOrder = 0
	reply.LikedCount = 0
	reply.BookmarkCount = 0
	reply.StarCount = 0
	reply.BlockCount = 0
	reply.Digest = 0
	reply.Points = 0

	// 与post不同的地方
	reply.PostParent = post.Id
	reply.Floor = post.ReplyCount + 1

	if reply.ObjectType == "" {
		reply.ObjectType = "post"
	}
	if reply.SubType == "" {
		reply.SubType = "reply"
	}
	if reply.PostStatus == "" {
		reply.PostStatus = "normal"
	}
	if reply.ReplyStatus == "" {
		reply.ReplyStatus = "open"
	}
	if reply.PingStatus == "" {
		reply.PingStatus = "open"
	}
}

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
	q := db.Where("sub_type=?", "reply").
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
		q = q.Order("floor asc")
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
	setPostDefaultValue(reply)

	reply.PostParent = post.Id
	reply.AuthorId = user.Id
	reply.AuthorName = user.Name

	err = db.Create(reply).Error
	if err == nil {
		post.ReplyCount += 1
		post.LastReplyAt = time.Now()
		if err2 := db.Save(post).Error; err2 != nil {
			glog.Error("update post %s reply info failed, reply id: %s\n", post.Id, reply.Id)
		}
	} else {
		return err
	}

	return
}

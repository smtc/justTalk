package main

import (
	"errors"
	"sort"
	"time"

	"github.com/jinzhu/gorm"
	//"github.com/satori/go.uuid"
	"github.com/smtc/glog"
)

// Post table
type Post struct {
	Id          string    `gorm:"primary_key" json:"id"`
	SiteId      int64     `json:"site_id"` // not used currently
	AuthorId    string    `json:"author_id"`
	AuthorName  string    `sql:"size:40" json:"author_name"`
	Title       string    `sql:"size:60" json:"title"`
	Content     string    `sql:"type:TEXT" json:"content"`
	ObjectType  string    `sql:"size:30" json:"obj_type"` // post, reply, etc...
	SubType     string    `sql:"size:20" json:"sub_type"` // 二级类型 问答，投票等
	PostStatus  string    `sql:"size:20" json:"post_status"`
	Excerpt     string    `sql:"size:500" json:"excerpt"`
	PostAt      time.Time `json:"post_at"`
	CreatedAt   time.Time `json:"created_at"`
	ModifyAt    time.Time `json:"modify_at"`
	ClosedAt    int64     `json:"closed_at"`
	CloseReason string    `json:"close_reason"`

	ReplyStatus     string    `sql:"size:20;default:'open'" json:"reply_status"`
	PingStatus      string    `sql:"size:20;default:'open'" json:"ping_status"`
	PostName        string    `sql:"size:200" json:"post_name"`
	PostPassword    string    `sql:"size:20" json:"post_password"`
	ToPing          string    `sql:"type:text" json:"to_ping"`
	Pinged          string    `sql:"type:text" json:"pinged"`
	ContentFiltered string    `sql:"type:text" json:"content_filtered"`
	PostParent      string    `sql:"size:20;index" json:"post_parent"`
	MenuOrder       int       `json:"menu_order"`
	PostMimeType    string    `sql:"size:200" json:"post_mime_type"`
	ReplyCount      int64     `json:"reply_count"`
	LastReplyAt     time.Time `json:"last_reply_at"`
	LikedCount      int       `json:"liked_count"`
	BookmarkCount   int       `json:"bookmark_count"`
	StarCount       int       `json:"star_count"`
	BlockCount      int       `json:"block_count"`
	Points          int       `json:"point"`
	Digest          int       `json:"digest"`

	Replies []*Post `sql:"-" json:"replies"`
}

func setPostDefaultValue(post *Post) {
	post.ReplyCount = 0
	post.MenuOrder = 0
	post.LikedCount = 0
	post.BookmarkCount = 0
	post.StarCount = 0
	post.BlockCount = 0
	post.Digest = 0
	post.Points = 0

	if post.ObjectType == "" {
		post.ObjectType = "post"
	}
	if post.SubType == "" {
		post.SubType = "post"
	}
	if post.PostStatus == "" {
		post.PostStatus = "normal"
	}
	if post.ReplyStatus == "" {
		post.ReplyStatus = "open"
	}
	if post.PingStatus == "" {
		post.PingStatus = "open"
	}
}

// 补充内容
type AppendContent struct {
	AppendAt time.Time `json:"append_at"`
	Content  string    `json:"content"`
}

// 获得用户的文章
// todo: 处理options参数
// options not used.
func getPostsByUser(u *User, start, count int, options ...interface{}) (posts []*Post, err error) {
	q := db.Where("author_id=?", u.Id)
	if len(options) > 0 {
		opt := options[0].(map[string]interface{})
		if typ, ok := opt["object_type"]; ok {
			q = q.Where("object_type=?", typ.(string))
		}
	}
	err = q.Order("publish_at desc").Offset(start).Limit(count).Find(&posts).Error
	if err == gorm.RecordNotFound {
		err = nil
	}

	return
}

func getPostById(id string) (*Post, error) {
	var (
		post Post
		err  error
	)

	if id == "" {
		err = errors.New("id is empty")
		return &post, err
	}
	err = db.Where("id=?", id).First(&post).Error
	if err != nil {
		post.Replies, err = getReplyByPid(post.Id, 0, 20)
	}

	return &post, err
}

// 更加不同的角度，查找posts
// point: default, digest, rocket, latest, controversy
func getPostsByPoint(point string, start, count int) (posts []*Post, ret int, err error) {
	q := db.Where("object_type=post")
	switch point {
	case "default":
		q.Order("points desc")
	case "digest":
		q.Where("digest=1").Order("publish_at desc")
	case "latest":
		q.Order("post_at desc")
	case "rocket":
	// Todo
	case "controversy":
	// Todo
	default:
		glog.Warn("getPostsByPoint: unknown piont: %s\n", point)
	}

	err = q.Offset(start).Limit(count).Find(&posts).Error
	if err != gorm.RecordNotFound {
		err = nil
	}

	return
}

// 创建新的post
// 鉴权
// sanitize
func createPost(post *Post) (err error) {
	/*
		var (
			npost Post
		)

		npost.Id = post.Id
		npost.Title = post.Title
		npost.Content = post.Content
		npost.AuthorId = post.AuthorId
		npost.AuthorName = post.AuthorName

		err = db.Save(&npost).Error
	*/
	if post.Id == "" {
		return errors.New("createPost: post id should supply.")
	}
	setPostDefaultValue(post)
	err = db.Save(post).Error
	return
}

// 修改post的title
func (post *Post) modifyTitle(title string) {
	// Todo: sanitize
	post.Title = title
	return
}

// 修改post的内容
func (post *Post) modifyContent(content string) {
	// Todo: sanitize
	post.Content = content
	return
}

// 补充post的content内容
func (post *Post) appendContent(content string) {
	return
}

// 写回数据库中
func (post *Post) flush() error {
	return db.Save(post).Error
}

// taxonomy
func getTopicsByTaxonomy(tax string, start, count int) (posts []*Post, err error) {
	term, err := getTaxByName(tax, "category")
	if err != nil {
		if err == gorm.RecordNotFound {
			err = nil
		}
		return
	}
	// 写不出join语句，先一个个查吧...
	ids, err := getObjectsByTerm(term, start, count)
	if err != nil {
		if err == gorm.RecordNotFound {
			err = nil
		}
		return
	}

	for _, id := range ids {
		if post, err := getPostById(id.ObjectId); err == nil {
			posts = append(posts, post)
		} else {
			glog.Warn("getTopicsByTaxonomy: Not found post by id: %s tax: %s\n", id.ObjectId, tax)
		}
	}

	postSlice(posts).orderedBy("points", true)
	return
}

// 对posts按照不同的键值排序
// sort
type postSlice []*Post
type lessFunc func(p1, p2 *Post) bool

// multiSorter implements the Sort interface, sorting the posts within.
type multiSorter struct {
	posts []*Post
	less  []lessFunc
}

// Len is part of sort.Interface.
func (ms *multiSorter) Len() int {
	return len(ms.posts)
}

// Swap is part of sort.Interface.
func (ms *multiSorter) Swap(i, j int) {
	ms.posts[i], ms.posts[j] = ms.posts[j], ms.posts[i]
}

// Less is part of sort.Interface. It is implemented by looping along the
// less functions until it finds a comparison that is either Less or
// !Less. Note that it can call the less functions twice per call. We
// could change the functions to return -1, 0, 1 and reduce the
// number of calls for greater efficiency: an exercise for the reader.
func (ms *multiSorter) Less(i, j int) bool {
	p, q := ms.posts[i], ms.posts[j]
	// Try all but the last comparison.
	var k int
	for k = 0; k < len(ms.less)-1; k++ {
		less := ms.less[k]
		switch {
		case less(p, q):
			// p < q, so we have a decision.
			return true
		case less(q, p):
			// p > q, so we have a decision.
			return false
		}
		// p == q; try the next comparison.
	}
	// All comparisons to here said "equal", so just return whatever
	// the final comparison reports.
	return ms.less[k](p, q)
}

func orderedBy(less ...lessFunc) *multiSorter {
	return &multiSorter{
		less: less,
	}
}

func (ps postSlice) orderedBy(key string, desc bool) {
	var ms multiSorter

	byPoints := func(p1, p2 *Post) bool {
		return p1.Points < p2.Points
	}
	byPosted := func(p1, p2 *Post) bool {
		return p2.PostAt.After(p1.PostAt)
	}

	switch key {
	case "points":
		ms = multiSorter{less: []lessFunc{byPoints}}
	case "post_at":
		ms = multiSorter{less: []lessFunc{byPosted}}
	default:
		glog.Warn("orderedBy: posts sort by key %s not implement yet.\n", key)
		return
	}

	ms.posts = ps
	if desc {
		sort.Sort(sort.Reverse(&ms))
	} else {
		sort.Sort(&ms)
	}
}

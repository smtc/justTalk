package main

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"
	"github.com/smtc/glog"
)

type User struct {
	//Id         int64  `json:"id"`
	Id         string `gorm:"primary_key"` //`sql:"size:64;unique_index" json:"object_id"`
	SiteId     int64  `json:"site_id"`
	Name       string `sql:"size:40;unique_index" json:"name"`
	Email      string `sql:"size:100;index" json:"email"`
	Msisdn     string `sql:"size:20;index" json:"msisdn"`
	Password   string `sql:"size:200" json:"-"`
	MainId     int64  `json:"main_id"` // 如果不是主用户, 主用户id；否则为0
	Approved   bool   `json:"approved"`
	Activing   bool   `json:"acitiving"`
	Banned     bool   `json:"banned"`
	ApprovedBy string `sql:"size:20" json:"approved_by"`
	IpAddr     string `sql:"size:30" json:"ipaddr"`
	DaysLogin  int    `json:"days_login"`

	Birthday   string    `json:"birthday"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	ApprovedAt time.Time `json:"approved_at"`
	LastLogin  time.Time `json:"last_login"`
	LastPost   int64     `json:"last_post"`

	Avatar string `sql:"size:120" json:"avatar"`
	City   string `sql:"size:40" json:"city"`
	Blog   string `sql:"size:200"`
	Qq     string `sql:"size:40"`
	Weibo  string `sql:"size:40"`
	Weixin string `sql:"size:40"`

	Goldcoins     int64 `json:"gold_coins"`   // 金币
	SilverCoins   int64 `json:"silver_coins"` // 银币
	Coppercoins   int64 `json:"copper_coins"` // 铜币
	Reputation    int   `json:"reputation"`   // 声望值
	Credits       int   `json:"credits"`      // 信用等级
	Experience    int   `json:"experience"`   // 经验值
	Activation    int   `json:"activation"`   // 活跃度
	Redcards      int   `json:"red_cards"`
	Yellowcards   int   `json:"yellow_cards"`
	Notifications int   `json:"notifications"`

	// 第三方登录的数据
	ProviderName string `json:"provider_name"` // 第三方登录名称
	ProviderId   string `json:"provider_id"`   // 用户在第三方的唯一id

	// user capability
	capability map[string]bool `sql:"-"`
	capParsed  bool            `sql:"-"`
	roles      string          `sql:"type:text"` // 这是一个string数组, 以,分割
	// other meta data
	metaData map[string]interface{} `sql:"-"`
}

// 用户收藏，喜欢，反对，打赏等
type UserAction struct {
	Id        string    `gorm:"primary_key" sql:"size:120" json:"id"`
	UserId    string    `sql:"size:120" json:"user_id"`
	Action    string    `sql:"size:40" json:"action"`
	TargetTyp string    `sql:"size:40" json:"target_typ"`
	TargetId  string    `sql:"size:120" json:"target_id"`
	Data      string    `sql:"size:200" json:"data"`
	CreatedAt time.Time `json:"created_at"`
}

// 创建用户
// 必须字段：
//    id: uuid, if empty, create it
//    name: username, must be uniqe
func createProviderUser(u *User) (err error) {
	var exist bool

	u.Id = strings.TrimSpace(u.Id)
	u.Name = strings.TrimSpace(u.Name)
	u.Email = strings.TrimSpace(u.Email)
	u.ProviderName = strings.TrimSpace(u.ProviderName)
	u.ProviderId = strings.TrimSpace(u.ProviderId)
	if u.ProviderId == "" || u.ProviderName == "" {
		err = errors.New("provider name and provider id should NOT be empty.")
		return
	}

	if u.Id == "" {
		u.Id = uuid.NewV4().String()
		fmt.Println("Id:", u.Id)
	} else {
		if exist, err = idHasExist(u.Id); exist || err != nil {
			err = fmt.Errorf("create user failed: id %s exist or db query failed: %s", u.Id, err.Error())
			return
		}
	}
	if u.Name == "" {
		err = errors.New("create user failed: name should not be empty.")
		return
	} else if exist, err = nameHasExist(u.Name); exist || err != nil {
		err = fmt.Errorf("create user failed: name %s exist or db query failed: %s", u.Name, err.Error())
		return
	}
	if u.Email != "" {
		if exist, err = emailHasExist(u.Email); exist || err != nil {
			err = fmt.Errorf("create user failed: email %s exist or db query failed: %s", u.Email, err.Error())
			return
		}
	}

	u.MainId = 0
	u.Approved = true
	u.Activing = true
	u.Banned = false
	u.ApprovedBy = "auto"
	u.IpAddr = ""
	u.DaysLogin = 0
	u.Birthday = ""

	now := time.Now()
	u.CreatedAt = now
	u.UpdatedAt = now
	u.ApprovedAt = now
	u.LastPost = 0

	u.Goldcoins, u.SilverCoins, u.Coppercoins = 0, 0, 0
	u.Reputation, u.Credits, u.Experience, u.Activation = 0, 0, 0, 0
	u.Redcards, u.Yellowcards, u.Notifications = 0, 0, 0

	err = db.Create(u).Error

	return nil
}

// 根据ObjectId来查找用户
func getUserById(oid string) *User {
	var u User

	if err := db.Where("id=?", oid).First(&u).Error; err != nil {
		if err != gorm.RecordNotFound {
			glog.Error("getUserById failed: oid=%s err=%v\n", oid, err)
		}
		return nil
	}

	return &u
}

func getUserByEmail(email string) (*User, error) {
	var u User

	if err := db.Where("email=?", email).First(&u).Error; err != nil {
		if err != gorm.RecordNotFound {
			glog.Error("getUserByEmail failed: email=%s err=%v\n", email, err)
		}
		return nil, err
	}

	return &u, nil
}

// 根据用户名查找用户
//
func getUserByName(name string) *User {
	var u User

	if err := db.Where("name=?", name).First(&u).Error; err != nil {
		if err != gorm.RecordNotFound {
			glog.Error("getUserByName failed: name=%s err=%v\n", name, err)
		}
		return nil
	}

	return &u
}

// 根据provider name和provider id来查找该用户是否存在
func providerUserHasExist(pname, pid string) (bool, *User, error) {
	var (
		err  error
		user User
	)

	err = db.Where("provider_name=? AND provider_id=?", pname, pid).Find(&user).Error
	if err == gorm.RecordNotFound {
		return false, nil, nil
	}
	if err != nil {
		return false, nil, err
	}

	return true, &user, nil
}

// id是否已经存在
func idHasExist(id string) (exist bool, err error) {
	var u User

	err = db.Where("id=?", id).Find(&u).Error
	if err == gorm.RecordNotFound {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

// email是否已经存在
func emailHasExist(email string) (exist bool, err error) {
	var u User

	err = db.Where("email=?", email).Find(&u).Error
	if err == gorm.RecordNotFound {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

// 名字是否存在
func nameHasExist(name string) (exist bool, err error) {
	var u User

	err = db.Where("name=?", name).Find(&u).Error
	if err == gorm.RecordNotFound {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

// 更新用户资料
// 目前只有city和avatar可以变更
// 哪些资料可以变更？哪些资料不能变更？
func updateUserById(u *User) (err error) {
	var (
		ou      *User
		nu      User
		changed bool
	)

	ou = getUserById(u.Id)
	if ou == nil {
		return errors.New("Not found user by id " + u.Id)
	}

	if u.Avatar != "" && u.Avatar != ou.Avatar {
		nu.Avatar = u.Avatar
		changed = true
	}
	if u.City != "" {
		nu.City = u.City
		changed = true
	}
	if changed {
		err = db.Model(ou).Updates(nu).Error
	}

	return err
}

func getPostsByUserAction(u *User, action string, start, count int) ([]*Post, error) {
	var (
		err     error
		actions []UserAction
		posts   []*Post
	)

	err = db.Where("user_id=? and action=?", u.Id, action).Offset(start).Limit(count).Find(&actions).Error
	if err != nil {
		if err == gorm.RecordNotFound {
			return []*Post{}, nil
		}
		return []*Post{}, err
	}

	for _, a := range actions {
		var p = new(Post)
		p.Id = a.TargetId
		err = db.Find(p).Error
		if err != nil {
			glog.Error("Find Post by Id %s failed: %s\n", p.Id, err.Error())
			continue
		}
		posts = append(posts, p)
	}

	return posts, nil
}

// 2015-08-18 Todo:
// user的动作，例如收藏，star，up，down，pay等

// 根据用户的角色，分析、填充用户的权限
func (u *User) parseUserCap() {
	// todo: parse user's role, fill user capability

	u.capability = make(map[string]bool)
	// 2015-08-29
	// 为了测试：设置用户administrator的权限
	if u.Id == "administrator" {
		u.capability["create_taxonomy"] = true
	}

	u.capParsed = true
}

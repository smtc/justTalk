package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"runtime"
	"testing"

	"github.com/satori/go.uuid"
)

type userTest struct {
	User
	ctoken string // 再其他需要认证的测试函数中用到user的cookie
}

var (
	_             = fmt.Printf
	_             = runtime.Caller
	users         = make(map[string]*userTest)
	anonymous     = &userTest{User{Id: "anonymous", Name: "nobody", City: "nowhere"}, ""}
	administrator = &userTest{User{Id: "administrator", Name: "root"}, ""}
	usersCookie   = make(map[string]string)
)

func init() {
	var _users = []*userTest{
		&userTest{User{Id: uuid.NewV1().String(), Name: "guotie", City: "Nanjing", Msisdn: "15651889188"}, ""},
		&userTest{User{Id: uuid.NewV4().String(), Name: "铁哥"}, ""},
		&userTest{User{Id: uuid.NewV4().String(), Name: "天津人民很伤心"}, ""},
		&userTest{User{Id: uuid.NewV4().String(), Name: "luomessi"}, ""},
		&userTest{User{Id: uuid.NewV4().String(), Name: "somebody"}, ""},
	}

	for _, u := range _users {
		users[u.Name] = u
	}
}

func (ut *userTest) createSelf() error {
	return db.Create(&ut.User).Error
}

// 以user的身份访问特定的资源
func (ut *userTest) fetch(method, path string, data interface{}) (rr *RespResult, err error) {
	var (
		body   []byte
		rbody  []byte
		bodyRd io.Reader
		req    *http.Request
		resp   *http.Response
	)

	if data != nil {
		if body, err = json.Marshal(data); err != nil {
			return
		}
	}

	if method == "GET" {
		bodyRd = nil
	} else {
		bodyRd = bytes.NewBuffer(body)
	}
	req, err = http.NewRequest(method, "http://127.0.0.1:8008"+path, bodyRd)
	if err != nil {
		return nil, err
	}
	if ut.isVisitor() == false {
		req.Header.Add(XAuthToken, ut.ctoken)
	}

	// 请求
	client := &http.Client{}
	if resp, err = client.Do(req); err != nil {
		return
	}

	// 返回的数据
	defer resp.Body.Close()
	rbody, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	rr = new(RespResult)
	err = json.Unmarshal(rbody, rr)
	if err != nil {
		fmt.Println("unmarshal failed:", string(rbody))
	}
	//fmt.Println("fetch response: ", rr)
	return
}

// 假装登陆，就是获得一个ctoken，并把这个token保存在userTest结构体中，其他需要登陆验证的地方，可以在http头部
// 携带xtokenAuth字段
func (ut *userTest) mockLogin() error {
	if ut.Id == "anonymous" {
		return nil
	}

	ctoken, err := setUserCache(&cachedUserData{ut.Id, nil})
	if err != nil {
		return err
	}
	ut.ctoken = ctoken
	return nil
}

// 是否已经登陆过
func (ut *userTest) isVisitor() bool {
	return ut.ctoken == ""
}

// 准备用户数据
func prepareUserData() {
	for _, u := range users {
		err := u.createSelf()
		if err != nil {
			log.Printf("create user %s failed: %v\n", u.Name, err)
		}
	}
}

func cleanDatabase() {
	db.Exec("truncate users;")
	db.Exec("truncate posts;")
	db.Exec("truncate taxonomies;")
	db.Exec("truncate term_relations;")
}

// todo: 20150828 需要完善，让几个测试帐号全部登陆
// done: 20150829
func testUserCurrent(t *testing.T) {
	var err error

	for _, u := range users {
		err = u.mockLogin()
		if err != nil {
			t.Fail()
			return
		}
		rr, err := u.fetch("GET", "/api/user/current", nil)
		assert(t, err == nil, "response should have no errors")
		assert(t, rr.Code == 0, "response should be ok")
		res := rr.Data.(map[string]interface{})
		if res["name"].(string) != u.Name {
			t.Fail()
		}
	}

	// 超级用户
	administrator.createSelf()
	administrator.mockLogin()
	administrator.capability = map[string]bool{"create_taxonomy": true}

	rr, err := anonymous.fetch("GET", "/api/user/current", nil)
	assert(t, rr.Code == ErrCodeVisitor, "should be visitor")
}

func assert(t *testing.T, cond bool, fmt string, args ...interface{}) {
	if !cond {
		//buf := make([]byte, 500000)
		//runtime.Stack(buf, true)
		//t.Logf(string(buf))
		_, file, line, _ := runtime.Caller(1)
		t.Logf("%s:%d:", file, line)
		t.Errorf(fmt, args...)
	}
}

func testNameExist(t *testing.T, path string, res bool) {
	rr, err := anonymous.fetch("GET", path, "")
	assert(t, err == nil, "request should NOT fail.")

	assert(t, rr.Code == 0, "resp code should be 0")
	exist := rr.Data.(bool)
	assert(t, exist == res, "name should has %v", res)
}

func testUserNameExist(t *testing.T) {
	testNameExist(t, "/api/user/exist?name=guotie", true)
	testNameExist(t, "/api/user/exist?name=guotie9", false)
}

func testGetUserInfo(t *testing.T) {
	rr, err := anonymous.fetch("GET", "/api/user/info/guotie", "")
	assert(t, err == nil, "resp should not fail")
	u := rr.Data.(map[string]interface{})
	assert(t, u["city"].(string) == "Nanjing", "city should be nanjing")
	assert(t, u["msisdn"].(string) == "15651889188", "msisdn should be 15651889188")

	rr, err = anonymous.fetch("GET", "/api/user/info/nobody", "")
	assert(t, rr.Code == ErrCodeNotFound, "nobody should not exist")
}

func testUserApi(t *testing.T) {
	testUserCurrent(t)
	testUserNameExist(t)
	testGetUserInfo(t)
}

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"testing"

	"github.com/satori/go.uuid"
)

var (
	_     = fmt.Printf
	users = []User{
		{Id: uuid.NewV1().String(), Name: "guotie", City: "Nanjing", Msisdn: "15651889188"},
		{Id: uuid.NewV4().String(), Name: "铁哥"},
		{Id: uuid.NewV4().String(), Name: "天津人民很伤心"},
		{Id: uuid.NewV4().String(), Name: "luomessi"},
		{Id: uuid.NewV4().String(), Name: "somebody"},
	}
)

// 准备用户数据
func prepareUserData() {
	for _, u := range users {
		err := db.Create(&u).Error
		if err != nil {
			log.Printf("save user %s failed: %v\n", u.Name, err)
		}
	}
}

func cleanUserData() {
	db.Exec("truncate users;")
}

// 生成一个request请求，请求头部包含XAuthToken字段，并请求path，返回结果
func makeGetRequest(path string, ctoken string) (*RespResult, error) {
	var (
		rr RespResult
	)

	client := &http.Client{}
	req, err := http.NewRequest("GET", "http://127.0.0.1:8008"+path, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add(XAuthToken, ctoken)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &rr)

	return &rr, err
}

func testUserCurrent(t *testing.T) {
	// 将user保存在
	ctoken, err := setUserCache(&cachedUserData{users[0].Id, nil})
	if err != nil {
		t.Fail()
		return
	}

	rr, err := makeGetRequest("/api/user/current", ctoken)
	//fmt.Println(rr.Data)
	u := rr.Data.(map[string]interface{})
	if u["name"].(string) != users[0].Name {
		t.Fail()
	}
}

func assert(t *testing.T, cond bool, fmt string, args ...interface{}) {
	if !cond {
		t.Errorf(fmt, args...)
	}
}

func testNameExist(t *testing.T, path string, res bool) {
	rr, err := makeGetRequest(path, "")
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
	rr, err := makeGetRequest("/api/user/info/guotie", "")
	assert(t, err == nil, "resp should not fail")
	u := rr.Data.(map[string]interface{})
	assert(t, u["city"].(string) == "Nanjing", "city should be nanjing")
	assert(t, u["msisdn"].(string) == "15651889188", "msisdn should be 15651889188")

	rr, err = makeGetRequest("/api/user/info/nobody", "")
	assert(t, rr.Code == ErrCodeNotFound, "nobody should not exist")
}

func TestUserApi(t *testing.T) {
	testUserCurrent(t)
	testUserNameExist(t)
	testGetUserInfo(t)
}

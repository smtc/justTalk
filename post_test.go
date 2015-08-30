package main

import (
	"testing"

	"github.com/satori/go.uuid"
)

var _posts = make(map[string]*Post)

// 12个post
var _testPosts = []map[string]interface{}{
	map[string]interface{}{"id": uuid.NewV4().String(),
		"title":         "golang库：io",
		"category_info": []TaxInfo{{TaxName: "golang", Taxonomy: "category"}},
		"content":       "golang库：io使用说明\nio.Reader\nio.Writer"},
	map[string]interface{}{"id": uuid.NewV4().String(),
		"title": "golang库：net", "tax_name": "golang",
		"category_info": []TaxInfo{{TaxName: "golang", Taxonomy: "category"}},
		"content":       "golang库：net使用说明"},
	map[string]interface{}{"id": uuid.NewV4().String(),
		"title":         "golang库：json",
		"category_info": []TaxInfo{{TaxName: "golang", Taxonomy: "category"}},
		"content":       "golang库：json使用说明\njson.Marshal\njosn.Unmarshal"},
	map[string]interface{}{"id": uuid.NewV4().String(),
		"title":         "golang库：runtime",
		"category_info": []TaxInfo{{TaxId: testTaxes[0].Id}},
		"content":       "golang库：runtime使用说明\nruntime.Stack"},
	map[string]interface{}{"id": uuid.NewV4().String(),
		"title":         "golang库：os",
		"category_info": []TaxInfo{{TaxName: "golang", Taxonomy: "category"}},
		"content":       "golang库：os使用说明\nos.Sytem\nos.Exec"},
	map[string]interface{}{"id": uuid.NewV4().String(),
		"title":         "golang库：time",
		"category_info": []TaxInfo{{TaxId: testTaxes[0].Id}},
		"content":       "golang库：time使用说明\ntime.Unix\ntime.Now\nTimer"},

	map[string]interface{}{"id": uuid.NewV4().String(),
		"title":         "切尔西球迷",
		"category_info": []TaxInfo{{TaxId: testTaxes[1].Id}, {TaxId: testTaxes[4].Id}, {TaxId: testTaxes[6].Id}},
		"content":       "切尔西球迷\nchelsa\nlondon"},
	map[string]interface{}{"id": uuid.NewV4().String(),
		"title":         "巴萨俱乐部",
		"category_info": []TaxInfo{{TaxId: testTaxes[4].Id}, {TaxId: testTaxes[5].Id}},
		"content":       "梅西，马拉多纳，欧冠王，球王球队"},
	map[string]interface{}{"id": uuid.NewV4().String(),
		"title":         "messi",
		"category_info": []TaxInfo{{TaxId: testTaxes[1].Id}, {TaxId: testTaxes[5].Id}},
		"content":       "the best player of all time"},

	map[string]interface{}{"id": uuid.NewV4().String(),
		"title":         "redis用法说明",
		"category_info": []TaxInfo{{TaxId: testTaxes[2].Id}},
		"content":       "redis command list"},
	map[string]interface{}{"id": uuid.NewV4().String(),
		"title":         "如何删除所有键值",
		"category_info": []TaxInfo{{TaxId: testTaxes[2].Id}},
		"content":       "flushall"},

	map[string]interface{}{"id": uuid.NewV4().String(),
		"title":         "佩刀如何？",
		"category_info": []TaxInfo{{TaxId: testTaxes[1].Id}, {TaxName: testTaxes[4].Name, Taxonomy: testTaxes[4].Taxonomy}},
		"content":       "这场比赛水平不行啊\n难道是打法问题？\n水晶宫真是如有神助啊...."},
}

func testOrderBy(t *testing.T) {
	var (
		p1 = Post{Title: "1", Points: 1}
		p2 = Post{Title: "2", Points: 2}
		p3 = Post{Title: "3", Points: 3}
		p4 = Post{Title: "4", Points: 4}
		p5 = Post{Title: "5", Points: 5}
		p6 = Post{Title: "6", Points: 6}

		ps1 = []*Post{&p3, &p4, &p1, &p6, &p5, &p2}
	)

	postSlice(ps1).orderedBy("points", true)
	if ps1[0].Title != "6" || ps1[1].Title != "5" || ps1[2].Title != "4" ||
		ps1[3].Title != "3" || ps1[4].Title != "2" || ps1[5].Title != "1" {
		t.Fail()
	}

	postSlice(ps1).orderedBy("points", false)
	if ps1[0].Title != "1" || ps1[1].Title != "2" || ps1[2].Title != "3" ||
		ps1[3].Title != "4" || ps1[4].Title != "5" || ps1[5].Title != "6" {
		t.Fail()
	}
}

func testPostTaxApi(t *testing.T) {
	testTaxCreate(t)
	testGetAllTax(t)

	// 创建post
	testPostCreate(t)
	// 测试category下的post
	testPostCategory(t)
	// 删除post
	testDeleteTopics(t)

}

// 准备post数据
// 创建post
// 创建12个post，分别属于3个不同的category：6，4，2
func testPostCreate(t *testing.T) {
	var (
		guotie   = users["guotie"]
		tiege    = users["铁哥"]
		tianj    = users["天津人民很伤心"]
		messi    = users["luomessi"]
		somebody = users["somebody"]

		path = "/api/topic/"
	)

	assert(t, guotie != nil, "user must exist")
	assert(t, tiege != nil, "user must exist")
	assert(t, tianj != nil, "user must exist")
	assert(t, messi != nil, "user must exist")
	assert(t, somebody != nil, "user must exist")

	guotie.fetch("POST", path+_testPosts[0]["id"].(string), _testPosts[0])
	guotie.fetch("POST", path+_testPosts[2]["id"].(string), _testPosts[2])
	guotie.fetch("POST", path+_testPosts[4]["id"].(string), _testPosts[4])
	testGetTopics(t, _testPosts[0], guotie)
	testGetTopics(t, _testPosts[2], guotie)
	testGetTopics(t, _testPosts[4], guotie)

	tiege.fetch("POST", path+_testPosts[1]["id"].(string), _testPosts[1])
	tiege.fetch("POST", path+_testPosts[3]["id"].(string), _testPosts[3])
	tiege.fetch("POST", path+_testPosts[5]["id"].(string), _testPosts[5])
	testGetTopics(t, _testPosts[1], tiege)
	testGetTopics(t, _testPosts[3], tiege)
	testGetTopics(t, _testPosts[5], tiege)

	tianj.fetch("POST", path+_testPosts[6]["id"].(string), _testPosts[6])
	tianj.fetch("POST", path+_testPosts[7]["id"].(string), _testPosts[7])
	testGetTopics(t, _testPosts[6], tianj)
	testGetTopics(t, _testPosts[7], tianj)

	messi.fetch("POST", path+_testPosts[8]["id"].(string), _testPosts[8])
	messi.fetch("POST", path+_testPosts[10]["id"].(string), _testPosts[10])
	testGetTopics(t, _testPosts[8], messi)
	testGetTopics(t, _testPosts[10], messi)

	somebody.fetch("POST", path+_testPosts[9]["id"].(string), _testPosts[9])
	testGetTopics(t, _testPosts[9], somebody)

	rr, _ := anonymous.fetch("POST", path+_testPosts[11]["id"].(string), _testPosts[11])
	assert(t, rr.Code != 0, "should need perm")

	somebody.fetch("POST", path+_testPosts[11]["id"].(string), _testPosts[11])
	testGetTopics(t, _testPosts[11], somebody)
}

func testGetTopics(t *testing.T, p map[string]interface{}, ut *userTest) {
	id := p["id"].(string)
	rr, err := anonymous.fetch("GET", "/api/topic/"+id, nil)
	assert(t, err == nil, "response should success")
	assert(t, rr.Code == 0, "result code should be 0")

	post := rr.Data.(map[string]interface{})
	assert(t, p["id"] == post["id"].(string), "id should equal")
	assert(t, p["title"] == post["title"].(string), "title should equal")
	assert(t, post["author_id"].(string) == ut.Id, "author should equal")
	assert(t, post["author_name"].(string) == ut.Name, "author should equal")
}

func testModifyTopic(t *testing.T) {
	// to be complete
}

func testDeleteTopic(t *testing.T, id string, ut *userTest) {
	//
	rr, err := ut.fetch("DELETE", "/api/topic/"+id, nil)
	assert(t, err == nil, "http should be ok")
	assert(t, rr.Code == 0, "resp result should be 0: %d", rr.Code)
}

func testDeleteTopics(t *testing.T) {
	rr, err := anonymous.fetch("DELETE", "/api/topic/"+_testPosts[0]["id"].(string), nil)
	assert(t, err == nil, "http response should be ok")
	assert(t, rr.Code == ErrCodeNeePerm, "anonymous should not delete post")

	guotie := users["guotie"]
	testDeleteTopic(t, _testPosts[0]["id"].(string), guotie)
	testDeleteTopic(t, _testPosts[1]["id"].(string), administrator)
}

func testGetTaxTopics(t *testing.T, cat string, count int) {
	rr, err := anonymous.fetch("GET", "/api/cat/topics?catname="+cat, nil)
	assert(t, err == nil, "response should ok")
	assert(t, rr.Code == 0, "result code not 0: %d", rr.Code)
	posts := rr.Data.([]interface{})
	assert(t, len(posts) == count, "category golang should has 6 posts: %d", len(posts))
}

func testPostCategory(t *testing.T) {
	testGetTaxTopics(t, "golang", 6)
	testGetTaxTopics(t, "足球", 3) // messi不属于足球分类
	testGetTaxTopics(t, "redis", 2)
}

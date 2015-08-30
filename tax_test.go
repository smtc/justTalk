package main

import (
	"fmt"
	"testing"

	"github.com/satori/go.uuid"
)

var testTaxes = []Taxonomy{
	Taxonomy{Id: uuid.NewV4().String(), Name: "golang", Slug: "golang", Taxonomy: "category"},
	Taxonomy{Id: uuid.NewV4().String(), Name: "足球", Slug: "football", Taxonomy: "category"},
	Taxonomy{Id: uuid.NewV4().String(), Name: "redis", Slug: "redis", Taxonomy: "category"},
	Taxonomy{Id: uuid.NewV4().String(), Name: "golang", Slug: "golang", Taxonomy: "tag"},
	Taxonomy{Id: uuid.NewV4().String(), Name: "足球", Slug: "football", Taxonomy: "tag"}, // 4
	Taxonomy{Id: uuid.NewV4().String(), Name: "梅西", Slug: "mesii", Taxonomy: "tag"},    // 5
	Taxonomy{Id: uuid.NewV4().String(), Name: "英超", Slug: "esl", Taxonomy: "tag"},      // 6
}

func init() {

}

// 创建3个不同分类
func testTaxCreate(t *testing.T) {
	for _, tax := range testTaxes {
		//fmt.Println("admin:", administrator)
		rr, err := administrator.fetch("POST", "/api/tax/create/"+tax.Id, &tax)
		//fmt.Println(err, rr)
		assert(t, err == nil, fmt.Sprintf("response should be ok: code=%d error=%v", rr.Code, err))
		assert(t, rr.Code == 0, "response result should be 0, code: %d msg: %s", rr.Code, rr.Message)
	}
}

func testGetAllTax(t *testing.T) {
	rr, err := anonymous.fetch("GET", "/api/tax/list", nil)
	assert(t, err == nil, "response should be ok")
	assert(t, rr.Code == 0, "resp result code should be 0")
	//fmt.Println(rr.Data)
}

package main

import (
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/smtc/glog"
	"io/ioutil"
	//"net/http"
)

// GET /api/tax/list?taxonomy=xxxx
// todo: modify this api
func getCategoryList(c *gin.Context) {
	taxes, err := getAllCategory()
	if err != nil {
		c.JSON(200, RespResult{ErrCodeDBQuery, err.Error(), nil})
		return
	}

	c.JSON(200, RespResult{0, "ok", taxes})
	return
}

// 新建taxonomy
// POST /api/tax/create/:id
/*
{
	"id":,
	"name":,
	"slug":,
	"taxonomy":,
	"description":,
	"parent",
}
*/
func createCategory(c *gin.Context) {
	var tax Taxonomy

	user := getCurrent(c)
	if user == nil {
		c.JSON(200, RespResult{ErrCodeNeeAuthen, "need login to create category", nil})
		return
	}
	id := c.Param("id")

	defer c.Request.Body.Close()
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(200, RespResult{ErrCodeReadRequest, err.Error(), nil})
		return
	}
	if err = json.Unmarshal(body, &tax); err != nil {
		c.JSON(200, RespResult{ErrCodeUnmarshaLTax, err.Error(), nil})
		return
	}
	if id != tax.Id {
		c.JSON(200, RespResult{ErrCodeParam,
			fmt.Sprintf("path param id %s Not equal with json param id %s", id, tax.Id), nil})
		return
	}

	err = createTaxonomy(&tax, user)
	if err != nil {
		c.JSON(200, RespResult{ErrCodeDBQuery, err.Error(), nil})
		return
	}

	c.JSON(200, RespResult{0, "ok", nil})
	return
}

/*
获取特定taxonomy的文章
目前使用join来获取文章

GET /api/cat/topics?catname=xxx&start=xx&count=xx
*/
func getTaxonomyTopics(c *gin.Context) {
	tax := c.Query("catname")
	if tax == "" {
		glog.Warn("no categroy param found!\n")
		getTopics(c)
		return
	}
	start := getQueryIntDefault(c, "start", 0)
	count := getQueryIntDefault(c, "count", 20)

	posts, err := getTopicsByTaxonomy(tax, start, count)
	if err != nil {
		c.JSON(200, RespResult{ErrCodeDBQuery, err.Error(), nil})
		return
	}
	c.JSON(200, RespResult{0, "ok", posts})
}

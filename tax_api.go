package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	//"github.com/smtc/glog"
	"io/ioutil"
	//"net/http"
)

// GET /api/tax/category
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
// POST /api/tax/category/:id
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

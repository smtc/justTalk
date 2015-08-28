package main

import (
	"testing"
)

func TestOrderBy(t *testing.T) {
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

// 准备post数据
func preparePostData() {

}

func cleanPostData() {

}

func TestPostTaxApi(t *testing.T) {
	testTaxCreate(t)
	testPostCreate(t)
}

// 创建post
// 创建10个post，分别属于3个不同的category：5，3，2
func testPostCreate(t *testing.T) {

}
func testGetTopics(t *testing.T) {

}

func testGetTheTopic(t *testing.T) {

}

func testModifyTopic(t *testing.T) {

}
func testDeleteTopic(t *testing.T) {

}

func testGetTaxTopics(t *testing.T) {

}

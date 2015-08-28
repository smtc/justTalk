package main

import (
	"testing"
)

// 创建3个不同分类
func testTaxCreate(t *testing.T) {

}

//
func makePostRequest(path string, ctoken string) (*RespResult, error) {
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

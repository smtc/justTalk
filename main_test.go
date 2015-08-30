package main

import (
	"os"
	"testing"

	"github.com/guotie/config"
	"github.com/guotie/deferinit"
)

func TestMain(m *testing.M) {
	config.ReadCfg("./config.json")

	deferinit.InitAll()
	go runServer()

	cleanDatabase()

	prepareUserData()
	code := m.Run()

	//cleanDatabase()

	// 清理redis，生成环境下慎用！！！
	//FlushAll()
	os.Exit(code)
}

func TestApi(t *testing.T) {
	testUserApi(t)
	testPostTaxApi(t)
}

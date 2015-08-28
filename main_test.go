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

	cleanUserData()
	prepareUserData()
	code := m.Run()

	cleanUserData()

	os.Exit(code)
}

package main

import (
	"fmt"
	"strings"

	"github.com/guotie/config"
	"github.com/guotie/deferinit"
	"github.com/smtc/glog"
)

var (
	port     int
	domain   string
	hostname string
)

func setGlobalVars() {
	domain = strings.TrimSpace(config.GetStringDefault("domain", "http://127.0.0.1"))
	if domain[len(domain)-1] == '/' {
		domain = domain[0 : len(domain)-1]
	}

	port = config.GetIntDefault("port", 80)
	if port != 80 {
		hostname = fmt.Sprintf("%s:%d", domain, port)
	} else {
		hostname = domain
	}
}

func runServer() {
	r := router()

	if err := r.Run(fmt.Sprintf(":%d", config.GetIntDefault("port", 8000))); err != nil {
		glog.Fatal("run failed:", err)
	}
}

func main() {
	config.ReadCfg("./config.json")
	setGlobalVars()
	deferinit.InitAll()
	runServer()
}

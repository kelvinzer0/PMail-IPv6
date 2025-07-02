package main

import (
	"embed" // Import the embed package
	"github.com/Jinnrry/pmail/config"
	"github.com/Jinnrry/pmail/listen/cron_server"
	"github.com/Jinnrry/pmail/res_init"
	log "github.com/sirupsen/logrus"
)

//go:embed fe/dist
var embeddedFS embed.FS

var (
	gitHash   string
	buildTime string
	goVersion string
	version   string
)

func main() {

	config.Init()

	if version == "" {
		version = "TestVersion"
	}

	log.Infoln("*******************************************************************")
	log.Infof("***\tServer Start Success \n")
	log.Infof("***\tServer Version: %s \n", version)
	log.Infof("***\tGit Commit Hash: %s ", gitHash)
	log.Infof("***\tBuild Date: %s ", buildTime)
	log.Infof("***\tBuild GoLang Version: %s ", goVersion)
	log.Infoln("*******************************************************************")

	// 定时任务启动
	go cron_server.Start()

	// 核心服务启动
	res_init.Init(version, embeddedFS)

	log.Warnf("Server Stoped \n")

}


package http_server

import (
	"errors"
	"fmt"
	"github.com/Jinnrry/pmail/config"
	"github.com/Jinnrry/pmail/controllers"
	"github.com/Jinnrry/pmail/controllers/email"
	"github.com/Jinnrry/pmail/session"
	"github.com/Jinnrry/pmail/utils/ip"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

var httpServer *http.Server

func HttpStop() {
	if httpServer != nil {
		httpServer.Close()
	}
}

func router(mux *http.ServeMux) {
	mux.HandleFunc("/.well-known/", controllers.AcmeChallenge)
	mux.HandleFunc("/api/ping", controllers.Ping)
	mux.HandleFunc("/api/login", contextIterceptor(controllers.Login))
	mux.HandleFunc("/api/logout", contextIterceptor(controllers.Logout))
	mux.HandleFunc("/api/group", contextIterceptor(controllers.GetUserGroup))
	mux.HandleFunc("/api/group/list", contextIterceptor(controllers.GetUserGroupList))
	mux.HandleFunc("/api/group/add", contextIterceptor(controllers.AddGroup))
	mux.HandleFunc("/api/group/del", contextIterceptor(controllers.DelGroup))
	mux.HandleFunc("/api/email/list", contextIterceptor(email.EmailList))
	mux.HandleFunc("/api/email/del", contextIterceptor(email.EmailDelete))
	mux.HandleFunc("/api/email/read", contextIterceptor(email.MarkRead))
	mux.HandleFunc("/api/email/detail", contextIterceptor(email.EmailDetail))
	mux.HandleFunc("/api/email/move", contextIterceptor(email.Move))
	mux.HandleFunc("/api/email/send", contextIterceptor(email.Send))
	mux.HandleFunc("/api/settings/modify_password", contextIterceptor(controllers.ModifyPassword))
	mux.HandleFunc("/api/rule/get", contextIterceptor(controllers.GetRule))
	mux.HandleFunc("/api/rule/add", contextIterceptor(controllers.UpsertRule))
	mux.HandleFunc("/api/rule/update", contextIterceptor(controllers.UpsertRule))
	mux.HandleFunc("/api/rule/del", contextIterceptor(controllers.DelRule))
	mux.HandleFunc("/attachments/", contextIterceptor(controllers.GetAttachments))
	mux.HandleFunc("/attachments/download/", contextIterceptor(controllers.Download))
	mux.HandleFunc("/api/user/create", contextIterceptor(controllers.CreateUser))
	mux.HandleFunc("/api/user/edit", contextIterceptor(controllers.EditUser))
	mux.HandleFunc("/api/user/info", contextIterceptor(controllers.Info))
	mux.HandleFunc("/api/user/list", contextIterceptor(controllers.UserList))
	mux.HandleFunc("/api/plugin/settings/", contextIterceptor(controllers.SettingsHtml))
	mux.HandleFunc("/api/plugin/list", contextIterceptor(controllers.GetPluginList))
	mux.HandleFunc("/api/config", controllers.GetAppConfigHttp)
}

func HttpStart() {
	mux := http.NewServeMux()
	router(mux)

	HttpPort := 80
	if config.Instance.HttpPort > 0 {
		HttpPort = config.Instance.HttpPort
	}

	bindingHost := ip.GetIp()
	if bindingHost == "" {
		bindingHost = "0.0.0.0"
	}

	addr := fmt.Sprintf("%s:%d", bindingHost, HttpPort)

	if config.Instance.HttpsEnabled != 2 {
		log.Infof("HttpServer Start On %s", addr)
		httpServer = &http.Server{
			Addr:         addr,
			Handler:      mux,
			ReadTimeout:  time.Second * 90,
			WriteTimeout: time.Second * 90,
		}
	} else {
		log.Infof("HttpServer Start On %s", addr)
		httpServer = &http.Server{
			Addr:         addr,
			Handler:      session.Instance.LoadAndSave(mux),
			ReadTimeout:  time.Second * 90,
			WriteTimeout: time.Second * 90,
		}
	}

	err := httpServer.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		panic(err)
	}
}

package http_server

import (
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Jinnrry/pmail/config"
	"github.com/Jinnrry/pmail/controllers"
	"github.com/Jinnrry/pmail/controllers/email"
	"github.com/Jinnrry/pmail/dto/response"
	"github.com/Jinnrry/pmail/i18n"
	"github.com/Jinnrry/pmail/models"
	"github.com/Jinnrry/pmail/session"
	"github.com/Jinnrry/pmail/utils/context"
	"github.com/Jinnrry/pmail/utils/id"
	"github.com/Jinnrry/pmail/utils/ip"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"io/fs"
	"net/http"
	"path"
	"strings"
	"time"
)

var httpServer *http.Server
var httpsServer *http.Server

func HttpStop() {
	if httpServer != nil {
		httpServer.Close()
	}
}

func HttpsStop() {
	if httpsServer != nil {
		httpsServer.Close()
	}
}

// 注入context
func contextIterceptor(h controllers.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if w.Header().Get("Content-Type") == "" {
			w.Header().Set("Content-Type", "application/json")
		}

		ctx := &context.Context{}
		ctx.Context = r.Context()
		ctx.SetValue(context.LogID, id.GenLogID())
		lang := r.Header.Get("Lang")
		if lang == "" {
			lang = "en"
		}
		ctx.Lang = lang

		if config.IsInit {
			user := cast.ToString(session.Instance.Get(ctx, "user"))
			var userInfo *models.User
			if user != "" {
				_ = json.Unmarshal([]byte(user), &userInfo)
			}
			if userInfo != nil && userInfo.ID > 0 {
				ctx.UserID = userInfo.ID
				ctx.UserName = userInfo.Name
				ctx.UserAccount = userInfo.Account
				ctx.IsAdmin = userInfo.IsAdmin == 1
			}

			if ctx.UserID == 0 {
				if r.URL.Path != "/api/ping" && r.URL.Path != "/api/login" {
					response.NewErrorResponse(response.NeedLogin, i18n.GetText(ctx.Lang, "login_exp"), "").FPrint(w)
					return
				}
			}
		} else if r.URL.Path != "/api/setup" {
			response.NewErrorResponse(response.NeedSetup, "", "").FPrint(w)
			return
		}
		h(ctx, w, r)
	}
}

func router(mux *http.ServeMux, embeddedFS embed.FS) {
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

	// Create a sub-filesystem rooted at "fe/dist"
	distFS, err := fs.Sub(embeddedFS, "fe/dist")
	if err != nil {
		log.Fatalf("failed to create sub file system for frontend: %v", err)
	}

	fileServer := http.FileServer(http.FS(distFS))

	// Handle all non-api requests with this handler
	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fsPath := strings.TrimPrefix(r.URL.Path, "/")

		// If the path is empty, it's the root.
		if fsPath == "" {
			fsPath = "index.html"
		}

		f, err := distFS.Open(fsPath)
		if err != nil {
			// File doesn't exist.
			// If it has an extension, it's a missing asset -> 404
			// If it has no extension, it's a route -> serve index.html
			if path.Ext(r.URL.Path) != "" {
				http.NotFound(w, r)
				return
			}

			// Serve index.html
			r.URL.Path = "/index.html" // rewrite for fileServer
			fileServer.ServeHTTP(w, r)
			return
		}
		f.Close()

		// File exists, serve it.
		fileServer.ServeHTTP(w, r)
	}))
}

func HttpStart(embeddedFS embed.FS) {
	mux := http.NewServeMux()
	router(mux, embeddedFS)

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

func HttpsStart(embeddedFS embed.FS) {
	mux := http.NewServeMux()
	router(mux, embeddedFS)

	HttpsPort := 443
	if config.Instance.HttpsPort > 0 {
		HttpsPort = config.Instance.HttpsPort
	}

	bindingHost := ip.GetIp()
	if bindingHost == "" {
		bindingHost = "0.0.0.0"
	}

	addr := fmt.Sprintf("%s:%d", bindingHost, HttpsPort)

	log.Infof("HttpsServer Start On %s", addr)
	httpsServer = &http.Server{
		Addr:         addr,
		Handler:      session.Instance.LoadAndSave(mux),
		ReadTimeout:  time.Second * 90,
		WriteTimeout: time.Second * 90,
	}

	err := httpsServer.ListenAndServeTLS(config.Instance.SslCert, config.Instance.SslKey)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		panic(err)
	}
}

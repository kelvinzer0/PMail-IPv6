package imap_server

import (
	"crypto/tls"
	"fmt"
	"github.com/Jinnrry/pmail/config"
	"github.com/Jinnrry/pmail/utils/ip"
	"github.com/emersion/go-imap/v2"
	"github.com/emersion/go-imap/v2/imapserver"
	log "github.com/sirupsen/logrus"
	"os"
)

var instanceTLS *imapserver.Server

func Stop() {
	if instanceTLS != nil {
		instanceTLS.Close()
		instanceTLS = nil
	}
}

// StarTLS 启动TLS端口监听，不加密的代码就懒得写了
func StarTLS() {

	crt, err := tls.LoadX509KeyPair(config.Instance.SSLPublicKeyPath, config.Instance.SSLPrivateKeyPath)
	if err != nil {
		panic(err)
	}
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{crt},
	}

	memServer := NewServer()

	option := &imapserver.Options{
		NewSession: func(conn *imapserver.Conn) (imapserver.Session, *imapserver.GreetingData, error) {
			return memServer.NewSession(), nil, nil
		},
		Caps: imap.CapSet{
			imap.CapIMAP4rev1: {},
		},
		TLSConfig:    tlsConfig,
		InsecureAuth: false,
	}

	if config.Instance.LogLevel == "debug" {
		option.DebugWriter = os.Stdout
	}

	instanceTLS = imapserver.New(option)

	bindingHost := ip.GetIp()
	if bindingHost == "" {
		bindingHost = "0.0.0.0"
	}

	addr := fmt.Sprintf("%s:%d", bindingHost, 993)

	log.Infof("IMAP With TLS Server Start On %s", addr)
	if err := instanceTLS.ListenAndServeTLS(addr); err != nil {
		panic(err)
	}
}

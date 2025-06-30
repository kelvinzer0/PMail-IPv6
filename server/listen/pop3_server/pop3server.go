package pop3_server

import (
	"crypto/rand"
	"crypto/tls"
	"github.com/Jinnrry/gopop"
	"github.com/Jinnrry/pmail/config"
	log "github.com/sirupsen/logrus"
	"time"
)

var instance *gopop.Server
var instanceTls *gopop.Server

func StartWithTls() {
	crt, err := tls.LoadX509KeyPair(config.Instance.SSLPublicKeyPath, config.Instance.SSLPrivateKeyPath)
	if err != nil {
		panic(err)
	}
	tlsConfig := &tls.Config{}
	tlsConfig.Certificates = []tls.Certificate{crt}
	tlsConfig.Time = time.Now
	tlsConfig.Rand = rand.Reader

	bindingHost := config.Instance.BindingHost
	if bindingHost == "" {
		bindingHost = "0.0.0.0"
	}

	instanceTls = gopop.NewPop3Server(995, bindingHost, true, tlsConfig, action{})
	instanceTls.ConnectAliveTime = 5 * time.Minute

	log.Infof("POP3 With TLS Server Start On %s:995", bindingHost)

	err = instanceTls.Start()
	if err != nil {
		panic(err)
	}
}

func Start() {
	crt, err := tls.LoadX509KeyPair(config.Instance.SSLPublicKeyPath, config.Instance.SSLPrivateKeyPath)
	if err != nil {
		panic(err)
	}
	tlsConfig := &tls.Config{}
	tlsConfig.Certificates = []tls.Certificate{crt}
	tlsConfig.Time = time.Now
	tlsConfig.Rand = rand.Reader

	bindingHost := config.Instance.BindingHost
	if bindingHost == "" {
		bindingHost = "0.0.0.0"
	}

	instance = gopop.NewPop3Server(110, bindingHost, false, tlsConfig, action{})
	instance.ConnectAliveTime = 5 * time.Minute
	log.Infof("POP3 Server Start On %s:110", bindingHost)

	err = instance.Start()
	if err != nil {
		panic(err)
	}
}

func Stop() {
	if instance != nil {
		instance.Stop()
	}

	if instanceTls != nil {
		instanceTls.Stop()
	}
}

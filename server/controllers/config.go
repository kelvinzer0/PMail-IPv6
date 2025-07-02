package controllers

import (
	"encoding/json"
	"github.com/Jinnrry/pmail/config"
	"net/http"
)

type AppConfig struct {
	PusherBeamsInstanceId string `json:"pusherBeamsInstanceId"`
}

// GetAppConfigHttp returns public application configuration, including Pusher Beams Instance ID.
func GetAppConfigHttp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	appConfig := AppConfig{
		PusherBeamsInstanceId: config.Instance.PusherBeamsInstanceId,
	}
	json.NewEncoder(w).Encode(appConfig)
}

package controllers

import (
	"github.com/Jinnrry/pmail/config"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AppConfig struct {
	PusherBeamsInstanceId string `json:"pusherBeamsInstanceId"`
}

// GetAppConfig returns public application configuration, including Pusher Beams Instance ID.
func GetAppConfig(c *gin.Context) {
	appConfig := AppConfig{
		PusherBeamsInstanceId: config.Instance.PusherBeamsInstanceId,
	}
	c.JSON(http.StatusOK, appConfig)
}

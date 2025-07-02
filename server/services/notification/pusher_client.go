package notification

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Jinnrry/pmail/config"
	"github.com/pusher/pusher-http-go/v5"
	log "github.com/sirupsen/logrus"
	"net/http"
)

var client *pusher.Client

func InitPusher() {
	client = &pusher.Client{
		AppID:   config.Instance.PusherAppID,
		Key:     config.Instance.PusherKey,
		Secret:  config.Instance.PusherSecret,
		Cluster: config.Instance.PusherCluster,
		Secure:  true,
	}
	log.Infoln("Pusher client initialized.")
}

func SendNewEmailNotification(userID int, subject, sender, emailID string) {
	channelName := fmt.Sprintf("private-user-%d", userID)
	eventName := "new-email"
	data := map[string]string{
		"subject": subject,
		"sender":  sender,
		"id":      emailID,
	}

	err := client.Trigger(channelName, eventName, data)
	if err != nil {
		log.Errorf("Error triggering Pusher event for user %d: %v", userID, err)
	}
	log.Debugf("Pusher event '%s' triggered for channel '%s' with data: %+v", eventName, channelName, data)

	// Send web push notification
	interest := fmt.Sprintf("user-%d", userID)
	title := "New Email from " + sender
	body := subject
	SendWebPushNotification(interest, title, body)
}

func SendWebPushNotification(interest string, title, body string) {
	if config.Instance.PusherBeamsInstanceId == "" || config.Instance.PusherBeamsSecretKey == "" {
		log.Warnln("Pusher Beams not configured. Skipping web push notification.")
		return
	}

	url := fmt.Sprintf("https://%s.pushnotifications.pusher.com/publish_api/v1/instances/%s/publishes", config.Instance.PusherBeamsInstanceId, config.Instance.PusherBeamsInstanceId)

	payload := map[string]interface{}{
		"interests": []string{interest},
		"web": map[string]interface{}{
			"notification": map[string]string{
				"title": title,
				"body":  body,
			},
		},
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		log.Errorf("Error marshalling web push payload: %v", err)
		return
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		log.Errorf("Error creating web push request: %v", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+config.Instance.PusherBeamsSecretKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Errorf("Error sending web push notification: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Errorf("Failed to send web push notification. Status: %s", resp.Status)
	} else {
		log.Debugf("Web push notification sent successfully to interest '%s'", interest)
	}
}

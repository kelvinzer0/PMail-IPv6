# Pusher Beams Integration

This document outlines the steps to integrate Pusher Beams for web push notifications into the PMail-IPv6 project.

## Overview

Pusher Beams is a hosted API that allows you to send push notifications to web browsers. This integration enables real-time notifications for events such as new email arrivals.

## Backend Configuration (Go)

To enable Pusher Beams in the backend, you need to add your Beams Instance ID and Secret Key to the `config.json` (or `config.dev.json` for development) file.

1.  **Update `server/config/config.go`**:
    Ensure the `Config` struct includes the following fields:

    ```go
    type Config struct {
        // ... existing fields ...
        PusherBeamsInstanceId string `json:"pusherBeamsInstanceId"`
        PusherBeamsSecretKey string `json:"pusherBeamsSecretKey"`
        // ... other fields ...
    }
    ```

2.  **Update `config.json` (or `config.dev.json`)**:
    Add your Pusher Beams Instance ID and Secret Key to your configuration file:

    ```json
    {
      "bindingHost": "0.0.0.0",
      // ... other configurations ...
      "pusherBeamsInstanceId": "YOUR_PUSHER_BEAMS_INSTANCE_ID",
      "pusherBeamsSecretKey": "YOUR_PUSHER_BEAMS_SECRET_KEY"
    }
    ```
    **Replace `YOUR_PUSHER_BEAMS_INSTANCE_ID` and `YOUR_PUSHER_BEAMS_SECRET_KEY` with your actual credentials.**

3.  **Backend Code Changes (`server/services/notification/pusher_client.go`)**:
    A new function `SendWebPushNotification` has been added to handle sending notifications via Pusher Beams. This function makes an HTTP POST request to the Pusher Beams Publish API.

    ```go
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
    ```

## Frontend Integration (Vue.js)

The frontend needs to register a service worker, load the Pusher Beams SDK, and subscribe to user-specific interests.

1.  **Create `fe/public/service-worker.js`**:
    This file is essential for web push notifications and must be served from the root of your web application.

    ```javascript
    importScripts("https://js.pusher.com/beams/service-worker.js");
    ```

2.  **Modify `fe/src/main.js`**:
    The Pusher Beams SDK is now dynamically loaded, and the client is initialized after the user logs in, subscribing to a unique interest based on the user's ID.

    ```javascript
    import './assets/main.css'

    import { createApp, watch } from 'vue'
    import { createPinia } from 'pinia'
    import App from './App.vue'
    import {router} from './router'

    import ElementPlus from 'element-plus'
    import 'element-plus/dist/index.css'

    import { lang } from './i18n/i18n'
    import { useGlobalStatusStore } from '@/stores/useGlobalStatusStore'

    const app = createApp(App)
    app.use(router)
    app.use(createPinia())
    app.use(ElementPlus)
    app.config.globalProperties.$lang = lang
    app.mount('#app')

    // Dynamically load Pusher Beams SDK
    const script = document.createElement('script');
    script.src = "https://js.pusher.com/beams/2.1.0/push-notifications-cdn.js";
    script.onload = () => {
      const globalStatus = useGlobalStatusStore();
      watch(() => globalStatus.userInfos.id, (newId) => {
        if (newId) {
          const beamsClient = new PusherPushNotifications.Client({
            instanceId: 'YOUR_PUSHER_BEAMS_INSTANCE_ID', // Replace with your actual instance ID
          });

          beamsClient.start()
            .then(() => beamsClient.addDeviceInterest(`user-${newId}`))
            .then(() => console.log('Successfully registered and subscribed to user interest!'))
            .catch(console.error);
        }
      }, { immediate: true });
    };
    document.head.appendChild(script);
    ```
    **Replace `YOUR_PUSHER_BEAMS_INSTANCE_ID` with your actual Instance ID.**

## Important Notes

*   **Serving `service-worker.js`**: Ensure your web server is configured to serve the `fe/public/service-worker.js` file directly from the root of your domain (e.g., `http://yourdomain.com/service-worker.js`).
*   **Instance ID Consistency**: The `instanceId` used in the frontend (`fe/src/main.js`) must match the `pusherBeamsInstanceId` configured in your backend's `config.json`.
*   **Secret Key Security**: The `pusherBeamsSecretKey` must be kept confidential and only used on your backend server. Never expose it in your frontend code.
*   **Notification Permissions**: Users will be prompted to grant notification permissions in their browser. The push notifications will only be delivered if permission is granted.

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

// Fetch config from backend and then dynamically load Pusher Beams SDK
fetch('/api/config')
  .then(response => response.json())
  .then(data => {
    const pusherBeamsInstanceId = data.pusherBeamsInstanceId;

    const script = document.createElement('script');
    script.src = "https://js.pusher.com/beams/2.1.0/push-notifications-cdn.js";
    script.onload = () => {
      const globalStatus = useGlobalStatusStore();
      watch(() => globalStatus.userInfos.id, (newId) => {
        if (newId) {
          const beamsClient = new PusherPushNotifications.Client({
            instanceId: pusherBeamsInstanceId,
          });

          beamsClient.start()
            .then(() => beamsClient.addDeviceInterest(`user-${newId}`))
            .then(() => console.log('Successfully registered and subscribed to user interest!'))
            .catch(console.error);
        }
      }, { immediate: true });
    };
    document.head.appendChild(script);
  })
  .catch(error => {
    console.error('Error fetching app config:', error);
  });

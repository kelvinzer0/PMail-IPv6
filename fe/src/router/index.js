import { createRouter, createWebHashHistory } from 'vue-router'
import ListView from '../views/ListView.vue'
import EditerView from '../views/EditerView.vue'
import LoginView from '../views/LoginView.vue'
import EmailDetailView from '../views/EmailDetailView.vue'
import SetupView from '../views/SetupView.vue'


const router = createRouter({
  history: createWebHashHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: ListView
    },
    {
      path: '/list',
      name: 'list',
      component: ListView
    },
    {
      path: '/editer',
      name: "editer",
      component: EditerView
    },
    {
      path: '/login',
      name: "login",
      component: LoginView
    },
    {
      path: '/setup',
      name: "setup",
      component: SetupView
    },
    {
      path: '/detail/:id',
      name: "detail",
      component: EmailDetailView
    }
  ]
})





import { useGlobalStatusStore } from "@/stores/useGlobalStatusStore";

router.beforeEach(async (to, from, next) => {
  const globalStatus = useGlobalStatusStore();
  console.log(`Navigating from ${from.name || from.path} to ${to.name || to.path}`);
  console.log("UserInfos before guard check:", globalStatus.userInfos);

  if (Object.keys(globalStatus.userInfos).length === 0 && to.name !== 'login' && to.name !== 'setup') {
    console.log("User not logged in, attempting to initialize user info...");
    await globalStatus.init(() => {});
    console.log("UserInfos after init attempt:", globalStatus.userInfos);

    if (Object.keys(globalStatus.userInfos).length === 0) {
      console.log("User info still empty after init, redirecting to login.");
      next({ name: 'login' });
    } else {
      console.log("User info successfully initialized, proceeding.");
      next();
    }
  } else {
    console.log("User already logged in or navigating to login/setup, proceeding.");
    next();
  }
});

export {router};

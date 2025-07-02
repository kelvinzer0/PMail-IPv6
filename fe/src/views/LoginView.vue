<template>
  <div id="main">
    <div class="login-container">
      <h1 class="login-title">Login</h1>
      <form @submit.prevent="onSubmit" class="login-form">
        <div class="form-group">
          <label for="account" class="form-label">{{ lang.account }}</label>
          <input
            type="text"
            id="account"
            v-model="form.account"
            placeholder="User Name"
            class="form-input"
            aria-required="true"
          />
        </div>
        <div class="form-group">
          <label for="password" class="form-label">{{ lang.password }}</label>
          <input
            type="password"
            id="password"
            v-model="form.password"
            placeholder="Password"
            class="form-input"
            aria-required="true"
          />
        </div>
        <button type="submit" class="login-button">{{ lang.login }}</button>
      </form>
      <div class="contact-info" style="margin-top: 20px; font-size: 14px; color: #5f6368;">
        <p v-html="lang.login_contact_info"></p>
      </div>
    </div>
  </div>
</template>

<script setup>

import {reactive} from 'vue'
import {ElMessage} from 'element-plus'
import {router} from "@/router"; //根路由对象
import { getCurrentInstance } from 'vue';
import {http} from "@/utils/axios";
import {useGlobalStatusStore} from "@/stores/useGlobalStatusStore";

const globalStatus = useGlobalStatusStore();
const { proxy } = getCurrentInstance();
const lang = proxy.$lang;

const form = reactive({
  account: '',
  password: '',
})

const onSubmit = () => {
  http.post("/api/login", form).then(res => {
    if (res.errorNo !== 0) {
      ElMessage.error(res.errorMsg)
      console.error("Login failed:", res.errorMsg);
    } else {
      Object.assign(globalStatus.userInfos , res.data) 
      console.log("Login successful, userInfos:", globalStatus.userInfos); 
      router.replace({
        path: '/',
        query: {
          redirect: router.currentRoute.fullPath
        }
      })
    }
  })

}
</script>


<style scoped>
#main {
  width: 100%;
  height: 100vh;
  background-color: #f8f9fa; /* Very light background */
  display: flex;
  justify-content: center;
  align-items: center;
  font-family: 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif; /* Modern font */
}

.login-container {
  background-color: #ffffff;
  padding: 48px;
  border-radius: 0; /* No border-radius */
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.12), 0 1px 2px rgba(0, 0, 0, 0.24); /* Subtle shadow */
  width: 100%;
  max-width: 450px; /* Slightly wider for a more spacious feel */
  text-align: center;
  box-sizing: border-box;
  border: 1px solid #dadce0; /* Light border */
}

.login-title {
  font-size: 24px;
  color: #202124; /* Darker text for better contrast */
  margin-bottom: 32px;
  font-weight: 400; /* Lighter font weight */
}

.login-form {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.form-group {
  text-align: left;
}

.form-label {
  display: block;
  margin-bottom: 8px;
  color: #5f6368; /* Google-like grey */
  font-size: 14px;
  font-weight: 500;
}

.form-input {
  width: 100%;
  padding: 14px 16px;
  border: 1px solid #dadce0;
  border-radius: 0; /* No border-radius */
  font-size: 16px;
  color: #202124;
  box-sizing: border-box;
  transition: border-color 0.2s ease, box-shadow 0.2s ease;
}

.form-input:focus {
  border-color: #004e8c; /* Darker Microsoft blue */
  outline: none;
  box-shadow: 0 0 0 2px rgba(0, 78, 140, 0.5); /* Darker Microsoft-style glow */
}

.login-button {
  background-color: #004e8c; /* Darker Microsoft blue */
  color: #ffffff;
  padding: 12px 24px;
  border: none;
  border-radius: 0; /* No border-radius */
  font-size: 16px;
  font-weight: 500;
  cursor: pointer;
  transition: background-color 0.3s ease, box-shadow 0.3s ease;
  width: 100%;
}

.login-button:hover {
  background-color: #003d6b; /* Even darker Microsoft blue for hover */
  box-shadow: 0 0 0 2px rgba(0, 78, 140, 0.5); /* Subtle glow on hover */
}

.login-button:focus {
  outline: none;
  box-shadow: 0 0 0 3px rgba(0, 78, 140, 0.7); /* More pronounced glow on focus */
}

/* Responsive adjustments */
@media (max-width: 600px) {
  .login-container {
    margin: 20px;
    padding: 32px;
  }

  .login-title {
    font-size: 22px;
  }

  .form-input, .login-button {
    padding: 12px;
    font-size: 15px;
  }
}
</style>
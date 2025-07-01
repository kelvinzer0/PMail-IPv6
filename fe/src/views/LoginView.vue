<template>
  <div id="main">
    <div id="form">
      <img src="https://codeberg.org/zcdns/mediakit/raw/branch/main/brand/logo/zcdns-dark-logo.svg" alt="Logo" class="logo">
      <el-form :model="form" label-width="120px" @keyup.enter="onSubmit">
        <el-form-item :label="lang.account">
          <el-input v-model="form.account" placeholder="User Name"/>
        </el-form-item>
        <el-form-item :label="lang.password">
          <el-input v-model="form.password" placeholder="Password" type="password"/>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="onSubmit">{{ lang.login }}</el-button>
        </el-form-item>
      </el-form>

    </div>
  </div>
</template>

<script setup>

import {reactive} from 'vue'
import {ElMessage} from 'element-plus'
import {router} from "@/router"; //根路由对象
import lang from '../i18n/i18n';
import {http} from "@/utils/axios";
import {useGlobalStatusStore} from "@/stores/useGlobalStatusStore";

const globalStatus = useGlobalStatusStore();
// eslint-disable-next-line no-unused-vars

const form = reactive({
  account: '',
  password: '',
})

const onSubmit = () => {
  http.post("/api/login", form).then(res => {
    if (res.errorNo !== 0) {
      ElMessage.error(res.errorMsg)
    } else {
      Object.assign(globalStatus.userInfos , res.data) 
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
  height: 100%;
  background-color: #f1f1f1;
  display: flex;
  justify-content: center;
  /* 水平居中 */
  align-items: center;
  /* 垂直居中 */
}

#form {
  background-color: #fff;
  padding: 40px;
  border-radius: 8px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
  display: flex;
  flex-direction: column;
  align-items: center;
}

.logo {
  width: 150px; /* Adjust as needed */
  margin-bottom: 20px;
}
</style>
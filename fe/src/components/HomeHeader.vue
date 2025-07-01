<template>
  <div id="header_main">
    <div id="logo">
      <router-link to="/">
        <img src="/zcdns-dark-logo.svg" alt="PMail Logo" class="logo-img"/>
      </router-link>
    </div>
    <div id="search-bar">
      <input type="text" placeholder="Search mail"/>
    </div>
    <div id="settings" @click="settings" v-if="isLogin">
      <el-icon style="font-size: 25px;">
        <TbSettings style="color:#FFFFFF"/>
      </el-icon>
    </div>
    <el-drawer v-model="openSettings" size="80%" :title="lang.settings">
      <el-tabs tab-position="left">
        <el-tab-pane :label="lang.security">
          <SecuritySettings/>
        </el-tab-pane>

        <el-tab-pane :label="lang.group_settings">
          <GroupSettings/>
        </el-tab-pane>

        <el-tab-pane :label="lang.rule_setting">
          <RuleSettings/>
        </el-tab-pane>

        <el-tab-pane v-if="userInfos.is_admin" :label="lang.user_management">
          <UserManagement/>
        </el-tab-pane>

        <el-tab-pane :label="lang.plugin_settings">
          <PluginSettings/>
        </el-tab-pane>

      </el-tabs>
    </el-drawer>

  </div>
</template>

<script setup>
import {TbSettings} from "vue-icons-plus/tb";
import {ref} from 'vue'
import SecuritySettings from '@/components/SecuritySettings.vue'
import lang from '../i18n/i18n';
import GroupSettings from './GroupSettings.vue';
import RuleSettings from './RuleSettings.vue';
import UserManagement from './UserManagement.vue';
import PluginSettings from './PluginSettings.vue';
import {useGlobalStatusStore} from "@/stores/useGlobalStatusStore";

const globalStatus = useGlobalStatusStore();
const isLogin = globalStatus.isLogin;
const userInfos = globalStatus.userInfos;


const openSettings = ref(false)
const settings = function () {
  if (Object.keys(userInfos).length === 0) {
    globalStatus.init(()=>{
      Object.assign(userInfos,globalStatus.userInfos)
      openSettings.value = true;
    })
  } else {
    openSettings.value = true;
  }


}

</script>


<style scoped>

#header_main {
  height: 60px;
  background-color: #f2f2f2;
  display: flex;
  align-items: center;
  padding: 0 20px;
  border-bottom: 1px solid #e0e0e0;
}

#logo {
  height: 40px;
  display: flex;
  align-items: center;
}

.logo-img {
  height: 40px;
  width: auto;
}

#search-bar {
  flex-grow: 1;
  margin: 0 20px;
}

#search-bar input {
  width: 100%;
  padding: 10px 15px;
  border: 1px solid #ccc;
  border-radius: 20px;
  font-size: 16px;
  outline: none;
}

#search-bar input:focus {
  border-color: #4285f4;
  box-shadow: 0 0 0 2px rgba(66, 133, 244, 0.2);
}

#settings {
  display: flex;
  justify-content: center;
  align-items: center;
  padding-right: 0;
}

#settings .el-icon {
  font-size: 24px;
  color: #5f6368;
}
</style>
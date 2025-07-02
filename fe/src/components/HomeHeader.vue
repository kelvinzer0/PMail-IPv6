<template>
  <div id="header_main">
    <div id="logo">
      <router-link to="/" style="text-decoration: none; display: flex; align-items: center; height: 100%;">
        <img src="../assets/logo.svg" alt="ZCMail Logo" style="height: 40px;">
      </router-link>
    </div>
    <div class="language-selector">
      <el-dropdown @command="handleLanguageChange">
        <span class="el-dropdown-link">
          <el-icon style="font-size: 25px;">
            <TbLanguage style="color:#FFFFFF"/>
          </el-icon>
        </span>
        <template #dropdown>
          <el-dropdown-menu>
            <el-dropdown-item command="en">English</el-dropdown-item>
            <el-dropdown-item command="zh">简体中文</el-dropdown-item>
            <el-dropdown-item command="id">Bahasa Indonesia</el-dropdown-item>
          </el-dropdown-menu>
        </template>
      </el-dropdown>
    </div>
    <div id="menu-toggle" @click="toggleSidebar">
      <el-icon style="font-size: 25px;">
        <TbMenu2 style="color:#FFFFFF"/>
      </el-icon>
    </div>
    <div id="settings" @click="settings" v-if="isLogin">
      <el-icon style="font-size: 25px;">
        <TbSettings style="color:#FFFFFF"/>
      </el-icon>
    </div>
    <el-drawer v-model="openSettings" :size="drawerSize" :title="lang.settings">
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
import {TbSettings, TbMenu2, TbLanguage} from "vue-icons-plus/tb";
import {ref, onMounted, onUnmounted} from 'vue'
import { storeToRefs } from 'pinia';
import SecuritySettings from '@/components/SecuritySettings.vue'
import { lang, setLanguage } from '../i18n/i18n';
import GroupSettings from './GroupSettings.vue';
import RuleSettings from './RuleSettings.vue';
import UserManagement from './UserManagement.vue';
import PluginSettings from './PluginSettings.vue';
import {useGlobalStatusStore} from "@/stores/useGlobalStatusStore";

const globalStatus = useGlobalStatusStore();
const { isLogin, userInfos } = storeToRefs(globalStatus);

const openSettings = ref(false)
const drawerSize = ref('80%');

const updateDrawerSize = () => {
  drawerSize.value = window.innerWidth < 768 ? '100%' : '80%';
};

onMounted(() => {
  updateDrawerSize();
  window.addEventListener('resize', updateDrawerSize);
});

onUnmounted(() => {
  window.removeEventListener('resize', updateDrawerSize);
});

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

const toggleSidebar = () => {
  globalStatus.toggleSidebar();
}

const handleLanguageChange = (command) => {
  setLanguage(command);
  // No need to reload the page, as the `lang` object is now reactive
  // window.location.reload();
};

</script>


<style scoped>

#header_main {
  height: 56px; /* Standard header height */
  background: linear-gradient(to right, #004e8c, #0067b8); /* Dark blue gradient */
  display: flex;
  align-items: center;
  padding: 0 20px;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.2), 0 0 15px rgba(0, 103, 184, 0.7); /* More pronounced glow */
  border-bottom: none; /* Remove border-bottom as gradient handles it */
}

#logo {
  height: 100%;
  display: flex;
  align-items: center;
  font-size: 24px; /* Larger font size for prominence */
  font-weight: 600; /* Bolder font for the brand name */
  color: #ffffff; /* White text for contrast */
  text-align: left;
  flex-grow: 1;
}

#logo h1 {
  margin: 0; /* Remove default margin */
  padding: 0;
  color: inherit; /* Inherit color from parent */
  font-size: 20px; /* Smaller font size for the logo */
  text-shadow: none; /* Remove text shadow/glow */
}

#settings, #menu-toggle, .language-selector {
  display: flex;
  justify-content: center;
  align-items: center;
  cursor: pointer;
  padding: 8px; /* Add some padding for click area */
  border-radius: 4px; /* Slightly rounded corners for interactive elements */
  transition: background-color 0.2s ease; /* Smooth transition for hover */
}

#settings:hover, #menu-toggle:hover, .language-selector:hover {
  background-color: rgba(255, 255, 255, 0.1); /* Light white background on hover */
}

#settings .el-icon, #menu-toggle .el-icon, .language-selector .el-icon {
  font-size: 24px; /* Adjust icon size */
  color: #ffffff; /* White for icons */
}

/* Microsoft-like Drawer Styles */
.el-drawer__header {
  margin-bottom: 20px; /* Adjust spacing */
  padding: 20px; /* Add padding */
  border-bottom: 1px solid #e0e0e0; /* Subtle separator */
  font-size: 24px; /* Larger title */
  font-weight: 600;
  color: #333;
}

.el-drawer__body {
  padding: 0; /* Remove default padding to allow tabs to fill */
}

.el-tabs--left .el-tabs__header {
  margin-right: 0; /* Remove default margin */
  background-color: #f8f8f8; /* Light background for tab navigation */
  border-right: 1px solid #e0e0e0; /* Separator for tab navigation */
}

.el-tabs__item {
  height: 50px; /* Consistent height for tab items */
  line-height: 50px;
  font-size: 16px;
  color: #555; /* Default tab text color */
  padding: 0 20px;
  transition: background-color 0.2s ease, color 0.2s ease;
}

.el-tabs__item.is-active {
  background-color: #e6f2fa; /* Light blue background for active tab */
  color: #0078d4; /* Microsoft blue for active tab text */
  font-weight: 600;
}

.el-tabs__item:hover {
  background-color: #f0f0f0; /* Light grey on hover */
  color: #333;
}

.el-tabs__active-bar {
  background-color: #0078d4; /* Microsoft blue active bar */
}

.el-tabs__content {
  padding: 20px; /* Padding for the content area */
}
</style>
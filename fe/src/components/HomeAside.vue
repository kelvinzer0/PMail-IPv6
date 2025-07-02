<template>
  <div id="main">
    <div class="search-container">
      <input id="search" :placeholder="lang.search" v-model="searchQuery" @keyup.enter="performSearch" />
      <el-icon class="search-icon"><TbSearch /></el-icon>
    </div>
    <el-tree
      :data="data"
      :defaultExpandAll="true"
      @node-click="handleNodeClick"
    >
      <template #default="{ node, data }">
        <span class="custom-tree-node">
          <el-icon v-if="data.icon"><component :is="data.icon" /></el-icon>
          <span>{{ node.label }}</span>
        </span>
      </template>
    </el-tree>
  </div>
</template>

<script setup>
import { useRouter } from "vue-router";
import { ref, markRaw } from "vue";
import useGroupStore from "../stores/group";
import { lang } from "../i18n/i18n";
import { http } from "@/utils/axios";
import { TbFolder, TbSearch, TbMail } from "vue-icons-plus/tb";

const groupStore = useGroupStore();
const router = useRouter();
const data = ref([]);
const searchQuery = ref('');

const performSearch = () => {
  groupStore.searchQuery = searchQuery.value;
  router.push({
    name: "list",
  });
};

http.get("/api/group").then((res) => {
  if (res.data) {
    data.value = addFolderIconRecursively(res.data, true);
  } else {
    data.value = [];
  }
});

const handleNodeClick = function (data) {
  if (data.tag != null) {
    groupStore.name = data.label;
    groupStore.tag = data.tag;
    router.push({
      name: "list",
    });
  }
};

const addFolderIconRecursively = (items, isRoot = false) => {
  return items.map((item, index) => {
    let iconToUse = markRaw(TbFolder);
    if (isRoot && index === 0) {
      iconToUse = markRaw(TbMail);
    }
    const newItem = { ...item, icon: iconToUse };
    if (newItem.children && newItem.children.length > 0) {
      newItem.children = addFolderIconRecursively(newItem.children);
    }
    return newItem;
  });
};
</script>

<style scoped>
#main {
  width: 243px;
  background-color: #f1f1f1;
  height: 100%;
}

.search-container {
  position: relative;
  width: 100%;
  margin-bottom: 10px;
}

#search {
  background-color: #d6e7f7;
  width: 100%;
  height: 40px;
  padding-left: 40px; /* Make space for the icon */
  border: none;
  outline: none;
  font-size: 16px;
  box-sizing: border-box;
}

.search-icon {
  position: absolute;
  left: 10px;
  top: 50%;
  transform: translateY(-50%);
  color: #555;
}

.el-tree {
  background-color: #f1f1f1;
}

.custom-tree-node {
  flex: 1;
  display: flex;
  align-items: center;
  /* Removed justify-content: space-between; */
  font-size: 14px;
  padding-right: 8px;
}

.custom-tree-node .el-icon {
  margin-right: 8px;
}

.add_group {
  font-size: 14px;
  text-align: left;
  padding-left: 15px;
}
</style>
<template>
  <div id="main">
    <el-tree
      :data="data"
      :defaultExpandAll="true"
      @node-click="handleNodeClick"
    />
  </div>
</template>

<script setup>
import { useRouter } from "vue-router";
import { ref } from "vue";
import useGroupStore from "../stores/group";
import lang from "../i18n/i18n";
import { http } from "@/utils/axios";

const groupStore = useGroupStore();
const router = useRouter();
const data = ref([]);

http.get("/api/group").then((res) => {
  if (res.data) data.value = res.data;
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
</script>

<style scoped>
#main {
  width: 250px;
  background-color: #f9f9f9;
  height: 100%;
  padding-top: 10px;
  border-right: 1px solid #e0e0e0;
}

.el-tree {
  background-color: #f9f9f9;
  padding: 10px 0;
}

.el-tree-node__content {
  height: 36px;
  line-height: 36px;
  padding-left: 15px !important;
}

.el-tree-node__content:hover {
  background-color: #e8e8e8;
}

.el-tree-node.is-current > .el-tree-node__content {
  background-color: #e8e8e8;
  color: #4285f4;
}
</style>
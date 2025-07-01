<script setup>
import {RouterView, useRoute} from 'vue-router'
import HomeHeader from '@/components/HomeHeader.vue'
import HomeAside from '@/components/HomeAside.vue';
import {ref, watch} from 'vue'

const route = useRoute()
const pageName = ref(route.name)




watch(
    () => route.fullPath,
    () => {
      pageName.value = route.name
    }
)

</script>

<template>
  <div id="app-container">
    <HomeHeader/>
    <div id="content-wrapper">
      <aside id="aside" v-if="pageName !== 'login' && pageName !== 'setup'">
        <HomeAside/>
      </aside>
      <main id="main-content">
        <RouterView/>
      </main>
      <section id="detail-view" aria-label="Email Detail View">
        <!-- This will be used for email detail view later -->
      </section>
    </div>
  </div>
</template>


<style scoped>
#aside {
  background-color: #F1F1F1;
}

#main-content {
  flex-grow: 1;
  height: 100%;
}

#detail-view {
  width: 0;
  height: 100%;
  overflow: hidden;
  transition: width 0.3s ease-in-out;
}

#content.show-detail #detail-view {
  width: 50%; /* Adjust as needed */
}

#content.show-detail #main-content {
  width: 50%; /* Adjust as needed */
}

#content {
  display: flex;
  height: 100%;
}

#main {
  height: 100%;
  display: flex;
  flex-direction: column;
}
</style>

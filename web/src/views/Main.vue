<script setup lang="ts">
  import { RouterView } from 'vue-router';
</script>

<template>
  <Header />
  <n-layout has-sider position="absolute" style="top: 3.75rem">
    <DefaultSidebar />
    <n-layout-content class="p-6 bg-white">
      <router-view />
    </n-layout-content>
  </n-layout>
</template>

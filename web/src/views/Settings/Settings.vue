<script setup lang="ts">
  import Header from '@/components/Header.vue';
</script>

<template>
  <Header />
  <n-layout has-sider position="absolute" style="top: 3.75rem">
    <SettingsSidebar />
    <n-layout-content class="p-6 bg-white">
      <router-view />
    </n-layout-content>
  </n-layout>
</template>

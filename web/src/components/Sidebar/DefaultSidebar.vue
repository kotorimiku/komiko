<template>
  <Sidebar :menuItems="menuItems" :onNavigate="handleClick" />
</template>

<script lang="ts" setup>
  import Sidebar from './Sidebar.vue';
  import { ref, onMounted } from 'vue';
  import { useRouter } from 'vue-router';
  import { libraryApi } from '@/api';

  const router = useRouter();
  const handleClick = (path: string) => {
    router.push(path);
  };

  const menuItems = ref([
    { label: '主页', path: '/home', icon: 'i-heroicons:home' },
    { label: '所有系列', path: '/library', icon: 'i-heroicons-book-open' },
  ]);

  onMounted(async () => {
    const res = await libraryApi.fetchLibraries();
    const libraryItems: any = {
      label: '库',
      path: '/library',
      children: [],
      isOpen: true,
      icon: 'i-heroicons:book-open',
    };
    libraryItems.children.push(
      ...res.map((item) => ({ label: item.name, path: `/library/${item.id}`}))
    );
    menuItems.value.push(libraryItems);
  });
</script>

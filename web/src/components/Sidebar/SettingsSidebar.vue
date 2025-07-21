<template>
  <Sidebar :menuItems="menuItems" :onNavigate="handleClick" />
</template>

<script lang="ts" setup>
  import Sidebar from './Sidebar.vue';
  import { useRouter } from 'vue-router';
  import { useUserStore } from '@/stores/user';
  import type { MenuItem } from './SidebarMenuItem.vue';
  const router = useRouter();
  const handleClick = (path: string) => {
    router.push(path);
  };

  const userStore = useUserStore();
  const isAdmin = computed(() => {
    return userStore.user?.role === 'admin';
  });

  const menuItems = ref<MenuItem[]>([
    { label: '用户', path: '/settings/user', icon: 'i-heroicons:user' },
    { label: '关于', path: '/settings/about', icon: 'i-heroicons:information-circle' },
  ]);

  const adminItems = [
    { label: '库', path: '/settings/library', icon: 'i-heroicons:book-open' },
    { label: '用户管理', path: '/settings/user-manage', icon: 'i-heroicons:user-group' },
  ];

  watch(isAdmin, (newVal) => {
    if (newVal) {
      menuItems.value.push(...adminItems);
    }
  });
  onMounted(() => {
    if (isAdmin.value) {
      menuItems.value.push(...adminItems);
    }
  });
</script>

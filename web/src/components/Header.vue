<template>
  <header
    class="fixed top-0 left-0 w-full z-50 bg-white shadow flex items-center justify-between h-15"
  >
    <div class="flex-1 ml-6">
      <span
        class="text-xl font-bold text-gray-800 cursor-pointer"
        @click="goHome"
        >Komiko</span
      >
    </div>
    <div class="flex-2 flex items-center justify-center">
      <n-input
        v-model:value="searchQuery"
        placeholder="搜索..."
        class="w-55"
        @keyup.enter="onSearch"
      >
        <template #suffix>
          <n-button quaternary size="small" @click="onSearch">
            <template #icon>
              <div class="i-heroicons:magnifying-glass"></div>
            </template>
          </n-button>
        </template>
      </n-input>
    </div>
    <div class="flex-1 flex items-center justify-end mr-6">
      <n-space>
        <!-- 管理员可见的任务管理按钮 -->
        <n-button v-if="isAdmin" quaternary size="small" @click="goTasks" title="任务管理">
          <template #icon>
            <div class="i-heroicons:queue-list"></div>
          </template>
        </n-button>
        <n-button quaternary size="small" @click="goSettings">
          <template #icon>
            <div class="i-heroicons:cog-6-tooth"></div>
          </template>
        </n-button>
        <template v-if="user">
          <n-dropdown :options="dropdownOptions" @select="handleDropdownSelect">
            <n-text class="cursor-pointer">{{ user.username }}</n-text>
          </n-dropdown>
        </template>
      </n-space>
    </div>
  </header>
</template>

<script setup lang="ts">
  import { ref, computed } from 'vue';
  import { useRouter } from 'vue-router';
  import { NInput, NButton, NSpace, NText, NDropdown } from 'naive-ui';
  import { useUserStore } from '@/stores/user';

  const searchQuery = ref('');
  const router = useRouter();

  const userStore = useUserStore();
  const user = computed(() => userStore.user);
  
  // 判断当前用户是否为管理员
  const isAdmin = computed(() => user.value?.role === 'admin');

  const dropdownOptions = [
    { label: '个人中心', key: 'profile' },
    { label: '退出登录', key: 'logout' },
  ];

  function handleDropdownSelect(key: string) {
    switch (key) {
      case 'profile':
        router.push('/settings/user');
        break;
      case 'logout':
        logout();
        break;
    }
  }

  function onSearch() {

  }

  function logout() {
    userStore.logout();
    router.push('/login');
  }

  function goHome() {
    router.push('/');
  }

  function goSettings() {
    router.push('/settings');
  }
  
  function goTasks() {
    router.push('/tasks');
  }
</script>

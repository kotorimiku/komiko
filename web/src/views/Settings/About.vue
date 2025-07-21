<template>
  <div class="max-w-2xl mx-auto">
    <n-card title="关于 Komiko" class="mb-6">
      <div class="space-y-4">
        <div class="flex items-center space-x-3">
          <div class="w-16 h-16 bg-gradient-to-br from-blue-500 to-purple-600 rounded-lg flex items-center justify-center">
            <span class="text-white text-2xl font-bold">K</span>
          </div>
          <div>
            <h2 class="text-xl font-semibold">Komiko</h2>
            <p class="text-gray-600">漫画和小说阅读平台</p>
          </div>
        </div>

        <div class="border-t pt-4">
          <h3 class="text-lg font-medium mb-3">项目信息</h3>
          <div class="space-y-2 text-sm">
            <div class="flex justify-between">
              <span class="text-gray-600">版本：</span>
              <span>{{ version }}</span>
            </div>
          </div>
        </div>

        <div class="border-t pt-4">
          <h3 class="text-lg font-medium mb-3">相关链接</h3>
          <div class="space-y-2">
            <n-button
              type="primary"
              ghost
              class="w-full justify-start"
              @click="openGitHub"
            >
              <template #icon>
                <i class="i-simple-icons:github text-lg"></i>
              </template>
              GitHub 仓库
            </n-button>
            <n-button
              type="info"
              ghost
              class="w-full justify-start"
              @click="openIssues"
            >
              <template #icon>
                <i class="i-heroicons:exclamation-triangle text-lg"></i>
              </template>
              问题反馈
            </n-button>
          </div>
        </div>

        <div class="border-t pt-4">
          <h3 class="text-lg font-medium mb-3">功能特性</h3>
          <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
            <div class="flex items-center space-x-2">
              <i class="i-heroicons:check-circle text-green-500"></i>
              <span class="text-sm">漫画阅读</span>
            </div>
            <div class="flex items-center space-x-2">
              <i class="i-heroicons:check-circle text-green-500"></i>
              <span class="text-sm">小说阅读</span>
            </div>
            <div class="flex items-center space-x-2">
              <i class="i-heroicons:check-circle text-green-500"></i>
              <span class="text-sm">进度同步</span>
            </div>
            <div class="flex items-center space-x-2">
              <i class="i-heroicons:check-circle text-green-500"></i>
              <span class="text-sm">多用户支持</span>
            </div>
          </div>
        </div>
      </div>
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { getVersion } from '@/api/version';

const version = ref('');

onMounted(async () => {
  const res = await getVersion();
  version.value = res;
});

const openGitHub = () => {
  window.open('https://github.com/kotorimiku/komiko', '_blank');
};

const openIssues = () => {
  window.open('https://github.com/kotorimiku/komiko/issues', '_blank');
};
</script>
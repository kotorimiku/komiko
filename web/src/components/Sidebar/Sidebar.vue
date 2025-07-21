<template>
  <n-layout-sider
    width="200"
    bordered
    class="bg-gradient-to-b from-white to-gray-100 shadow-lg rounded-xl p-0 m-3"
  >
    <div
      class="flex flex-col gap-3 p-4 border border-gray-200 rounded-lg bg-white/80 backdrop-blur"
    >
      <template v-for="item in menuItems" :key="item.path">
        <SidebarMenuItem
          v-if="!item.hide"
          :item="item"
          @navigate="onNavigate"
          class="   "
        />
      </template>
    </div>
  </n-layout-sider>
</template>

<script lang="ts" setup>
  import { NLayoutSider } from 'naive-ui';
  import SidebarMenuItem from '@/components/Sidebar/SidebarMenuItem.vue';
  import { defineProps } from 'vue';

  interface MenuItem {
    label: string;
    path: string;
    isOpen?: boolean;
    hide?: boolean;
    children?: MenuItem[];
  }

  const props = defineProps<{
    menuItems: MenuItem[];
    onNavigate: (path: string) => void;
  }>();
</script>

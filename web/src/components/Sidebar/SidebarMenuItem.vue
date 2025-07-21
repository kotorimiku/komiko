<script setup lang="ts">
  import { ref } from 'vue';

  export interface MenuItem {
    label: string;
    path: string;
    icon?: string;
    isOpen?: boolean;
    hide?: boolean;
    children?: MenuItem[];
  }

  const props = defineProps<{ item: MenuItem }>();
  const emit = defineEmits(['navigate']);

  const isOpen = ref(props.item.isOpen || false);

  function handleClick() {
    if (props.item.children) {
      isOpen.value = !isOpen.value;
    } else {
      emit('navigate', props.item.path);
    }
  }
</script>

<template>
  <div>
    <div
      class="flex items-center rounded hover:bg-gray-200 cursor-pointer text-lg gap-3"
      @click="handleClick"
    >
      <span v-if="item.icon" :class="item.icon"></span>
      <span>{{ item.label }}</span>
      <span v-if="item.children" class="ml-auto">
        <n-icon>
          <div class="i-heroicons:chevron-right" v-if="!isOpen"></div>
          <div class="i-heroicons:chevron-down" v-else></div>
        </n-icon>
      </span>
    </div>
    <div v-if="item.children && isOpen" class="pl-4">
      <SidebarMenuItem
        v-for="child in item.children"
        :key="child.path"
        :item="child"
        @navigate="emit('navigate', $event)"
      />
    </div>
  </div>
</template>

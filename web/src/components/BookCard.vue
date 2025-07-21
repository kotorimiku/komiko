<script setup lang="ts">
  import { NCard, NDropdown } from 'naive-ui';

  const props = withDefaults(
    defineProps<{
      src: string;
      title?: string;
      count?: number;
      dropdownOptions?: { label: string; key: string }[];
      height?: number;
    }>(),
    {
      title: '',
      count: 1,
      dropdownOptions: () => [],
      height: 200,
    }
  );

  const emit = defineEmits<{
    (e: 'click'): void;
    (e: 'menu-click', key: string): void;
  }>();

  const handleDropdown = (key: string) => {
    emit('menu-click', key);
  };

  const onImgLoad = (e: Event) => {
    const target = e.target as HTMLImageElement | null;
    if (target && target.classList) {
      target.classList.add('opacity-100');
    }
  };
</script>

<template>
  <n-card
    class="book-card flex-shrink-0 cursor-pointer relative transition-all duration-200 hover:shadow-2xl hover:-translate-y-1 p-0 overflow-hidden"
    :content-style="{ padding: '0', background: 'transparent' }"
    :bordered="false"
    @click="emit('click')"
  >
    <div
      class="justify-center relative"
      :style="{
        height: `${height}px`,
      }"
    >
      <div class="w-full h-full">
        <transition name="fade">
          <div
            v-if="count > 1"
            class="absolute top-2 right-2 bg-black/70 text-white text-xs px-1 py-0.5 rounded-full z-10"
          >
            {{ count }}
          </div>
        </transition>
        <img
          :src="src"
          class="h-full w-full object-cover object-[25%_center]"
          @load="onImgLoad"
        />
      </div>
    </div>

    <div
      class="flex items-center justify-between px-2 text-sm text-gray-700 mt-2 mb-1"
    >
      <div class="truncate font-medium max-w-[80%]">{{ title }}</div>
      <n-dropdown
        trigger="click"
        :options="dropdownOptions"
        @select="handleDropdown"
      >
        <n-button
          quaternary
          @click.stop
          class="hover:text-black min-w-0 p-1 text-gray-500"
        >
          <template #icon>
            <span class="i-heroicons:list-bullet" />
          </template>
        </n-button>
      </n-dropdown>
    </div>
  </n-card>
</template>

<style scoped>
  .book-card {
    border-radius: 12px;
    box-shadow: 0 2px 8px 0 rgba(0, 0, 0, 0.08);
    background: #fff;
  }
</style>

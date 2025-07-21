<script setup lang="ts">
  import {
    ref,
    computed,
    onMounted,
    watch,
    nextTick,
    onBeforeUnmount,
  } from 'vue';
  import BookCard from './BookCard.vue';
  import { NButton } from 'naive-ui';

  interface Book {
    id: string | number;
    src: string;
    title?: string;
    count?: number;
    dropdownOptions?: { label: string; key: string }[];
    height?: number;
  }

  const props = defineProps<{
    bookList: Book[];
    cardHeight?: number;
    title?: string;
  }>();

  const emit = defineEmits<{
    (e: 'click', value: Book): void;
  }>();

  const handleClick = (value: Book) => {
    emit('click', value);
  };

  const cardWidth = 150; // px, 与模板 class 保持一致
  const gap = 16; // px, gap-4 = 16px
  const cardHeight = computed(() => props.cardHeight ?? 200);
  const currentIndex = ref(0);
  const visibleCount = ref(1);
  const containerRef = ref<HTMLElement | null>(null);

  function updateVisibleCount() {
    if (!containerRef.value) return;
    const width = containerRef.value.offsetWidth;
    visibleCount.value = Math.max(
      1,
      Math.floor((width + gap) / (cardWidth + gap))
    );
    // 修正 currentIndex 防止越界
    if (currentIndex.value > maxIndex.value) {
      currentIndex.value = maxIndex.value;
    }
  }

  const maxIndex = computed(() => {
    return Math.max(0, props.bookList.length - visibleCount.value);
  });

  const canPrev = computed(() => currentIndex.value > 0);
  const canNext = computed(() => currentIndex.value < maxIndex.value);

  function prev() {
    if (canPrev.value) {
      currentIndex.value = Math.max(
        0,
        currentIndex.value - visibleCount.value / 2
      );
    }
  }
  function next() {
    if (canNext.value) {
      currentIndex.value = Math.min(
        maxIndex.value,
        currentIndex.value + visibleCount.value / 2
      );
    }
  }

  const trackStyle = computed(() => {
    return {
      transform: `translateX(-${(cardWidth + gap) * currentIndex.value}px)`,
      transition: 'transform 0.4s cubic-bezier(0.4,0,0.2,1)',
    };
  });

  onMounted(() => {
    updateVisibleCount();
    window.addEventListener('resize', updateVisibleCount);
    nextTick(updateVisibleCount);
  });
  onBeforeUnmount(() => {
    window.removeEventListener('resize', updateVisibleCount);
  });
  watch(() => props.bookList.length, updateVisibleCount);
</script>

<template>
  <n-h2 v-if="title">{{ title }}</n-h2>
  <div class="relative w-full">
    <div
      class="flex justify-end items-center gap-2 absolute right-0 top-0 z-10"
      style="padding: 8px"
    >
      <n-button
        quaternary
        circle
        size="small"
        :disabled="!canPrev"
        @click="prev"
        :class="{ '!cursor-auto': !canPrev }"
      >
        <template #icon> <span class="i-heroicons:chevron-left" /> </template>
      </n-button>
      <n-button
        quaternary
        circle
        size="small"
        :disabled="!canNext"
        @click="next"
        :class="{ '!cursor-auto': !canNext }"
      >
        <template #icon> <span class="i-heroicons:chevron-right" /> </template>
      </n-button>
    </div>
    <div class="w-full" ref="containerRef">
      <div class="flex gap-4" :style="trackStyle">
        <BookCard
          v-for="book in props.bookList"
          :key="book.id"
          v-bind="book"
          :height="cardHeight"
          class="w-[150px] flex-shrink-0"
          @click="handleClick(book)"
        />
      </div>
    </div>
  </div>
</template>

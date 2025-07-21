<script setup lang="ts">
  import { ref, computed } from 'vue';

  interface image {
    src: string;
    title: string;
    value: any;
    pageCount: number;
  }

  const props = withDefaults(
    defineProps<{
      images: image[];
      dropdownOptions?: { label: string; key: string }[];
      width?: number;
      xGap?: number;
    }>(),
    {
      dropdownOptions: () => [],
      width: 150,
      xGap: 20,
    }
  );

  const start = ref(0);
  const end = ref(0);
  const page = defineModel<number>('page', { default: 1 });

  const pagedImages = computed(() => {
    return props.images.slice(start.value, end.value);
  });

  const emit = defineEmits<{
    (e: 'click', value: any): void;
    (e: 'menu-click', key: string, value: any): void;
  }>();

  function handleClick(value: any) {
    emit('click', value);
  }

  const handleDropdown = (key: string, value: any) => {
    emit('menu-click', key, value);
  };
</script>

<template>
  <div ref="gridRef" class="space-y-4">
    <div
      :style="{
        gridTemplateColumns: `repeat(auto-fit, minmax(${width}px, ${width}px))`,
        columnGap: `${xGap}px`,
      }"
      class="justify-center grid gap-y-[10px]"
    >
      <div v-for="(image, index) in pagedImages" :key="index">
        <BookCard
          :src="image.src"
          :title="image.title"
          @click="() => handleClick(image.value)"
          @menu-click="(key: string) => handleDropdown(key, image.value)"
          :count="image.pageCount"
          :dropdown-options="dropdownOptions"
          :height="200"
        />
      </div>
    </div>

    <Pagination
      :min-col-width="width + xGap"
      v-model:start="start"
      v-model:end="end"
      :row-count="4"
      :item-count="images.length"
      v-model:page="page"
    />
  </div>
</template>

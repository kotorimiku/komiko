<script setup lang="ts">
  import { bookApi, progressApi, seriesApi } from '@/api';
  import type { Book, Progress, Series } from '@/types';
  import { useRoute, useRouter } from 'vue-router';

  const route = useRoute();
  const router = useRouter();

  const seriesId = route.params.id as string;
  const series = ref<Series>();
  const books = ref<Book[]>([]);
  const progresses = ref<Progress[]>([]);

  const images = computed(() => {
    return books.value.map((item) => ({
      src: item.cover ? `/api/cover/${item.cover}` : '',
      title: item.title,
      value: item,
      pageCount: item.pageCount || 0,
    }));
  });

  const dropdownOptions = [
    {
      label: '查看图片',
      key: 'image',
    },
  ];

  const handleDropdown = (key: any, value: any) => {
    switch (key) {
      case 'image':
        router.push(`/series/${seriesId}/comic/${value.id}`);
    }
  };

  const handleClick = (value: Book) => {
    const progress = progresses.value.find((item) => item.bookId === value.id);
    if (value.type === 'comic') {
      router.push(
        `/series/${seriesId}/comic/${value.id}/${progress?.page || 1}`
      );
    } else if (value.type === 'novel') {
      router.push(
        `/series/${seriesId}/novel/${value.id}/${progress?.page || 1}`
      );
    }
  };

  onMounted(async () => {
    bookApi.fetchBooksBySeries(seriesId).then((res) => {
      books.value = res;
    });

    seriesApi.fetchSeriesById(seriesId).then((res) => {
      series.value = res;
    });

    progressApi.fetchProgress(seriesId).then((res) => {
      progresses.value = res;
    });
  });
</script>

<template>
  <div v-if="series" class="series-info flex items-start gap-6 mb-6">
    <img
      v-if="series.cover"
      :src="`/api/cover/${series.cover}`"
      alt="系列封面"
      class="w-40 h-55 object-cover rounded-xl object-[25%_center] shadow-md"
    />
    <div>
      <h2>{{ series.title }}</h2>
      <p v-if="series.description" class="my-2 text-gray-600">
        {{ series.description }}
      </p>
      <p v-if="series.author" class="text-gray-500">
        作者：{{ series.author }}
      </p>
    </div>
  </div>
  <BookCardGrid
    :images="images"
    @click="handleClick"
    @menu-click="handleDropdown"
    :dropdown-options="dropdownOptions"
    :x-gap="20"
  />
</template>

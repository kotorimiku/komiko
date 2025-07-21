<script setup lang="ts">
  import { seriesApi } from '@/api';
  import type { Series } from '@/types';
  import { ref, computed, watch, onMounted } from 'vue';
  import { useRoute, useRouter } from 'vue-router';

  const router = useRouter();
  const route = useRoute();

  const id = computed(() => route.params.id as string);
  const page = ref(1);
  const series = ref<Series[]>([]);

  const images = computed(() => {
    return series.value.map((item) => ({
      src: item.cover ? '/api/cover/' + item.cover : '',
      title: item.title,
      value: item,
      pageCount: 0,
    }));
  });

  watch(id, async () => {
    await getSeries();
    page.value = Number(route.params.page) || 1;
  });

  const getSeries = async () => {
    if (id.value) {
      const res = await seriesApi.fetchSeriesByLibrary(id.value);
      series.value = res;
    } else {
      seriesApi.fetchSeriesAll().then((response) => {
        series.value = response;
      });
    }
  };

  const handleClick = (value: any) => {
    router.push(`/series/${value.id}`);
  };

  onMounted(async () => {
    getSeries();
  });
</script>

<template>
  <BookCardGrid :images="images" @click="handleClick" v-model:page="page" />
</template>

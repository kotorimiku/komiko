<script setup lang="ts">
  import { onMounted } from 'vue';
  import { getSeriesProgress } from '@/api/progress';
  import type { Progress, Series } from '@/types';
  import BookCardSlider from '@/components/BookCardSlider.vue';
  import { fetchSeriesNew, fetchSeriesUpdated } from '@/api/series';
  import { useRouter } from 'vue-router';

  const seriesProgress = ref<Progress[]>([]);
  const updatedSeries = ref<Series[]>([]);
  const newSeries = ref<Series[]>([]);

  const router = useRouter();

  const handleClick = (value: any) => {
    router.push(`/series/${value.id}`);
  };

  onMounted(() => {
    getSeriesProgress(15, 0).then((response) => {
      seriesProgress.value = response;
    });
    fetchSeriesUpdated(15, 0).then((response) => {
      updatedSeries.value = response;
    });
    fetchSeriesNew(15, 0).then((response) => {
      newSeries.value = response;
    });
  });
</script>

<template>
  <n-card class="reading-list">
    <n-space vertical :size="16">
      <!-- Recently Read Section -->
      <section>
        <BookCardSlider
          @click="handleClick"
          :bookList="
            seriesProgress.map((progress) => ({
              id: progress.series?.id ?? progress.id,
              src: '/api/cover/' + progress.series?.cover || '',
              title: progress.series?.title || '',
            }))
          "
          title="最近阅读"
        />
      </section>

      <n-divider />

      <!-- Recently Updated Series -->
      <section>
        <BookCardSlider
          @click="handleClick"
          :bookList="
            updatedSeries.map((series) => ({
              id: series.id,
              src: '/api/cover/' + series.cover || '',
              title: series.title || '',
            }))
          "
          title="最近更新的系列"
        />
      </section>

      <n-divider />

      <!-- New Series -->
      <section>
        <BookCardSlider
          @click="handleClick"
          :bookList="
            newSeries.map((series) => ({
              id: series.id,
              src: '/api/cover/' + series.cover || '',
              title: series.title || '',
            }))
          "
          title="新增系列"
        />
      </section>
    </n-space>
  </n-card>
</template>

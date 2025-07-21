<script setup lang="ts">
  import type { ComicReadConfig } from '@/utils/read-config';
  import ComicPanel from './ComicPanel.vue';

  enum PageType {
    First,
    Second,
    Single,
    Null,
  }

  interface Props {
    imageUrl: string;
    imageFirstUrl: string;
    imageSecondUrl: string;
    pageCount: number;
    pageType?: PageType;
  }

  const imageLeftUrl = computed(() => {
    if (readConfig.value.direction === 'rightToLeft') {
      return props.imageSecondUrl;
    } else {
      return props.imageFirstUrl;
    }
  });

  const imageRightUrl = computed(() => {
    if (readConfig.value.direction === 'rightToLeft') {
      return props.imageFirstUrl;
    } else {
      return props.imageSecondUrl;
    }
  });

  const props = withDefaults(defineProps<Props>(), {
    pageType: PageType.Single,
  });

  const page = defineModel<number>('page', { required: true });
  const back = defineModel<() => void>('back', { required: true });
  const readConfig = defineModel<ComicReadConfig>('readConfig', {
    default: () => ({
      layoutMode: 'single',
      direction: 'rightToLeft',
      isSplit: false,
      matchSpreadPages: false,
    }),
  });

  const isDouble = computed(() => {
    return (
      readConfig.value.layoutMode === 'double' &&
      props.pageType !== PageType.Single
    );
  });

  watch(page, () => {
    if (page.value < 1) {
      page.value = 1;
    }
    if (page.value > props.pageCount) {
      page.value = props.pageCount;
    }

    canvasImg.value.src = props.imageUrl;
  });

  function handlePageClick(event: MouseEvent) {
    const width = window.innerWidth;
    const x = event.clientX;

    if (x < width / 4) {
      readConfig.value.direction === 'leftToRight' ? prevPage() : nextPage();
    } else if (x > (3 * width) / 4) {
      readConfig.value.direction === 'leftToRight' ? nextPage() : prevPage();
    } else {
      showPanel.value = !showPanel.value;
    }
  }

  watch(readConfig.value, () => {
    if (useCanvas.value) {
      canvasIndex.value = 0;
      drawSplit(canvasSide.value);
    }
  });

  const showPanel = ref(false);
  const canvasRef = ref<HTMLCanvasElement | null>(null);
  const useCanvas = computed(() => {
    return (
      readConfig.value.layoutMode === 'single' &&
      imgHeight.value < imgWidth.value &&
      readConfig.value.isSplit
    );
  });
  const canvasSide = computed(() => {
    if (readConfig.value.direction === 'rightToLeft') {
      return canvasIndex.value === 0 ? 'right' : 'left';
    } else {
      return canvasIndex.value === 0 ? 'left' : 'right';
    }
  });
  const canvasIndex = ref(0);
  const imgWidth = ref(0);
  const imgHeight = ref(0);
  const canvasImg = ref(new Image());
  canvasImg.value.onload = () => {
    imgWidth.value = canvasImg.value.width;
    imgHeight.value = canvasImg.value.height;

    if (useCanvas.value) {
      nextTick(() => {
        drawSplit(canvasSide.value);
      });
    }
  };

  function drawSplit(side: string) {
    const canvas = canvasRef.value;

    const ctx = canvas!.getContext('2d');
    if (!ctx) return;

    const drawWidth = canvasImg.value.width / 2;

    canvas!.width = drawWidth;
    canvas!.height = canvasImg.value.height;

    ctx.clearRect(0, 0, canvas!.width, canvas!.height);
    ctx.drawImage(canvasImg.value, 0, 0, drawWidth, canvasImg.value.height);
    if (side == 'left') {
      ctx.drawImage(
        canvasImg.value,
        0,
        0,
        drawWidth,
        canvasImg.value.height,
        0,
        0,
        drawWidth,
        canvasImg.value.height
      );
    } else {
      ctx.drawImage(
        canvasImg.value,
        drawWidth,
        0,
        drawWidth,
        canvasImg.value.height,
        0,
        0,
        drawWidth,
        canvasImg.value.height
      );
    }
  }

  const nextPage = () => {
    if (readConfig.value.layoutMode === 'single') {
      if (useCanvas.value) {
        if (canvasIndex.value === 0) {
          canvasIndex.value = 1;
          drawSplit(canvasSide.value);
          return;
        }
      }
      canvasIndex.value = 0;
      if (page.value < props.pageCount) page.value++;
    } else {
      if (page.value < props.pageCount) {
        if (props.pageType == PageType.Single) {
          page.value++;
        } else if (props.pageType == PageType.First) {
          page.value += 2;
        } else if (props.pageType == PageType.Second) {
          page.value++;
        }
      }
    }
  };
  const prevPage = () => {
    if (readConfig.value.layoutMode === 'single') {
      if (useCanvas.value) {
        if (canvasIndex.value === 1) {
          canvasIndex.value = 0;
          drawSplit(canvasSide.value);
          return;
        }
      }
      canvasIndex.value = 1;
      if (page.value > 1) page.value--;
    } else {
      if (page.value > 1) {
        if (props.pageType == PageType.Single) {
          page.value--;
        } else if (props.pageType == PageType.First) {
          page.value--;
        } else if (props.pageType == PageType.Second) {
          page.value -= 2;
        }
      }
    }
  };

  onMounted(() => {
    canvasImg.value.src = props.imageUrl;
  });
</script>

<template>
  <div
    class="absolute top-0 left-0 h-full w-full z-20"
    @click="handlePageClick($event)"
  >
    <div
      class="flex justify-center items-center fixed inset-0 bg-black select-none pointer-events-auto h-full w-full"
    >
      <div v-if="isDouble" class="flex">
        <img
          class="max-h-screen max-w-screen"
          v-if="imageLeftUrl"
          :src="imageLeftUrl"
        />
        <img
          class="max-h-screen max-w-screen"
          v-if="imageRightUrl"
          :src="imageRightUrl"
        />
      </div>

      <div v-else class="h-full">
        <canvas v-show="useCanvas" ref="canvasRef" class="h-full" />
        <img
          v-if="!useCanvas && imageUrl"
          class="max-h-screen max-w-screen"
          :src="imageUrl"
        />
      </div>
    </div>
  </div>
  <ComicPanel
    v-show="showPanel"
    :page-count="pageCount"
    v-model:page="page"
    :next-page="nextPage"
    :prev-page="prevPage"
    :back="back"
    v-model:read-config="readConfig"
    class="z-30"
  />
</template>

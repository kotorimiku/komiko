<script setup lang="ts">
  import { ref, watch } from 'vue';
  import { useRoute, useRouter } from 'vue-router';
  import { bookApi, progressApi } from '@/api';
  import type { Progress } from '@/types';
  import {
    getComicReadConfig,
    type ComicReadConfig,
  } from '@/utils/read-config';

  const route = useRoute();
  const router = useRouter();

  const readConfig = ref<ComicReadConfig>(getComicReadConfig());

  const bookId = ref(route.params.id as string);
  const seriesId = ref(route.params.seriesID as string);
  const book = ref<any>(null);
  const page = ref(Number(route.params.page) || 1);
  const cacheImages = ref<{ [key: string]: HTMLImageElement }>({});

  watch(page, () => {
    pageUpdate();
  });

  const pageUpdate = () => {
    router.push(
      `/series/${seriesId.value}/comic/${bookId.value}/${page.value}`
    );

    preloadImages(page.value - 3, 3 * 2 + 1);

    updateProgress();
  };

  const updateProgress = () => {
    const progress = {
      seriesId: book.value.seriesId,
      bookId: book.value.id,
      page: page.value,
    };

    progressApi.updateProgress(progress as Progress);
  };

  const imageUrl = computed(() => {
    return `/api/book/${bookId.value}/comic/${page.value}`;
  });

  enum PageType {
    First,
    Second,
    Single,
    Null,
  }

  const doublePages = computed(() => {
    let images: PageType[] = [];
    images.push(PageType.Null);
    if (book.value) {
      if (readConfig.value.matchSpreadPages) {
        for (let i = 1; i <= book.value.pageCount; i++) {
          if (
            book.value.images[i - 1].width < book.value.images[i - 1].height
          ) {
            if (
              i + 1 <= book.value.pageCount &&
              book.value.images[i].width < book.value.images[i].height
            ) {
              images.push(PageType.First);
              images.push(PageType.Second);
              i++;
            } else {
              images.push(PageType.Single);
            }
          } else {
            images.push(PageType.Single);
          }
        }
      } else {
        images.push(PageType.Single);
        for (let i = 2; i <= book.value.pageCount; i++) {
          if (
            book.value.images[i - 1].width < book.value.images[i - 1].height
          ) {
            if (
              i + 1 <= book.value.pageCount &&
              book.value.images[i].width < book.value.images[i].height
            ) {
              images.push(PageType.First);
              images.push(PageType.Second);
              i++;
            } else {
              images.push(PageType.Single);
            }
          } else {
            images.push(PageType.Single);
          }
        }
      }
    }

    return images;
  });

  const imageFirstUrl = computed(() => {
    if (doublePages.value) {
      if (doublePages.value[page.value] === PageType.First) {
        return `/api/book/${bookId.value}/comic/${page.value}`;
      } else if (doublePages.value[page.value] === PageType.Second) {
        return `/api/book/${bookId.value}/comic/${page.value - 1}`;
      }
    }
    return '';
  });
  const imageSecondUrl = computed(() => {
    if (doublePages.value) {
      if (doublePages.value[page.value] === PageType.First) {
        return `/api/book/${bookId.value}/comic/${page.value + 1}`;
      } else if (doublePages.value[page.value] === PageType.Second) {
        return `/api/book/${bookId.value}/comic/${page.value}`;
      }
    }
    return '';
  });

  const preloadImages = (startPage: number, count: number) => {
    if (book.value) {
      var urls = [];
      for (let i = startPage; i < startPage + count; i++) {
        if (i > 1 && i <= book.value.pageCount) {
          var src = `/api/book/${bookId.value}/comic/${i}`;
          urls.push(src);
          if (!cacheImages.value[src]) {
            const img = new Image();
            img.src = src;
            cacheImages.value[src] = img;
          }
        }
      }

      for (let url in cacheImages.value) {
        if (!urls.includes(url)) {
          delete cacheImages.value[url];
        }
      }
    }
  };

  const getBook = async () => {
    const res = await bookApi.fetchBook(bookId.value);
    res.pageCount = res.images.length;
    book.value = res;
  };

  const back = () => {
    router.push(`/series/${seriesId.value}`);
  };

  onMounted(async () => {
    await getBook();

    pageUpdate();
  });
</script>

<template>
  <ComicReader
    :image-url="imageUrl"
    :image-first-url="imageFirstUrl"
    :image-second-url="imageSecondUrl"
    :page-count="book?.pageCount ?? 0"
    :page-type="doublePages[page]"
    :back="back"
    v-model:page="page"
    v-model:read-config="readConfig"
  />
</template>

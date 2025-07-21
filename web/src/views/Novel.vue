<script setup lang="ts">
  import { bookApi } from '@/api';
  import { useRoute, useRouter } from 'vue-router';

  const route = useRoute();
  const router = useRouter();

  const bookId = route.params.id as string;
  const page = ref(Number(route.params.page) || 1);
  const chapters = ref<string[]>([]);
  const book = ref<any>(null);
  const html = ref('');
  const src = computed(
    () => `/api/book/${bookId}/novel/file/${chapters.value[page.value - 1]}`
  );

  const direction = ref<'rightToLeft' | 'leftToRight'>('leftToRight');
  const showPanel = ref(false);

  watch(page, async () => {
    getHtml();
  });

  const prevPage = () => {
    if (page.value > 1) page.value--;
    router.push(`./${page.value}`);
  };

  const nextPage = () => {
    if (page.value < book.value.pageCount) page.value++;
    router.push(`./${page.value}`);
  };

  function handlePageClick(event: MouseEvent) {
    const width = window.innerWidth;
    const x = event.clientX;

    if (x < width / 4) {
      direction.value === 'leftToRight' ? prevPage() : nextPage();
    } else if (x > (3 * width) / 4) {
      direction.value === 'leftToRight' ? nextPage() : prevPage();
    } else {
      showPanel.value = !showPanel.value;
    }
  }

  const getBook = async () => {
    const res = await bookApi.fetchBook(bookId);
    book.value = res;
    chapters.value = res.pages;
  };

  const readerIframe = ref<HTMLIFrameElement | null>(null);

  const getHtml = async () => {
    let res = await bookApi.fetchNovel(bookId, page.value);
    res = res.replace(
      '<body',
      `<div class="h-full w-full top-0 left-0 absolute"`
    );
    res = res.replace('</body', '</div');
    res = res.replace(
      /src=(["'])([^"']+)\1/g,
      `src="/api/book/${bookId}/novel/$2"`
    );
    res = res.replace(
      /href=(["'])file\/([^"']+)\1/g,
      `href="/api/book/${bookId}/novel/file/$2"`
    );

    html.value = res;

    // const iframe = readerIframe.value;
    // const doc = iframe!.contentDocument!;

    // doc.body.addEventListener("click", handlePageClick);

    // function loadScript(src: any) {
    //   return new Promise((resolve, reject) => {
    //     const script = document.createElement("script");
    //     script.src = src;
    //     script.onload = resolve;
    //     script.onerror = reject;
    //     document.head.appendChild(script);
    //   });
    // }

    // await loadScript("http://localhost:5173/api/book/10/novel/file/OEBPS/Misc/notereplace.js");
  };

  onMounted(async () => {
    await getBook();
    // await getChapters();
    await getHtml();
  });
</script>

<template>
  <div
    class="absolute top-0 left-0 h-full w-full z-20"
    @click="handlePageClick($event)"
  >
    <div ref="content" class="h-full w-full" v-html="html"></div>
    <!-- <iframe
      :src="src"
      style="width: 100%; height: 600px; border: none"
    ></iframe> -->
    <!-- <iframe
      @click="handlePageClick($event)"
      ref="readerIframe"
      :srcdoc="html"
      style="width: 100%; height: 600px; border: none"
    ></iframe> -->
  </div>
</template>

<style>
  img {
    max-width: 100%;
    max-height: 100%;
    /* display: block;
  margin: 0 auto; */
    object-fit: scale-down;
    height: 100vh;
    width: 100vw;
  }
</style>

<template>
  <div ref="gridRef">
    <div class="flex justify-center pt-4">
      <n-pagination
        v-model:page="page"
        :item-count="count"
        :page-size="pageSize"
        :show-size-picker="false"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
  const gridRef = ref<HTMLElement | null>(null);
  const columnCount = ref(1); // 动态列数

  const pageSize = computed(() => columnCount.value * props.rowCount);

  const props = withDefaults(
    defineProps<{
      content?: Array<any>; // 内容列表
      itemCount?: number; // 总条数,如果content为空，则使用itemCount
      minWidth?: number;
      minColWidth: number; // minWidth + x-gap
      rowCount?: number; // 每页固定行数
    }>(),
    {
      content: () => [],
      rowCount: 6,
    }
  );

  const count = computed(() => {
    if (props.content.length > 0) {
      return props.content.length;
    }
    return props.itemCount;
  });

  const page = defineModel<number>('page', { default: 1 });
  const startOffset = defineModel<number>('start', { default: 0 });
  const endOffset = defineModel<number>('end', { default: 0 });
  const start = computed(() => (page.value - 1) * pageSize.value);
  const end = computed(() => start.value + pageSize.value);
  watchEffect(() => {
    startOffset.value = start.value;
    endOffset.value = end.value;
  });

  const pagedContent = defineModel<Array<any>>('pagedContent', { default: [] });
  const paged = computed(() => {
    return props.content.slice(start.value, end.value);
  });
  watch([paged, pageSize], () => {
    pagedContent.value = paged.value;
  });

  // 自动计算列数
  function calculateColumns() {
    if (gridRef.value && isVisible.value) {
      const width = gridRef.value.clientWidth;
      // const minColWidth = props.minWidth + 10; // min-width + gap

      columnCount.value = Math.floor(width / props.minColWidth) || 1;
    }
  }

  // 使用 ResizeObserver 监听容器宽度变化
  onMounted(() => {
    calculateColumns();
    const resizeObserver = new ResizeObserver(() => {
      calculateColumns();
    });
    if (gridRef.value) resizeObserver.observe(gridRef.value);
    if (props.content.length > 0) {
      pagedContent.value = paged.value;
    }
  });

  // 监听内容变化
  watch(
    () => props.content,
    () => {
      if (props.content.length > 0) {
        pagedContent.value = paged.value;
      }
    }
  );

  // 容器变化时重置分页
  watch(pageSize, () => {
    if (isVisible.value) {
      page.value = 1;
    }
  });

  const isVisible = ref(true);

  onActivated(() => {
    isVisible.value = true;
  });

  onDeactivated(() => {
    isVisible.value = false;
  });
</script>

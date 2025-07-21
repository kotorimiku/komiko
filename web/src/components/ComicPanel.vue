<template>
  <div>
    <!-- 顶部控制栏 -->
    <div
      class="fixed top-0 inset-x-0 z-40 bg-black/80 backdrop-blur-md text-white px-4 py-3 shadow-lg rounded-b-xl flex items-center justify-between transition-opacity duration-300"
    >
      <div class="flex items-center gap-3">
        <n-button
          color="white"
          quaternary
          size="medium"
          circle
          @click="back"
          title="返回"
        >
          <template #icon>
            <div class="i-material-symbols-arrow-back text-lg"></div>
          </template>
        </n-button>
        <div class="text-sm font-medium opacity-80 hidden md:block">
          {{ bookTitle }}
        </div>
      </div>
      <div class="flex items-center justify-end gap-3">
        <n-button
          color="white"
          quaternary
          size="medium"
          circle
          @click="showSettingsDrawer = true"
          title="设置"
        >
          <template #icon>
            <div class="i-material-symbols-settings text-lg"></div>
          </template>
        </n-button>
      </div>
    </div>

    <!-- 底部控制栏 -->
    <div
      class="fixed bottom-0 inset-x-0 z-40 bg-black/80 backdrop-blur-md text-white px-4 py-4 shadow-lg rounded-t-xl transition-opacity duration-300"
    >
      <div class="flex items-center gap-3 mb-2">
        <n-button
          color="white"
          quaternary
          size="large"
          circle
          @click="prevPage"
          title="上一页"
        >
          <template #icon>
            <div class="i-material-symbols-arrow-back text-lg"></div>
          </template>
        </n-button>
        <div class="w-16 text-center text-sm font-medium">
          <span class="text-base">{{ page }}</span>
          <span class="opacity-70">/{{ pageCount }}</span>
        </div>
        <n-slider
          v-model:value="page"
          :min="1"
          :max="pageCount"
          class="flex-grow mx-2"
          tooltip
          :step="1"
        />
        <n-button
          color="white"
          quaternary
          size="large"
          circle
          @click="nextPage"
          title="下一页"
        >
          <template #icon>
            <div class="i-material-symbols-arrow-forward text-lg"></div>
          </template>
        </n-button>
      </div>
    </div>

    <!-- 设置抽屉 -->
    <n-drawer v-model:show="showSettingsDrawer" :width="320" placement="right">
      <n-drawer-content title="阅读设置" closable>
        <div class="flex flex-col gap-6 py-2">
          <div>
            <div class="mb-3 font-medium text-base">布局模式</div>
            <n-radio-group v-model:value="readConfig.layoutMode" class="w-full">
              <n-space justify="space-around" class="w-full">
                <n-radio-button value="single">
                  <div class="flex items-center gap-2">
                    <div class="i-material-symbols-book"></div>
                    <span>单页</span>
                  </div>
                </n-radio-button>
                <n-radio-button value="double">
                  <div class="flex items-center gap-2">
                    <div class="i-material-symbols-book-2"></div>
                    <span>双页</span>
                  </div>
                </n-radio-button>
              </n-space>
            </n-radio-group>
          </div>

          <div>
            <div class="mb-3 font-medium text-base">阅读方向</div>
            <n-radio-group v-model:value="readConfig.direction" class="w-full">
              <n-space justify="space-around" class="w-full">
                <n-radio-button value="rightToLeft">
                  <div class="flex items-center gap-2">
                    <div class="i-material-symbols-arrow-left-alt"></div>
                    <span>从右到左</span>
                  </div>
                </n-radio-button>
                <n-radio-button value="leftToRight">
                  <div class="flex items-center gap-2">
                    <div class="i-material-symbols-arrow-right-alt"></div>
                    <span>从左到右</span>
                  </div>
                </n-radio-button>
              </n-space>
            </n-radio-group>
          </div>

          <div class="flex items-center justify-between py-2">
            <div class="font-medium">分割页面</div>
            <n-switch v-model:value="readConfig.isSplit"></n-switch>
          </div>

          <div class="flex items-center justify-between py-2">
            <div class="font-medium">更改跨页匹配</div>
            <n-switch v-model:value="readConfig.matchSpreadPages"></n-switch>
          </div>

          <n-button type="primary" class="w-full mt-4" @click="saveSettings">
            保存设置
          </n-button>
        </div>
      </n-drawer-content>
    </n-drawer>
  </div>
</template>

<script setup lang="ts">
  import { ref, onMounted, onUnmounted } from 'vue';
  import {
    NSlider,
    NButton,
    NRadioGroup,
    NRadioButton,
    NSpace,
    NSwitch,
    NDrawer,
    NDrawerContent,
  } from 'naive-ui';
  import {
    setComicReadConfig,
    type ComicReadConfig,
  } from '../utils/read-config';

  const props = defineProps<{
    pageCount: number;
    nextPage: () => void;
    prevPage: () => void;
    back: () => void;
    bookTitle?: string;
  }>();

  const page = defineModel<number>('page', { default: 1 });
  const readConfig = defineModel<ComicReadConfig>('readConfig', {
    default: () => ({
      layoutMode: 'single',
      direction: 'rightToLeft',
      isSplit: false,
      matchSpreadPages: false,
    }),
  });

  const showSettingsDrawer = ref(false);
  const controlsTimer = ref<number | null>(null);

  function saveSettings() {
    setComicReadConfig(readConfig.value);
    showSettingsDrawer.value = false;
  }

  onMounted(() => {
    window.addEventListener('keydown', handleKeyDown);
  });

  onUnmounted(() => {
    if (controlsTimer.value) {
      window.clearTimeout(controlsTimer.value);
    }
    window.removeEventListener('keydown', handleKeyDown);
  });

  function handleKeyDown(event: KeyboardEvent) {
    if (event.key === 'ArrowLeft') {
      readConfig.value.direction === 'rightToLeft'
        ? props.nextPage()
        : props.prevPage();
    } else if (event.key === 'ArrowRight') {
      readConfig.value.direction === 'rightToLeft'
        ? props.prevPage()
        : props.nextPage();
    } else if (event.key === 'Escape') {
      props.back();
    }
  }
</script>

<script setup lang="tsx">
  import { ref, reactive, onMounted } from 'vue';
  import { libraryApi } from '@/api';
  import {
    NButton,
    NDataTable,
    NModal,
    NForm,
    NFormItem,
    NInput,
    NSelect,
    NUpload,
    useMessage,
  } from 'naive-ui';
  import type { Library } from '@/types/index';

  const message = useMessage();

  const libraries = ref<Library[]>([]);
  const loading = ref(false);

  const showModal = ref(false);
  const isEdit = ref(false);
  const form = reactive<Partial<Library>>({
    name: '',
    path: '',
    type: 'comic',
  });
  const currentId = ref<number | null>(null);
  const coverFile = ref<File | null>(null);

  const typeOptions = [
    { label: '漫画', value: 'comic' },
    { label: '小说', value: 'novel' },
  ];

  async function fetchLibraries() {
    loading.value = true;
    try {
      libraries.value = await libraryApi.fetchLibraries();
    } finally {
      loading.value = false;
    }
  }

  function openAdd() {
    isEdit.value = false;
    Object.assign(form, { name: '', path: '', type: 'comic' });
    coverFile.value = null;
    showModal.value = true;
  }

  function openEdit(row: Library) {
    isEdit.value = true;
    Object.assign(form, row);
    currentId.value = row.id;
    coverFile.value = null;
    showModal.value = true;
  }

  async function handleDelete(row: Library) {
    if (confirm(`确定要删除库：${row.name} 吗？`)) {
      await libraryApi.deleteLibraryById(row.id.toString());
      message.success('删除成功');
      fetchLibraries();
    }
  }

  async function handleScanUpdate(row: Library) {
    await libraryApi.scanUpdateLibrary(row.id.toString());
    message.success('已触发扫描更新');
  }

  async function handleScanCreate(row: Library) {
    await libraryApi.scanCreateLibrary(row.id.toString());
    message.success('已触发扫描新建');
  }

  async function handleUpdateCover(row: Library) {
    await libraryApi.updateCover(row.id.toString());
    message.success('已触发更新封面');
  }

  async function handleCoverUpload(options: any) {
    // options.fileList 是数组，取第一个原始文件
    const rawFile =
      options.fileList && options.fileList[0] && options.fileList[0].file;
    coverFile.value = rawFile || null;
  }

  async function handleSubmit() {
    if (!form.name || !form.path || !form.type) {
      message.error('请填写完整信息');
      return;
    }
    let lib: Library;
    if (isEdit.value && currentId.value != null) {
      lib = await libraryApi.updateLibrary({ ...form, id: currentId.value } as Library);
      if (coverFile.value) {
        await libraryApi.updateLibraryCover(currentId.value.toString(), coverFile.value);
      }
      message.success('修改成功');
    } else {
      lib = await libraryApi.createLibrary(form as Library);
      if (coverFile.value) {
        await libraryApi.updateLibraryCover(lib.id.toString(), coverFile.value);
      }
      message.success('添加成功');
    }
    showModal.value = false;
    fetchLibraries();
  }

  const columns = [
    { title: 'ID', key: 'id', width: 60 },
    { title: '名称', key: 'name' },
    { title: '路径', key: 'path' },
    {
      title: '类型',
      key: 'type',
      render(row: Library) {
        return row.type === 'comic' ? '漫画' : '小说';
      },
    },
    {
      title: '操作',
      key: 'actions',
      width: 260,
      render(row: Library) {
        return (
          <div class="flex ">
            <n-tooltip placement="top">
              {{
                trigger: () => (
                  <n-button quaternary onClick={() => openEdit(row)}>
                    <n-icon>
                      <div class="i-heroicons:pencil-square"></div>
                    </n-icon>
                  </n-button>
                ),
                default: () => '编辑'
              }}
            </n-tooltip>

            <n-tooltip placement="top">
              {{
                trigger: () => (
                  <n-button quaternary onClick={() => handleDelete(row)}>
                    <n-icon>
                      <div class="i-heroicons:trash"></div>
                    </n-icon>
                  </n-button>
                ),
                default: () => '删除'
              }}
            </n-tooltip>

            <n-tooltip placement="top">
              {{
                trigger: () => (
                  <n-button quaternary onClick={() => handleScanUpdate(row)}>
                    <n-icon>
                      <div class="i-heroicons:arrow-path"></div>
                    </n-icon>
                  </n-button>
                ),
                default: () => '扫描库(更新元数据)'
              }}
            </n-tooltip>

            <n-tooltip placement="top">
              {{
                trigger: () => (
                  <n-button quaternary onClick={() => handleScanCreate(row)}>
                    <n-icon>
                      <div class="i-heroicons:plus-circle"></div>
                    </n-icon>
                  </n-button>
                ),
                default: () => '扫描库'
              }}
            </n-tooltip>

            <n-tooltip placement="top">
              {{
                trigger: () => (
                  <n-button quaternary onClick={() => handleUpdateCover(row)}>
                    <n-icon>
                      <div class="i-heroicons:photo"></div>
                    </n-icon>
                  </n-button>
                ),
                default: () => '更新封面'
              }}
            </n-tooltip>
          </div>
        );
      },
    },
  ];

  onMounted(fetchLibraries);
</script>

<template>
  <div>
    <div class="flex justify-between items-center mb-4">
      <h2 class="text-2xl font-bold">库管理</h2>
      <n-button type="primary" @click="openAdd">新增库</n-button>
    </div>
    <n-data-table :columns="columns" :data="libraries" :loading="loading" />

    <n-modal
      v-model:show="showModal"
      preset="dialog"
      title="库信息"
      style="width: 400px"
    >
      <n-form label-placement="top" class="space-y-2">
        <n-form-item label="名称" required>
          <n-input v-model:value="form.name" placeholder="请输入库名称">
            <template #password-visible-icon></template>
          </n-input>
        </n-form-item>
        <n-form-item label="路径" required>
          <n-input v-model:value="form.path" placeholder="请输入路径" />
        </n-form-item>
        <n-form-item label="类型" required>
          <n-select v-model:value="form.type" :options="typeOptions" />
        </n-form-item>
        <n-form-item label="封面">
          <n-upload
            :max="1"
            :default-upload="false"
            @change="handleCoverUpload"
          >
            <n-button>选择图片</n-button>
          </n-upload>
        </n-form-item>
      </n-form>
      <template #action>
        <n-button @click="showModal = false">取消</n-button>
        <n-button type="primary" @click="handleSubmit" style="margin-left: 8px"
          >保存</n-button
        >
      </template>
    </n-modal>
  </div>
</template>

<style scoped>
  .n-data-table .n-button + .n-button {
    margin-left: 8px;
  }
</style>

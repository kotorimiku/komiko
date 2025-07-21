<script setup lang="ts">
  import { ref, reactive, onMounted } from 'vue';
  import {
    fetchAllUsers,
    createUser,
    updateUser,
    deleteUserById,
  } from '@/api/user';
  import {
    NButton,
    NDataTable,
    NModal,
    NForm,
    NFormItem,
    NInput,
    useMessage,
    NSelect,
    NUpload,
  } from 'naive-ui';
  import type { User } from '@/types/index';
  import { h } from 'vue';

  const message = useMessage();

  const users = ref<User[]>([]);
  const loading = ref(false);

  const showModal = ref(false);
  const isEdit = ref(false);
  const form = reactive<Partial<User>>({
    username: '',
    password: '',
    name: '',
    email: '',
    cover: '',
    role: '',
  });
  const currentId = ref<number | null>(null);

  const roleOptions = [
    { label: '普通用户', value: 'user' },
    { label: '管理员', value: 'admin' },
  ];

  async function fetchUsers() {
    loading.value = true;
    try {
      users.value = await fetchAllUsers();
    } finally {
      loading.value = false;
    }
  }

  function openAdd() {
    isEdit.value = false;
    Object.assign(form, {
      username: '',
      password: '',
      name: '',
      email: '',
      cover: '',
      role: '',
    });
    currentId.value = null;
    showModal.value = true;
  }

  function openEdit(row: User) {
    isEdit.value = true;
    Object.assign(form, row);
    form.password = '';
    currentId.value = row.id;
    showModal.value = true;
  }

  async function handleDelete(row: User) {
    if (confirm(`确定要删除用户：${row.username} 吗？`)) {
      await deleteUserById(row.id);
      message.success('删除成功');
      fetchUsers();
    }
  }

  async function handleSubmit() {
    if (!form.username || (!isEdit.value && !form.password) || !form.role) {
      message.error('请填写完整信息');
      return;
    }
    if (isEdit.value && currentId.value != null) {
      await updateUser({ ...form, id: currentId.value } as User);
      message.success('修改成功');
    } else {
      await createUser(form as User);
      message.success('添加成功');
    }
    showModal.value = false;
    fetchUsers();
  }

  const columns = [
    { title: 'ID', key: 'id', width: 60 },
    { title: '用户名', key: 'username' },
    { title: '姓名', key: 'name' },
    { title: '邮箱', key: 'email' },
    {
      title: '角色',
      key: 'role',
      render(row: User) {
        return row.role === 'admin' ? '管理员' : '普通用户';
      },
    },
    {
      title: '操作',
      key: 'actions',
      width: 180,
      render(row: User) {
        return [
          h(
            NButton,
            { size: 'small', onClick: () => openEdit(row) },
            { default: () => '编辑' }
          ),
          h(
            NButton,
            {
              size: 'small',
              type: 'error',
              onClick: () => handleDelete(row),
              style: 'margin-left:8px',
            },
            { default: () => '删除' }
          ),
        ];
      },
    },
  ];

  onMounted(fetchUsers);
</script>

<template>
  <div>
    <div class="flex justify-between items-center mb-4">
      <h2 class="text-2xl font-bold">用户管理</h2>
      <n-button type="primary" @click="openAdd">新增用户</n-button>
    </div>
    <n-data-table :columns="columns" :data="users" :loading="loading" />

    <n-modal
      v-model:show="showModal"
      preset="dialog"
      title="用户信息"
      style="width: 400px"
    >
      <n-form label-placement="top" class="space-y-2">
        <n-form-item label="用户名" required>
          <n-input v-model:value="form.username" placeholder="请输入用户名" />
        </n-form-item>
        <n-form-item label="密码" :required="!isEdit">
          <n-input
            v-model:value="form.password"
            placeholder="请输入密码"
            type="password"
          />
        </n-form-item>
        <n-form-item label="姓名">
          <n-input v-model:value="form.name" placeholder="请输入姓名" />
        </n-form-item>
        <n-form-item label="邮箱">
          <n-input v-model:value="form.email" placeholder="请输入邮箱" />
        </n-form-item>
        <n-form-item label="头像URL">
          <n-input
            v-model:value="form.cover"
            placeholder="请输入头像图片链接"
          />
        </n-form-item>
        <n-form-item label="角色" required>
          <n-select
            v-model:value="form.role"
            :options="roleOptions"
            placeholder="请选择角色"
          />
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

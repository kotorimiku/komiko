<template>
  <div class="p-8 max-w-2xl mx-auto">
    <n-card title="用户信息" class="mb-6">
      <!-- <div class="flex items-center gap-4 mb-4">
        <n-avatar :src="form.cover" size="large" round />
        <n-upload
          :show-file-list="false"
          :custom-request="handleAvatarUpload"
        >
          <n-button size="small">更换头像</n-button>
        </n-upload>
      </div> -->
      <n-form :model="form" label-width="80">
        <n-form-item label="昵称">
          <n-input v-model:value="form.name" placeholder="请输入昵称" />
        </n-form-item>
        <n-form-item label="邮箱">
          <n-input v-model:value="form.email" placeholder="请输入邮箱" />
        </n-form-item>
        <n-form-item>
          <n-button
            type="primary"
            @click="onSaveUserInfo"
            :loading="userInfoLoading"
            >保存</n-button
          >
        </n-form-item>
      </n-form>
    </n-card>

    <n-card title="修改密码">
      <n-form :model="passwordForm" label-width="80" ref="passwordFormRef">
        <n-form-item label="原密码" path="oldPassword">
          <n-input
            v-model:value="passwordForm.oldPassword"
            type="password"
            placeholder="请输入原密码"
          />
        </n-form-item>
        <n-form-item label="新密码" path="newPassword">
          <n-input
            v-model:value="passwordForm.newPassword"
            type="password"
            placeholder="请输入新密码"
          />
        </n-form-item>
        <n-form-item label="确认新密码" path="confirmPassword">
          <n-input
            v-model:value="passwordForm.confirmPassword"
            type="password"
            placeholder="请再次输入新密码"
          />
        </n-form-item>
        <n-form-item>
          <n-button
            type="primary"
            @click="onChangePassword"
            :loading="passwordLoading"
            >修改密码</n-button
          >
        </n-form-item>
      </n-form>
    </n-card>
  </div>
</template>

<script setup lang="ts">
  import { ref, reactive, watch } from 'vue';
  import { useUserStore } from '@/stores/user';
  import { updateUser } from '@/api/user';
  import { useMessage } from 'naive-ui';
  import type { UploadCustomRequestOptions } from 'naive-ui';
  import type { User } from '@/types';

  const userStore = useUserStore();
  const message = useMessage();

  const form = reactive({
    cover: '',
    name: '',
    email: '',
  });

  const userInfoLoading = ref(false);

  watch(
    () => userStore.user,
    (user) => {
      if (user) {
        form.cover = user.cover;
        form.name = user.name;
        form.email = user.email;
      }
    },
    { immediate: true }
  );

  async function onSaveUserInfo() {
    userInfoLoading.value = true;
    try {
      // 只传递已存在的user字段，避免undefined
      const user = userStore.user;
      if (!user) throw new Error('用户未登录');
      await updateUser({
        id: user.id,
        username: user.username,
        password: user.password,
        name: form.name,
        email: form.email,
        cover: form.cover,
        role: user.role,
      });
      await userStore.getUser();
      message.success('用户信息已更新');
    } catch (e) {
      message.error('更新失败');
    } finally {
      userInfoLoading.value = false;
    }
  }

  function handleAvatarUpload(options: UploadCustomRequestOptions) {
    const file = options.file.file;
    if (!file) return;
    const reader = new FileReader();
    reader.onload = (e) => {
      form.cover = e.target?.result as string;
      options.onFinish();
    };
    reader.readAsDataURL(file);
  }

  // 修改密码表单
  const passwordForm = reactive({
    oldPassword: '',
    newPassword: '',
    confirmPassword: '',
  });
  const passwordLoading = ref(false);
  const passwordFormRef = ref();

  async function onChangePassword() {
    if (
      !passwordForm.oldPassword ||
      !passwordForm.newPassword ||
      !passwordForm.confirmPassword
    ) {
      message.warning('请填写完整');
      return;
    }
    if (passwordForm.newPassword !== passwordForm.confirmPassword) {
      message.warning('两次输入的新密码不一致');
      return;
    }
    passwordLoading.value = true;
    try {
      await updateUser({
        id: userStore.user!.id,
        password: passwordForm.newPassword,
      } as User);
      message.success('密码修改成功');
      passwordForm.oldPassword = '';
      passwordForm.newPassword = '';
      passwordForm.confirmPassword = '';
    } catch (e) {
      message.error('密码修改失败');
    } finally {
      passwordLoading.value = false;
    }
  }
</script>

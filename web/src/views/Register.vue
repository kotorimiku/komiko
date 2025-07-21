<template>
  <div class="w-full h-screen flex items-center justify-center bg-gray-100">
    <n-card class="w-[320px] shadow-lg">
      <div class="text-center mb-4 text-xl font-bold">注册</div>
      <n-form :model="form" :rules="rules" ref="formRef">
        <n-form-item path="username" label="用户名">
          <n-input v-model:value="form.username" placeholder="请输入用户名" />
        </n-form-item>
        <n-form-item path="password" label="密码">
          <n-input
            v-model:value="form.password"
            type="password"
            placeholder="请输入密码"
          />
        </n-form-item>
        <n-form-item path="confirmPassword" label="确认密码">
          <n-input
            v-model:value="form.confirmPassword"
            type="password"
            placeholder="请再次输入密码"
          />
        </n-form-item>
        <n-button type="primary" block @click="handleRegister">注册</n-button>
        <n-button text type="default" class="mt-2 w-full" @click="goToLogin">
          已有账号？登录
        </n-button>
      </n-form>
    </n-card>
  </div>
</template>

<script setup lang="ts">
  import { ref, reactive } from 'vue';
  import { useRouter } from 'vue-router';
  import { useMessage, type FormInst } from 'naive-ui';
  import { userApi } from '@/api';
  import type { User } from '@/types';

  const router = useRouter();
  const message = useMessage();
  const formRef = ref<FormInst | null>();

  const form = reactive({
    username: '',
    password: '',
    confirmPassword: '',
  });

  const rules = {
    username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
    password: [
      {
        required: true,
        trigger: ['blur', 'change'],
        validator: (_rule: any, value: string) => {
          if (!value) {
            return new Error('请输入密码');
          }
          if (value.length < 6) {
            return new Error('密码长度不能小于6位');
          }
          return true;
        },
      },
    ],
    confirmPassword: [
      {
        validator: (_rule: any, value: string) =>
          value !== form.password ? new Error('两次密码不一致') : true,
        trigger: 'blur',
      },
    ],
  };

  const handleRegister = async () => {
    try {
      await formRef.value?.validate();
    } catch {
      return;
    }
    try {
      await userApi.register({
        username: form.username,
        password: form.password,
      } as User);
      message.success(`注册成功`);
      router.push('/login');
    } catch (err: any) {
      message.error(err?.message || '注册失败，请稍后重试');
      console.log(err);
    }
  };

  const goToLogin = () => {
    router.push('/login');
  };

  onMounted(async () => {
    const allowRegister = await userApi.allowRegister();
    if (!allowRegister) {
      router.push('/login');
      return;
    }
    const token = localStorage.getItem('token');
    if (token) {
      router.push('/');
    }
  });
</script>

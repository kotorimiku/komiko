<template>
  <div class="w-full h-screen flex items-center justify-center bg-gray-100">
    <n-card class="w-[320px] shadow-lg">
      <div class="text-center mb-4 text-xl font-bold">登录</div>
      <n-form :model="form" :rules="rules" ref="formRef">
        <n-form-item path="username" label="用户名">
          <n-input v-model:value="form.username" placeholder="请输入用户名" />
        </n-form-item>
        <n-form-item path="password" label="密码">
          <n-input
            v-model:value="form.password"
            type="password"
            placeholder="请输入密码"
            @keydown.enter.prevent="handleLogin"
          />
        </n-form-item>
        <n-button type="primary" block @click="handleLogin">登录</n-button>
        <n-button
          v-if="allowRegister"
          text
          type="default"
          class="mt-2 w-full"
          @click="goToRegister"
        >
          还没有账号？注册
        </n-button>
      </n-form>
    </n-card>
  </div>
</template>

<script setup lang="ts">
  import { ref, reactive } from 'vue';
  import { useRouter } from 'vue-router';
  import { useMessage } from 'naive-ui';
  import type { User } from '@/types';
  import { useUserStore } from '@/stores/user';
  import { userApi } from '@/api';

  const router = useRouter();
  const message = useMessage();
  const formRef = ref();
  const allowRegister = ref(false);

  const form = reactive({
    username: '',
    password: '',
  });

  const rules = {
    username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
    password: [
      {
        required: true,
        trigger: 'blur',
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
  };

  const handleLogin = async () => {
    try {
      await formRef.value?.validate();
    } catch {
      return;
    }
    try {
      await useUserStore().login({
        username: form.username,
        password: form.password,
      } as User);
      message.success('登录成功');
      router.push('/');
    } catch (err: any) {
      message.error(err?.message || '登录失败');
      console.log(err);
    }
  };

  const goToRegister = () => {
    router.push('/register');
  };

  onMounted(async () => {
    allowRegister.value = await userApi.allowRegister();
    const token = localStorage.getItem('token');
    if (token) {
      router.push('/');
    }
  });
</script>

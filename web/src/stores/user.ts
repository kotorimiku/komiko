import { ref, computed } from 'vue';
import { defineStore } from 'pinia';
import type { User } from '@/types';
import { fetchUserByToken } from '@/api/user';
import { userApi } from '@/api';

export const useUserStore = defineStore('user', () => {
  const user = ref<User | null>(null);

  const getUser = async () => {
    const res = await fetchUserByToken();
    user.value = res;
  };

  const login = async (user: User): Promise<string> => {
    const res = await userApi.login(user);
    localStorage.setItem('token', res);
    getUser();
    return res;
  };

  const logout = () => {
    localStorage.removeItem('token');
    user.value = null;
  };

  onMounted(() => {
    getUser();
  });

  return { user, getUser, login, logout };
});

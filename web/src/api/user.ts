import { apiGetStr, apiPost, apiGet } from '@/client';
import type { User } from '@/types/index';

export async function register(user: User): Promise<void> {
  return await apiPost(`/user/register`, user);
}

export async function login(user: User): Promise<string> {
  return await apiPost(`/user/login`, user);
}

export async function allowRegister(): Promise<boolean> {
  return await apiGet(`/user/allow-register`);
}

export async function fetchAllUsers(): Promise<User[]> {
  return await apiGet<User[]>(`/user`);
}

export async function fetchUserById(id: number): Promise<User> {
  return await apiGet<User>(`/user/${id}`);
}

export async function createUser(user: User): Promise<string> {
  return await apiPost(`/user`, user);
}

export async function updateUser(user: User): Promise<string> {
  return await apiPost(`/user/update`, user);
}

export async function deleteUserById(id: number): Promise<string> {
  return await apiPost(`/user/${id}/delete`, {});
}

export async function fetchUserByToken(): Promise<User> {
  return await apiGet<User>(`/user/current`);
}

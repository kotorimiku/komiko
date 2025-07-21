import { apiGet, apiPost } from '@/client';
import type { Task } from '@/types';

export function fetchTasks(): Promise<Task[]> {
  return apiGet('/task');
}

export function fetchTask(id: string): Promise<Task> {
  return apiGet(`/task/${id}`);
}

export function stopTask(id: string): Promise<void> {
  return apiPost(`/task/${id}/stop`, {});
}
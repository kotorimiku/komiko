import { apiGet } from '@/client';

export const getVersion = async (): Promise<string> => {
  return await apiGet('/version');
};
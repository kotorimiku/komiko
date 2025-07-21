import { apiGet, apiPost } from '@/client';
import type { Progress } from '@/types/index';

export async function fetchProgress(seriesId: string): Promise<Progress[]> {
  return await apiGet(`/progress?series=${seriesId}`);
}

export async function updateProgress(progress: Progress): Promise<Progress> {
  return await apiPost(`/progress`, progress);
}

export async function getSeriesProgress(
  limit: number,
  offset: number
): Promise<Progress[]> {
  return await apiGet(`/progress/series?limit=${limit}&offset=${offset}`);
}

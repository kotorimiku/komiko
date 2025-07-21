import { apiGet } from '@/client';
import type { Series } from '@/types/index';

export async function fetchSeriesByLibrary(id: string): Promise<Series[]> {
  const params = new URLSearchParams({
    query: JSON.stringify({
      libraryId: Number(id),
    }),
  });
  return await apiGet(`/series?${params.toString()}`);
}

export async function fetchSeriesById(id: string): Promise<Series> {
  return await apiGet(`/series/${id}`);
}

export async function fetchSeriesAll(): Promise<Series[]> {
  return await apiGet(`/series`);
}

export async function fetchSeriesUpdated(
  limit: number,
  offset: number
): Promise<Series[]> {
  return await apiGet(
    `/series?sort=updated_at&desc=true&limit=${limit}&offset=${offset}`
  );
}

export async function fetchSeriesNew(
  limit: number,
  offset: number
): Promise<Series[]> {
  return await apiGet(
    `/series?sort=created_at&desc=true&limit=${limit}&offset=${offset}`
  );
}

export async function fetchSeries(
  sort: string,
  desc: boolean,
  limit: number,
  offset: number
): Promise<Series[]> {
  return await apiGet(
    `/series?sort=${sort}&desc=${desc}&limit=${limit}&offset=${offset}`
  );
}

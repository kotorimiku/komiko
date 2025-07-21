import { apiGet, apiGetStr } from '@/client';
import type { Book } from '@/types/index';

export async function fetchBook(id: string): Promise<Book> {
  return await apiGet(`/book/${id}`);
}

export async function fetchBooksBySeries(id: string): Promise<Book[]> {
  return await apiGet(`/book?series=${id}`);
}

export async function fetchNovel(bookId: string, page: number): Promise<string> {
  return await apiGetStr(`/book/${bookId}/novel/${page}`);
}

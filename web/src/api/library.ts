import { apiGet, apiPost } from '@/client';
import type { Library } from '@/types/index';

export const fetchLibraries = async (): Promise<Library[]> => {
  return await apiGet(`/library`);
};

export const fetchLibraryById = async (id: string): Promise<Library> => {
  return await apiGet(`/library/${id}`);
};

export const createLibrary = async (library: Library): Promise<Library> => {
  return await apiPost(`/library`, library);
};

export const updateLibrary = async (library: Library): Promise<Library> => {
  return await apiPost(`/library/update`, library);
};

export const deleteLibraryById = async (id: string): Promise<void> => {
  return await apiPost(`/library/${id}/delete`);
};

export const updateLibraryCover = async (
  id: string,
  cover: File | Blob
): Promise<Library> => {
  const formData = new FormData();
  formData.append('cover', cover);
  return await apiPost(`/library/${id}/update-cover`, formData);
};

export const scanUpdateLibrary = async (id: string): Promise<any> => {
  return await apiPost(`/library/${id}/scan-update`);
};

export const scanCreateLibrary = async (id: string): Promise<any> => {
  return await apiPost(`/library/${id}/scan-create`);
};

export const updateCover = async (id: string,): Promise<Library> => {
  return await apiPost(`/library/${id}/update-cover`);
};

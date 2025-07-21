import axios, { type AxiosRequestConfig } from 'axios';

interface ApiResponse<T> {
  code: number;
  data: T;
  msg: string;
}

export const client = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || '/api',
  timeout: 10000,
});

client.interceptors.request.use((config) => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers.Authorization = `${token}`;
  }
  return config;
});

// 响应拦截器：处理 JWT 过期
client.interceptors.response.use(
  (res) => res,
  (err) => {
    if (err.response && err.response.status === 401) {
      localStorage.removeItem('token');
      window.location.href = '/login';
    }
    return Promise.reject(err);
  }
);

export async function apiGet<T>(
  url: string,
  config?: AxiosRequestConfig
): Promise<T> {
  try {
    const res = (await client.get<ApiResponse<T>>(url, config)).data;
    if (res.code !== 200) {
      throw new Error(res.msg || '请求失败');
    }
    return res.data;
  } catch (error) {
    console.error('apiGet 请求失败:', error);
    throw error;
  }
}

export async function apiGetStr(
  url: string,
  config?: AxiosRequestConfig
): Promise<string> {
  try {
    const res = (await client.get<string>(url, config)).data;
    return res;
  } catch (error) {
    console.error('apiGet 请求失败:', error);
    throw error;
  }
}

export async function apiPost<T>(
  url: string,
  data?: any,
  config?: AxiosRequestConfig
): Promise<T> {
  try {
    const res = (await client.post<ApiResponse<T>>(url, data, config)).data;
    if (res.code !== 200) {
      throw new Error(res.msg || '请求失败');
    }
    return res.data;
  } catch (error) {
    console.error('apiGet 请求失败:', error);
    throw error;
  }
}

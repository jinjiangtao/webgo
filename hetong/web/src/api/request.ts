import axios from 'axios'
import { ElMessage } from 'element-plus'
import type { AxiosResponse, AxiosRequestConfig } from 'axios'
import type { ApiResponse } from '../types'

const request = axios.create({
  baseURL: 'http://localhost:8080/api',
  timeout: 30000,
})

request.interceptors.response.use(
  (response: AxiosResponse) => {
    const res = response.data as ApiResponse
    if (res.code !== 200 && res.code !== 0) {
      ElMessage.error(res.message || '请求失败')
      return Promise.reject(new Error(res.message || 'Error'))
    }
    return res as any
  },
  (error) => {
    ElMessage.error(error.message || '网络错误')
    return Promise.reject(error)
  }
)

export function get<T = any>(url: string, config?: AxiosRequestConfig): Promise<T> {
  return request.get(url, config) as any
}

export function post<T = any>(url: string, data?: any, config?: AxiosRequestConfig): Promise<T> {
  return request.post(url, data, config) as any
}

export function postBlob(url: string, data?: any, config?: AxiosRequestConfig): Promise<Blob> {
  return request.post(url, data, {
    ...config,
    responseType: 'blob',
  }) as any
}

export default request

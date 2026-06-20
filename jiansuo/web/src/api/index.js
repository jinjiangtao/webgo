import request from '@/utils/request'

export function search(params) {
  return request({
    url: '/search',
    method: 'get',
    params
  })
}

export function getSuggestions(params) {
  return request({
    url: '/search/suggest',
    method: 'get',
    params
  })
}

export function getHotKeywords(params) {
  return request({
    url: '/search/hot',
    method: 'get',
    params
  })
}

export function getStatistics() {
  return request({
    url: '/search/stats',
    method: 'get'
  })
}

export function incrementView(id) {
  return request({
    url: `/keywords/${id}/view`,
    method: 'post'
  })
}

export function getKeywordDetail(id) {
  return request({
    url: `/keywords/${id}`,
    method: 'get'
  })
}

export function listKeywords(params) {
  return request({
    url: '/keywords',
    method: 'get',
    params
  })
}

export function createKeyword(data) {
  return request({
    url: '/keywords',
    method: 'post',
    data
  })
}

export function updateKeyword(id, data) {
  return request({
    url: `/keywords/${id}`,
    method: 'put',
    data
  })
}

export function deleteKeyword(id) {
  return request({
    url: `/keywords/${id}`,
    method: 'delete'
  })
}

export function setKeywordStatus(id, status) {
  return request({
    url: `/keywords/${id}/status`,
    method: 'patch',
    data: { status }
  })
}

export function batchKeywords(data) {
  return request({
    url: '/keywords/batch',
    method: 'post',
    data
  })
}

export function importCSV(file) {
  const formData = new FormData()
  formData.append('file', file)
  return request({
    url: '/keywords/import',
    method: 'post',
    data: formData,
    headers: { 'Content-Type': 'multipart/form-data' }
  })
}

export function listCategories(params = {}) {
  return request({
    url: '/categories',
    method: 'get',
    params
  })
}

export function createCategory(data) {
  return request({
    url: '/categories',
    method: 'post',
    data
  })
}

export function updateCategory(id, data) {
  return request({
    url: `/categories/${id}`,
    method: 'put',
    data
  })
}

export function deleteCategory(id) {
  return request({
    url: `/categories/${id}`,
    method: 'delete'
  })
}

export function setCategoryStatus(id, status) {
  return request({
    url: `/categories/${id}/status`,
    method: 'patch',
    data: { status }
  })
}

export function listSearchLogs(params) {
  return request({
    url: '/logs',
    method: 'get',
    params
  })
}

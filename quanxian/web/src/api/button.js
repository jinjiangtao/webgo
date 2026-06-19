import request from '@/utils/request'

export function getButtonList(params) {
  return request({
    url: '/buttons',
    method: 'get',
    params
  })
}

export function getButton(id) {
  return request({
    url: `/buttons/${id}`,
    method: 'get'
  })
}

export function createButton(data) {
  return request({
    url: '/buttons',
    method: 'post',
    data
  })
}

export function updateButton(id, data) {
  return request({
    url: `/buttons/${id}`,
    method: 'put',
    data
  })
}

export function deleteButton(id) {
  return request({
    url: `/buttons/${id}`,
    method: 'delete'
  })
}

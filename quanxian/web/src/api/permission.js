import request from '@/utils/request'

export function getPermissionsTree() {
  return request({
    url: '/permissions/tree',
    method: 'get'
  })
}

export function getPermissionPreview(params) {
  return request({
    url: '/permissions/preview',
    method: 'get',
    params
  })
}

export function batchAssignPermissions(data) {
  return request({
    url: '/permissions/batch-assign',
    method: 'post',
    data
  })
}

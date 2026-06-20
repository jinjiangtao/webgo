import request from '../utils/request'

export const createActivity = (data) => {
  return request({
    url: '/activity',
    method: 'post',
    data
  })
}

export const getActivityList = () => {
  return request({
    url: '/activity',
    method: 'get'
  })
}

export const getActivityDetail = (id) => {
  return request({
    url: `/activity/${id}`,
    method: 'get'
  })
}

export const toggleActivityStatus = (id, status) => {
  return request({
    url: `/activity/${id}/status`,
    method: 'put',
    data: { status }
  })
}

export const deleteActivity = (id) => {
  return request({
    url: `/activity/${id}`,
    method: 'delete'
  })
}

export const getDashboardStats = () => {
  return request({
    url: '/dashboard',
    method: 'get'
  })
}

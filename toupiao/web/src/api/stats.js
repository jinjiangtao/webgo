import request from '../utils/request'

export const getVoteStats = (id) => {
  return request({
    url: `/stats/${id}`,
    method: 'get'
  })
}

export const getVoteRecords = (id, page = 1, pageSize = 20) => {
  return request({
    url: `/stats/${id}/records`,
    method: 'get',
    params: {
      page,
      page_size: pageSize
    }
  })
}

import request from '../utils/request'

export const submitVote = (activityId, optionIds) => {
  return request({
    url: '/vote',
    method: 'post',
    data: {
      activity_id: activityId,
      option_ids: optionIds
    }
  })
}

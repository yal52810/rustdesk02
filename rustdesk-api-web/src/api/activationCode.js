import request from '@/utils/request'

export function list(params) {
  return request({
    url: '/activation_code/list',
    method: 'get',
    params
  })
}

export function create(data) {
  return request({
    url: '/activation_code/create',
    method: 'post',
    data
  })
}

export function batchCreate(data) {
  return request({
    url: '/activation_code/batchCreate',
    method: 'post',
    data
  })
}

export function remove(id) {
  return request({
    url: `/activation_code/delete/${id}`,
    method: 'post'
  })
}

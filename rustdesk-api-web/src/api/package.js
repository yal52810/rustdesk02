import request from '@/utils/request'

export function list (params) {
  return request({
    url: '/package/list',
    params,
  })
}

export function detail (id) {
  return request({
    url: `/package/detail/${id}`,
  })
}

export function create (data) {
  return request({
    url: '/package/create',
    method: 'post',
    data,
  })
}

export function update (data) {
  return request({
    url: '/package/update',
    method: 'post',
    data,
  })
}

export function remove (data) {
  return request({
    url: '/package/delete',
    method: 'post',
    data,
  })
}

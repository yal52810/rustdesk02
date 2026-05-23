import request from '@/utils/request'

export function list (params) {
  return request({
    url: '/server/list',
    params,
  })
}

export function detail (id) {
  return request({
    url: `/server/detail/${id}`,
  })
}

export function create (data) {
  return request({
    url: '/server/create',
    method: 'post',
    data,
  })
}

export function update (data) {
  return request({
    url: '/server/update',
    method: 'post',
    data,
  })
}

export function remove (data) {
  return request({
    url: '/server/delete',
    method: 'post',
    data,
  })
}

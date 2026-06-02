import request from '@/utils/request'

export function server () {
  return request({
    url: '/config/server',
    method: 'get',
  })
}

export function app () {
  return request({
    url: '/config/app',
    method: 'get',
  })
}

export function admin () {
  return request({
    url: '/config/admin',
    method: 'get',
  })
}

export function mail () {
  return request({
    url: '/config/mail',
    method: 'get',
  })
}

export function updateMail (data) {
  return request({
    url: '/config/mail',
    method: 'post',
    data,
  })
}

export function testMail (data) {
  return request({
    url: '/config/mail/test',
    method: 'post',
    data,
  })
}

export function getServerConfig () {
  return request({
    url: '/config/server',
    method: 'get',
  })
}

export function updateServerKey (data) {
  return request({
    url: '/config/server-key',
    method: 'post',
    data,
  })
}

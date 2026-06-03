import axios from 'axios'
import { getToken } from '@/utils/auth'

function rawPortalRequest (config) {
  const token = getToken()
  const headers = {
    ...(config.headers || {}),
  }
  if (token) {
    headers['api-token'] = token
  }
  return axios({
    baseURL: '/api',
    timeout: 50000,
    withCredentials: true,
    ...config,
    headers,
  }).then(res => res.data)
}

export function listVipServers () {
  return rawPortalRequest({
    url: '/vip/servers',
    method: 'get',
  })
}

export function listVipPackages () {
  return rawPortalRequest({
    url: '/vip/packages',
    method: 'get',
  })
}

export function getClientConfig () {
  return rawPortalRequest({
    url: '/vip/client-config',
    method: 'get',
  })
}

export function vipRegister (data) {
  return rawPortalRequest({
    url: '/register',
    method: 'post',
    data,
  })
}

export function redeemActivationCode (data) {
  return rawPortalRequest({
    url: '/vip/redeem',
    method: 'post',
    data,
  })
}

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
    baseURL: import.meta.env.VITE_SERVER_API,
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

export function redeemActivationCode (data) {
  return rawPortalRequest({
    url: '/vip/redeem',
    method: 'post',
    data,
  })
}

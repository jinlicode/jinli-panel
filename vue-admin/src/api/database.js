import request from '@/utils/request'

export function fetchList(query) {
  return request({
    url: '/database/list',
    method: 'get',
    params: query
  })
}

export function fetchSite(id) {
  return request({
    url: '/database/detail',
    method: 'get',
    params: { id }
  })
}

export function createSite(data) {
  return request({
    url: '/database/create',
    method: 'post',
    data
  })
}

export function updateSite(data) {
  return request({
    url: '/database/update',
    method: 'post',
    data
  })
}

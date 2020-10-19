import request from '@/utils/request'

export function fetchList(query) {
  return request({
    url: '/site/list',
    method: 'get',
    params: query
  })
}

export function fetchSite(id) {
  return request({
    url: '/site/detail',
    method: 'get',
    params: { id }
  })
}

export function createSite(data) {
  return request({
    url: '/site/create',
    method: 'post',
    data
  })
}

export function updateSite(data) {
  return request({
    url: '/site/update',
    method: 'post',
    data
  })
}

export function deleteSite(data) {
  return request({
    url: '/site/delete',
    method: 'post',
    data
  })
}

import request from '@/utils/request'

export function fetchList(query) {
  return request({
    url: '/vue-element-admin/site/list',
    method: 'get',
    params: query
  })
}

export function fetchSite(id) {
  return request({
    url: '/vue-element-admin/site/detail',
    method: 'get',
    params: { id }
  })
}

export function createSite(data) {
  return request({
    url: '/vue-element-admin/site/create',
    method: 'post',
    data
  })
}

export function updateSite(data) {
  return request({
    url: '/vue-element-admin/site/update',
    method: 'post',
    data
  })
}

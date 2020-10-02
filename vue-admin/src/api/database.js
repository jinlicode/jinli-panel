import request from '@/utils/request'

export function fetchList(query) {
  return request({
    url: '/vue-element-admin/database/list',
    method: 'get',
    params: query
  })
}

export function fetchSite(id) {
  return request({
    url: '/vue-element-admin/database/detail',
    method: 'get',
    params: { id }
  })
}

export function createSite(data) {
  return request({
    url: '/vue-element-admin/database/create',
    method: 'post',
    data
  })
}

export function updateSite(data) {
  return request({
    url: '/vue-element-admin/database/update',
    method: 'post',
    data
  })
}

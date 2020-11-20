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

export function deleteSite(data) {
  return request({
    url: '/site/delete',
    method: 'post',
    data
  })
}

export function getSiteConf(id) {
  return request({
    url: '/site/get_conf',
    method: 'get',
    params: { id }
  })
}

export function updateSiteConf(data) {
  return request({
    url: '/site/update_conf',
    method: 'post',
    data
  })
}

export function getSiteRewrite(id) {
  return request({
    url: '/site/get_rewrite',
    method: 'get',
    params: { id }
  })
}

export function updateSiteRewrite(data) {
  return request({
    url: '/site/update_rewrite',
    method: 'post',
    data
  })
}

export function getSitePhp(id) {
  return request({
    url: '/site/get_php',
    method: 'get',
    params: { id }
  })
}

export function updateSitePhp(data) {
  return request({
    url: '/site/update_php',
    method: 'post',
    data
  })
}

export function getSiteDomain(id) {
  return request({
    url: '/site/get_domain',
    method: 'get',
    params: { id }
  })
}

export function updateSiteDomain(data) {
  return request({
    url: '/site/update_domain',
    method: 'post',
    data
  })
}

export function delSiteDomain(data) {
  return request({
    url: '/site/del_domain',
    method: 'post',
    data
  })
}

export function getSiteBasepath(id) {
  return request({
    url: '/site/get_basepath',
    method: 'get',
    params: { id }
  })
}

export function updateSiteBasepath(data) {
  return request({
    url: '/site/update_basepath',
    method: 'post',
    data
  })
}

export function updateSiteStatus(data) {
  return request({
    url: '/site/update_status',
    method: 'post',
    data
  })
}

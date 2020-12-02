import request from '@/utils/request'

export function fetchList() {
  return request({
    url: '/soft/list',
    method: 'get'
  })
}

export function getPHPList() {
  return request({
    url: '/soft/phplist',
    method: 'get'
  })
}

export function InstallData() {
  return request({
    url: '/soft/install',
    method: 'post'
  })
}

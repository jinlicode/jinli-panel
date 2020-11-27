import request from '@/utils/request'

export function getPHPList() {
  return request({
    url: '/soft/phplist',
    method: 'get'
  })
}

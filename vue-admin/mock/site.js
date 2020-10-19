const Mock = require('mockjs')

const List = []
const count = 100
const phpVersion = ['5.6', '7.0', '7.1', '7.2', '7.3', '5.6-sec', '7.0-sec', '7.1-sec', '7.2-sec', '7.3-sec']

for (let i = 0; i < count; i++) {
  List.push(Mock.mock({
    id: '@increment',
    php_version: phpVersion[i % 10],
    domain: '@domain',
    email: '@email',
    status: '@integer(0,1,2)',
    is_ssl: '@integer(0,1)'
  }))
}

module.exports = [
  {
    url: '/site/list',
    type: 'get',
    response: config => {
      const { importance, type, title, page = 1, limit = 20, sort } = config.query

      let mockList = List.filter(item => {
        if (importance && item.importance !== +importance) return false
        if (type && item.type !== type) return false
        if (title && item.title.indexOf(title) < 0) return false
        return true
      })

      if (sort === '-id') {
        mockList = mockList.reverse()
      }

      const pageList = mockList.filter((item, index) => index < limit * page && index >= limit * (page - 1))

      return {
        code: 20000,
        data: {
          total: mockList.length,
          list: pageList
        }
      }
    }
  },
  {
    url: '/site/create',
    type: 'post',
    response: _ => {
      return {
        code: 20000,
        data: 'success'
      }
    }
  },
  {
    url: '/site/update',
    type: 'post',
    response: _ => {
      return {
        code: 20000,
        data: 'success'
      }
    }
  }
]


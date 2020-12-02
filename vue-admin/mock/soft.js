const Mock = require('mockjs')

const List = []
const count = 20

for (let i = 0; i < count; i++) {
  List.push(Mock.mock({
    id: '@increment',
    name: '@name',
    desc: '@name',
    status: '@integer(0,1,2)'
  }))
}

module.exports = [
  {
    url: '/soft/list',
    type: 'get',
    response: config => {
      const { importance, type, title, sort } = config.query

      let mockList = List.filter(item => {
        if (importance && item.importance !== +importance) return false
        if (type && item.type !== type) return false
        if (title && item.title.indexOf(title) < 0) return false
        return true
      })

      if (sort === '-id') {
        mockList = mockList.reverse()
      }

      return {
        code: 20000,
        data: {
          list: mockList
        }
      }
    }
  },
  {
    url: '/soft/phplist',
    type: 'get',
    response: config => {
      return {
        code: 20000,
        data: {
          list: [
            { name: 'php:v5.6', desc: 'php-5.6' },
            { name: 'php:v5.6-sec', desc: 'php-5.6 安全版本' },
            { name: 'php:v7.0', desc: 'php-7.0' },
            { name: 'php:v7.0-sec', desc: 'php-7.0 安全版本' },
            { name: 'php:v7.1', desc: 'php-7.1' },
            { name: 'php:v7.1-sec', desc: 'php-7.1 安全版本' },
            { name: 'php:v7.2', desc: 'php-7.2' },
            { name: 'php:v7.2-sec', desc: 'php-7.2 安全版本' },
            { name: 'php:v7.3', desc: 'php-7.3' },
            { name: 'php:v7.3-sec', desc: 'php-7.3 安全版本' }
          ]
        }
      }
    }
  },
  {
    url: '/soft/install',
    type: 'post',
    response: _ => {
      return {
        code: 20000,
        data: 'success'
      }
    }
  }
]


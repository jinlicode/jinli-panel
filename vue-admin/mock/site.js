const Mock = require('mockjs')

const List = []
const count = 100
const phpVersion = ['php:v5.6', 'php:v7.0', 'php:v7.1', 'php:v7.2', 'php:v7.3', 'php:v5.6-sec', 'php:v7.0-sec', 'php:v7.1-sec', 'php:v7.2-sec', 'php:v7.3-sec']

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
  },
  {
    url: '/site/delete',
    type: 'post',
    response: _ => {
      return {
        code: 20000,
        data: 'success'
      }
    }
  },
  {
    url: '/site/get_conf',
    type: 'get',
    response: _ => {
      return {
        code: 20000,
        data: {
          text: '这个是get_conf内容'
        }
      }
    }
  },
  {
    url: '/site/update_conf',
    type: 'post',
    response: _ => {
      return {
        code: 20000,
        data: 'success'
      }
    }
  },
  {
    url: '/site/get_rewrite',
    type: 'get',
    response: _ => {
      return {
        code: 20000,
        data: {
          text: '这个是get_rewrite内容'
        }
      }
    }
  },
  {
    url: '/site/update_rewrite',
    type: 'post',
    response: _ => {
      return {
        code: 20000,
        data: 'success'
      }
    }
  },
  {
    url: '/site/get_php',
    type: 'get',
    response: _ => {
      return {
        code: 20000,
        data: {
          text: '7.1'
        }
      }
    }
  },
  {
    url: '/site/update_php',
    type: 'post',
    response: _ => {
      return {
        code: 20000,
        data: 'success'
      }
    }
  },
  {
    url: '/site/get_domain',
    type: 'get',
    response: _ => {
      return {
        code: 20000,
        data: {
          list: [
            { id: 1, name:"www.baidu.com", pid: 1 },
            { id: 2, name:"www.baidu2.com", pid: 1 },
            { id: 3, name:"www.baidu3.com", pid: 1 }
          ]
        }
      }
    }
  },
  {
    url: '/site/update_domain',
    type: 'post',
    response: _ => {
      return {
        code: 20000,
        data: 'success'
      }
    }
  },
  {
    url: '/site/del_domain',
    type: 'post',
    response: _ => {
      return {
        code: 20000,
        data: 'success'
      }
    }
  },
  {
    url: '/site/get_basepath',
    type: 'get',
    response: _ => {
      return {
        code: 20000,
        data: {
          list: [
            '/',
            '/public',
            '/conf'
          ]
        }
      }
    }
  },
  {
    url: '/site/update_basepath',
    type: 'post',
    response: _ => {
      return {
        code: 20000,
        data: 'success'
      }
    }
  },
  {
    url: '/site/update_status',
    type: 'post',
    response: _ => {
      return {
        code: 20000,
        data: 'success'
      }
    }
  }
]


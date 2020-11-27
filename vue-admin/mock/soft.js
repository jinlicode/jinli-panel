module.exports = [
  {
    url: '/soft/phplist',
    type: 'get',
    response: config => {
      return {
        code: 20000,
        data: {
          list: [
            { key: '5.6', display_name: 'php-5.6' },
            { key: '5.6-sec', display_name: 'php-5.6 安全版本' },
            { key: '7.0', display_name: 'php-7.0' },
            { key: '7.0-sec', display_name: 'php-7.0 安全版本' },
            { key: '7.1', display_name: 'php-7.1' },
            { key: '7.1-sec', display_name: 'php-7.1 安全版本' },
            { key: '7.2', display_name: 'php-7.2' },
            { key: '7.2-sec', display_name: 'php-7.2 安全版本' },
            { key: '7.3', display_name: 'php-7.3' },
            { key: '7.3-sec', display_name: 'php-7.3 安全版本' }
          ]
        }
      }
    }
  }
]


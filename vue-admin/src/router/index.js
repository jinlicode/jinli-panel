import Vue from 'vue'
import Router from 'vue-router'

Vue.use(Router)

/* Layout */
import Layout from '@/layout'

/**
 * Note: sub-menu only appear when route children.length >= 1
 * Detail see: https://panjiachen.github.io/vue-element-admin-site/guide/essentials/router-and-nav.html
 *
 * hidden: true                   if set true, item will not show in the sidebar(default is false)
 * alwaysShow: true               if set true, will always show the root menu
 *                                if not set alwaysShow, when item has more than one children route,
 *                                it will becomes nested mode, otherwise not show the root menu
 * redirect: noRedirect           if set noRedirect will no redirect in the breadcrumb
 * name:'router-name'             the name is used by <keep-alive> (must set!!!)
 * meta : {
    roles: ['admin','editor']    control the page roles (you can set multiple roles)
    title: 'title'               the name show in sidebar and breadcrumb (recommend set)
    icon: 'svg-name'/'el-icon-x' the icon show in the sidebar
    noCache: true                if set true, the page will no be cached(default is false)
    affix: true                  if set true, the tag will affix in the tags-view
    breadcrumb: false            if set false, the item will hidden in breadcrumb(default is true)
    activeMenu: '/example/list'  if set path, the sidebar will highlight the path you set
  }
 */

/**
 * constantRoutes
 * a base page that does not have permission requirements
 * all roles can be accessed
 */
export const constantRoutes = [
  {
    path: '/redirect',
    component: Layout,
    hidden: true,
    children: [
      {
        path: '/redirect/:path(.*)',
        component: () => import('@/views/redirect/index')
      }
    ]
  },
  {
    path: '/login',
    component: () => import('@/views/login/index'),
    hidden: true
  },
  {
    path: '/auth-redirect',
    component: () => import('@/views/login/auth-redirect'),
    hidden: true
  },
  {
    path: '/',
    component: Layout,
    redirect: '/dashboard',
    children: [
      {
        path: 'dashboard',
        component: () => import('@/views/dashboard/index'),
        name: 'Dashboard',
        meta: { title: '统计', icon: 'dashboard', affix: true }
      }
    ]
  },
  {
    path: '/profile',
    component: Layout,
    redirect: '/profile/index',
    hidden: true,
    children: [
      {
        path: 'index',
        component: () => import('@/views/profile/index'),
        name: 'Profile',
        meta: { title: 'Profile', icon: 'user', noCache: true }
      }
    ]
  }
]

/**
 * asyncRoutes
 * the routes that need to be dynamically loaded based on user roles
 */
export const asyncRoutes = [
  {
    path: '/site',
    component: Layout,
    redirect: '/site/list',
    name: 'Site',
    meta: {
      title: '网站',
      icon: 'el-icon-s-help'
    },
    children: [
      {
        path: 'create',
        name: 'CreateSite',
        meta: { title: '创建网站', icon: 'edit' },
        hidden: true
      },
      {
        path: 'edit/:id(\\d+)',
        name: 'EditSite',
        meta: { title: '修改网站', noCache: true, activeMenu: '/site/list' },
        hidden: true
      },
      {
        path: 'list',
        component: () => import('@/views/site/list'),
        name: 'SiteList',
        meta: { title: '网站列表', icon: 'list' }
      },
      {
        path: 'delete',
        name: 'DeleteSite',
        meta: { title: '删除网站', icon: 'del' },
        hidden: true
      },
      {
        path: 'get_conf',
        name: 'GETSiteConf',
        meta: { title: '获取配置', icon: 'edit' },
        hidden: true
      },
      {
        path: 'update_conf',
        name: 'UpdateSiteConf',
        meta: { title: '更新配置', icon: 'edit' },
        hidden: true
      },
      {
        path: 'get_rewrite',
        name: 'GETSiteRewrite',
        meta: { title: '获取伪静态', icon: 'edit' },
        hidden: true
      },
      {
        path: 'update_rewrite',
        name: 'UpdateSiteRewrite',
        meta: { title: '更新伪静态', icon: 'edit' },
        hidden: true
      },
      {
        path: 'get_php',
        name: 'GETSitePhp',
        meta: { title: '获取php', icon: 'edit' },
        hidden: true
      },
      {
        path: 'update_php',
        name: 'UpdateSitePhp',
        meta: { title: '设置php版本', icon: 'edit' },
        hidden: true
      },
      {
        path: 'get_domain',
        name: 'GETSiteDomain',
        meta: { title: '获取域名', icon: 'edit' },
        hidden: true
      },
      {
        path: 'update_domain',
        name: 'UpdateSiteDomain',
        meta: { title: '设置域名', icon: 'edit' },
        hidden: true
      },
      {
        path: 'del_domain',
        name: 'DelSiteDomain',
        meta: { title: '删除域名', icon: 'edit' },
        hidden: true
      },
      {
        path: 'get_basepath',
        name: 'GETSiteBasepath',
        meta: { title: '获取根目录', icon: 'edit' },
        hidden: true
      },
      {
        path: 'update_basepath',
        name: 'UpdateSiteBasepath',
        meta: { title: '设置根目录', icon: 'edit' },
        hidden: true
      },
      {
        path: 'update_status',
        name: 'UpdateSiteStatus',
        meta: { title: '设置网站状态', icon: 'edit' },
        hidden: true
      }
    ]
  },
  {
    path: '/database',
    component: Layout,
    redirect: '/database/list',
    name: 'Database',
    meta: {
      title: '数据库',
      icon: 'el-icon-s-help'
    },
    children: [
      {
        path: 'create',
        name: 'CreateDatabase',
        meta: { title: '创建数据库', icon: 'edit' },
        hidden: true
      },
      {
        path: 'edit/:id(\\d+)',
        name: 'EditDatabase',
        meta: { title: '修改数据库', noCache: true, activeMenu: '/database/list' },
        hidden: true
      },
      {
        path: 'list',
        component: () => import('@/views/database/list'),
        name: 'DatabaseList',
        meta: { title: '数据库列表', icon: 'list' }
      }
    ]
  },

  {
    path: '/setting',
    component: Layout,
    children: [
      {
        path: 'index',
        name: 'Tab',
        meta: { title: '配置', icon: 'tab' }
      }
    ]
  },
  // 404 page must be placed at the end !!!
  { path: '*', redirect: '/404', hidden: true }
]

const createRouter = () => new Router({
  // mode: 'history', // require service support
  scrollBehavior: () => ({ y: 0 }),
  routes: constantRoutes
})

const router = createRouter()

// Detail see: https://github.com/vuejs/vue-router/issues/1234#issuecomment-357941465
export function resetRouter() {
  const newRouter = createRouter()
  router.matcher = newRouter.matcher // reset router
}

export default router

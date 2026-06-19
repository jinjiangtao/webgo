import { createRouter, createWebHashHistory } from 'vue-router'
import { useUserStore } from '@/stores/user'

const Layout = () => import('@/layout/index.vue')

const constantRoutes = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/login/index.vue'),
    meta: { hidden: true }
  },
  {
    path: '/404',
    name: '404',
    component: () => import('@/views/error/404.vue'),
    meta: { hidden: true }
  }
]

const router = createRouter({
  history: createWebHashHistory(),
  routes: constantRoutes
})

const WHITE_LIST = ['/login', '/404']

router.beforeEach(async (to, from, next) => {
  const userStore = useUserStore()
  const token = userStore.token

  if (token) {
    if (to.path === '/login') {
      next('/')
    } else {
      if (!userStore.userInfo || !userStore.userInfo.id) {
        try {
          await userStore.getUserInfo()
          const dynamicRoutes = generateDynamicRoutes(userStore.menus)
          dynamicRoutes.forEach(route => {
            router.addRoute(route)
          })
          router.addRoute({ path: '/:pathMatch(.*)*', redirect: '/404' })
          next({ ...to, replace: true })
        } catch (error) {
          userStore.logout()
          next('/login')
        }
      } else {
        next()
      }
    }
  } else {
    if (WHITE_LIST.includes(to.path)) {
      next()
    } else {
      next('/login')
    }
  }
})

function generateDynamicRoutes(menus) {
  const routes = []
  const topRoute = {
    path: '/',
    component: Layout,
    redirect: '/dashboard',
    children: [
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: () => import('@/views/dashboard/index.vue'),
        meta: { title: '首页', icon: 'HomeFilled' }
      }
    ]
  }

  function processMenus(menuList, parentPath = '') {
    menuList.forEach(menu => {
      if (menu.children && menu.children.length > 0) {
        processMenus(menu.children, menu.path)
      } else {
        const routePath = menu.path.replace(/^\//, '')
        const componentPath = menu.component
        const route = {
          path: routePath,
          name: `Menu_${menu.id}`,
          component: () => import(`@/views/${componentPath}.vue`).catch(() => import('@/views/error/404.vue')),
          meta: {
            title: menu.name,
            icon: menu.icon,
            menuId: menu.id
          }
        }
        topRoute.children.push(route)
      }
    })
  }

  processMenus(menus)
  routes.push(topRoute)
  return routes
}

export default router

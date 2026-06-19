import { createRouter, createWebHashHistory } from 'vue-router'
import { useUserStore } from '@/stores/user'

const Layout = () => import('@/layout/index.vue')

export const constantRoutes = [
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
  },
  {
    path: '/',
    name: 'Layout',
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
]

const router = createRouter({
  history: createWebHashHistory(),
  routes: constantRoutes
})

const WHITE_LIST = ['/login', '/404']

let dynamicRoutesAdded = false

router.beforeEach(async (to, from, next) => {
  const userStore = useUserStore()
  const token = userStore.token

  if (token) {
    if (to.path === '/login') {
      next('/dashboard')
    } else {
      if (!dynamicRoutesAdded) {
        try {
          await userStore.getUserInfo()
          const menuRoutes = generateMenuRoutes(userStore.menus)
          menuRoutes.forEach(route => {
            router.addRoute('Layout', route)
          })
          router.addRoute({
            path: '/:pathMatch(.*)*',
            redirect: '/404'
          })
          dynamicRoutesAdded = true
          next({ ...to, replace: true })
        } catch (error) {
          console.error('获取用户信息失败:', error)
          userStore.logout()
          dynamicRoutesAdded = false
          next('/login')
        }
      } else {
        next()
      }
    }
  } else {
    dynamicRoutesAdded = false
    if (WHITE_LIST.includes(to.path)) {
      next()
    } else {
      next('/login')
    }
  }
})

function generateMenuRoutes(menus) {
  const routes = []

  function findLeafMenus(menuList) {
    menuList.forEach(menu => {
      if (menu.children && menu.children.length > 0) {
        findLeafMenus(menu.children)
      } else {
        let routePath = menu.path
        if (routePath.startsWith('/')) {
          routePath = routePath.substring(1)
        }
        const componentPath = menu.component
        const route = {
          path: routePath,
          name: `Menu_${menu.id}`,
          component: () => {
            return import(`@/views/${componentPath}.vue`).catch(err => {
              console.warn(`无法加载组件: ${componentPath}`, err)
              return import('@/views/error/404.vue')
            })
          },
          meta: {
            title: menu.name,
            icon: menu.icon,
            menuId: menu.id
          }
        }
        routes.push(route)
      }
    })
  }

  findLeafMenus(menus)
  return routes
}

export function resetRouter() {
  dynamicRoutesAdded = false
  const newRouter = createRouter({
    history: createWebHashHistory(),
    routes: constantRoutes
  })
  router.matcher = newRouter.matcher
}

export default router

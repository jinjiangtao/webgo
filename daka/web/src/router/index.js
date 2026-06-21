import { createRouter, createWebHashHistory } from 'vue-router'

const routes = [
  {
    path: '/',
    name: 'home',
    component: () => import('../views/Home.vue')
  },
  {
    path: '/stats',
    name: 'stats',
    component: () => import('../views/Stats.vue')
  }
]

const router = createRouter({
  history: createWebHashHistory(),
  routes
})

export default router

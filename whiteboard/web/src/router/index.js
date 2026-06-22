import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  {
    path: '/',
    name: 'Whiteboard',
    component: () => import('@/views/WhiteboardView.vue')
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router

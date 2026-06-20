import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  {
    path: '/',
    name: 'Home',
    component: () => import('../views/Home.vue')
  },
  {
    path: '/create',
    name: 'CreateActivity',
    component: () => import('../views/CreateActivity.vue')
  },
  {
    path: '/vote/:id',
    name: 'Vote',
    component: () => import('../views/Vote.vue')
  },
  {
    path: '/records/:id',
    name: 'VoteRecords',
    component: () => import('../views/VoteRecords.vue')
  },
  {
    path: '/admin',
    name: 'Admin',
    component: () => import('../views/Admin.vue')
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router

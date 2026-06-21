import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'
import Dashboard from '@/pages/Dashboard.vue'
import ReportCenter from '@/pages/ReportCenter.vue'
import DataTrace from '@/pages/DataTrace.vue'
import SnapshotList from '@/pages/SnapshotList.vue'

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    redirect: '/dashboard',
  } as RouteRecordRaw,
  {
    path: '/dashboard',
    name: 'dashboard',
    component: Dashboard,
    meta: { title: '数据大屏', icon: '📊' },
  },
  {
    path: '/reports',
    name: 'reports',
    component: ReportCenter,
    meta: { title: '报表中心', icon: '📋' },
  },
  {
    path: '/trace',
    name: 'trace',
    component: DataTrace,
    meta: { title: '数据溯源', icon: '🔍' },
  },
  {
    path: '/snapshots',
    name: 'snapshots',
    component: SnapshotList,
    meta: { title: '快照管理', icon: '📸' },
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

router.beforeEach((to, _from, next) => {
  document.title = to.meta.title ? `${to.meta.title} - 企业数据分析平台` : '企业数据分析平台'
  next()
})

export default router

import { createRouter, createWebHistory } from 'vue-router'
import { useUserStore } from '@/stores/user'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      component: () => import('@/views/layout/Layout.vue'),
      redirect: '/dashboard',
      meta: { requiresAuth: true },
      children: [
        {
          path: 'dashboard',
          name: 'dashboard',
          component: () => import('@/views/DashboardView.vue')
        },
        {
          path: 'stats/:formId',
          name: 'stats',
          component: () => import('@/views/StatsView.vue')
        }
      ]
    },
    // --- 编辑器路由 ---
    {
      path: '/editor/new',
      name: 'newForm',
      meta: { requiresAuth: true },
      component: () => import('@/views/EditorView.vue')
    },
    {
      path: '/editor/:formId',
      name: 'editForm',
      meta: { requiresAuth: true },
      component: () => import('@/views/EditorView.vue')
    },
    // --- 独立页面 ---
    {
      path: '/login',
      name: 'login',
      component: () => import('@/views/LoginView.vue')
    },
    {
      path: '/form/:formKey',
      name: 'formFill',
      component: () => import('@/views/FormFillView.vue')
    },
    {
      path: '/preview',
      name: 'preview',
      component: () => import('@/views/PreviewView.vue')
    },
    // 新增提交成功页面路由
    {
      path: '/success',
      name: 'success',
      component: () => import('@/views/SuccessView.vue')
    }
  ]
})

// 导航守卫 beforeEach
router.beforeEach((to, from, next) => {
  const userStore = useUserStore()
  if (to.meta.requiresAuth) {
    if (!userStore.isAuthenticated) {
      next({
        path: '/login',
        query: { redirect: to.fullPath }
      })
    } else {
      next()
    }
  } else {
    if (to.name === 'login' && userStore.isAuthenticated) {
      next({ path: '/dashboard' })
    } else {
      next()
    }
  }
})

export default router

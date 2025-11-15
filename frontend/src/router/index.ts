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
        },
        // 个人中心路由
        {
          path: 'profile',
          name: 'profile',
          component: () => import('@/views/ProfileView.vue')
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
  // 检查路由是否需要认证
  if (to.meta.requiresAuth) {
    // 如果需要认证，但用户未登录
    if (!userStore.isAuthenticated) {
      // 重定向到登录页
      next({
        path: '/login',
        // 保存用户想访问的页面路径，以便登录后重定向回去
        query: { redirect: to.fullPath }
      })
    } else {
      // 用户已登录，正常放行
      next()
    }
  } else {
    // 如果路由不需要认证
    // 特殊处理：如果用户已登录，但访问的是登录页，则直接跳转到首页
    if (to.name === 'login' && userStore.isAuthenticated) {
      next({ path: '/dashboard' })
    } else {
      // 其他情况（未登录访问公开页，或已登录访问公开页），正常放行
      next()
    }
  }
})

export default router

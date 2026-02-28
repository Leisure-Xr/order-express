import { createRouter, createWebHistory } from 'vue-router'
import { adminRoutes } from './admin'
import { customerRoutes } from './customer'
import { useAuthStore } from '@/stores/auth'
import { useAppStore } from '@/stores/app'
import { i18n } from '@/i18n'

declare module 'vue-router' {
  interface RouteMeta {
    requiresAuth?: boolean
    requiresAdmin?: boolean
    title?: string
    layout?: 'admin' | 'customer' | 'blank'
    showBack?: boolean
  }
}

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    ...customerRoutes,
    ...adminRoutes,
    {
      path: '/:pathMatch(.*)*',
      name: 'NotFound',
      component: () => import('@/views/NotFoundView.vue'),
      meta: { title: 'routes.notFound', layout: 'blank' },
    },
  ],
  scrollBehavior(_to, _from, savedPosition) {
    if (savedPosition) {
      return savedPosition
    }
    return { top: 0 }
  },
})

router.beforeEach(async (to, _from) => {
  const authStore = useAuthStore()
  const appStore = useAppStore()
  const t = i18n.global.t

  // Initialise auth once per app load (validates token if present)
  if (!authStore.initialized) {
    await authStore.initAuth()
  }

  // Capture table number from QR code links like /menu?table=T05
  const tableQuery = to.query.table
  const tableId = typeof tableQuery === 'string' ? tableQuery.trim() : ''
  if (tableId) {
    appStore.setTableFromQuery(tableId)
  }

  // Set document title
  const titleKey = to.meta.title
  if (titleKey) {
    const translatedTitle = t(titleKey)
    const appName = t('common.appName', 'Order Express')
    document.title = `${translatedTitle} - ${appName}`
  } else {
    document.title = 'Order Express'
  }

  // Check authentication for admin routes
  if (to.meta.requiresAuth && !authStore.isAuthenticated) {
    return {
      name: 'AdminLogin',
      query: { redirect: to.fullPath },
    }
  }

  // Check admin role for admin routes
  if (to.meta.requiresAdmin && authStore.user?.role !== 'admin') {
    return {
      name: 'AdminLogin',
      query: { redirect: to.fullPath },
    }
  }

  // If already logged in and going to admin login, redirect to dashboard
  if (to.name === 'AdminLogin' && authStore.isAuthenticated && authStore.user?.role === 'admin') {
    return { name: 'AdminDashboard' }
  }
})

export default router

import type { RouteRecordRaw } from 'vue-router'

export const adminRoutes: RouteRecordRaw[] = [
  {
    path: '/admin/login',
    name: 'AdminLogin',
    component: () => import('@/views/admin/LoginView.vue'),
    meta: { layout: 'blank', title: 'routes.adminLogin', showBack: false },
  },
  {
    path: '/admin',
    meta: { requiresAuth: true, requiresAdmin: true },
    redirect: '/admin/dashboard',
  },
  {
    path: '/admin/dashboard',
    name: 'AdminDashboard',
    component: () => import('@/views/admin/DashboardView.vue'),
    meta: { title: 'routes.dashboard', layout: 'admin', requiresAuth: true, requiresAdmin: true },
  },
  {
    path: '/admin/menu',
    name: 'AdminMenuList',
    component: () => import('@/views/admin/menu/MenuListView.vue'),
    meta: { title: 'routes.menuManagement', layout: 'admin', requiresAuth: true, requiresAdmin: true },
  },
  {
    path: '/admin/menu/categories',
    name: 'AdminCategories',
    component: () => import('@/views/admin/menu/CategoryManageView.vue'),
    meta: { title: 'routes.categoryManagement', layout: 'admin', requiresAuth: true, requiresAdmin: true },
  },
  {
    path: '/admin/menu/dish/create',
    name: 'AdminDishCreate',
    component: () => import('@/views/admin/menu/DishEditView.vue'),
    meta: { title: 'routes.addDish', layout: 'admin', requiresAuth: true, requiresAdmin: true },
  },
  {
    path: '/admin/menu/dish/:id/edit',
    name: 'AdminDishEdit',
    component: () => import('@/views/admin/menu/DishEditView.vue'),
    meta: { title: 'routes.editDish', layout: 'admin', requiresAuth: true, requiresAdmin: true },
    props: true,
  },
  {
    path: '/admin/orders',
    name: 'AdminOrderList',
    component: () => import('@/views/admin/order/OrderListView.vue'),
    meta: { title: 'routes.orderManagement', layout: 'admin', requiresAuth: true, requiresAdmin: true },
  },
  {
    path: '/admin/orders/:id',
    name: 'AdminOrderDetail',
    component: () => import('@/views/admin/order/OrderDetailView.vue'),
    meta: { title: 'routes.orderDetail', layout: 'admin', requiresAuth: true, requiresAdmin: true },
    props: true,
  },
  {
    path: '/admin/tables',
    name: 'AdminTableManage',
    component: () => import('@/views/admin/table/TableManageView.vue'),
    meta: { title: 'routes.tableManagement', layout: 'admin', requiresAuth: true, requiresAdmin: true },
  },
  {
    path: '/admin/settings',
    name: 'AdminStoreSettings',
    component: () => import('@/views/admin/settings/StoreSettingsView.vue'),
    meta: { title: 'routes.storeSettings', layout: 'admin', requiresAuth: true, requiresAdmin: true },
  },
]

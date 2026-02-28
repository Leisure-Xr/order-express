import type { RouteRecordRaw } from 'vue-router'

export const customerRoutes: RouteRecordRaw[] = [
  {
    path: '/',
    name: 'CustomerHome',
    component: () => import('@/views/customer/HomeView.vue'),
    meta: { title: 'routes.home', layout: 'customer' },
  },
  {
    path: '/menu',
    name: 'CustomerMenu',
    component: () => import('@/views/customer/MenuView.vue'),
    meta: { title: 'routes.menu', layout: 'customer' },
  },
  {
    path: '/dish/:id',
    name: 'CustomerDishDetail',
    component: () => import('@/views/customer/DishDetailView.vue'),
    meta: { title: 'routes.dishDetail', layout: 'customer' },
    props: true,
  },
  {
    path: '/cart',
    name: 'CustomerCart',
    component: () => import('@/views/customer/CartView.vue'),
    meta: { title: 'routes.cart', layout: 'customer' },
  },
  {
    path: '/orders',
    name: 'CustomerOrderHistory',
    component: () => import('@/views/customer/OrderHistoryView.vue'),
    meta: { title: 'routes.orderHistory', layout: 'customer' },
  },
  {
    path: '/checkout',
    name: 'CustomerCheckout',
    component: () => import('@/views/customer/CheckoutView.vue'),
    meta: { title: 'routes.checkout', layout: 'blank' },
  },
  {
    path: '/payment/:orderId',
    name: 'CustomerPayment',
    component: () => import('@/views/customer/PaymentView.vue'),
    meta: { title: 'routes.payment', layout: 'blank' },
    props: true,
  },
  {
    path: '/payment-result/:orderId',
    name: 'CustomerPaymentResult',
    component: () => import('@/views/customer/PaymentResultView.vue'),
    meta: { title: 'routes.paymentResult', layout: 'blank', showBack: false },
    props: true,
  },
  {
    path: '/order/:id/tracking',
    name: 'CustomerOrderTracking',
    component: () => import('@/views/customer/OrderTrackingView.vue'),
    meta: { title: 'routes.orderTracking', layout: 'blank' },
    props: true,
  },
]

import { createApp } from 'vue'
import { createPinia } from 'pinia'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'

import App from './App.vue'
import router from './router'
import { i18n } from './i18n'
import { configureApiClient } from './api/client'
import { useAuthStore } from '@/stores/auth'

import './styles/index.scss'
import './styles/element-plus.scss'
import './styles/admin.scss'
import './styles/customer.scss'
import './styles/responsive.scss'

const app = createApp(App)

const pinia = createPinia()

app.use(pinia)
app.use(router)
app.use(i18n)
app.use(ElementPlus)

configureApiClient({
  getToken: () => useAuthStore(pinia).token,
  onUnauthorized: async () => {
    const authStore = useAuthStore(pinia)
    await authStore.logout()

    const current = router.currentRoute.value
    if (typeof current.path === 'string' && current.path.startsWith('/admin') && current.name !== 'AdminLogin') {
      await router.replace({
        name: 'AdminLogin',
        query: { redirect: current.fullPath },
      })
    }
  },
})

app.mount('#app')

<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import { reactive } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { useAuthStore } from '@/stores/auth'

const { t } = useI18n()
const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

const formModel = reactive({
  username: import.meta.env.DEV ? 'admin' : '',
  password: import.meta.env.DEV ? 'admin123' : '',
  rememberMe: true,
})

const rules = {
  username: [{ required: true, message: t('auth.username'), trigger: 'blur' }],
  password: [{ required: true, message: t('auth.password'), trigger: 'blur' }],
}

async function handleLogin() {
  const ok = await authStore.login(
    { username: formModel.username, password: formModel.password },
    { rememberMe: formModel.rememberMe },
  )

  if (!ok) {
    ElMessage.error(t('auth.loginFailed'))
    return
  }

  ElMessage.success(t('auth.loginSuccess'))

  const redirectRaw = typeof route.query.redirect === 'string' ? route.query.redirect : ''
  const redirect = redirectRaw.startsWith('/admin') ? redirectRaw : '/admin/dashboard'
  router.replace(redirect)
}
</script>

<template>
  <div class="login-page">
    <el-card class="login-card" shadow="hover">
      <template #header>
        <div class="login-title">{{ t('routes.adminLogin') }}</div>
      </template>

      <el-form :model="formModel" :rules="rules" label-position="top">
        <el-form-item :label="t('auth.username')" prop="username">
          <el-input v-model="formModel.username" autocomplete="username" />
        </el-form-item>
        <el-form-item :label="t('auth.password')" prop="password">
          <el-input v-model="formModel.password" type="password" autocomplete="current-password" show-password />
        </el-form-item>

        <div class="login-options">
          <el-checkbox v-model="formModel.rememberMe">
            {{ t('auth.rememberMe') }}
          </el-checkbox>
        </div>

        <el-button
          type="primary"
          class="login-btn"
          :loading="authStore.loginLoading"
          @click="handleLogin"
        >
          {{ t('auth.login') }}
        </el-button>
      </el-form>
    </el-card>
  </div>
</template>

<style scoped lang="scss">
.login-page {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  padding: 24px;

  background: radial-gradient(circle at 20% 20%, rgba(231, 76, 60, 0.14), transparent 45%),
    radial-gradient(circle at 80% 0%, rgba(52, 152, 219, 0.12), transparent 50%),
    radial-gradient(circle at 50% 100%, rgba(243, 156, 18, 0.08), transparent 40%),
    #f5f5f5;
}

.login-card {
  width: 100%;
  max-width: 420px;
  border-radius: 18px;
  backdrop-filter: blur(12px);
  background: rgba(255, 255, 255, 0.92);
  border: 1px solid rgba(255, 255, 255, 0.5);
  box-shadow: 0 16px 48px rgba(0, 0, 0, 0.10);
  animation: fadeSlideUp 0.5s cubic-bezier(0.4, 0, 0.2, 1);
}

.login-title {
  font-size: 22px;
  font-weight: 800;
  background: linear-gradient(135deg, #303133, #e74c3c);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.login-options {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.login-btn {
  width: 100%;
  height: 44px;
  font-size: 15px;
  font-weight: 700;
}
</style>

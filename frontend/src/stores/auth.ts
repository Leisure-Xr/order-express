import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { LoginCredentials, User } from '@/types'
import { getCurrentUserApi, loginApi, logoutApi } from '@/api/auth'
import {
  type AuthStorageKind,
  readAuthFromStorage,
  setAuthToStorage,
  clearAuthFromStorage,
  clearAllAuthStorage,
} from '@/utils/authStorage'

export const useAuthStore = defineStore('auth', () => {
  const cached = readAuthFromStorage()

  const initialized = ref(false)
  const storageKind = ref<AuthStorageKind | null>(cached.kind)
  const token = ref<string | null>(cached.token)
  const user = ref<User | null>(cached.user)
  const loginLoading = ref(false)

  const isAuthenticated = computed(() => !!token.value && !!user.value)
  const isAdmin = computed(() => user.value?.role === 'admin')
  const userName = computed(() => user.value?.name ?? '')

  async function initAuth(): Promise<void> {
    if (initialized.value) return

    const { token: storedToken, kind } = readAuthFromStorage()
    if (!storedToken || !kind) {
      token.value = null
      user.value = null
      storageKind.value = null
      initialized.value = true
      return
    }

    const res = await getCurrentUserApi(storedToken)
    if (res.code !== 200 || !res.data) {
      token.value = null
      user.value = null
      storageKind.value = null
      clearAllAuthStorage()
      initialized.value = true
      return
    }

    token.value = storedToken
    user.value = res.data
    storageKind.value = kind
    setAuthToStorage(kind, storedToken, res.data)

    // Ensure the other storage doesn't contain stale auth
    clearAuthFromStorage(kind === 'local' ? sessionStorage : localStorage)
    initialized.value = true
  }

  async function login(credentials: LoginCredentials, options?: { rememberMe?: boolean }): Promise<boolean> {
    loginLoading.value = true
    try {
      const res = await loginApi(credentials.username, credentials.password)
      if (res.code !== 200 || !res.data) {
        return false
      }

      const rememberMe = options?.rememberMe ?? true
      const kind: AuthStorageKind = rememberMe ? 'local' : 'session'

      token.value = res.data.token
      user.value = res.data.user
      storageKind.value = kind
      initialized.value = true

      setAuthToStorage(kind, res.data.token, res.data.user)
      clearAuthFromStorage(kind === 'local' ? sessionStorage : localStorage)

      return true
    } finally {
      loginLoading.value = false
    }
  }

  async function logout(): Promise<void> {
    const currentToken = token.value
    token.value = null
    user.value = null
    storageKind.value = null
    initialized.value = true
    clearAllAuthStorage()

    if (!currentToken) return
    try {
      await logoutApi(currentToken)
    } catch {
      // ignore mock failures
    }
  }

  return {
    initialized,
    storageKind,
    token,
    user,
    loginLoading,
    isAuthenticated,
    isAdmin,
    userName,
    initAuth,
    login,
    logout,
  }
})

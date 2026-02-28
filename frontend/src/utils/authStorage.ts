import type { User } from '@/types'

export const AUTH_TOKEN_KEY = 'auth_token'
export const AUTH_USER_KEY = 'auth_user'

export type AuthStorageKind = 'session' | 'local'

function safeParseUser(raw: string | null): User | null {
  if (!raw) return null
  try {
    return JSON.parse(raw) as User
  } catch {
    return null
  }
}

export function readAuthFromStorage(): {
  token: string | null
  user: User | null
  kind: AuthStorageKind | null
} {
  const sessionToken = sessionStorage.getItem(AUTH_TOKEN_KEY)
  if (sessionToken) {
    return {
      token: sessionToken,
      user: safeParseUser(sessionStorage.getItem(AUTH_USER_KEY)),
      kind: 'session',
    }
  }

  const localToken = localStorage.getItem(AUTH_TOKEN_KEY)
  if (localToken) {
    return {
      token: localToken,
      user: safeParseUser(localStorage.getItem(AUTH_USER_KEY)),
      kind: 'local',
    }
  }

  return { token: null, user: null, kind: null }
}

export function setAuthToStorage(kind: AuthStorageKind, token: string, user: User) {
  const storage = kind === 'local' ? localStorage : sessionStorage
  storage.setItem(AUTH_TOKEN_KEY, token)
  storage.setItem(AUTH_USER_KEY, JSON.stringify(user))
}

export function clearAuthFromStorage(storage: Storage) {
  storage.removeItem(AUTH_TOKEN_KEY)
  storage.removeItem(AUTH_USER_KEY)
}

export function clearAllAuthStorage() {
  clearAuthFromStorage(localStorage)
  clearAuthFromStorage(sessionStorage)
}


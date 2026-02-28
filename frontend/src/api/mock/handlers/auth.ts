import type { User } from '@/types'
import { initialUsers } from '@/api/mock/data/users'

const users = [...initialUsers]

function generateToken(userId: string): string {
  const ts = Date.now().toString(36)
  const rand = Math.random().toString(36).substring(2, 10)
  return `token_${userId}_${ts}_${rand}`
}

function parseToken(token: string): string | null {
  const parts = token.split('_')
  if (parts.length < 3) return null
  if (parts[0] !== 'token') return null
  return parts[1] || null
}

function login(username: string, password: string): { token: string; user: User } | { error: string } {
  const found = users.find(
    (u) => u.user.username === username && u.password === password,
  )

  if (!found) {
    return { error: 'Invalid username or password' }
  }

  const token = generateToken(found.user.id)

  return { token, user: { ...found.user } }
}

function logout(token: string): boolean {
  // Stateless token: nothing to invalidate in mock mode
  return !!token
}

function getCurrentUser(token: string): User | null {
  const userId = parseToken(token)
  if (!userId) return null

  const found = users.find((u) => u.user.id === userId)
  return found ? { ...found.user } : null
}

export const authHandlers = {
  login,
  logout,
  getCurrentUser,
}

import { mockDelay, successResponse, errorResponse } from './index'
import type { ApiResponse } from './index'
import type { User } from '@/types'
import { authHandlers } from './mock/handlers/auth'
import { apiRequest, isMockApiEnabled } from './client'
import { endpoints } from './endpoints'

export interface LoginResult {
  token: string
  user: User
}

export async function loginApi(username: string, password: string): Promise<ApiResponse<LoginResult | null>> {
  if (isMockApiEnabled()) {
    await mockDelay()

    if (!username || !password) {
      return errorResponse(400, 'Username and password are required')
    }

    const result = authHandlers.login(username, password)

    if ('error' in result) {
      return errorResponse(401, result.error)
    }

    return successResponse<LoginResult>(result)
  }

  if (!username || !password) {
    return errorResponse(400, 'Username and password are required')
  }

  return apiRequest<LoginResult>({
    path: endpoints.auth.login,
    method: 'POST',
    auth: false,
    body: { username, password },
  })
}

export async function logoutApi(token: string): Promise<ApiResponse<null>> {
  if (isMockApiEnabled()) {
    await mockDelay()

    authHandlers.logout(token)
    return successResponse(null)
  }

  if (!token) {
    return errorResponse(401, 'Token is required')
  }

  return apiRequest<null>({
    path: endpoints.auth.logout,
    method: 'POST',
    token,
    handleUnauthorized: false,
  })
}

export async function getCurrentUserApi(token: string): Promise<ApiResponse<User | null>> {
  if (isMockApiEnabled()) {
    await mockDelay()

    if (!token) {
      return errorResponse(401, 'Token is required')
    }

    const user = authHandlers.getCurrentUser(token)
    if (!user) {
      return errorResponse(401, 'Invalid or expired token')
    }

    return successResponse(user)
  }

  if (!token) {
    return errorResponse(401, 'Token is required')
  }

  return apiRequest<User>({
    path: endpoints.auth.me,
    token,
  })
}

import type { ApiResponse } from './index'
import { successResponse } from './index'

export interface ApiClientConfig {
  getToken?: () => string | null | undefined
  onUnauthorized?: () => void | Promise<void>
}

export type HttpMethod = 'GET' | 'POST' | 'PUT' | 'PATCH' | 'DELETE'

export interface ApiRequestOptions {
  path: string
  method?: HttpMethod
  query?: Record<string, string | number | boolean | null | undefined> | object
  body?: unknown
  headers?: HeadersInit
  signal?: AbortSignal
  auth?: boolean
  token?: string | null
  handleUnauthorized?: boolean
}

let apiClientConfig: ApiClientConfig = {}
let unauthorizedHandling: Promise<void> | null = null

export function configureApiClient(config: ApiClientConfig) {
  apiClientConfig = config
}

export function isMockApiEnabled(): boolean {
  return import.meta.env.VITE_USE_MOCK_API === 'true'
}

function buildUrl(path: string, query?: ApiRequestOptions['query']): string {
  if (!query) return path
  const params = new URLSearchParams()
  for (const [key, value] of Object.entries(query as Record<string, unknown>)) {
    if (value === undefined || value === null) continue
    params.set(key, String(value))
  }
  const qs = params.toString()
  return qs ? `${path}?${qs}` : path
}

function isApiEnvelope(value: unknown): value is ApiResponse<unknown> {
  if (!value || typeof value !== 'object') return false
  const maybe = value as Record<string, unknown>
  return typeof maybe.code === 'number' && 'data' in maybe && typeof maybe.message === 'string'
}

async function triggerUnauthorizedHandler() {
  if (!apiClientConfig.onUnauthorized) return

  if (!unauthorizedHandling) {
    unauthorizedHandling = Promise.resolve(apiClientConfig.onUnauthorized())
      .catch(() => undefined)
      .finally(() => {
        unauthorizedHandling = null
      })
  }

  await unauthorizedHandling
}

function normalizeError<T>(code: number, message: string): ApiResponse<T> {
  return { code, data: null, message }
}

export async function apiRequest<T>(options: ApiRequestOptions): Promise<ApiResponse<T>> {
  const method: HttpMethod = options.method ?? 'GET'
  const url = buildUrl(options.path, options.query)

  const headers = new Headers(options.headers)
  if (!headers.has('Accept')) {
    headers.set('Accept', 'application/json')
  }

  if (options.auth !== false) {
    const token = options.token ?? apiClientConfig.getToken?.()
    if (token) {
      headers.set('Authorization', `Bearer ${token}`)
    }
  }

  const init: RequestInit = {
    method,
    headers,
    signal: options.signal,
  }

  if (options.body !== undefined && options.body !== null) {
    if (options.body instanceof FormData) {
      init.body = options.body
    } else {
      if (!headers.has('Content-Type')) {
        headers.set('Content-Type', 'application/json')
      }
      init.body = JSON.stringify(options.body)
    }
  }

  try {
    const res = await fetch(url, init)

    if (res.status === 401 && options.handleUnauthorized !== false) {
      await triggerUnauthorizedHandler()
    }

    if (res.status === 204) {
      return successResponse(null as unknown as T)
    }

    const contentType = res.headers.get('content-type') ?? ''
    let payload: unknown = null

    if (contentType.includes('application/json')) {
      payload = await res.json().catch(() => null)
    } else {
      payload = await res.text().catch(() => '')
    }

    if (isApiEnvelope(payload)) {
      if (payload.code === 401 && options.handleUnauthorized !== false) {
        await triggerUnauthorizedHandler()
      }
      return payload as ApiResponse<T>
    }

    if (res.ok) {
      return successResponse(payload as T)
    }

    const message =
      typeof payload === 'string'
        ? payload || res.statusText || 'Request failed'
        : ((payload as any)?.message as string | undefined) || res.statusText || 'Request failed'

    return normalizeError<T>(res.status, message)
  } catch (e: any) {
    return normalizeError<T>(0, e?.message || 'Network error')
  }
}

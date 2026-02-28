export function mockDelay(ms: number = 300 + Math.random() * 500): Promise<void> {
  return new Promise((resolve) => setTimeout(resolve, ms))
}

export interface ApiResponse<T> {
  code: number
  data: T | null
  message: string
}

export function successResponse<T>(data: T): ApiResponse<T> {
  return { code: 200, data, message: 'success' }
}

export function errorResponse(code: number, message: string): ApiResponse<null> {
  return { code, data: null, message }
}

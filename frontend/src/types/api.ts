export interface ApiResponse<T> {
  code: number
  data: T | null
  message: string
}

export interface PaginatedResult<T> {
  items: T[]
  total: number
  page: number
  pageSize: number
}

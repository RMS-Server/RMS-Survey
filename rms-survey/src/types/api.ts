// API response wrapper
export interface ApiResponse<T> {
  code: number
  message: string
  data: T
}

// Pagination response matching backend PageResponse[T]
export interface PageResponse<T> {
  list: T[]
  total: number
}

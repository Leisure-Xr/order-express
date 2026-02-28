export type UserRole = 'admin' | 'customer'

export interface User {
  id: string
  username: string
  name: string
  role: UserRole
  phone?: string
  avatar?: string
}

export interface LoginCredentials {
  username: string
  password: string
}

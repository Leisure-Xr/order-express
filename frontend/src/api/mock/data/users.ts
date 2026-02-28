import type { User } from '@/types'

export interface UserWithPassword {
  user: User
  password: string
}

export const initialUsers: UserWithPassword[] = [
  {
    user: {
      id: '1',
      username: 'admin',
      name: '管理员',
      role: 'admin',
      phone: '13800138000',
      avatar: 'https://placehold.co/100x100/e74c3c/white?text=Admin',
    },
    password: 'admin123',
  },
  {
    user: {
      id: '2',
      username: 'customer1',
      name: '张三',
      role: 'customer',
      phone: '13900139001',
      avatar: 'https://placehold.co/100x100/3498db/white?text=ZS',
    },
    password: '123456',
  },
  {
    user: {
      id: '3',
      username: 'customer2',
      name: '李四',
      role: 'customer',
      phone: '13900139002',
      avatar: 'https://placehold.co/100x100/2ecc71/white?text=LS',
    },
    password: '123456',
  },
  {
    user: {
      id: '4',
      username: 'customer3',
      name: '王五',
      role: 'customer',
      phone: '13900139003',
      avatar: 'https://placehold.co/100x100/9b59b6/white?text=WW',
    },
    password: '123456',
  },
]

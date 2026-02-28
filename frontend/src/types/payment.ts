import type { PaymentMethod, PaymentStatus } from './order'

export interface Payment {
  paymentId: string
  orderId: string
  method: PaymentMethod
  status: PaymentStatus
  amount: number
  createdAt: string
  paidAt?: string
}


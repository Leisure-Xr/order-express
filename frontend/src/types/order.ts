export type OrderStatus = 'pending' | 'confirmed' | 'preparing' | 'ready' | 'completed' | 'delivered' | 'cancelled'
export type OrderType = 'dine_in' | 'takeaway' | 'pickup'
export type PaymentMethod = 'wechat' | 'alipay' | 'cash'
export type PaymentStatus = 'unpaid' | 'processing' | 'paid' | 'refunded'

export interface OrderItem {
  dishId: string
  dishName: { zh: string; en: string }
  quantity: number
  unitPrice: number
  selectedOptions: { optionName: string; valueName: string; priceAdjustment: number }[]
  subtotal: number
}

export interface StatusHistoryEntry {
  status: OrderStatus
  timestamp: string
  note?: string
}

export interface Order {
  id: string
  orderNumber: string
  type: OrderType
  status: OrderStatus
  tableId?: string
  items: OrderItem[]
  subtotal: number
  deliveryFee: number
  discount: number
  total: number
  remarks?: string
  payment: {
    method: PaymentMethod
    status: PaymentStatus
    paidAt?: string
  }
  deliveryAddress?: string
  contactPhone?: string
  estimatedDeliveryTime?: string
  statusHistory: StatusHistoryEntry[]
  createdAt: string
  updatedAt: string
}

export interface CreateOrderPayload {
  type: OrderType
  tableId?: string
  items: { dishId: string; quantity: number; selectedOptions: { optionName: string; valueName: string; priceAdjustment: number }[] }[]
  remarks?: string
  deliveryAddress?: string
  contactPhone?: string
  paymentMethod: PaymentMethod
}

export interface OrderQueryParams {
  status?: OrderStatus | 'all'
  type?: OrderType
  page?: number
  pageSize?: number
}

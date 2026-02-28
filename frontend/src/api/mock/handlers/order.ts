import { ref } from 'vue'
import type {
  Order,
  OrderStatus,
  OrderQueryParams,
  CreateOrderPayload,
  OrderItem,
  StatusHistoryEntry,
} from '@/types'
import type { PaginatedResult } from '@/types'
import type { PaymentMethod } from '@/types'
import { initialOrders } from '@/api/mock/data/orders'
import { menuHandlers } from './menu'

// Reactive order list seeded from initial data, persists across the session
const orders = ref<Order[]>(JSON.parse(JSON.stringify(initialOrders)))

let orderCounter = 200

function nextOrderId(): string {
  orderCounter++
  return `order-${orderCounter}`
}

function nextOrderNumber(): string {
  orderCounter // already incremented
  const datePart = new Date().toISOString().slice(0, 10).replace(/-/g, '')
  return `OE${datePart}${String(orderCounter).padStart(4, '0')}`
}

function now(): string {
  return new Date().toISOString()
}

// Valid status transitions
const VALID_TRANSITIONS: Record<string, string[]> = {
  pending: ['confirmed', 'cancelled'],
  confirmed: ['preparing', 'cancelled'],
  preparing: ['ready', 'cancelled'],
  ready: ['completed', 'delivered'],
  completed: [],
  delivered: [],
  cancelled: [],
}

// -----------------------------------------------------------
// Auto-progression timers (for demo purposes)
// -----------------------------------------------------------
const autoTimers = new Map<string, ReturnType<typeof setTimeout>[]>()

function scheduleAutoProgression(orderId: string): void {
  const timers: ReturnType<typeof setTimeout>[] = []

  // 3s: pending -> confirmed
  timers.push(
    setTimeout(() => {
      const order = orders.value.find((o) => o.id === orderId)
      if (order && order.status === 'pending') {
        applyStatusChange(order, 'confirmed', 'Auto-confirmed by system')
      }
    }, 3000),
  )

  // 3s + 8s = 11s: confirmed -> preparing
  timers.push(
    setTimeout(() => {
      const order = orders.value.find((o) => o.id === orderId)
      if (order && order.status === 'confirmed') {
        applyStatusChange(order, 'preparing', 'Kitchen started preparation')
      }
    }, 11000),
  )

  // 3s + 8s + 20s = 31s: preparing -> ready
  timers.push(
    setTimeout(() => {
      const order = orders.value.find((o) => o.id === orderId)
      if (order && order.status === 'preparing') {
        applyStatusChange(order, 'ready', 'Order is ready for pickup')
      }
    }, 31000),
  )

  autoTimers.set(orderId, timers)
}

function cancelAutoProgression(orderId: string): void {
  const timers = autoTimers.get(orderId)
  if (timers) {
    timers.forEach(clearTimeout)
    autoTimers.delete(orderId)
  }
}

function applyStatusChange(order: Order, newStatus: OrderStatus, note?: string): void {
  order.status = newStatus
  order.updatedAt = now()

  const entry: StatusHistoryEntry = {
    status: newStatus,
    timestamp: now(),
    note,
  }
  order.statusHistory.push(entry)
}

// -----------------------------------------------------------
// Handler functions
// -----------------------------------------------------------

function getOrders(params?: OrderQueryParams): PaginatedResult<Order> {
  let result = [...orders.value]

  // Filter by status
  if (params?.status && params.status !== 'all') {
    result = result.filter((o) => o.status === params.status)
  }

  // Filter by type
  if (params?.type) {
    result = result.filter((o) => o.type === params.type)
  }

  // Sort by creation date, newest first
  result.sort((a, b) => new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime())

  // Pagination
  const page = params?.page ?? 1
  const pageSize = params?.pageSize ?? 20
  const total = result.length
  const start = (page - 1) * pageSize
  const items = result.slice(start, start + pageSize)

  return { items, total, page, pageSize }
}

function getOrderById(id: string): Order | null {
  return orders.value.find((o) => o.id === id) ?? null
}

function createOrder(payload: CreateOrderPayload): Order {
  const orderId = nextOrderId()
  const orderNumber = nextOrderNumber()
  const timestamp = now()

  // Resolve order items from the dish catalog
  const items: OrderItem[] = payload.items.map((item) => {
    const dish = menuHandlers.getDishById(item.dishId)
    const basePrice = dish?.price ?? 0
    const optionsTotal = item.selectedOptions.reduce((sum, opt) => sum + opt.priceAdjustment, 0)
    const unitPrice = basePrice + optionsTotal
    const subtotal = unitPrice * item.quantity

    return {
      dishId: item.dishId,
      dishName: dish?.name ?? { zh: '未知菜品', en: 'Unknown Dish' },
      quantity: item.quantity,
      unitPrice,
      selectedOptions: item.selectedOptions,
      subtotal,
    }
  })

  const subtotal = items.reduce((sum, item) => sum + item.subtotal, 0)
  const deliveryFee = payload.type === 'takeaway' ? 5 : 0
  const discount = 0
  const total = subtotal + deliveryFee - discount

  const order: Order = {
    id: orderId,
    orderNumber,
    type: payload.type,
    status: 'pending',
    tableId: payload.tableId,
    items,
    subtotal,
    deliveryFee,
    discount,
    total,
    remarks: payload.remarks,
    payment: {
      method: payload.paymentMethod,
      status: 'unpaid',
    },
    deliveryAddress: payload.deliveryAddress,
    contactPhone: payload.contactPhone,
    statusHistory: [
      {
        status: 'pending',
        timestamp,
        note: 'Order placed',
      },
    ],
    createdAt: timestamp,
    updatedAt: timestamp,
  }

  orders.value.push(order)

  // Schedule auto-progression for demo
  scheduleAutoProgression(orderId)

  return { ...order }
}

function updateOrderStatus(
  id: string,
  newStatus: OrderStatus,
  note?: string,
): { success: true; order: Order } | { success: false; error: string } {
  const order = orders.value.find((o) => o.id === id)
  if (!order) {
    return { success: false, error: 'Order not found' }
  }

  const allowed = VALID_TRANSITIONS[order.status] ?? []
  if (!allowed.includes(newStatus)) {
    return {
      success: false,
      error: `Cannot transition from "${order.status}" to "${newStatus}". Allowed: [${allowed.join(', ')}]`,
    }
  }

  // If cancelling, stop auto-progression
  if (newStatus === 'cancelled') {
    cancelAutoProgression(id)
  }

  applyStatusChange(order, newStatus, note)

  return { success: true, order: { ...order } }
}

function getOrderHistory(): Order[] {
  return orders.value
    .filter((o) => ['completed', 'delivered', 'cancelled'].includes(o.status))
    .sort((a, b) => new Date(b.updatedAt).getTime() - new Date(a.updatedAt).getTime())
}

function markOrderPaid(id: string, method: PaymentMethod, paidAt?: string): Order | null {
  const order = orders.value.find((o) => o.id === id)
  if (!order) return null

  order.payment.method = method
  order.payment.status = 'paid'
  order.payment.paidAt = paidAt ?? now()
  order.updatedAt = now()

  return { ...order }
}

export const orderHandlers = {
  getOrders,
  getOrderById,
  createOrder,
  updateOrderStatus,
  getOrderHistory,
  markOrderPaid,
}

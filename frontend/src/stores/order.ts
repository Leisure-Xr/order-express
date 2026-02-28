import { defineStore } from 'pinia'
import { ref } from 'vue'
import type {
  Order,
  OrderQueryParams,
  CreateOrderPayload,
  OrderStatus,
  PaginatedResult,
} from '@/types'
import {
  getOrdersApi,
  getOrderByIdApi,
  createOrderApi,
  updateOrderStatusApi,
  getOrderHistoryApi,
} from '@/api/order'
import { addOrderIdToHistory } from '@/utils/orderHistory'

function assertOk(code: number, message: string): asserts code is 200 {
  if (code !== 200) throw new Error(message || 'Request failed')
}

export const useOrderStore = defineStore('order', () => {
  const orders = ref<Order[]>([])
  const total = ref(0)
  const page = ref(1)
  const pageSize = ref(20)

  const currentOrder = ref<Order | null>(null)
  const orderHistory = ref<Order[]>([])

  const listLoading = ref(false)
  const detailLoading = ref(false)
  const historyLoading = ref(false)

  async function fetchOrders(params?: OrderQueryParams): Promise<PaginatedResult<Order>> {
    listLoading.value = true
    try {
      const res = await getOrdersApi(params)
      assertOk(res.code, res.message)
      const data = res.data ?? { items: [], total: 0, page: 1, pageSize: 20 }
      orders.value = data.items
      total.value = data.total
      page.value = data.page
      pageSize.value = data.pageSize
      return data
    } finally {
      listLoading.value = false
    }
  }

  async function fetchOrderById(id: string): Promise<Order | null> {
    detailLoading.value = true
    try {
      const res = await getOrderByIdApi(id)
      if (res.code !== 200) {
        currentOrder.value = null
        return null
      }
      currentOrder.value = res.data ?? null
      return currentOrder.value
    } finally {
      detailLoading.value = false
    }
  }

  async function createOrder(payload: CreateOrderPayload): Promise<Order | null> {
    const res = await createOrderApi(payload)
    if (res.code !== 200) throw new Error(res.message || 'Create order failed')
    const order = res.data ?? null
    if (order?.id) {
      addOrderIdToHistory(order.id)
    }
    return order
  }

  async function updateOrderStatus(id: string, status: OrderStatus, note?: string): Promise<Order | null> {
    const res = await updateOrderStatusApi(id, status, note)
    if (res.code !== 200) throw new Error(res.message || 'Update status failed')
    const updated = res.data ?? null
    if (!updated) return null

    const idx = orders.value.findIndex((o) => o.id === id)
    if (idx !== -1) orders.value[idx] = updated
    if (currentOrder.value?.id === id) currentOrder.value = updated
    return updated
  }

  async function fetchOrderHistory(): Promise<Order[]> {
    historyLoading.value = true
    try {
      const res = await getOrderHistoryApi()
      assertOk(res.code, res.message)
      orderHistory.value = res.data ?? []
      return orderHistory.value
    } finally {
      historyLoading.value = false
    }
  }

  return {
    // State
    orders,
    total,
    page,
    pageSize,
    currentOrder,
    orderHistory,
    listLoading,
    detailLoading,
    historyLoading,
    // Actions
    fetchOrders,
    fetchOrderById,
    createOrder,
    updateOrderStatus,
    fetchOrderHistory,
  }
})

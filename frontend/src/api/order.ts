import { mockDelay, successResponse, errorResponse } from './index'
import type { ApiResponse } from './index'
import type { Order, OrderStatus, OrderQueryParams, CreateOrderPayload, PaginatedResult } from '@/types'
import { orderHandlers } from './mock/handlers/order'
import { apiRequest, isMockApiEnabled } from './client'
import { endpoints } from './endpoints'

export async function getOrdersApi(
  params?: OrderQueryParams,
): Promise<ApiResponse<PaginatedResult<Order>>> {
  if (isMockApiEnabled()) {
    await mockDelay()
    return successResponse(orderHandlers.getOrders(params))
  }

  return apiRequest<PaginatedResult<Order>>({
    path: endpoints.orders.list,
    method: 'GET',
    query: params,
  })
}

export async function getOrderByIdApi(id: string): Promise<ApiResponse<Order | null>> {
  if (isMockApiEnabled()) {
    await mockDelay()
    const order = orderHandlers.getOrderById(id)
    if (!order) {
      return errorResponse(404, 'Order not found')
    }
    return successResponse(order)
  }

  return apiRequest<Order>({
    path: endpoints.orders.get(id),
    method: 'GET',
  })
}

export async function createOrderApi(payload: CreateOrderPayload): Promise<ApiResponse<Order | null>> {
  if (isMockApiEnabled()) {
    await mockDelay()

    // Basic validation
    if (!payload.items || payload.items.length === 0) {
      return errorResponse(400, 'Order must contain at least one item')
    }

    if (payload.type === 'dine_in' && !payload.tableId) {
      return errorResponse(400, 'Table ID is required for dine-in orders')
    }

    if (payload.type === 'takeaway') {
      if (!payload.deliveryAddress) return errorResponse(400, 'Delivery address is required for takeaway orders')
      if (!payload.contactPhone) return errorResponse(400, 'Contact phone is required for takeaway orders')
    }

    if (payload.type === 'pickup') {
      if (!payload.contactPhone) return errorResponse(400, 'Contact phone is required for pickup orders')
    }

    const order = orderHandlers.createOrder(payload)
    return successResponse(order)
  }

  // Keep basic client-side validation to fail fast
  if (!payload.items || payload.items.length === 0) {
    return errorResponse(400, 'Order must contain at least one item')
  }

  if (payload.type === 'dine_in' && !payload.tableId) {
    return errorResponse(400, 'Table ID is required for dine-in orders')
  }

  if (payload.type === 'takeaway') {
    if (!payload.deliveryAddress) return errorResponse(400, 'Delivery address is required for takeaway orders')
    if (!payload.contactPhone) return errorResponse(400, 'Contact phone is required for takeaway orders')
  }

  if (payload.type === 'pickup') {
    if (!payload.contactPhone) return errorResponse(400, 'Contact phone is required for pickup orders')
  }

  return apiRequest<Order>({
    path: endpoints.orders.create,
    method: 'POST',
    body: payload,
  })
}

export async function updateOrderStatusApi(
  id: string,
  status: OrderStatus,
  note?: string,
): Promise<ApiResponse<Order | null>> {
  if (isMockApiEnabled()) {
    await mockDelay()

    const result = orderHandlers.updateOrderStatus(id, status, note)

    if (!result.success) {
      return errorResponse(400, result.error)
    }

    return successResponse(result.order)
  }

  return apiRequest<Order>({
    path: endpoints.orders.updateStatus(id),
    method: 'PATCH',
    body: { status, note },
  })
}

export async function getOrderHistoryApi(): Promise<ApiResponse<Order[]>> {
  if (isMockApiEnabled()) {
    await mockDelay()
    return successResponse(orderHandlers.getOrderHistory())
  }

  return apiRequest<Order[]>({
    path: endpoints.orders.history,
    method: 'GET',
  })
}

import { mockDelay, successResponse, errorResponse } from './index'
import type { ApiResponse } from './index'
import type { Payment, PaymentMethod } from '@/types'
import { paymentHandlers } from './mock/handlers/payment'
import { apiRequest, isMockApiEnabled } from './client'
import { endpoints } from './endpoints'

export async function initiatePaymentApi(
  orderId: string,
  method: PaymentMethod,
  amount: number,
): Promise<ApiResponse<Payment | null>> {
  if (isMockApiEnabled()) {
    await mockDelay()

    if (!orderId) {
      return errorResponse(400, 'Order ID is required')
    }

    const result = paymentHandlers.initiatePayment(orderId, method, amount)
    return successResponse<Payment>(result)
  }

  if (!orderId) {
    return errorResponse(400, 'Order ID is required')
  }

  return apiRequest<Payment>({
    path: endpoints.payments.initiate,
    method: 'POST',
    body: { orderId, method, amount },
  })
}

export async function checkPaymentStatusApi(
  paymentId: string,
): Promise<ApiResponse<Payment | null>> {
  if (isMockApiEnabled()) {
    await mockDelay()

    const payment = paymentHandlers.checkPaymentStatus(paymentId)
    if (!payment) {
      return errorResponse(404, 'Payment not found')
    }

    return successResponse(payment)
  }

  return apiRequest<Payment>({
    path: endpoints.payments.get(paymentId),
    method: 'GET',
  })
}

export async function getPaymentsByOrderIdApi(orderId: string): Promise<ApiResponse<Payment[]>> {
  if (isMockApiEnabled()) {
    await mockDelay()
    return successResponse(paymentHandlers.getPaymentsByOrderId(orderId))
  }

  return apiRequest<Payment[]>({
    path: endpoints.payments.byOrderId(orderId),
    method: 'GET',
  })
}

export async function refundPaymentApi(paymentId: string): Promise<ApiResponse<Payment | null>> {
  if (isMockApiEnabled()) {
    await mockDelay()

    const payment = paymentHandlers.refundPayment(paymentId)
    if (!payment) {
      return errorResponse(404, 'Payment not found')
    }

    return successResponse(payment)
  }

  return apiRequest<Payment>({
    path: endpoints.payments.refund(paymentId),
    method: 'POST',
  })
}


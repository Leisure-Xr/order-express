import { ref } from 'vue'
import type { Payment, PaymentMethod } from '@/types'
import { orderHandlers } from './order'

// Reactive state: all payment records
const payments = ref<Payment[]>([])

let paymentCounter = 1000

function nextPaymentId(): string {
  paymentCounter++
  return `pay-${Date.now().toString(36)}-${paymentCounter}`
}

function now(): string {
  return new Date().toISOString()
}

function initiatePayment(orderId: string, method: PaymentMethod, amount: number): Payment {
  const paymentId = nextPaymentId()

  const record: Payment = {
    paymentId,
    orderId,
    method,
    status: 'processing',
    amount,
    createdAt: now(),
  }

  payments.value.push(record)

  // Simulate: auto-complete payment after a short delay
  setTimeout(() => {
    const payment = payments.value.find((p) => p.paymentId === paymentId)
    if (payment && payment.status === 'processing') {
      payment.status = 'paid'
      payment.paidAt = now()
      orderHandlers.markOrderPaid(payment.orderId, payment.method, payment.paidAt)
    }
  }, 2000)

  return { ...record }
}

function checkPaymentStatus(paymentId: string): Payment | null {
  const payment = payments.value.find((p) => p.paymentId === paymentId)
  return payment ? { ...payment } : null
}

function getPaymentsByOrderId(orderId: string): Payment[] {
  return payments.value.filter((p) => p.orderId === orderId).map((p) => ({ ...p }))
}

function refundPayment(paymentId: string): Payment | null {
  const payment = payments.value.find((p) => p.paymentId === paymentId)
  if (!payment) return null

  payment.status = 'refunded'
  return { ...payment }
}

export const paymentHandlers = {
  initiatePayment,
  checkPaymentStatus,
  getPaymentsByOrderId,
  refundPayment,
}

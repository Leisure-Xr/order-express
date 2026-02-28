import { getStorageItem, setStorageItem } from './storage'

const ORDER_HISTORY_IDS_KEY = 'order_history_ids'
const MAX_HISTORY_IDS = 50

export function addOrderIdToHistory(orderId: string): void {
  if (!orderId) return

  const ids = getStorageItem<string[]>(ORDER_HISTORY_IDS_KEY, [])
  const next = [orderId, ...ids.filter((id) => id !== orderId)].slice(0, MAX_HISTORY_IDS)
  setStorageItem(ORDER_HISTORY_IDS_KEY, next)
}

export function getOrderHistoryIds(): string[] {
  return getStorageItem<string[]>(ORDER_HISTORY_IDS_KEY, [])
}

export function clearOrderHistoryIds(): void {
  setStorageItem<string[]>(ORDER_HISTORY_IDS_KEY, [])
}


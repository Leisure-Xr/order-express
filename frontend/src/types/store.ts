export interface StoreInfo {
  name: { zh: string; en: string }
  address: { zh: string; en: string }
  phone: string
  logo: string
  description: { zh: string; en: string }
}

export interface BusinessHours {
  dayOfWeek: number
  openTime: string
  closeTime: string
  isClosed: boolean
}

export interface DeliverySettings {
  enabled: boolean
  minimumOrder: number
  deliveryFee: number
  freeDeliveryThreshold: number
  estimatedMinutes: number
  deliveryRadius: number
}

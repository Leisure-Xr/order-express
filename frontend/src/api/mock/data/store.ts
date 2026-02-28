import type { StoreInfo, BusinessHours, DeliverySettings } from '@/types'

export const initialStoreInfo: StoreInfo = {
  name: { zh: '食光小馆', en: 'Order Express' },
  address: { zh: '北京市朝阳区建国路88号', en: '88 Jianguo Road, Chaoyang, Beijing' },
  phone: '010-88886666',
  logo: '/images/logo.png',
  description: {
    zh: '精选时令食材，匠心烹饪美味',
    en: 'Fresh seasonal ingredients, crafted with care',
  },
}

export const initialBusinessHours: BusinessHours[] = [
  { dayOfWeek: 0, openTime: '', closeTime: '', isClosed: true },
  { dayOfWeek: 1, openTime: '10:00', closeTime: '22:00', isClosed: false },
  { dayOfWeek: 2, openTime: '10:00', closeTime: '22:00', isClosed: false },
  { dayOfWeek: 3, openTime: '10:00', closeTime: '22:00', isClosed: false },
  { dayOfWeek: 4, openTime: '10:00', closeTime: '22:00', isClosed: false },
  { dayOfWeek: 5, openTime: '10:00', closeTime: '23:00', isClosed: false },
  { dayOfWeek: 6, openTime: '10:00', closeTime: '23:00', isClosed: false },
]

export const initialDeliverySettings: DeliverySettings = {
  enabled: true,
  minimumOrder: 30,
  deliveryFee: 5,
  freeDeliveryThreshold: 80,
  estimatedMinutes: 30,
  deliveryRadius: 5,
}

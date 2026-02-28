import { ref } from 'vue'
import type { StoreInfo, BusinessHours, DeliverySettings } from '@/types'
import {
  initialStoreInfo,
  initialBusinessHours,
  initialDeliverySettings,
} from '@/api/mock/data/store'

// Reactive state persists in-session
const storeInfo = ref<StoreInfo>(JSON.parse(JSON.stringify(initialStoreInfo)))
const businessHours = ref<BusinessHours[]>(JSON.parse(JSON.stringify(initialBusinessHours)))
const deliverySettings = ref<DeliverySettings>(JSON.parse(JSON.stringify(initialDeliverySettings)))

function getStoreInfo(): StoreInfo {
  return { ...storeInfo.value }
}

function updateStoreInfo(data: Partial<StoreInfo>): StoreInfo {
  storeInfo.value = { ...storeInfo.value, ...data }
  return { ...storeInfo.value }
}

function getBusinessHours(): BusinessHours[] {
  return [...businessHours.value]
}

function updateBusinessHours(data: BusinessHours[]): BusinessHours[] {
  businessHours.value = [...data]
  return [...businessHours.value]
}

function getDeliverySettings(): DeliverySettings {
  return { ...deliverySettings.value }
}

function updateDeliverySettings(data: Partial<DeliverySettings>): DeliverySettings {
  deliverySettings.value = { ...deliverySettings.value, ...data }
  return { ...deliverySettings.value }
}

export const storeHandlers = {
  getStoreInfo,
  updateStoreInfo,
  getBusinessHours,
  updateBusinessHours,
  getDeliverySettings,
  updateDeliverySettings,
}

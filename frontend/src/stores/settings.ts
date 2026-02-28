import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { StoreInfo, BusinessHours, DeliverySettings } from '@/types'
import {
  getStoreInfoApi,
  updateStoreInfoApi,
  getBusinessHoursApi,
  updateBusinessHoursApi,
  getDeliverySettingsApi,
  updateDeliverySettingsApi,
} from '@/api/store'

function assertOk(code: number, message: string): asserts code is 200 {
  if (code !== 200) throw new Error(message || 'Request failed')
}

export const useSettingsStore = defineStore('settings', () => {
  const storeInfo = ref<StoreInfo | null>(null)
  const businessHours = ref<BusinessHours[]>([])
  const deliverySettings = ref<DeliverySettings | null>(null)
  const loading = ref(false)

  async function fetchStoreInfo(): Promise<StoreInfo> {
    loading.value = true
    try {
      const res = await getStoreInfoApi()
      assertOk(res.code, res.message)
      storeInfo.value = res.data
      return res.data as StoreInfo
    } finally {
      loading.value = false
    }
  }

  async function updateStoreInfo(data: Partial<StoreInfo>): Promise<StoreInfo> {
    loading.value = true
    try {
      const res = await updateStoreInfoApi(data)
      assertOk(res.code, res.message)
      storeInfo.value = res.data
      return res.data as StoreInfo
    } finally {
      loading.value = false
    }
  }

  async function fetchBusinessHours(): Promise<BusinessHours[]> {
    loading.value = true
    try {
      const res = await getBusinessHoursApi()
      assertOk(res.code, res.message)
      businessHours.value = res.data ?? []
      return businessHours.value
    } finally {
      loading.value = false
    }
  }

  async function updateBusinessHours(hours: BusinessHours[]): Promise<BusinessHours[]> {
    loading.value = true
    try {
      const res = await updateBusinessHoursApi(hours)
      assertOk(res.code, res.message)
      businessHours.value = res.data ?? []
      return businessHours.value
    } finally {
      loading.value = false
    }
  }

  async function fetchDeliverySettings(): Promise<DeliverySettings> {
    loading.value = true
    try {
      const res = await getDeliverySettingsApi()
      assertOk(res.code, res.message)
      deliverySettings.value = res.data
      return res.data as DeliverySettings
    } finally {
      loading.value = false
    }
  }

  async function updateDeliverySettings(settings: Partial<DeliverySettings>): Promise<DeliverySettings> {
    loading.value = true
    try {
      const res = await updateDeliverySettingsApi(settings)
      assertOk(res.code, res.message)
      deliverySettings.value = res.data
      return res.data as DeliverySettings
    } finally {
      loading.value = false
    }
  }

  async function fetchAllSettings(): Promise<void> {
    loading.value = true
    try {
      await Promise.all([fetchStoreInfo(), fetchBusinessHours(), fetchDeliverySettings()])
    } finally {
      loading.value = false
    }
  }

  function isStoreOpen(): boolean {
    if (!businessHours.value.length) return true

    const now = new Date()
    const dayOfWeek = now.getDay()
    const currentHour = businessHours.value.find((h) => h.dayOfWeek === dayOfWeek)

    if (!currentHour || currentHour.isClosed) return false

    const timeStr = `${String(now.getHours()).padStart(2, '0')}:${String(now.getMinutes()).padStart(2, '0')}`
    return timeStr >= currentHour.openTime && timeStr <= currentHour.closeTime
  }

  return {
    storeInfo,
    businessHours,
    deliverySettings,
    loading,
    fetchStoreInfo,
    updateStoreInfo,
    fetchBusinessHours,
    updateBusinessHours,
    fetchDeliverySettings,
    updateDeliverySettings,
    fetchAllSettings,
    isStoreOpen,
  }
})


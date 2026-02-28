import { mockDelay, successResponse } from './index'
import type { ApiResponse } from './index'
import type { StoreInfo, BusinessHours, DeliverySettings } from '@/types'
import { storeHandlers } from './mock/handlers/store'
import { apiRequest, isMockApiEnabled } from './client'
import { endpoints } from './endpoints'

export async function getStoreInfoApi(): Promise<ApiResponse<StoreInfo>> {
  if (isMockApiEnabled()) {
    await mockDelay()
    return successResponse(storeHandlers.getStoreInfo())
  }

  return apiRequest<StoreInfo>({
    path: endpoints.store.info,
    method: 'GET',
  })
}

export async function updateStoreInfoApi(
  data: Partial<StoreInfo>,
): Promise<ApiResponse<StoreInfo>> {
  if (isMockApiEnabled()) {
    await mockDelay()
    return successResponse(storeHandlers.updateStoreInfo(data))
  }

  return apiRequest<StoreInfo>({
    path: endpoints.store.info,
    method: 'PATCH',
    body: data,
  })
}

export async function getBusinessHoursApi(): Promise<ApiResponse<BusinessHours[]>> {
  if (isMockApiEnabled()) {
    await mockDelay()
    return successResponse(storeHandlers.getBusinessHours())
  }

  return apiRequest<BusinessHours[]>({
    path: endpoints.store.businessHours,
    method: 'GET',
  })
}

export async function updateBusinessHoursApi(
  data: BusinessHours[],
): Promise<ApiResponse<BusinessHours[]>> {
  if (isMockApiEnabled()) {
    await mockDelay()
    return successResponse(storeHandlers.updateBusinessHours(data))
  }

  return apiRequest<BusinessHours[]>({
    path: endpoints.store.businessHours,
    method: 'PUT',
    body: data,
  })
}

export async function getDeliverySettingsApi(): Promise<ApiResponse<DeliverySettings>> {
  if (isMockApiEnabled()) {
    await mockDelay()
    return successResponse(storeHandlers.getDeliverySettings())
  }

  return apiRequest<DeliverySettings>({
    path: endpoints.store.deliverySettings,
    method: 'GET',
  })
}

export async function updateDeliverySettingsApi(
  data: Partial<DeliverySettings>,
): Promise<ApiResponse<DeliverySettings>> {
  if (isMockApiEnabled()) {
    await mockDelay()
    return successResponse(storeHandlers.updateDeliverySettings(data))
  }

  return apiRequest<DeliverySettings>({
    path: endpoints.store.deliverySettings,
    method: 'PATCH',
    body: data,
  })
}

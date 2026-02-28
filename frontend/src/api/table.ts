import { mockDelay, successResponse, errorResponse } from './index'
import type { ApiResponse } from './index'
import type { Table, TableStatus } from '@/types'
import { tableHandlers } from './mock/handlers/table'
import { apiRequest, isMockApiEnabled } from './client'
import { endpoints } from './endpoints'

export async function getTablesApi(): Promise<ApiResponse<Table[]>> {
  if (isMockApiEnabled()) {
    await mockDelay()
    return successResponse(tableHandlers.getTables())
  }

  return apiRequest<Table[]>({
    path: endpoints.tables.list,
    method: 'GET',
  })
}

export async function getPublicTablesApi(): Promise<ApiResponse<Table[]>> {
  if (isMockApiEnabled()) {
    await mockDelay()
    return successResponse(tableHandlers.getTables())
  }

  return apiRequest<Table[]>({
    path: endpoints.tables.publicList,
    method: 'GET',
    auth: false,
  })
}

export async function getTableByIdApi(id: string): Promise<ApiResponse<Table | null>> {
  if (isMockApiEnabled()) {
    await mockDelay()
    const table = tableHandlers.getTableById(id)
    if (!table) {
      return errorResponse(404, 'Table not found')
    }
    return successResponse(table)
  }

  return apiRequest<Table>({
    path: endpoints.tables.get(id),
    method: 'GET',
  })
}

export async function getTablesByAreaApi(area: string): Promise<ApiResponse<Table[]>> {
  if (isMockApiEnabled()) {
    await mockDelay()
    return successResponse(tableHandlers.getTablesByArea(area))
  }

  return apiRequest<Table[]>({
    path: endpoints.tables.list,
    method: 'GET',
    query: { area },
  })
}

export async function getTablesByStatusApi(status: TableStatus): Promise<ApiResponse<Table[]>> {
  if (isMockApiEnabled()) {
    await mockDelay()
    return successResponse(tableHandlers.getTablesByStatus(status))
  }

  return apiRequest<Table[]>({
    path: endpoints.tables.list,
    method: 'GET',
    query: { status },
  })
}

export async function createTableApi(
  data: Omit<Table, 'id' | 'qrCodeUrl'>,
): Promise<ApiResponse<Table>> {
  if (isMockApiEnabled()) {
    await mockDelay()
    const table = tableHandlers.createTable(data)
    return successResponse(table)
  }

  return apiRequest<Table>({
    path: endpoints.tables.create,
    method: 'POST',
    body: data,
  })
}

export async function updateTableApi(
  id: string,
  data: Partial<Omit<Table, 'id'>>,
): Promise<ApiResponse<Table | null>> {
  if (isMockApiEnabled()) {
    await mockDelay()
    const table = tableHandlers.updateTable(id, data)
    if (!table) {
      return errorResponse(404, 'Table not found')
    }
    return successResponse(table)
  }

  return apiRequest<Table>({
    path: endpoints.tables.update(id),
    method: 'PATCH',
    body: data,
  })
}

export async function deleteTableApi(id: string): Promise<ApiResponse<null>> {
  if (isMockApiEnabled()) {
    await mockDelay()
    const success = tableHandlers.deleteTable(id)
    if (!success) {
      return errorResponse(404, 'Table not found')
    }
    return successResponse(null)
  }

  return apiRequest<null>({
    path: endpoints.tables.delete(id),
    method: 'DELETE',
  })
}

export async function updateTableStatusApi(
  id: string,
  status: TableStatus,
  currentOrderId?: string,
): Promise<ApiResponse<Table | null>> {
  if (isMockApiEnabled()) {
    await mockDelay()
    const table = tableHandlers.updateTableStatus(id, status, currentOrderId)
    if (!table) {
      return errorResponse(404, 'Table not found')
    }
    return successResponse(table)
  }

  return apiRequest<Table>({
    path: endpoints.tables.updateStatus(id),
    method: 'PATCH',
    body: { status, currentOrderId },
  })
}

export async function getAreasApi(): Promise<ApiResponse<string[]>> {
  if (isMockApiEnabled()) {
    await mockDelay()
    return successResponse(tableHandlers.getAreas())
  }

  return apiRequest<string[]>({
    path: endpoints.tables.areas,
    method: 'GET',
  })
}

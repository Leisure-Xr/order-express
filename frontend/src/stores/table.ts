import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { Table, TableStatus } from '@/types'
import {
  getTablesApi,
  createTableApi,
  updateTableApi,
  deleteTableApi,
  updateTableStatusApi,
} from '@/api/table'

function assertOk(code: number, message: string): asserts code is 200 {
  if (code !== 200) throw new Error(message || 'Request failed')
}

export const useTableStore = defineStore('table', () => {
  const tables = ref<Table[]>([])
  const loading = ref(false)

  const availableTables = computed(() => tables.value.filter((t) => t.status === 'available'))
  const occupiedTables = computed(() => tables.value.filter((t) => t.status === 'occupied'))
  const reservedTables = computed(() => tables.value.filter((t) => t.status === 'reserved'))
  const disabledTables = computed(() => tables.value.filter((t) => t.status === 'disabled'))

  const tablesByArea = computed(() => {
    const map = new Map<string, Table[]>()
    for (const table of tables.value) {
      const area = table.area || 'Default'
      const list = map.get(area) ?? []
      list.push(table)
      map.set(area, list)
    }
    for (const [key, list] of map.entries()) {
      list.sort((a, b) => a.number.localeCompare(b.number))
      map.set(key, list)
    }
    return map
  })

  async function fetchTables(): Promise<Table[]> {
    loading.value = true
    try {
      const res = await getTablesApi()
      assertOk(res.code, res.message)
      tables.value = res.data ?? []
      return tables.value
    } finally {
      loading.value = false
    }
  }

  async function createTable(payload: Omit<Table, 'id' | 'qrCodeUrl'>): Promise<Table> {
    const res = await createTableApi(payload)
    assertOk(res.code, res.message)
    if (!res.data) throw new Error('Create table failed')
    tables.value.push(res.data)
    return res.data
  }

  async function updateTable(id: string, payload: Partial<Omit<Table, 'id'>>): Promise<Table | null> {
    const res = await updateTableApi(id, payload)
    assertOk(res.code, res.message)
    const updated = res.data ?? null
    if (!updated) return null
    const idx = tables.value.findIndex((t) => t.id === id)
    if (idx !== -1) tables.value[idx] = updated
    return updated
  }

  async function deleteTable(id: string): Promise<boolean> {
    const res = await deleteTableApi(id)
    if (res.code !== 200) return false
    tables.value = tables.value.filter((t) => t.id !== id)
    return true
  }

  async function updateTableStatus(id: string, status: TableStatus, currentOrderId?: string): Promise<Table | null> {
    const res = await updateTableStatusApi(id, status, currentOrderId)
    assertOk(res.code, res.message)
    const updated = res.data ?? null
    if (!updated) return null
    const idx = tables.value.findIndex((t) => t.id === id)
    if (idx !== -1) tables.value[idx] = updated
    return updated
  }

  function getTableById(id: string): Table | undefined {
    return tables.value.find((t) => t.id === id)
  }

  function getTableByNumber(number: string): Table | undefined {
    return tables.value.find((t) => t.number === number)
  }

  function getOrderLinkForTableNumber(tableNumber: string): string {
    const baseUrl = window.location.origin
    return `${baseUrl}/menu?table=${encodeURIComponent(tableNumber)}`
  }

  return {
    tables,
    loading,
    availableTables,
    occupiedTables,
    reservedTables,
    disabledTables,
    tablesByArea,
    fetchTables,
    createTable,
    updateTable,
    deleteTable,
    updateTableStatus,
    getTableById,
    getTableByNumber,
    getOrderLinkForTableNumber,
  }
})


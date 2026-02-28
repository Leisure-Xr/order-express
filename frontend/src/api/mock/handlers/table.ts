import { ref } from 'vue'
import type { Table, TableStatus } from '@/types'
import { initialTables } from '@/api/mock/data/tables'

// Reactive state persists in-session
const tables = ref<Table[]>(JSON.parse(JSON.stringify(initialTables)))

let idCounter = 100

function nextId(): string {
  idCounter++
  return `table-${idCounter}`
}

function getTables(): Table[] {
  return [...tables.value]
}

function getTableById(id: string): Table | null {
  return tables.value.find((t) => t.id === id) ?? null
}

function getTablesByArea(area: string): Table[] {
  return tables.value.filter((t) => t.area === area)
}

function getTablesByStatus(status: TableStatus): Table[] {
  return tables.value.filter((t) => t.status === status)
}

function createTable(data: Omit<Table, 'id' | 'qrCodeUrl'>): Table {
  const id = nextId()
  const newTable: Table = {
    ...data,
    id,
    qrCodeUrl: `/qr/${id}`,
  }
  tables.value.push(newTable)
  return { ...newTable }
}

function updateTable(id: string, data: Partial<Omit<Table, 'id'>>): Table | null {
  const index = tables.value.findIndex((t) => t.id === id)
  if (index === -1) return null

  const existing = tables.value[index]
  if (!existing) return null

  tables.value[index] = { ...existing, ...data }
  return { ...tables.value[index]! }
}

function deleteTable(id: string): boolean {
  const index = tables.value.findIndex((t) => t.id === id)
  if (index === -1) return false

  tables.value.splice(index, 1)
  return true
}

function updateTableStatus(id: string, status: TableStatus, currentOrderId?: string): Table | null {
  const table = tables.value.find((t) => t.id === id)
  if (!table) return null

  table.status = status
  table.currentOrderId = currentOrderId
  return { ...table }
}

function getAreas(): string[] {
  const areaSet = new Set(tables.value.map((t) => t.area).filter(Boolean) as string[])
  return [...areaSet]
}

export const tableHandlers = {
  getTables,
  getTableById,
  getTablesByArea,
  getTablesByStatus,
  createTable,
  updateTable,
  deleteTable,
  updateTableStatus,
  getAreas,
}

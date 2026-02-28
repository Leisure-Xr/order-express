export type TableStatus = 'available' | 'occupied' | 'reserved' | 'disabled'

export interface Table {
  id: string
  number: string
  seats: number
  status: TableStatus
  currentOrderId?: string
  qrCodeUrl?: string
  area?: string
}

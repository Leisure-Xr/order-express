export interface Category {
  id: string
  name: { zh: string; en: string }
  icon?: string
  image?: string
  sortOrder: number
  status: 'active' | 'inactive'
  dishCount?: number
}

export interface DishOptionValue {
  id: string
  label: { zh: string; en: string }
  priceAdjustment: number
}

export interface DishOption {
  id: string
  name: { zh: string; en: string }
  values: DishOptionValue[]
  required: boolean
  multiSelect: boolean
}

export interface Dish {
  id: string
  categoryId: string
  name: { zh: string; en: string }
  description: { zh: string; en: string }
  price: number
  originalPrice?: number
  image: string
  images?: string[]
  status: 'on_sale' | 'off_sale' | 'sold_out'
  options: DishOption[]
  tags: string[]
  preparationTime?: number
  createdAt: string
  updatedAt: string
}

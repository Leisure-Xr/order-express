import { ref } from 'vue'
import type { Category, Dish } from '@/types'
import { initialCategories } from '@/api/mock/data/categories'
import { initialDishes } from '@/api/mock/data/dishes'

// Reactive state so mutations persist in-session
const categories = ref<Category[]>(JSON.parse(JSON.stringify(initialCategories)))
const dishes = ref<Dish[]>(JSON.parse(JSON.stringify(initialDishes)))

let idCounter = 1000

function nextId(prefix: string): string {
  idCounter++
  return `${prefix}-${idCounter}`
}

function now(): string {
  return new Date().toISOString()
}

// -----------------------------------------------------------
// Category operations
// -----------------------------------------------------------

function getCategories(): Category[] {
  return [...categories.value].sort((a, b) => a.sortOrder - b.sortOrder)
}

function getCategoryById(id: string): Category | null {
  return categories.value.find((c) => c.id === id) ?? null
}

function createCategory(data: Omit<Category, 'id' | 'dishCount'>): Category {
  const newCategory: Category = {
    ...data,
    id: nextId('cat'),
    dishCount: 0,
  }
  categories.value.push(newCategory)
  return { ...newCategory }
}

function updateCategory(id: string, data: Partial<Omit<Category, 'id'>>): Category | null {
  const index = categories.value.findIndex((c) => c.id === id)
  if (index === -1) return null

  const existing = categories.value[index]
  if (!existing) return null

  categories.value[index] = { ...existing, ...data }
  return { ...categories.value[index]! }
}

function deleteCategory(id: string): boolean {
  const index = categories.value.findIndex((c) => c.id === id)
  if (index === -1) return false

  // Also remove any dishes in that category
  dishes.value = dishes.value.filter((d) => d.categoryId !== id)
  categories.value.splice(index, 1)
  return true
}

function reorderCategories(orderedIds: string[]): Category[] {
  orderedIds.forEach((id, idx) => {
    const cat = categories.value.find((c) => c.id === id)
    if (cat) {
      cat.sortOrder = idx + 1
    }
  })
  return getCategories()
}

// -----------------------------------------------------------
// Dish operations
// -----------------------------------------------------------

function getDishes(params?: { categoryId?: string; status?: string }): Dish[] {
  let result = [...dishes.value]

  if (params?.categoryId) {
    result = result.filter((d) => d.categoryId === params.categoryId)
  }
  if (params?.status) {
    result = result.filter((d) => d.status === params.status)
  }

  return result
}

function getDishById(id: string): Dish | null {
  return dishes.value.find((d) => d.id === id) ?? null
}

function createDish(data: Omit<Dish, 'id' | 'createdAt' | 'updatedAt'>): Dish {
  const timestamp = now()
  const newDish: Dish = {
    ...data,
    id: nextId('dish'),
    createdAt: timestamp,
    updatedAt: timestamp,
  }
  dishes.value.push(newDish)

  // Increment category dish count
  const cat = categories.value.find((c) => c.id === data.categoryId)
  if (cat) {
    cat.dishCount = (cat.dishCount ?? 0) + 1
  }

  return { ...newDish }
}

function updateDish(id: string, data: Partial<Omit<Dish, 'id' | 'createdAt'>>): Dish | null {
  const index = dishes.value.findIndex((d) => d.id === id)
  if (index === -1) return null

  const existing = dishes.value[index]
  if (!existing) return null

  dishes.value[index] = {
    ...existing,
    ...data,
    updatedAt: now(),
  }

  return { ...dishes.value[index]! }
}

function deleteDish(id: string): boolean {
  const index = dishes.value.findIndex((d) => d.id === id)
  if (index === -1) return false

  const dish = dishes.value[index]
  if (!dish) return false

  // Decrement category dish count
  const cat = categories.value.find((c) => c.id === dish.categoryId)
  if (cat && cat.dishCount && cat.dishCount > 0) {
    cat.dishCount--
  }

  dishes.value.splice(index, 1)
  return true
}

function toggleDishStatus(id: string): Dish | null {
  const dish = dishes.value.find((d) => d.id === id)
  if (!dish) return null

  dish.status = dish.status === 'on_sale' ? 'off_sale' : 'on_sale'
  dish.updatedAt = now()

  return { ...dish }
}

export const menuHandlers = {
  // Categories
  getCategories,
  getCategoryById,
  createCategory,
  updateCategory,
  deleteCategory,
  reorderCategories,
  // Dishes
  getDishes,
  getDishById,
  createDish,
  updateDish,
  deleteDish,
  toggleDishStatus,
}

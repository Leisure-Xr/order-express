import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { Category, Dish } from '@/types'
import {
  getCategoriesApi,
  createCategoryApi,
  updateCategoryApi,
  deleteCategoryApi,
  reorderCategoriesApi,
  getDishesApi,
  getDishByIdApi,
  createDishApi,
  updateDishApi,
  deleteDishApi,
  toggleDishStatusApi,
} from '@/api/menu'

function assertOk(code: number, message: string): asserts code is 200 {
  if (code !== 200) {
    throw new Error(message || 'Request failed')
  }
}

export const useMenuStore = defineStore('menu', () => {
  const categories = ref<Category[]>([])
  const dishes = ref<Dish[]>([])

  const categoriesLoading = ref(false)
  const dishesLoading = ref(false)

  const activeCategories = computed(() =>
    [...categories.value]
      .filter((c) => c.status === 'active')
      .sort((a, b) => a.sortOrder - b.sortOrder),
  )

  const dishesByCategoryId = computed(() => {
    const map = new Map<string, Dish[]>()
    for (const dish of dishes.value) {
      const list = map.get(dish.categoryId) ?? []
      list.push(dish)
      map.set(dish.categoryId, list)
    }
    for (const [key, list] of map.entries()) {
      list.sort((a, b) => a.name.zh.localeCompare(b.name.zh))
      map.set(key, list)
    }
    return map
  })

  async function fetchCategories(): Promise<Category[]> {
    categoriesLoading.value = true
    try {
      const res = await getCategoriesApi()
      assertOk(res.code, res.message)
      categories.value = res.data ?? []
      return categories.value
    } finally {
      categoriesLoading.value = false
    }
  }

  async function fetchDishes(params?: { categoryId?: string; status?: string }): Promise<Dish[]> {
    dishesLoading.value = true
    try {
      const res = await getDishesApi(params)
      assertOk(res.code, res.message)
      dishes.value = res.data ?? []
      return dishes.value
    } finally {
      dishesLoading.value = false
    }
  }

  async function fetchAll(): Promise<void> {
    await Promise.all([fetchCategories(), fetchDishes()])
  }

  async function fetchDishById(id: string): Promise<Dish | null> {
    const res = await getDishByIdApi(id)
    if (res.code !== 200) return null
    return res.data ?? null
  }

  async function createCategory(payload: Omit<Category, 'id' | 'dishCount'>): Promise<Category> {
    const res = await createCategoryApi(payload)
    assertOk(res.code, res.message)
    const created = res.data
    if (!created) throw new Error('Category create failed')
    categories.value.push(created)
    categories.value.sort((a, b) => a.sortOrder - b.sortOrder)
    return created
  }

  async function updateCategory(id: string, payload: Partial<Omit<Category, 'id'>>): Promise<Category | null> {
    const res = await updateCategoryApi(id, payload)
    assertOk(res.code, res.message)
    const updated = res.data
    if (!updated) return null
    const idx = categories.value.findIndex((c) => c.id === id)
    if (idx !== -1) {
      categories.value[idx] = updated
      categories.value.sort((a, b) => a.sortOrder - b.sortOrder)
    }
    return updated
  }

  async function deleteCategory(id: string): Promise<boolean> {
    const res = await deleteCategoryApi(id)
    if (res.code !== 200) return false
    categories.value = categories.value.filter((c) => c.id !== id)
    dishes.value = dishes.value.filter((d) => d.categoryId !== id)
    return true
  }

  async function reorderCategories(orderedIds: string[]): Promise<Category[]> {
    const res = await reorderCategoriesApi(orderedIds)
    assertOk(res.code, res.message)
    categories.value = res.data ?? []
    return categories.value
  }

  async function createDish(payload: Omit<Dish, 'id' | 'createdAt' | 'updatedAt'>): Promise<Dish> {
    const res = await createDishApi(payload)
    assertOk(res.code, res.message)
    const created = res.data
    if (!created) throw new Error('Dish create failed')
    dishes.value.push(created)
    return created
  }

  async function updateDish(id: string, payload: Partial<Omit<Dish, 'id' | 'createdAt'>>): Promise<Dish | null> {
    const res = await updateDishApi(id, payload)
    assertOk(res.code, res.message)
    const updated = res.data
    if (!updated) return null
    const idx = dishes.value.findIndex((d) => d.id === id)
    if (idx !== -1) dishes.value[idx] = updated
    return updated
  }

  async function deleteDish(id: string): Promise<boolean> {
    const res = await deleteDishApi(id)
    if (res.code !== 200) return false
    dishes.value = dishes.value.filter((d) => d.id !== id)
    return true
  }

  async function toggleDishStatus(id: string): Promise<Dish | null> {
    const res = await toggleDishStatusApi(id)
    assertOk(res.code, res.message)
    const updated = res.data
    if (!updated) return null
    const idx = dishes.value.findIndex((d) => d.id === id)
    if (idx !== -1) dishes.value[idx] = updated
    return updated
  }

  return {
    categories,
    dishes,
    categoriesLoading,
    dishesLoading,
    activeCategories,
    dishesByCategoryId,
    fetchCategories,
    fetchDishes,
    fetchAll,
    fetchDishById,
    createCategory,
    updateCategory,
    deleteCategory,
    reorderCategories,
    createDish,
    updateDish,
    deleteDish,
    toggleDishStatus,
  }
})


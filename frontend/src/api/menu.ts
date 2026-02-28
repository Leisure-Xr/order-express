import { mockDelay, successResponse, errorResponse } from './index'
import type { ApiResponse } from './index'
import type { Category, Dish } from '@/types'
import { menuHandlers } from './mock/handlers/menu'
import { apiRequest, isMockApiEnabled } from './client'
import { endpoints } from './endpoints'

// -----------------------------------------------------------
// Category APIs
// -----------------------------------------------------------

export async function getCategoriesApi(): Promise<ApiResponse<Category[]>> {
  if (isMockApiEnabled()) {
    await mockDelay()
    return successResponse(menuHandlers.getCategories())
  }

  return apiRequest<Category[]>({
    path: endpoints.categories.list,
    method: 'GET',
  })
}

export async function getCategoryByIdApi(id: string): Promise<ApiResponse<Category | null>> {
  if (isMockApiEnabled()) {
    await mockDelay()
    const category = menuHandlers.getCategoryById(id)
    if (!category) {
      return errorResponse(404, 'Category not found')
    }
    return successResponse(category)
  }

  return apiRequest<Category>({
    path: endpoints.categories.get(id),
    method: 'GET',
  })
}

export async function createCategoryApi(
  data: Omit<Category, 'id' | 'dishCount'>,
): Promise<ApiResponse<Category>> {
  if (isMockApiEnabled()) {
    await mockDelay()
    const category = menuHandlers.createCategory(data)
    return successResponse(category)
  }

  return apiRequest<Category>({
    path: endpoints.categories.create,
    method: 'POST',
    body: data,
  })
}

export async function updateCategoryApi(
  id: string,
  data: Partial<Omit<Category, 'id'>>,
): Promise<ApiResponse<Category | null>> {
  if (isMockApiEnabled()) {
    await mockDelay()
    const category = menuHandlers.updateCategory(id, data)
    if (!category) {
      return errorResponse(404, 'Category not found')
    }
    return successResponse(category)
  }

  return apiRequest<Category>({
    path: endpoints.categories.update(id),
    method: 'PATCH',
    body: data,
  })
}

export async function deleteCategoryApi(id: string): Promise<ApiResponse<null>> {
  if (isMockApiEnabled()) {
    await mockDelay()
    const success = menuHandlers.deleteCategory(id)
    if (!success) {
      return errorResponse(404, 'Category not found')
    }
    return successResponse(null)
  }

  return apiRequest<null>({
    path: endpoints.categories.delete(id),
    method: 'DELETE',
  })
}

export async function reorderCategoriesApi(orderedIds: string[]): Promise<ApiResponse<Category[]>> {
  if (isMockApiEnabled()) {
    await mockDelay()
    const categories = menuHandlers.reorderCategories(orderedIds)
    return successResponse(categories)
  }

  return apiRequest<Category[]>({
    path: endpoints.categories.reorder,
    method: 'POST',
    body: { orderedIds },
  })
}

// -----------------------------------------------------------
// Dish APIs
// -----------------------------------------------------------

export async function getDishesApi(
  params?: { categoryId?: string; status?: string },
): Promise<ApiResponse<Dish[]>> {
  if (isMockApiEnabled()) {
    await mockDelay()
    return successResponse(menuHandlers.getDishes(params))
  }

  return apiRequest<Dish[]>({
    path: endpoints.dishes.list,
    method: 'GET',
    query: params,
  })
}

export async function getDishByIdApi(id: string): Promise<ApiResponse<Dish | null>> {
  if (isMockApiEnabled()) {
    await mockDelay()
    const dish = menuHandlers.getDishById(id)
    if (!dish) {
      return errorResponse(404, 'Dish not found')
    }
    return successResponse(dish)
  }

  return apiRequest<Dish>({
    path: endpoints.dishes.get(id),
    method: 'GET',
  })
}

export async function createDishApi(
  data: Omit<Dish, 'id' | 'createdAt' | 'updatedAt'>,
): Promise<ApiResponse<Dish>> {
  if (isMockApiEnabled()) {
    await mockDelay()
    const dish = menuHandlers.createDish(data)
    return successResponse(dish)
  }

  return apiRequest<Dish>({
    path: endpoints.dishes.create,
    method: 'POST',
    body: data,
  })
}

export async function updateDishApi(
  id: string,
  data: Partial<Omit<Dish, 'id' | 'createdAt'>>,
): Promise<ApiResponse<Dish | null>> {
  if (isMockApiEnabled()) {
    await mockDelay()
    const dish = menuHandlers.updateDish(id, data)
    if (!dish) {
      return errorResponse(404, 'Dish not found')
    }
    return successResponse(dish)
  }

  return apiRequest<Dish>({
    path: endpoints.dishes.update(id),
    method: 'PATCH',
    body: data,
  })
}

export async function deleteDishApi(id: string): Promise<ApiResponse<null>> {
  if (isMockApiEnabled()) {
    await mockDelay()
    const success = menuHandlers.deleteDish(id)
    if (!success) {
      return errorResponse(404, 'Dish not found')
    }
    return successResponse(null)
  }

  return apiRequest<null>({
    path: endpoints.dishes.delete(id),
    method: 'DELETE',
  })
}

export async function toggleDishStatusApi(id: string): Promise<ApiResponse<Dish | null>> {
  if (isMockApiEnabled()) {
    await mockDelay()
    const dish = menuHandlers.toggleDishStatus(id)
    if (!dish) {
      return errorResponse(404, 'Dish not found')
    }
    return successResponse(dish)
  }

  return apiRequest<Dish>({
    path: endpoints.dishes.toggleStatus(id),
    method: 'POST',
  })
}

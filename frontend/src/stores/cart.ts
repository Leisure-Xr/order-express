import { defineStore } from 'pinia'
import { ref, computed, watch } from 'vue'
import type { CartItem, Dish } from '@/types'

const CART_STORAGE_KEY = 'cart_items'

/**
 * Generate a unique key for a cart item based on dishId and selected options.
 * This ensures the same dish with different options occupies separate cart rows.
 */
function generateCartItemKey(
  dishId: string,
  selectedOptions: CartItem['selectedOptions'],
): string {
  if (!selectedOptions.length) {
    return dishId
  }
  const sorted = [...selectedOptions]
    .sort((a, b) => a.optionName.localeCompare(b.optionName))
    .map((o) => `${o.optionName}:${o.valueName}`)
    .join('|')
  return `${dishId}__${sorted}`
}

/** Load persisted cart from localStorage */
function loadCartFromStorage(): CartItem[] {
  try {
    const raw = localStorage.getItem(CART_STORAGE_KEY)
    return raw ? (JSON.parse(raw) as CartItem[]) : []
  } catch {
    return []
  }
}

export const useCartStore = defineStore('cart', () => {
  // ---- State ----
  const items = ref<CartItem[]>(loadCartFromStorage())

  // ---- Auto-persist to localStorage ----
  watch(
    items,
    (newItems) => {
      localStorage.setItem(CART_STORAGE_KEY, JSON.stringify(newItems))
    },
    { deep: true },
  )

  // ---- Computed ----

  /** Total number of items (sum of quantities) */
  const itemCount = computed(() =>
    items.value.reduce((sum, item) => sum + item.quantity, 0),
  )

  /** Subtotal price across all items (unitPrice already includes option adjustments) */
  const subtotal = computed(() =>
    items.value.reduce((sum, item) => sum + item.unitPrice * item.quantity, 0),
  )

  /** Whether the cart is empty */
  const isEmpty = computed(() => items.value.length === 0)

  // ---- Actions ----

  /**
   * Add a dish to the cart.
   * If the exact same dish + options combination already exists, increase its quantity.
   */
  function addItem(
    dish: Dish,
    quantity: number = 1,
    selectedOptions: CartItem['selectedOptions'] = [],
  ) {
    // Calculate unit price including option adjustments
    const optionExtra = selectedOptions.reduce(
      (sum, opt) => sum + opt.priceAdjustment,
      0,
    )
    const unitPrice = dish.price + optionExtra

    const key = generateCartItemKey(dish.id, selectedOptions)

    const existing = items.value.find((item) => item.key === key)
    if (existing) {
      existing.quantity += quantity
    } else {
      items.value.push({
        key,
        dishId: dish.id,
        dishName: { ...dish.name },
        image: dish.image,
        unitPrice,
        quantity,
        selectedOptions: [...selectedOptions],
      })
    }
  }

  /**
   * Update the quantity of a specific cart item by key.
   * If quantity is 0 or less, the item is removed.
   */
  function updateQuantity(key: string, quantity: number) {
    if (quantity <= 0) {
      removeItem(key)
      return
    }
    const item = items.value.find((i) => i.key === key)
    if (item) {
      item.quantity = quantity
    }
  }

  /** Remove a single item from the cart by key */
  function removeItem(key: string) {
    const index = items.value.findIndex((i) => i.key === key)
    if (index !== -1) {
      items.value.splice(index, 1)
    }
  }

  /** Clear all items from the cart */
  function clearCart() {
    items.value = []
  }

  return {
    items,
    itemCount,
    subtotal,
    isEmpty,
    addItem,
    updateQuantity,
    removeItem,
    clearCart,
  }
})

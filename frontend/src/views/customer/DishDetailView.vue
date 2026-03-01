<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { useI18n } from 'vue-i18n'
import { useMenuStore } from '@/stores/menu'
import { useCartStore } from '@/stores/cart'
import { useLocaleText } from '@/composables/useLocaleText'
import DishOptionPicker, { type OptionSelectionMap } from '@/components/customer/DishOptionPicker.vue'
import PriceTag from '@/components/common/PriceTag.vue'
import EmptyState from '@/components/common/EmptyState.vue'
import type { Dish } from '@/types'

const { t } = useI18n()
const route = useRoute()
const router = useRouter()
const menuStore = useMenuStore()
const cartStore = useCartStore()
const { localText } = useLocaleText()

const dishId = computed(() => (typeof route.params.id === 'string' ? route.params.id : ''))

const loading = ref(false)
const dish = ref<Dish | null>(null)
const quantity = ref(1)
const selection = ref<OptionSelectionMap>({})

onMounted(async () => {
  if (!dishId.value) return
  loading.value = true
  try {
    const inState = menuStore.dishes.find((d) => d.id === dishId.value)
    dish.value = inState || (await menuStore.fetchDishById(dishId.value))
    if (!dish.value) {
      ElMessage.error('Dish not found')
      return
    }

    // Init selection defaults (empty)
    const init: OptionSelectionMap = {}
    for (const opt of dish.value.options) {
      init[opt.id] = []
    }
    selection.value = init
  } catch (e: any) {
    ElMessage.error(e?.message || t('common.fail'))
  } finally {
    loading.value = false
  }
})

const selectedOptionDetails = computed(() => {
  if (!dish.value) return []
  const result: { optionName: string; valueName: string; priceAdjustment: number }[] = []

  for (const opt of dish.value.options) {
    const ids = selection.value[opt.id] ?? []
    for (const id of ids) {
      const v = opt.values.find((vv) => vv.id === id)
      if (!v) continue
      result.push({
        optionName: localText(opt.name),
        valueName: localText(v.label),
        priceAdjustment: v.priceAdjustment,
      })
    }
  }
  return result
})

const unitPrice = computed(() => {
  if (!dish.value) return 0
  const optionExtra = selectedOptionDetails.value.reduce((sum, o) => sum + o.priceAdjustment, 0)
  return dish.value.price + optionExtra
})

const totalPrice = computed(() => unitPrice.value * quantity.value)

function validateSelections(): boolean {
  if (!dish.value) return false
  for (const opt of dish.value.options) {
    if (!opt.required) continue
    const picked = selection.value[opt.id] ?? []
    if (!picked.length) {
      ElMessage.warning(`${t('menu.required')}: ${localText(opt.name)}`)
      return false
    }
  }
  return true
}

function addToCart() {
  if (!dish.value) return
  if (!validateSelections()) return

  cartStore.addItem(dish.value, quantity.value, selectedOptionDetails.value)
  ElMessage.success(t('common.success'))
  router.push({ name: 'CustomerCart' })
}
</script>

<template>
  <div class="dish-detail">
    <el-skeleton v-if="loading" :rows="8" animated />

    <EmptyState v-else-if="!dish" :title="t('common.noData')" />

    <div v-else class="content">
      <img class="hero" :src="dish.image" :alt="localText(dish.name)" />

      <div class="card">
        <div class="name">{{ localText(dish.name) }}</div>
        <div class="desc">{{ localText(dish.description) }}</div>
        <div class="price-row">
          <PriceTag :price="unitPrice" :original-price="dish.originalPrice" size="lg" />
          <div class="qty">
            <el-input-number v-model="quantity" :min="1" :max="99" />
          </div>
        </div>
      </div>

      <div v-if="dish.options.length" class="card">
        <DishOptionPicker v-model="selection" :options="dish.options" />
      </div>

      <div class="footer">
        <div class="total">
          {{ t('common.total') }}: <span class="total-price">¥{{ totalPrice.toFixed(2) }}</span>
        </div>
        <el-button type="primary" size="large" class="add-btn" @click="addToCart">
          {{ t('cart.goCheckout') }}
        </el-button>
      </div>
    </div>
  </div>
</template>

<style scoped lang="scss">
.dish-detail {
  padding: 12px;
}

.content {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.hero {
  width: 100%;
  height: 220px;
  object-fit: cover;
  border-radius: 16px;
  background: linear-gradient(110deg, #f0f0f0 30%, #fafafa 50%, #f0f0f0 70%);
  background-size: 200% 100%;
  animation: shimmer 1.5s infinite linear;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
}

.card {
  background: #ffffff;
  border-radius: 14px;
  padding: 14px;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.06);
  transition: box-shadow var(--app-transition-base);

  &:hover {
    box-shadow: 0 4px 16px rgba(0, 0, 0, 0.08);
  }
}

.name {
  font-size: 20px;
  font-weight: 800;
  color: #303133;
}

.desc {
  margin-top: 6px;
  font-size: 13px;
  color: #909399;
  line-height: 1.5;
}

.price-row {
  margin-top: 12px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.qty {
  flex: none;
}

.footer {
  position: sticky;
  bottom: calc(var(--customer-bottom-bar-height) + 8px);
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  background: var(--app-glass-bg-strong);
  backdrop-filter: var(--app-glass-blur);
  border-radius: 14px;
  padding: 12px 14px;
  box-shadow: 0 -4px 24px rgba(0, 0, 0, 0.08);
  border: 1px solid rgba(255, 255, 255, 0.3);
  animation: slideUp 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.total {
  font-size: 14px;
  color: #606266;
}

.total-price {
  color: #e74c3c;
  font-weight: 800;
}

.add-btn {
  flex: none;
}
</style>

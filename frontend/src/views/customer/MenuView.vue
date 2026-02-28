<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { useI18n } from 'vue-i18n'
import { useMenuStore } from '@/stores/menu'
import { useCartStore } from '@/stores/cart'
import CategoryNav from '@/components/customer/CategoryNav.vue'
import DishCard from '@/components/customer/DishCard.vue'
import EmptyState from '@/components/common/EmptyState.vue'
import type { Dish } from '@/types'

const { t } = useI18n()
const router = useRouter()
const menuStore = useMenuStore()
const cartStore = useCartStore()

const activeCategoryId = ref('')

const categories = computed(() => menuStore.activeCategories)

const dishesForActiveCategory = computed(() => {
  const list = menuStore.dishes.filter((d) => d.status === 'on_sale')
  if (!activeCategoryId.value) return list
  return list.filter((d) => d.categoryId === activeCategoryId.value)
})

onMounted(async () => {
  try {
    if (!menuStore.categories.length || !menuStore.dishes.length) {
      await menuStore.fetchAll()
    }
    activeCategoryId.value = categories.value[0]?.id || ''
  } catch (e: any) {
    ElMessage.error(e?.message || t('common.fail'))
  }
})

function onSelectCategory(categoryId: string) {
  activeCategoryId.value = categoryId
}

function openDishDetail(dish: Dish) {
  router.push({ name: 'CustomerDishDetail', params: { id: dish.id } })
}

function quickAdd(dish: Dish) {
  if (dish.options?.length) {
    openDishDetail(dish)
    return
  }
  cartStore.addItem(dish, 1, [])
  ElMessage.success(t('common.success'))
}
</script>

<template>
  <div class="menu-page">
    <CategoryNav :categories="categories" :active-id="activeCategoryId" @select="onSelectCategory" />

    <div class="dish-list" v-loading="menuStore.dishesLoading">
      <DishCard
        v-for="dish in dishesForActiveCategory"
        :key="dish.id"
        :dish="dish"
        @click="openDishDetail"
        @add="quickAdd"
      />

      <EmptyState
        v-if="!menuStore.dishesLoading && dishesForActiveCategory.length === 0"
        :title="t('common.noData')"
      />
    </div>
  </div>
</template>

<style scoped lang="scss">
.menu-page {
  padding: 0 12px 12px 12px;
}

.dish-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
  padding-top: 4px;
}
</style>


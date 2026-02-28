<script setup lang="ts">
import type { Dish } from '@/types'
import { useLocaleText } from '@/composables/useLocaleText'
import PriceTag from '@/components/common/PriceTag.vue'
import { Plus } from '@element-plus/icons-vue'

const props = withDefaults(
  defineProps<{
    dish: Dish
    disabled?: boolean
    showAdd?: boolean
  }>(),
  {
    disabled: false,
    showAdd: true,
  },
)

const emit = defineEmits<{
  (e: 'click', dish: Dish): void
  (e: 'add', dish: Dish): void
}>()

const { localText } = useLocaleText()

function handleClick() {
  emit('click', props.dish)
}

function handleAdd(e: MouseEvent) {
  e.stopPropagation()
  emit('add', props.dish)
}
</script>

<template>
  <div class="dish-card" :class="{ disabled }" @click="handleClick">
    <img class="dish-image" :src="dish.image" :alt="localText(dish.name)" />
    <div class="dish-body">
      <div class="dish-title">{{ localText(dish.name) }}</div>
      <div class="dish-desc">{{ localText(dish.description) }}</div>
      <div class="dish-bottom">
        <PriceTag :price="dish.price" :original-price="dish.originalPrice" size="sm" />
        <el-button
          v-if="showAdd"
          type="primary"
          circle
          :icon="Plus"
          size="small"
          :disabled="disabled"
          @click="handleAdd"
        />
      </div>
    </div>
  </div>
</template>

<style scoped lang="scss">
.dish-card {
  display: flex;
  gap: 12px;
  padding: 12px;
  background: #ffffff;
  border-radius: 12px;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.06);
  cursor: pointer;
  user-select: none;
  -webkit-tap-highlight-color: transparent;
}

.dish-card.disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.dish-image {
  width: 88px;
  height: 88px;
  object-fit: cover;
  border-radius: 10px;
  background: #f2f2f2;
  flex: none;
}

.dish-body {
  min-width: 0;
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.dish-title {
  font-size: 15px;
  font-weight: 700;
  color: #303133;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.dish-desc {
  font-size: 12px;
  color: #909399;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
}

.dish-bottom {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: auto;
}
</style>


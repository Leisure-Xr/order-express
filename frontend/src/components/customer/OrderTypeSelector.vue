<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import type { OrderType } from '@/types'

defineProps<{
  modelValue: OrderType
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: OrderType): void
}>()

const { t } = useI18n()

const types: { value: OrderType; icon: string; labelKey: string }[] = [
  { value: 'dine_in', icon: '🍽️', labelKey: 'order.dineIn' },
  { value: 'takeaway', icon: '🛵', labelKey: 'order.takeout' },
  { value: 'pickup', icon: '🏃', labelKey: 'order.pickup' },
]
</script>

<template>
  <div class="type-selector">
    <div
      v-for="type in types"
      :key="type.value"
      class="type-card"
      :class="{ active: modelValue === type.value }"
      @click="emit('update:modelValue', type.value)"
    >
      <span class="type-icon">{{ type.icon }}</span>
      <span class="type-label">{{ t(type.labelKey) }}</span>
    </div>
  </div>
</template>

<style scoped lang="scss">
.type-selector {
  display: flex;
  gap: 10px;
}

.type-card {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 6px;
  padding: 14px 8px;
  border-radius: 14px;
  border: 2px solid #e5e7eb;
  background: #ffffff;
  cursor: pointer;
  transition: all var(--app-transition-base);

  &:hover:not(.active) {
    border-color: rgba(231, 76, 60, 0.3);
    background: rgba(231, 76, 60, 0.02);
  }

  &.active {
    border-color: #e74c3c;
    background: rgba(231, 76, 60, 0.06);
    box-shadow: 0 2px 12px rgba(231, 76, 60, 0.15);
  }

  &:active {
    transform: scale(0.97);
  }
}

.type-icon {
  font-size: 24px;
}

.type-label {
  font-size: 13px;
  font-weight: 600;
  color: #606266;

  .active & {
    color: #e74c3c;
  }
}
</style>

<script setup lang="ts">
import type { Category } from '@/types'
import { useLocaleText } from '@/composables/useLocaleText'

const props = defineProps<{
  categories: Category[]
  activeId: string
}>()

const emit = defineEmits<{
  (e: 'select', categoryId: string): void
}>()

const { localText } = useLocaleText()

function selectCategory(categoryId: string) {
  emit('select', categoryId)
}
</script>

<template>
  <div class="category-nav">
    <div
      v-for="cat in categories"
      :key="cat.id"
      class="cat-pill"
      :class="{ active: cat.id === activeId }"
      @click="selectCategory(cat.id)"
    >
      <span class="cat-icon" v-if="cat.icon">{{ cat.icon }}</span>
      <span class="cat-name">{{ localText(cat.name) }}</span>
    </div>
  </div>
</template>

<style scoped lang="scss">
.category-nav {
  display: flex;
  gap: 8px;
  overflow-x: auto;
  padding: 12px 4px;
  -webkit-overflow-scrolling: touch;
  scrollbar-width: none;

  &::-webkit-scrollbar {
    display: none;
  }
}

.cat-pill {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 8px 12px;
  border-radius: 999px;
  background: #ffffff;
  color: #606266;
  font-size: 13px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.06);
  white-space: nowrap;
  cursor: pointer;
  user-select: none;
  border: 1.5px solid transparent;
  transition: background var(--app-transition-base),
              color var(--app-transition-base),
              transform var(--app-transition-fast),
              box-shadow var(--app-transition-base),
              border-color var(--app-transition-base);

  &:hover:not(.active) {
    background: #f7f8fb;
    border-color: rgba(0, 0, 0, 0.04);
  }

  &.active {
    background: linear-gradient(135deg, rgba(231, 76, 60, 0.12), rgba(243, 156, 18, 0.06));
    color: #e74c3c;
    font-weight: 700;
    border-color: rgba(231, 76, 60, 0.18);
    box-shadow: 0 2px 8px rgba(231, 76, 60, 0.12);
  }

  &:active {
    transform: scale(0.95);
  }
}

.cat-icon {
  font-size: 14px;
}
</style>


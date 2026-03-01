<script setup lang="ts">
import { computed } from 'vue'

const props = withDefaults(
  defineProps<{
    title: string
    description?: string
  }>(),
  {
    description: '',
  },
)

const showDescription = computed(() => !!props.description?.trim())
</script>

<template>
  <div class="admin-page-header">
    <div class="row">
      <div class="left">
        <div class="title">{{ title }}</div>
        <div v-if="showDescription" class="desc">{{ description }}</div>
      </div>
      <div class="right">
        <slot name="actions" />
      </div>
    </div>
  </div>
</template>

<style scoped lang="scss">
.admin-page-header {
  margin-bottom: 14px;
}

.row {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 12px;
  flex-wrap: wrap;
}

.left {
  min-width: 0;
}

.title {
  font-size: 18px;
  font-weight: 900;
  color: #111827;
  letter-spacing: 0.2px;
  position: relative;
  padding-left: 12px;

  &::before {
    content: '';
    position: absolute;
    left: 0;
    top: 3px;
    bottom: 3px;
    width: 4px;
    border-radius: 999px;
    background: var(--app-gradient-primary);
  }
}

.desc {
  margin-top: 6px;
  font-size: 13px;
  color: #6b7280;
  line-height: 1.5;
  max-width: 860px;
}

.right {
  display: flex;
  align-items: center;
  gap: 10px;
}
</style>


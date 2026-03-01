<script setup lang="ts">
import { computed } from 'vue'
import type { CartItem } from '@/types'
import { useLocaleText } from '@/composables/useLocaleText'
import PriceTag from '@/components/common/PriceTag.vue'
import { Delete } from '@element-plus/icons-vue'

const props = defineProps<{
  item: CartItem
}>()

const emit = defineEmits<{
  (e: 'update:quantity', payload: { key: string; quantity: number }): void
  (e: 'remove', key: string): void
}>()

const { localText } = useLocaleText()

const lineTotal = computed(() => props.item.unitPrice * props.item.quantity)

function updateQuantity(qty: number) {
  emit('update:quantity', { key: props.item.key, quantity: qty })
}

function remove() {
  emit('remove', props.item.key)
}
</script>

<template>
  <div class="cart-row">
    <img class="thumb" :src="item.image" :alt="localText(item.dishName)" />
    <div class="body">
      <div class="title">
        <span class="name">{{ localText(item.dishName) }}</span>
        <el-button text :icon="Delete" @click="remove" />
      </div>

      <div v-if="item.selectedOptions.length" class="options">
        <span
          v-for="(opt, idx) in item.selectedOptions"
          :key="idx"
          class="opt"
        >
          {{ opt.optionName }}: {{ opt.valueName }}
        </span>
      </div>

      <div class="bottom">
        <PriceTag :price="item.unitPrice" size="sm" />
        <div class="qty">
          <el-input-number
            :model-value="item.quantity"
            :min="1"
            :max="99"
            size="small"
            @update:model-value="(v: number | undefined) => updateQuantity(typeof v === 'number' ? v : item.quantity)"
          />
        </div>
      </div>

      <div class="line-total">¥{{ lineTotal.toFixed(2) }}</div>
    </div>
  </div>
</template>

<style scoped lang="scss">
.cart-row {
  display: flex;
  gap: 12px;
  padding: 12px;
  background: #ffffff;
  border-radius: 12px;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.06);
  transition: box-shadow var(--app-transition-base);

  &:hover {
    box-shadow: 0 4px 16px rgba(0, 0, 0, 0.08);
  }
}

.thumb {
  width: 72px;
  height: 72px;
  border-radius: 10px;
  object-fit: cover;
  background: #f2f2f2;
  flex: none;
}

.body {
  flex: 1;
  min-width: 0;
}

.title {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
}

.name {
  font-weight: 700;
  color: #303133;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.options {
  margin-top: 6px;
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  color: #909399;
  font-size: 12px;
}

.opt {
  background: linear-gradient(135deg, #f5f7fa, #eef2f7);
  border-radius: 999px;
  padding: 3px 8px;
  border: 1px solid rgba(0, 0, 0, 0.03);
}

.bottom {
  margin-top: 10px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.line-total {
  margin-top: 6px;
  color: #e74c3c;
  font-weight: 800;
  text-align: right;
}
</style>

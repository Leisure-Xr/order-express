<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import type { OrderType } from '@/types'

const props = defineProps<{
  modelValue: OrderType
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: OrderType): void
}>()

const { t } = useI18n()
</script>

<template>
  <el-radio-group
    :model-value="modelValue"
    @update:model-value="
      (v: string | number | boolean | undefined) => {
        if (v === 'dine_in' || v === 'takeaway' || v === 'pickup') emit('update:modelValue', v)
      }
    "
  >
    <el-radio value="dine_in">{{ t('order.dineIn') }}</el-radio>
    <el-radio value="takeaway">{{ t('order.takeout') }}</el-radio>
    <el-radio value="pickup">{{ t('order.pickup') }}</el-radio>
  </el-radio-group>
</template>

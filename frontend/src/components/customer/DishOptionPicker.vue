<script setup lang="ts">
import { ref, watch } from 'vue'
import type { DishOption } from '@/types'
import { useLocaleText } from '@/composables/useLocaleText'

export type OptionSelectionMap = Record<string, string[]>

const props = withDefaults(
  defineProps<{
    options: DishOption[]
    modelValue: OptionSelectionMap
    disabled?: boolean
  }>(),
  {
    disabled: false,
  },
)

const emit = defineEmits<{
  (e: 'update:modelValue', value: OptionSelectionMap): void
}>()

const { localText } = useLocaleText()

function clone<T>(value: T): T {
  return JSON.parse(JSON.stringify(value)) as T
}

const selection = ref<OptionSelectionMap>(clone(props.modelValue))

watch(
  () => props.modelValue,
  (value) => {
    selection.value = clone(value)
  },
  { deep: true },
)

watch(
  selection,
  (value) => {
    emit('update:modelValue', clone(value))
  },
  { deep: true },
)

function setSingle(optionId: string, valueId: string) {
  if (!valueId) {
    selection.value[optionId] = []
    return
  }
  selection.value[optionId] = [valueId]
}

function formatPriceAdj(adj: number): string {
  if (!adj) return ''
  return adj > 0 ? `+${adj}` : `${adj}`
}
</script>

<template>
  <div class="option-picker">
    <div v-for="opt in options" :key="opt.id" class="option-group">
      <div class="option-title">
        <span>{{ localText(opt.name) }}</span>
        <span v-if="opt.required" class="required">*</span>
      </div>

      <el-checkbox-group
        v-if="opt.multiSelect"
        v-model="selection[opt.id]"
        :disabled="disabled"
        class="option-values"
      >
        <el-checkbox v-for="v in opt.values" :key="v.id" :value="v.id" class="option-pill">
          {{ localText(v.label) }}
          <span v-if="v.priceAdjustment" class="adj">{{ formatPriceAdj(v.priceAdjustment) }}</span>
        </el-checkbox>
      </el-checkbox-group>

      <el-radio-group
        v-else
        :model-value="selection[opt.id]?.[0] || ''"
        :disabled="disabled"
        class="option-values"
        @update:model-value="(v: string | number | boolean | undefined) => setSingle(opt.id, typeof v === 'string' ? v : '')"
      >
        <el-radio v-for="v in opt.values" :key="v.id" :value="v.id" class="option-pill">
          {{ localText(v.label) }}
          <span v-if="v.priceAdjustment" class="adj">{{ formatPriceAdj(v.priceAdjustment) }}</span>
        </el-radio>
      </el-radio-group>
    </div>
  </div>
</template>

<style scoped lang="scss">
.option-group {
  margin-bottom: 16px;
}

.option-title {
  font-size: 13px;
  font-weight: 600;
  color: #303133;
  margin-bottom: 8px;
}

.required {
  color: #e74c3c;
  margin-left: 4px;
}

.option-values {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.option-pill {
  margin-right: 0 !important;
  :deep(.el-checkbox__label),
  :deep(.el-radio__label) {
    padding-left: 6px;
    font-size: 13px;
  }
}

.adj {
  margin-left: 6px;
  color: #e74c3c;
  font-weight: 600;
  font-size: 12px;
}
</style>

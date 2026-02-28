<script setup lang="ts">
import { ref, watch, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import type { Category, Dish, DishOption, DishOptionValue } from '@/types'

export type DishFormValue = Omit<Dish, 'id' | 'createdAt' | 'updatedAt'>

const props = withDefaults(
  defineProps<{
    modelValue: DishFormValue
    categories: Category[]
    disabled?: boolean
  }>(),
  {
    disabled: false,
  },
)

const emit = defineEmits<{
  (e: 'update:modelValue', value: DishFormValue): void
}>()

const { t } = useI18n()

function clone<T>(value: T): T {
  return JSON.parse(JSON.stringify(value)) as T
}

const model = ref<DishFormValue>(clone(props.modelValue))

watch(
  () => props.modelValue,
  (value) => {
    model.value = clone(value)
  },
  { deep: true },
)

watch(
  model,
  (value) => {
    emit('update:modelValue', clone(value))
  },
  { deep: true },
)

const statusOptions = computed(() => [
  { value: 'on_sale', label: t('menu.onSale') },
  { value: 'off_sale', label: t('menu.offShelf') },
  { value: 'sold_out', label: t('menu.soldOut') },
])

const tagPresets = computed(() => [
  { value: 'popular', label: t('menu.popular') },
  { value: 'new', label: t('menu.newItem') },
  { value: 'spicy', label: t('menu.spicy') },
  { value: 'vegetarian', label: t('menu.vegetarian') },
])

function generateId(prefix: string): string {
  return `${prefix}-${Date.now()}-${Math.random().toString(36).slice(2, 8)}`
}

function addOptionGroup() {
  const option: DishOption = {
    id: generateId('opt'),
    name: { zh: '', en: '' },
    values: [],
    required: false,
    multiSelect: false,
  }
  model.value.options.push(option)
}

function removeOptionGroup(index: number) {
  model.value.options.splice(index, 1)
}

function addOptionValue(optionIndex: number) {
  const value: DishOptionValue = {
    id: generateId('val'),
    label: { zh: '', en: '' },
    priceAdjustment: 0,
  }
  model.value.options[optionIndex]?.values.push(value)
}

function removeOptionValue(optionIndex: number, valueIndex: number) {
  model.value.options[optionIndex]?.values.splice(valueIndex, 1)
}
</script>

<template>
  <el-form :model="model" label-position="top" class="dish-form">
    <div class="grid">
      <el-form-item :label="t('menu.itemName') + ' (ZH)'" class="col-6">
        <el-input v-model="model.name.zh" :disabled="disabled" />
      </el-form-item>
      <el-form-item :label="t('menu.itemName') + ' (EN)'" class="col-6">
        <el-input v-model="model.name.en" :disabled="disabled" />
      </el-form-item>
    </div>

    <el-form-item :label="t('menu.category')">
      <el-select v-model="model.categoryId" :disabled="disabled" filterable>
        <el-option
          v-for="c in categories"
          :key="c.id"
          :value="c.id"
          :label="`${c.icon || ''} ${c.name.zh} / ${c.name.en}`"
        />
      </el-select>
    </el-form-item>

    <div class="grid">
      <el-form-item :label="t('menu.price')" class="col-4">
        <el-input-number v-model="model.price" :disabled="disabled" :min="0" :precision="2" :step="1" />
      </el-form-item>
      <el-form-item :label="t('menu.originalPrice')" class="col-4">
        <el-input-number v-model="model.originalPrice" :disabled="disabled" :min="0" :precision="2" :step="1" />
      </el-form-item>
      <el-form-item :label="t('menu.prepTime')" class="col-4">
        <el-input-number v-model="model.preparationTime" :disabled="disabled" :min="0" :step="1" />
      </el-form-item>
    </div>

    <el-form-item :label="t('common.status')">
      <el-radio-group v-model="model.status" :disabled="disabled">
        <el-radio v-for="s in statusOptions" :key="s.value" :value="s.value">
          {{ s.label }}
        </el-radio>
      </el-radio-group>
    </el-form-item>

    <el-form-item :label="t('menu.image')">
      <el-input v-model="model.image" :disabled="disabled" placeholder="https://..." />
    </el-form-item>

    <div class="grid">
      <el-form-item :label="t('menu.description') + ' (ZH)'" class="col-6">
        <el-input v-model="model.description.zh" :disabled="disabled" type="textarea" :rows="3" />
      </el-form-item>
      <el-form-item :label="t('menu.description') + ' (EN)'" class="col-6">
        <el-input v-model="model.description.en" :disabled="disabled" type="textarea" :rows="3" />
      </el-form-item>
    </div>

    <el-form-item :label="t('menu.tags')">
      <el-select
        v-model="model.tags"
        :disabled="disabled"
        multiple
        filterable
        allow-create
        default-first-option
      >
        <el-option v-for="p in tagPresets" :key="p.value" :label="p.label" :value="p.value" />
      </el-select>
    </el-form-item>

    <div class="options-header">
      <div class="options-title">{{ t('menu.specOptions') }}</div>
      <el-button type="primary" text :disabled="disabled" @click="addOptionGroup">
        {{ t('menu.addSpec') }}
      </el-button>
    </div>

    <el-card v-for="(opt, optIndex) in model.options" :key="opt.id" class="option-card" shadow="never">
      <div class="grid">
        <el-form-item :label="t('menu.specName') + ' (ZH)'" class="col-6">
          <el-input v-model="opt.name.zh" :disabled="disabled" />
        </el-form-item>
        <el-form-item :label="t('menu.specName') + ' (EN)'" class="col-6">
          <el-input v-model="opt.name.en" :disabled="disabled" />
        </el-form-item>
      </div>

      <div class="option-flags">
        <el-checkbox v-model="opt.required" :disabled="disabled">{{ t('menu.required') }}</el-checkbox>
        <el-checkbox v-model="opt.multiSelect" :disabled="disabled">Multi</el-checkbox>
        <el-button type="danger" text :disabled="disabled" @click="removeOptionGroup(optIndex)">
          {{ t('common.delete') }}
        </el-button>
      </div>

      <div class="values-header">
        <div class="values-title">{{ t('menu.optionValue') }}</div>
        <el-button type="primary" text :disabled="disabled" @click="addOptionValue(optIndex)">
          + {{ t('menu.optionValue') }}
        </el-button>
      </div>

      <div v-for="(val, valIndex) in opt.values" :key="val.id" class="value-row">
        <el-input v-model="val.label.zh" :disabled="disabled" placeholder="ZH" class="value-input" />
        <el-input v-model="val.label.en" :disabled="disabled" placeholder="EN" class="value-input" />
        <el-input-number v-model="val.priceAdjustment" :disabled="disabled" :step="1" />
        <el-button type="danger" text :disabled="disabled" @click="removeOptionValue(optIndex, valIndex)">
          {{ t('common.delete') }}
        </el-button>
      </div>
    </el-card>
  </el-form>
</template>

<style scoped lang="scss">
.dish-form {
  :deep(.el-form-item) {
    margin-bottom: 12px;
  }
}

.grid {
  display: grid;
  grid-template-columns: repeat(12, 1fr);
  gap: 12px;
}

.col-4 {
  grid-column: span 4;
}

.col-6 {
  grid-column: span 6;
}

.options-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin: 8px 0 12px 0;
}

.options-title {
  font-weight: 600;
  color: #303133;
}

.option-card {
  margin-bottom: 12px;
  border-radius: 10px;
}

.option-flags {
  display: flex;
  align-items: center;
  gap: 12px;
  margin: 4px 0 10px 0;
}

.values-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 8px;
}

.values-title {
  font-weight: 600;
  color: #606266;
  font-size: 13px;
}

.value-row {
  display: grid;
  grid-template-columns: 1fr 1fr 140px 80px;
  gap: 8px;
  align-items: center;
  margin-bottom: 8px;
}

.value-input {
  width: 100%;
}
</style>


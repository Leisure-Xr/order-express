<script setup lang="ts">
import { ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'

export type CategoryFormValue = {
  nameZh: string
  nameEn: string
  icon: string
  image: string
  status: 'active' | 'inactive'
}

const props = withDefaults(
  defineProps<{
    modelValue: CategoryFormValue
    disabled?: boolean
  }>(),
  {
    disabled: false,
  },
)

const emit = defineEmits<{
  (e: 'update:modelValue', value: CategoryFormValue): void
}>()

const { t } = useI18n()

function clone<T>(value: T): T {
  return JSON.parse(JSON.stringify(value)) as T
}

const model = ref<CategoryFormValue>(clone(props.modelValue))

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
</script>

<template>
  <el-form :model="model" label-position="top" class="category-form">
    <el-form-item :label="t('menu.category') + ' (ZH)'">
      <el-input v-model="model.nameZh" :disabled="disabled" />
    </el-form-item>
    <el-form-item :label="t('menu.category') + ' (EN)'">
      <el-input v-model="model.nameEn" :disabled="disabled" />
    </el-form-item>
    <el-form-item label="Icon">
      <el-input v-model="model.icon" :disabled="disabled" placeholder="🔥" />
    </el-form-item>
    <el-form-item :label="t('menu.image')">
      <el-input v-model="model.image" :disabled="disabled" placeholder="https://..." />
    </el-form-item>
    <el-form-item :label="t('common.status')">
      <el-radio-group v-model="model.status" :disabled="disabled">
        <el-radio value="active">{{ t('common.yes') }}</el-radio>
        <el-radio value="inactive">{{ t('common.no') }}</el-radio>
      </el-radio-group>
    </el-form-item>
  </el-form>
</template>

<style scoped lang="scss">
.category-form {
  :deep(.el-form-item) {
    margin-bottom: 12px;
  }
}
</style>

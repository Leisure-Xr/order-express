<script setup lang="ts">
import { useI18n } from 'vue-i18n'

export type CategoryFormValue = {
  nameZh: string
  nameEn: string
  icon: string
  image: string
  status: 'active' | 'inactive'
}

const model = defineModel<CategoryFormValue>({ required: true })

const props = withDefaults(
  defineProps<{
    disabled?: boolean
  }>(),
  {
    disabled: false,
  },
)

const { t } = useI18n()
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

<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import { DocumentDelete } from '@element-plus/icons-vue'

const { t } = useI18n()

withDefaults(
  defineProps<{
    icon?: any
    title?: string
    description?: string
  }>(),
  {
    title: '',
    description: '',
  },
)
</script>

<template>
  <div class="empty-state">
    <div class="empty-icon">
      <el-icon :size="64" color="#c0c4cc">
        <component :is="icon || DocumentDelete" />
      </el-icon>
    </div>
    <h3 v-if="title" class="empty-title">{{ title }}</h3>
    <p v-if="description" class="empty-description">{{ description }}</p>
    <p v-if="!title && !description" class="empty-description">{{ t('common.noData') }}</p>
    <div class="empty-action">
      <slot />
    </div>
  </div>
</template>

<style scoped lang="scss">
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 48px 24px;
  text-align: center;
}

.empty-icon {
  margin-bottom: 16px;
  opacity: 0.6;
}

.empty-title {
  font-size: 16px;
  font-weight: 500;
  color: #606266;
  margin: 0 0 8px 0;
}

.empty-description {
  font-size: 14px;
  color: #909399;
  margin: 0 0 16px 0;
  max-width: 300px;
  line-height: 1.5;
}

.empty-action {
  margin-top: 8px;

  &:empty {
    display: none;
  }
}
</style>

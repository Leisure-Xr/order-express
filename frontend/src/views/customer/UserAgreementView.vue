<script setup lang="ts">
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
</script>

<template>
  <div class="legal-page">
    <div class="title">{{ t('routes.userAgreement') }}</div>

    <el-card shadow="never" class="section-card">
      <div class="placeholder-title">{{ t('legal.placeholderTitle') }}</div>
      <div class="placeholder-body">{{ t('legal.userAgreementPlaceholder') }}</div>
    </el-card>
  </div>
</template>

<style scoped lang="scss">
.legal-page {
  padding: 12px;
}

.title {
  font-size: 18px;
  font-weight: 900;
  color: #303133;
  margin-bottom: 12px;
}

.section-card {
  border-radius: 14px;
  border: 1px solid rgba(0, 0, 0, 0.03);
}

.placeholder-title {
  font-weight: 900;
  color: #303133;
  margin-bottom: 8px;
}

.placeholder-body {
  color: #909399;
  font-size: 13px;
  line-height: 1.7;
}
</style>

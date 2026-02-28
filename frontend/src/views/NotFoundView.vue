<script setup lang="ts">
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { WarningFilled } from '@element-plus/icons-vue'

const router = useRouter()
const { t } = useI18n()

function goHome() {
  router.push('/')
}
</script>

<template>
  <div class="not-found-page">
    <el-icon :size="80" color="#c0c4cc">
      <WarningFilled />
    </el-icon>
    <h1 class="error-code">404</h1>
    <p class="error-message">{{ t('routes.notFound') }}</p>
    <el-button type="primary" @click="goHome">
      {{ t('common.back') }}
    </el-button>
  </div>
</template>

<style scoped lang="scss">
.not-found-page {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  background:
    radial-gradient(1200px 600px at 20% 0%, rgba(231, 76, 60, 0.10), transparent 55%),
    radial-gradient(1000px 500px at 100% 0%, rgba(52, 152, 219, 0.08), transparent 50%),
    var(--el-bg-color-page);
  text-align: center;
  padding: 24px;
}

.error-code {
  font-size: 72px;
  font-weight: 800;
  color: #303133;
  margin: 16px 0 8px 0;
  line-height: 1;
}

.error-message {
  font-size: 18px;
  color: #909399;
  margin: 0 0 32px 0;
}
</style>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import { ArrowLeft } from '@element-plus/icons-vue'

defineProps<{
  showBack?: boolean
}>()

const router = useRouter()

function goBack() {
  router.back()
}
</script>

<template>
  <div class="blank-layout">
    <div v-if="showBack" class="blank-header">
      <el-button :icon="ArrowLeft" text @click="goBack" class="back-btn">
      </el-button>
    </div>
    <slot />
  </div>
</template>

<style scoped lang="scss">
.blank-layout {
  min-height: 100vh;
  background:
    radial-gradient(1200px 600px at 20% 0%, rgba(231, 76, 60, 0.10), transparent 55%),
    radial-gradient(1000px 500px at 100% 0%, rgba(52, 152, 219, 0.08), transparent 50%),
    var(--el-bg-color-page);
}

.blank-header {
  position: sticky;
  top: 0;
  z-index: 100;
  display: flex;
  align-items: center;
  padding: 8px 12px;
  background: rgba(255, 255, 255, 0.92);
  backdrop-filter: blur(12px);
  border-bottom: 1px solid rgba(2, 6, 23, 0.06);
}

.back-btn {
  font-size: 16px;
  color: #333333;

  &:hover {
    color: #e74c3c;
  }
}
</style>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import { CircleCheckFilled } from '@element-plus/icons-vue'
import type { Order } from '@/types'
import { getOrderByIdApi } from '@/api/order'
import PriceTag from '@/components/common/PriceTag.vue'

const { t } = useI18n()
const route = useRoute()
const router = useRouter()

const orderId = typeof route.params.orderId === 'string' ? route.params.orderId : ''

const loading = ref(false)
const order = ref<Order | null>(null)

onMounted(async () => {
  if (!orderId) return
  loading.value = true
  try {
    const res = await getOrderByIdApi(orderId)
    if (res.code !== 200 || !res.data) {
      ElMessage.error(res.message || t('common.fail'))
      return
    }
    order.value = res.data
  } catch (e: any) {
    ElMessage.error(e?.message || t('common.fail'))
  } finally {
    loading.value = false
  }
})

function trackOrder() {
  router.push({ name: 'CustomerOrderTracking', params: { id: orderId } })
}

function backToMenu() {
  router.push({ name: 'CustomerMenu' })
}
</script>

<template>
  <div class="result-page">
    <el-skeleton v-if="loading" :rows="6" animated />

    <div v-else-if="!order" class="empty">{{ t('common.noData') }}</div>

    <div v-else class="content">
      <div class="icon-wrap">
        <el-icon :size="86" color="#2ecc71">
          <CircleCheckFilled />
        </el-icon>
      </div>

      <div class="title">{{ t('payment.success') }}</div>

      <el-card shadow="never" class="card">
        <div class="row">
          <span>{{ t('payment.orderNumber') }}</span>
          <strong>{{ order.orderNumber }}</strong>
        </div>
        <div class="row">
          <span>{{ t('payment.amount') }}</span>
          <PriceTag :price="order.total" size="lg" />
        </div>
      </el-card>

      <div class="btns">
        <el-button type="primary" size="large" @click="trackOrder">{{ t('payment.trackOrder') }}</el-button>
        <el-button size="large" @click="backToMenu">{{ t('payment.backToMenu') }}</el-button>
      </div>
    </div>
  </div>
</template>

<style scoped lang="scss">
.result-page {
  padding: 16px;
}

.content {
  min-height: calc(100vh - 32px);
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 14px;
  text-align: center;
}

.icon-wrap {
  position: relative;
  animation: successBounce 0.6s cubic-bezier(0.34, 1.56, 0.64, 1);

  &::before {
    content: '';
    position: absolute;
    inset: -20px;
    border-radius: 50%;
    background: radial-gradient(circle, rgba(46, 204, 113, 0.12) 0%, transparent 70%);
    pointer-events: none;
  }
}

.title {
  font-size: 22px;
  font-weight: 900;
  color: #303133;
  animation: fadeSlideUp 0.4s cubic-bezier(0.4, 0, 0.2, 1) 0.2s both;
}

.card {
  width: 100%;
  max-width: 460px;
  border-radius: 16px;
  border: 1px solid rgba(0, 0, 0, 0.03);
  animation: fadeSlideUp 0.4s cubic-bezier(0.4, 0, 0.2, 1) 0.3s both;
}

.row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  color: #606266;
  margin-bottom: 8px;
}

.row strong {
  color: #303133;
}

.btns {
  width: 100%;
  max-width: 460px;
  display: flex;
  gap: 12px;
  animation: fadeSlideUp 0.4s cubic-bezier(0.4, 0, 0.2, 1) 0.4s both;
}

.btns :deep(.el-button) {
  flex: 1;
}
</style>


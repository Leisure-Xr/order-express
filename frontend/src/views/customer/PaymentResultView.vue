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
      <el-icon :size="86" color="#2ecc71">
        <CircleCheckFilled />
      </el-icon>

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

.title {
  font-size: 20px;
  font-weight: 900;
  color: #303133;
}

.card {
  width: 100%;
  max-width: 460px;
  border-radius: 12px;
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
}

.btns :deep(.el-button) {
  flex: 1;
}
</style>


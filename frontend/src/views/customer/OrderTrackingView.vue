<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import type { Order, OrderStatus } from '@/types'
import { getOrderByIdApi } from '@/api/order'
import { useOrderPolling } from '@/composables/useOrderPolling'
import { formatDate } from '@/utils/format'

const { t } = useI18n()
const route = useRoute()
const { startPolling, stopPolling } = useOrderPolling()

const orderId = computed(() => (typeof route.params.id === 'string' ? route.params.id : ''))

const order = ref<Order | null>(null)
const loading = ref(false)

function statusTagType(status: OrderStatus): 'info' | 'success' | 'warning' | 'danger' {
  if (status === 'pending') return 'warning'
  if (status === 'confirmed' || status === 'preparing' || status === 'ready') return 'info'
  if (status === 'completed' || status === 'delivered') return 'success'
  return 'danger'
}

async function fetchOrder() {
  if (!orderId.value) return
  loading.value = true
  try {
    const res = await getOrderByIdApi(orderId.value)
    if (res.code !== 200 || !res.data) {
      ElMessage.error(res.message || t('common.fail'))
      order.value = null
      stopPolling()
      return
    }
    order.value = res.data
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  startPolling(fetchOrder, 3000)
})
</script>

<template>
  <div class="tracking-page">
    <div class="title">{{ t('order.tracking.title') }}</div>

    <el-skeleton v-if="loading && !order" :rows="8" animated />

    <div v-else-if="!order" class="empty">{{ t('common.noData') }}</div>

    <div v-else class="content">
      <el-card shadow="never" class="card">
        <div class="row">
          <span>{{ t('payment.orderNumber') }}</span>
          <strong>{{ order.orderNumber }}</strong>
        </div>
        <div class="row">
          <span>{{ t('common.status') }}</span>
          <el-tag :type="statusTagType(order.status)">{{ t(`order.status.${order.status}`) }}</el-tag>
        </div>
      </el-card>

      <el-card shadow="never" class="card">
        <el-timeline>
          <el-timeline-item
            v-for="(h, idx) in order.statusHistory"
            :key="idx"
            :timestamp="formatDate(h.timestamp)"
            :type="idx === order.statusHistory.length - 1 ? 'primary' : 'info'"
          >
            {{ t(`order.status.${h.status}`) }}
            <span v-if="h.note">- {{ h.note }}</span>
          </el-timeline-item>
        </el-timeline>
      </el-card>
    </div>
  </div>
</template>

<style scoped lang="scss">
.tracking-page {
  padding: 12px;
}

.title {
  font-size: 18px;
  font-weight: 900;
  color: #303133;
  margin-bottom: 12px;
}

.content {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.card {
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
</style>


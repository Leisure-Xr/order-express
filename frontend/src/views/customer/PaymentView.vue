<script setup lang="ts">
import { onMounted, onUnmounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import type { Order, PaymentMethod, PaymentStatus } from '@/types'
import { getOrderByIdApi } from '@/api/order'
import { initiatePaymentApi, checkPaymentStatusApi } from '@/api/payment'
import PriceTag from '@/components/common/PriceTag.vue'

const { t } = useI18n()
const route = useRoute()
const router = useRouter()

const orderId = typeof route.params.orderId === 'string' ? route.params.orderId : ''

const loading = ref(false)
const paying = ref(false)
const order = ref<Order | null>(null)
const method = ref<PaymentMethod>('wechat')
const paymentId = ref<string | null>(null)
const paymentStatus = ref<PaymentStatus | ''>('')

let pollTimer: number | null = null

function stopPolling() {
  if (pollTimer !== null) {
    window.clearInterval(pollTimer)
    pollTimer = null
  }
}

async function pollPayment() {
  if (!paymentId.value) return
  const res = await checkPaymentStatusApi(paymentId.value)
  if (res.code !== 200 || !res.data) return

  paymentStatus.value = res.data.status

  if (res.data.status === 'paid') {
    stopPolling()
    router.replace({ name: 'CustomerPaymentResult', params: { orderId } })
  }
}

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
    method.value = res.data.payment.method

    if (res.data.payment.status === 'paid') {
      router.replace({ name: 'CustomerPaymentResult', params: { orderId } })
    }
  } catch (e: any) {
    ElMessage.error(e?.message || t('common.fail'))
  } finally {
    loading.value = false
  }
})

onUnmounted(() => {
  stopPolling()
})

async function payNow() {
  if (!order.value) return

  paying.value = true
  try {
    const res = await initiatePaymentApi(order.value.id, method.value, order.value.total)
    if (res.code !== 200 || !res.data) {
      ElMessage.error(res.message || t('common.fail'))
      return
    }

    paymentId.value = res.data.paymentId
    paymentStatus.value = res.data.status

    if (res.data.status === 'paid') {
      stopPolling()
      router.replace({ name: 'CustomerPaymentResult', params: { orderId } })
      return
    }

    stopPolling()
    pollTimer = window.setInterval(() => {
      pollPayment()
    }, 500)
  } catch (e: any) {
    ElMessage.error(e?.message || t('common.fail'))
  } finally {
    paying.value = false
  }
}
</script>

<template>
  <div class="payment-page">
    <div class="title">{{ t('routes.payment') }}</div>

    <el-skeleton v-if="loading" :rows="6" animated />

    <div v-else-if="!order" class="empty">{{ t('common.noData') }}</div>

    <div v-else class="content">
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

      <el-card shadow="never" class="card">
        <div class="section-title">{{ t('payment.selectMethod') }}</div>
        <el-radio-group v-model="method">
          <el-radio value="wechat">{{ t('payment.wechat') }}</el-radio>
          <el-radio value="alipay">{{ t('payment.alipay') }}</el-radio>
          <el-radio value="cash">{{ t('payment.cash') }}</el-radio>
        </el-radio-group>

        <div v-if="paymentId" class="processing">
          <el-tag v-if="paymentStatus === 'processing'" type="warning">{{ t('payment.processing') }}</el-tag>
          <el-tag v-else-if="paymentStatus === 'paid'" type="success">{{ t('payment.success') }}</el-tag>
          <el-tag v-else type="info">{{ paymentStatus }}</el-tag>
        </div>
      </el-card>

      <el-button type="primary" size="large" class="pay-btn" :loading="paying" @click="payNow">
        {{ t('order.payNow') }}
      </el-button>
    </div>
  </div>
</template>

<style scoped lang="scss">
.payment-page {
  padding: 12px;
}

.title {
  font-size: 18px;
  font-weight: 800;
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
  margin-bottom: 8px;
  color: #606266;
}

.row strong {
  color: #303133;
}

.section-title {
  font-weight: 700;
  margin-bottom: 10px;
  color: #303133;
}

.processing {
  margin-top: 12px;
  display: flex;
  align-items: center;
  gap: 8px;
}

.status {
  color: #909399;
  font-size: 12px;
}

.pay-btn {
  width: 100%;
}
</style>

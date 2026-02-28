<script setup lang="ts">
import { computed, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { useCartStore } from '@/stores/cart'
import { useAppStore } from '@/stores/app'
import { useOrderStore } from '@/stores/order'
import type { CreateOrderPayload, OrderType, PaymentMethod } from '@/types'
import OrderTypeSelector from '@/components/customer/OrderTypeSelector.vue'
import TableSelector from '@/components/customer/TableSelector.vue'

const { t } = useI18n()
const router = useRouter()
const cartStore = useCartStore()
const appStore = useAppStore()
const orderStore = useOrderStore()

const orderType = ref<OrderType>('dine_in')
const tableNumber = ref(appStore.currentTableId || '')
const deliveryAddress = ref('')
const contactPhone = ref('')
const remarks = ref('')
const paymentMethod = ref<PaymentMethod>('wechat')

const submitting = ref(false)

const deliveryFee = computed(() => (orderType.value === 'takeaway' ? 5 : 0))
const total = computed(() => cartStore.subtotal + deliveryFee.value)

function validate(): string | null {
  if (cartStore.isEmpty) return t('cart.emptyCart')

  if (orderType.value === 'dine_in') {
    if (!tableNumber.value.trim()) return t('order.selectTable')
  }

  if (orderType.value === 'takeaway') {
    if (!deliveryAddress.value.trim()) return t('order.deliveryAddress')
    if (!contactPhone.value.trim()) return t('order.contactPhone')
  }

  if (orderType.value === 'pickup') {
    if (!contactPhone.value.trim()) return t('order.contactPhone')
  }

  return null
}

async function submit() {
  const err = validate()
  if (err) {
    ElMessage.warning(err)
    return
  }

  const payload: CreateOrderPayload = {
    type: orderType.value,
    tableId: orderType.value === 'dine_in' ? tableNumber.value.trim() : undefined,
    items: cartStore.items.map((i) => ({
      dishId: i.dishId,
      quantity: i.quantity,
      selectedOptions: i.selectedOptions,
    })),
    remarks: remarks.value.trim() || undefined,
    deliveryAddress: orderType.value === 'takeaway' ? deliveryAddress.value.trim() : undefined,
    contactPhone: orderType.value !== 'dine_in' ? contactPhone.value.trim() || undefined : undefined,
    paymentMethod: paymentMethod.value,
  }

  submitting.value = true
  try {
    const order = await orderStore.createOrder(payload)
    if (!order) return

    cartStore.clearCart()
    router.replace({ name: 'CustomerPayment', params: { orderId: order.id } })
  } catch (e: any) {
    ElMessage.error(e?.message || t('common.fail'))
  } finally {
    submitting.value = false
  }
}
</script>

<template>
  <div class="checkout-page">
    <div class="title">{{ t('routes.checkout') }}</div>

    <el-card shadow="never" class="card">
      <div class="section-title">{{ t('order.orderType') }}</div>
      <OrderTypeSelector v-model="orderType" />

      <div v-if="orderType === 'dine_in'" class="section">
        <div class="section-title">{{ t('order.tableNumber') }}</div>
        <TableSelector v-model="tableNumber" />
      </div>

      <div v-else-if="orderType === 'takeaway'" class="section">
        <div class="section-title">{{ t('order.deliveryAddress') }}</div>
        <el-input v-model="deliveryAddress" :placeholder="t('order.deliveryAddress')" />

        <div class="section-title" style="margin-top: 12px">{{ t('order.contactPhone') }}</div>
        <el-input v-model="contactPhone" :placeholder="t('order.contactPhone')" />
      </div>

      <div v-else class="section">
        <div class="section-title">{{ t('order.contactPhone') }}</div>
        <el-input v-model="contactPhone" :placeholder="t('order.contactPhone')" />
      </div>

      <div class="section">
        <div class="section-title">{{ t('order.remark') }}</div>
        <el-input v-model="remarks" type="textarea" :rows="3" :placeholder="t('order.remarkPlaceholder')" />
      </div>
    </el-card>

    <el-card shadow="never" class="card">
      <div class="section-title">{{ t('payment.selectMethod') }}</div>
      <el-radio-group v-model="paymentMethod">
        <el-radio value="wechat">{{ t('payment.wechat') }}</el-radio>
        <el-radio value="alipay">{{ t('payment.alipay') }}</el-radio>
        <el-radio value="cash">{{ t('payment.cash') }}</el-radio>
      </el-radio-group>
    </el-card>

    <div class="summary">
      <div class="row">
        <span>{{ t('cart.subtotal') }}</span>
        <span>¥{{ cartStore.subtotal.toFixed(2) }}</span>
      </div>
      <div class="row">
        <span>{{ t('order.deliveryFee') }}</span>
        <span>¥{{ deliveryFee.toFixed(2) }}</span>
      </div>
      <div class="row total">
        <span>{{ t('common.total') }}</span>
        <span class="total-price">¥{{ total.toFixed(2) }}</span>
      </div>

      <el-button type="primary" size="large" class="pay-btn" :loading="submitting" @click="submit">
        {{ t('order.payNow') }}
      </el-button>
    </div>
  </div>
</template>

<style scoped lang="scss">
.checkout-page {
  padding: 12px;
}

.title {
  font-size: 18px;
  font-weight: 800;
  color: #303133;
  margin-bottom: 12px;
}

.card {
  border-radius: 12px;
  margin-bottom: 12px;
}

.section {
  margin-top: 14px;
}

.section-title {
  font-size: 13px;
  font-weight: 700;
  color: #303133;
  margin-bottom: 8px;
}

.summary {
  background: #ffffff;
  border-radius: 12px;
  padding: 14px;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.06);
}

.row {
  display: flex;
  justify-content: space-between;
  color: #606266;
  font-size: 14px;
  margin-bottom: 8px;
}

.row.total {
  margin-top: 10px;
  padding-top: 10px;
  border-top: 1px dashed #e5e7eb;
  font-weight: 800;
  color: #303133;
}

.total-price {
  color: #e74c3c;
}

.pay-btn {
  width: 100%;
  margin-top: 10px;
}
</style>

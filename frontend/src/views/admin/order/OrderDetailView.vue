<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { useOrderStore } from '@/stores/order'
import { useLocaleText } from '@/composables/useLocaleText'
import AdminPageHeader from '@/components/admin/AdminPageHeader.vue'
import { formatDate } from '@/utils/format'
import type { Order, OrderStatus } from '@/types'

const { t } = useI18n()
const route = useRoute()
const router = useRouter()
const orderStore = useOrderStore()
const { localText } = useLocaleText()

const orderId = computed(() => (typeof route.params.id === 'string' ? route.params.id : ''))
const order = computed<Order | null>(() => orderStore.currentOrder)

const loading = ref(false)
const updating = ref(false)

const VALID_TRANSITIONS: Record<OrderStatus, OrderStatus[]> = {
  pending: ['confirmed', 'cancelled'],
  confirmed: ['preparing', 'cancelled'],
  preparing: ['ready', 'cancelled'],
  ready: ['completed', 'delivered'],
  completed: [],
  delivered: [],
  cancelled: [],
}

const availableActions = computed(() => {
  if (!order.value) return []
  return VALID_TRANSITIONS[order.value.status] ?? []
})

onMounted(async () => {
  if (!orderId.value) return
  loading.value = true
  try {
    const data = await orderStore.fetchOrderById(orderId.value)
    if (!data) {
      ElMessage.error(t('common.noData'))
      router.replace({ name: 'AdminOrderList' })
    }
  } catch (e: any) {
    ElMessage.error(e?.message || t('common.fail'))
  } finally {
    loading.value = false
  }
})

function statusTagType(status: OrderStatus): 'info' | 'success' | 'warning' | 'danger' {
  if (status === 'pending') return 'warning'
  if (status === 'confirmed' || status === 'preparing' || status === 'ready') return 'info'
  if (status === 'completed' || status === 'delivered') return 'success'
  return 'danger'
}

async function updateStatus(status: OrderStatus) {
  if (!order.value) return
  updating.value = true
  try {
    await orderStore.updateOrderStatus(order.value.id, status)
    ElMessage.success(t('common.success'))
  } catch (e: any) {
    ElMessage.error(e?.message || t('common.fail'))
  } finally {
    updating.value = false
  }
}
</script>

<template>
  <div class="order-detail admin-page">
    <AdminPageHeader :title="t('routes.orderDetail')">
      <template #actions>
        <el-button @click="router.back()">{{ t('common.back') }}</el-button>
      </template>
    </AdminPageHeader>

    <el-skeleton v-if="loading" :rows="8" animated />

    <div v-else-if="!order" class="empty">
      {{ t('common.noData') }}
    </div>

    <div v-else class="content">
      <el-card shadow="never" class="card">
        <div class="info-grid">
          <div class="info-item">
            <div class="label">{{ t('payment.orderNumber') }}</div>
            <div class="value">{{ order.orderNumber }}</div>
          </div>
          <div class="info-item">
            <div class="label">{{ t('order.orderType') }}</div>
            <div class="value">
              {{
                order.type === 'dine_in'
                  ? t('order.dineIn')
                  : order.type === 'takeaway'
                    ? t('order.takeout')
                    : t('order.pickup')
              }}
            </div>
          </div>
          <div class="info-item">
            <div class="label">{{ t('order.tableNumber') }}</div>
            <div class="value">{{ order.tableId || '-' }}</div>
          </div>
          <div class="info-item">
            <div class="label">{{ t('common.status') }}</div>
            <div class="value">
              <el-tag :type="statusTagType(order.status)">{{ t(`order.status.${order.status}`) }}</el-tag>
            </div>
          </div>
          <div class="info-item">
            <div class="label">{{ t('order.payment') }}</div>
            <div class="value">
              {{ order.payment.method }} /
              <el-tag :type="order.payment.status === 'paid' ? 'success' : 'info'">
                {{ order.payment.status }}
              </el-tag>
            </div>
          </div>
          <div class="info-item">
            <div class="label">{{ t('order.createdAt') }}</div>
            <div class="value">{{ formatDate(order.createdAt) }}</div>
          </div>
        </div>

        <div v-if="availableActions.length" class="status-actions">
          <div class="label">{{ t('common.action') }}</div>
          <div class="btns">
            <el-button
              v-for="s in availableActions"
              :key="s"
              type="primary"
              :loading="updating"
              @click="updateStatus(s)"
            >
              {{ t(`order.status.${s}`) }}
            </el-button>
          </div>
        </div>
      </el-card>

      <el-card shadow="never" class="card">
        <template #header>
          <div class="card-title">{{ t('order.items') }}</div>
        </template>

        <div class="admin-table-scroll">
          <el-table :data="order.items" stripe>
            <el-table-column :label="t('dashboard.menuItem')" min-width="220">
              <template #default="{ row }">
                <div class="dish-name">{{ localText(row.dishName) }}</div>
                <div v-if="row.selectedOptions.length" class="dish-opts">
                  <span v-for="(opt, idx) in row.selectedOptions" :key="idx" class="opt">
                    {{ opt.optionName }}: {{ opt.valueName }}
                  </span>
                </div>
              </template>
            </el-table-column>
            <el-table-column :label="t('order.qty')" width="90">
              <template #default="{ row }">{{ row.quantity }}</template>
            </el-table-column>
            <el-table-column :label="t('order.unitPrice')" width="120">
              <template #default="{ row }">¥{{ row.unitPrice.toFixed(2) }}</template>
            </el-table-column>
            <el-table-column :label="t('order.subtotal')" width="140">
              <template #default="{ row }">¥{{ row.subtotal.toFixed(2) }}</template>
            </el-table-column>
          </el-table>
        </div>

        <div class="totals">
          <div class="row">
            <span>{{ t('cart.subtotal') }}</span>
            <span>¥{{ order.subtotal.toFixed(2) }}</span>
          </div>
          <div class="row">
            <span>{{ t('order.deliveryFee') }}</span>
            <span>¥{{ order.deliveryFee.toFixed(2) }}</span>
          </div>
          <div class="row total">
            <span>{{ t('common.total') }}</span>
            <span class="price">¥{{ order.total.toFixed(2) }}</span>
          </div>
        </div>
      </el-card>

      <el-card shadow="never" class="card">
        <template #header>
          <div class="card-title">{{ t('order.tracking.title') }}</div>
        </template>

        <el-timeline>
          <el-timeline-item
            v-for="(h, idx) in order.statusHistory"
            :key="idx"
            :timestamp="formatDate(h.timestamp)"
          >
            {{ t(`order.status.${h.status}`) }} <span v-if="h.note">- {{ h.note }}</span>
          </el-timeline-item>
        </el-timeline>
      </el-card>
    </div>
  </div>
</template>

<style scoped lang="scss">
.content {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.card {
  border-radius: 14px;
  transition: box-shadow var(--app-transition-base);

  &:hover {
    box-shadow: 0 4px 16px rgba(0, 0, 0, 0.06);
  }
}

.card-title {
  font-weight: 700;
}

.info-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 12px;
}

.info-item {
  background: #f9fafb;
  border-radius: 10px;
  padding: 10px 12px;
}

.info-item .label {
  font-size: 11px;
  color: #909399;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.info-item .value {
  margin-top: 4px;
  font-weight: 700;
  color: #303133;
}

.status-actions {
  margin-top: 14px;
}

.status-actions .btns {
  margin-top: 8px;
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.dish-name {
  font-weight: 700;
  color: #303133;
}

.dish-opts {
  margin-top: 6px;
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  color: #909399;
  font-size: 12px;
}

.opt {
  background: #f5f7fa;
  border-radius: 999px;
  padding: 2px 8px;
}

.totals {
  margin-top: 12px;
  max-width: 360px;
  margin-left: auto;
}

.row {
  display: flex;
  justify-content: space-between;
  color: #606266;
  margin-bottom: 6px;
}

.row.total {
  padding-top: 8px;
  border-top: 1.5px dashed rgba(231, 76, 60, 0.2);
  font-weight: 900;
  color: #303133;
}

.price {
  color: #e74c3c;
  font-size: 16px;
}

@media (max-width: 900px) {
  .info-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 600px) {
  .info-grid {
    grid-template-columns: 1fr;
  }
}
</style>

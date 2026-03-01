<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { useOrderStore } from '@/stores/order'
import { useOrderPolling } from '@/composables/useOrderPolling'
import AdminPageHeader from '@/components/admin/AdminPageHeader.vue'
import { formatDate } from '@/utils/format'
import type { OrderStatus } from '@/types'

const { t } = useI18n()
const router = useRouter()
const orderStore = useOrderStore()
const { startPolling, stopPolling } = useOrderPolling()

const activeStatus = ref<OrderStatus | 'all'>('all')

const statusTabs = computed(() => [
  { key: 'all', label: t('common.all') },
  { key: 'pending', label: t('order.status.pending') },
  { key: 'confirmed', label: t('order.status.confirmed') },
  { key: 'preparing', label: t('order.status.preparing') },
  { key: 'ready', label: t('order.status.ready') },
  { key: 'completed', label: t('order.status.completed') },
  { key: 'delivered', label: t('order.status.delivered') },
  { key: 'cancelled', label: t('order.status.cancelled') },
])

function statusTagType(status: OrderStatus): 'info' | 'success' | 'warning' | 'danger' {
  if (status === 'pending') return 'warning'
  if (status === 'confirmed' || status === 'preparing' || status === 'ready') return 'info'
  if (status === 'completed' || status === 'delivered') return 'success'
  return 'danger'
}

async function fetchList() {
  try {
    await orderStore.fetchOrders({ status: activeStatus.value, page: 1, pageSize: 20 })
  } catch (e: any) {
    ElMessage.error(e?.message || t('common.fail'))
  }
}

onMounted(() => {
  startPolling(fetchList, 10000)
})

watch(activeStatus, () => {
  fetchList()
})

function openDetail(id: string) {
  stopPolling()
  router.push({ name: 'AdminOrderDetail', params: { id } })
}
</script>

<template>
  <div class="order-list admin-page">
    <AdminPageHeader :title="t('routes.orderManagement')" />

    <el-tabs v-model="activeStatus" class="tabs">
      <el-tab-pane v-for="tab in statusTabs" :key="tab.key" :label="tab.label" :name="tab.key" />
    </el-tabs>

    <el-card shadow="never" class="table-shell" style="--el-card-padding: 0px">
      <div class="admin-table-scroll">
        <el-table :data="orderStore.orders" v-loading="orderStore.listLoading" stripe>
          <template #empty>
            <el-empty :description="t('common.noData')" />
          </template>

          <el-table-column label="Order" min-width="180">
            <template #default="{ row }">
              <div class="order-no">{{ row.orderNumber }}</div>
              <div class="order-time">{{ formatDate(row.createdAt) }}</div>
            </template>
          </el-table-column>

          <el-table-column :label="t('order.orderType')" width="120">
            <template #default="{ row }">
              {{
                row.type === 'dine_in'
                  ? t('order.dineIn')
                  : row.type === 'takeaway'
                    ? t('order.takeout')
                    : t('order.pickup')
              }}
            </template>
          </el-table-column>

          <el-table-column :label="t('order.tableNumber')" width="140">
            <template #default="{ row }">
              {{ row.tableId || '-' }}
            </template>
          </el-table-column>

          <el-table-column :label="t('common.status')" width="140">
            <template #default="{ row }">
              <el-tag :type="statusTagType(row.status)">
                {{ t(`order.status.${row.status}`) }}
              </el-tag>
            </template>
          </el-table-column>

          <el-table-column :label="t('payment.amount')" width="140">
            <template #default="{ row }">¥{{ row.total.toFixed(2) }}</template>
          </el-table-column>

          <el-table-column label="Payment" width="140">
            <template #default="{ row }">
              <el-tag :type="row.payment.status === 'paid' ? 'success' : row.payment.status === 'processing' ? 'warning' : 'info'">
                {{ row.payment.status }}
              </el-tag>
            </template>
          </el-table-column>

          <el-table-column :label="t('common.action')" width="140" fixed="right">
            <template #default="{ row }">
              <el-button type="primary" link @click="openDetail(row.id)">
                {{ t('routes.orderDetail') }}
              </el-button>
            </template>
          </el-table-column>
        </el-table>
      </div>
    </el-card>
  </div>
</template>

<style scoped lang="scss">
.order-no {
  font-weight: 800;
  color: #303133;
}

.order-time {
  font-size: 12px;
  color: #909399;
}

.tabs {
  margin-bottom: 12px;
}

.table-shell {
  overflow: hidden;
}
</style>

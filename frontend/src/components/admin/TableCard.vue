<script setup lang="ts">
import type { Table } from '@/types'
import { useI18n } from 'vue-i18n'
import { Edit, Delete, Link } from '@element-plus/icons-vue'

const props = defineProps<{
  table: Table
}>()

const emit = defineEmits<{
  (e: 'edit', table: Table): void
  (e: 'delete', table: Table): void
  (e: 'qrcode', table: Table): void
}>()

const { t } = useI18n()

function statusTagType(status: Table['status']): 'success' | 'warning' | 'info' | 'danger' {
  if (status === 'available') return 'success'
  if (status === 'occupied') return 'warning'
  if (status === 'reserved') return 'info'
  return 'danger'
}
</script>

<template>
  <el-card shadow="hover" class="table-card">
    <div class="top">
      <div class="number">{{ table.number }}</div>
      <el-tag :type="statusTagType(table.status)">
        {{ t(`table.${table.status === 'available' ? 'free' : table.status === 'occupied' ? 'inUse' : table.status}`) }}
      </el-tag>
    </div>

    <div class="meta">
      <div class="line">
        <span class="label">{{ t('table.seats') }}</span>
        <span class="value">{{ table.seats }}</span>
      </div>
      <div class="line" v-if="table.area">
        <span class="label">{{ t('table.area') }}</span>
        <span class="value">{{ table.area }}</span>
      </div>
      <div class="line" v-if="table.currentOrderId">
        <span class="label">Order</span>
        <span class="value">{{ table.currentOrderId }}</span>
      </div>
    </div>

    <div class="actions">
      <el-button size="small" :icon="Link" @click="emit('qrcode', table)">{{ t('table.viewQRCode') }}</el-button>
      <el-button size="small" :icon="Edit" @click="emit('edit', table)">{{ t('common.edit') }}</el-button>
      <el-button size="small" type="danger" :icon="Delete" @click="emit('delete', table)">{{ t('common.delete') }}</el-button>
    </div>
  </el-card>
</template>

<style scoped lang="scss">
.table-card {
  border-radius: 12px;
}

.top {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10px;
}

.number {
  font-weight: 900;
  color: #303133;
  font-size: 18px;
}

.meta .line {
  display: flex;
  justify-content: space-between;
  color: #606266;
  font-size: 13px;
  margin-bottom: 6px;
}

.meta .label {
  color: #909399;
}

.actions {
  margin-top: 10px;
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}
</style>

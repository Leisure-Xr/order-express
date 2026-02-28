<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import { useTableStore } from '@/stores/table'
import TableCard from '@/components/admin/TableCard.vue'
import QRCodeDialog from '@/components/admin/QRCodeDialog.vue'
import type { Table } from '@/types'

const { t } = useI18n()
const tableStore = useTableStore()

const formVisible = ref(false)
const formSubmitting = ref(false)
const editingId = ref<string | null>(null)

const formModel = ref({
  number: '',
  seats: 4,
  area: 'Main Hall',
  status: 'available' as Table['status'],
})

const qrVisible = ref(false)
const qrTableNumber = ref('')
const qrLink = ref('')

const areaGroups = computed(() => {
  return Array.from(tableStore.tablesByArea.entries()).map(([area, tables]) => ({
    area,
    tables,
  }))
})

onMounted(async () => {
  try {
    await tableStore.fetchTables()
  } catch (e: any) {
    ElMessage.error(e?.message || t('common.fail'))
  }
})

function openCreate() {
  editingId.value = null
  formModel.value = { number: '', seats: 4, area: 'Main Hall', status: 'available' }
  formVisible.value = true
}

function openEdit(table: Table) {
  editingId.value = table.id
  formModel.value = {
    number: table.number,
    seats: table.seats,
    area: table.area || '',
    status: table.status,
  }
  formVisible.value = true
}

async function submit() {
  if (!formModel.value.number.trim()) {
    ElMessage.warning(t('table.tableNumber'))
    return
  }
  formSubmitting.value = true
  try {
    if (!editingId.value) {
      await tableStore.createTable({
        number: formModel.value.number.trim(),
        seats: formModel.value.seats,
        area: formModel.value.area.trim() || undefined,
        status: formModel.value.status,
      })
      ElMessage.success(t('common.success'))
    } else {
      await tableStore.updateTable(editingId.value, {
        number: formModel.value.number.trim(),
        seats: formModel.value.seats,
        area: formModel.value.area.trim() || undefined,
        status: formModel.value.status,
      })
      ElMessage.success(t('common.success'))
    }
    formVisible.value = false
  } catch (e: any) {
    ElMessage.error(e?.message || t('common.fail'))
  } finally {
    formSubmitting.value = false
  }
}

async function remove(table: Table) {
  try {
    await ElMessageBox.confirm(
      `${t('common.confirm')} ${t('common.delete')}?`,
      t('common.confirm'),
      { type: 'warning' },
    )
  } catch {
    return
  }

  const ok = await tableStore.deleteTable(table.id)
  if (ok) ElMessage.success(t('common.success'))
  else ElMessage.error(t('common.fail'))
}

function openQr(table: Table) {
  qrTableNumber.value = table.number
  qrLink.value = tableStore.getOrderLinkForTableNumber(table.number)
  qrVisible.value = true
}
</script>

<template>
  <div class="table-manage">
    <div class="toolbar">
      <div class="title">{{ t('routes.tableManagement') }}</div>
      <el-button type="primary" :icon="Plus" @click="openCreate">
        {{ t('table.addTable') }}
      </el-button>
    </div>

    <el-skeleton v-if="tableStore.loading && !tableStore.tables.length" :rows="8" animated />

    <div v-else class="areas">
      <div v-for="g in areaGroups" :key="g.area" class="area">
        <div class="area-title">{{ g.area }}</div>
        <div class="grid">
          <TableCard
            v-for="tb in g.tables"
            :key="tb.id"
            :table="tb"
            @edit="openEdit"
            @delete="remove"
            @qrcode="openQr"
          />
        </div>
      </div>
    </div>

    <el-dialog v-model="formVisible" :title="editingId ? t('common.edit') : t('table.addTable')" width="520px">
      <el-form label-position="top">
        <el-form-item :label="t('table.tableNumber')">
          <el-input v-model="formModel.number" placeholder="T05" />
        </el-form-item>
        <el-form-item :label="t('table.seats')">
          <el-input-number v-model="formModel.seats" :min="1" :max="20" />
        </el-form-item>
        <el-form-item :label="t('table.area')">
          <el-input v-model="formModel.area" />
        </el-form-item>
        <el-form-item :label="t('common.status')">
          <el-select v-model="formModel.status" style="width: 220px">
            <el-option value="available" :label="t('table.free')" />
            <el-option value="occupied" :label="t('table.inUse')" />
            <el-option value="reserved" :label="t('table.reserved')" />
            <el-option value="disabled" :label="t('table.disabled')" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="formVisible = false">{{ t('common.cancel') }}</el-button>
        <el-button type="primary" :loading="formSubmitting" @click="submit">
          {{ t('common.save') }}
        </el-button>
      </template>
    </el-dialog>

    <QRCodeDialog v-model="qrVisible" :table-number="qrTableNumber" :link="qrLink" />
  </div>
</template>

<style scoped lang="scss">
.toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 16px;
}

.title {
  font-size: 18px;
  font-weight: 800;
  color: #303133;
}

.areas {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.area-title {
  font-size: 14px;
  font-weight: 800;
  color: #303133;
  margin-bottom: 10px;
}

.grid {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 12px;
}

@media (max-width: 1200px) {
  .grid {
    grid-template-columns: repeat(3, minmax(0, 1fr));
  }
}

@media (max-width: 900px) {
  .grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}

@media (max-width: 600px) {
  .grid {
    grid-template-columns: 1fr;
  }
}
</style>


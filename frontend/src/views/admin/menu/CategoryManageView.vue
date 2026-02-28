<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ArrowUp, ArrowDown, Plus, Edit, Delete } from '@element-plus/icons-vue'
import { useMenuStore } from '@/stores/menu'
import CategoryForm, { type CategoryFormValue } from '@/components/admin/CategoryForm.vue'
import type { Category } from '@/types'

const { t } = useI18n()
const menuStore = useMenuStore()

const dialogVisible = ref(false)
const dialogSubmitting = ref(false)
const editingId = ref<string | null>(null)

const formValue = ref<CategoryFormValue>({
  nameZh: '',
  nameEn: '',
  icon: '',
  image: '',
  status: 'active',
})

const sortedCategories = computed(() =>
  [...menuStore.categories].sort((a, b) => a.sortOrder - b.sortOrder),
)

onMounted(async () => {
  try {
    await menuStore.fetchCategories()
  } catch (e: any) {
    ElMessage.error(e?.message || 'Failed to load categories')
  }
})

function openCreate() {
  editingId.value = null
  formValue.value = { nameZh: '', nameEn: '', icon: '', image: '', status: 'active' }
  dialogVisible.value = true
}

function openEdit(category: Category) {
  editingId.value = category.id
  formValue.value = {
    nameZh: category.name.zh,
    nameEn: category.name.en,
    icon: category.icon || '',
    image: category.image || '',
    status: category.status,
  }
  dialogVisible.value = true
}

async function submit() {
  dialogSubmitting.value = true
  try {
    const payload = {
      name: { zh: formValue.value.nameZh.trim(), en: formValue.value.nameEn.trim() },
      icon: formValue.value.icon.trim() || undefined,
      image: formValue.value.image.trim() || undefined,
      status: formValue.value.status,
    }

    if (!payload.name.zh || !payload.name.en) {
      ElMessage.warning('Please input both ZH/EN names')
      return
    }

    if (!editingId.value) {
      await menuStore.createCategory({
        ...payload,
        sortOrder: sortedCategories.value.length + 1,
      })
      ElMessage.success(t('common.success'))
    } else {
      await menuStore.updateCategory(editingId.value, payload)
      ElMessage.success(t('common.success'))
    }

    dialogVisible.value = false
  } catch (e: any) {
    ElMessage.error(e?.message || t('common.fail'))
  } finally {
    dialogSubmitting.value = false
  }
}

async function removeCategory(category: Category) {
  try {
    await ElMessageBox.confirm(
      `${t('common.confirm')} ${t('common.delete')}?`,
      t('common.confirm'),
      { type: 'warning' },
    )
  } catch {
    return
  }

  const ok = await menuStore.deleteCategory(category.id)
  if (ok) {
    ElMessage.success(t('common.success'))
  } else {
    ElMessage.error(t('common.fail'))
  }
}

async function move(category: Category, direction: 'up' | 'down') {
  const list = sortedCategories.value
  const idx = list.findIndex((c) => c.id === category.id)
  if (idx === -1) return

  const target = direction === 'up' ? idx - 1 : idx + 1
  if (target < 0 || target >= list.length) return

  const orderedIds = list.map((c) => c.id)
  ;[orderedIds[idx], orderedIds[target]] = [orderedIds[target]!, orderedIds[idx]!]

  try {
    await menuStore.reorderCategories(orderedIds)
    ElMessage.success(t('common.success'))
  } catch (e: any) {
    ElMessage.error(e?.message || t('common.fail'))
  }
}
</script>

<template>
  <div class="category-manage">
    <div class="toolbar">
      <div class="title">{{ t('routes.categoryManagement') }}</div>
      <el-button type="primary" :icon="Plus" @click="openCreate">
        {{ t('menu.addCategory') }}
      </el-button>
    </div>

    <el-card shadow="never" class="table-shell" style="--el-card-padding: 0px">
      <el-table :data="sortedCategories" v-loading="menuStore.categoriesLoading" stripe>
        <el-table-column prop="sortOrder" label="#" width="70" />
        <el-table-column label="Icon" width="90">
          <template #default="{ row }">
            <span style="font-size: 18px">{{ row.icon || '-' }}</span>
          </template>
        </el-table-column>
        <el-table-column label="Name">
          <template #default="{ row }">
            <div class="name-cell">
              <div class="zh">{{ row.name.zh }}</div>
              <div class="en">{{ row.name.en }}</div>
            </div>
          </template>
        </el-table-column>
        <el-table-column :label="t('common.status')" width="120">
          <template #default="{ row }">
            <el-tag :type="row.status === 'active' ? 'success' : 'info'">
              {{ row.status }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="Dishes" width="120">
          <template #default="{ row }">
            {{ row.dishCount ?? 0 }}
          </template>
        </el-table-column>
        <el-table-column :label="t('common.action')" width="240" fixed="right">
          <template #default="{ row }">
            <el-button-group>
              <el-button :icon="ArrowUp" @click="move(row, 'up')" :disabled="row.sortOrder === 1" />
              <el-button :icon="ArrowDown" @click="move(row, 'down')" :disabled="row.sortOrder === sortedCategories.length" />
              <el-button type="primary" :icon="Edit" @click="openEdit(row)" />
              <el-button type="danger" :icon="Delete" @click="removeCategory(row)" />
            </el-button-group>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog v-model="dialogVisible" :title="editingId ? t('common.edit') : t('common.add')" width="520px">
      <CategoryForm v-model="formValue" :disabled="dialogSubmitting" />
      <template #footer>
        <el-button @click="dialogVisible = false">{{ t('common.cancel') }}</el-button>
        <el-button type="primary" :loading="dialogSubmitting" @click="submit">{{ t('common.save') }}</el-button>
      </template>
    </el-dialog>
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
  font-weight: 700;
  color: #303133;
}

.name-cell {
  .zh {
    font-weight: 600;
    color: #303133;
  }

  .en {
    font-size: 12px;
    color: #909399;
  }
}

.table-shell {
  overflow: hidden;
}
</style>

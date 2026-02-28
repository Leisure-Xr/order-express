<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Edit, Delete, SwitchFilled } from '@element-plus/icons-vue'
import { useMenuStore } from '@/stores/menu'
import type { Dish } from '@/types'

const { t } = useI18n()
const router = useRouter()
const menuStore = useMenuStore()

const selectedCategoryId = ref<string>('all')
const selectedStatus = ref<'all' | Dish['status']>('all')
const keyword = ref('')

const categoryOptions = computed(() => [
  { id: 'all', label: t('common.all') },
  ...menuStore.categories.map((c) => ({ id: c.id, label: `${c.icon || ''} ${c.name.zh}` })),
])

const categoryNameMap = computed(() => {
  const map = new Map<string, string>()
  for (const c of menuStore.categories) map.set(c.id, `${c.icon || ''} ${c.name.zh}`)
  return map
})

const filteredDishes = computed(() => {
  const kw = keyword.value.trim().toLowerCase()
  return menuStore.dishes.filter((d) => {
    if (selectedCategoryId.value !== 'all' && d.categoryId !== selectedCategoryId.value) return false
    if (selectedStatus.value !== 'all' && d.status !== selectedStatus.value) return false
    if (!kw) return true
    return (
      d.name.zh.toLowerCase().includes(kw) ||
      d.name.en.toLowerCase().includes(kw) ||
      d.description.zh.toLowerCase().includes(kw) ||
      d.description.en.toLowerCase().includes(kw)
    )
  })
})

onMounted(async () => {
  try {
    await menuStore.fetchAll()
  } catch (e: any) {
    ElMessage.error(e?.message || 'Failed to load menu')
  }
})

function goCreate() {
  router.push({ name: 'AdminDishCreate' })
}

function goEdit(dish: Dish) {
  router.push({ name: 'AdminDishEdit', params: { id: dish.id } })
}

async function toggleStatus(dish: Dish) {
  try {
    await menuStore.toggleDishStatus(dish.id)
    ElMessage.success(t('common.success'))
  } catch (e: any) {
    ElMessage.error(e?.message || t('common.fail'))
  }
}

async function removeDish(dish: Dish) {
  try {
    await ElMessageBox.confirm(
      `${t('common.confirm')} ${t('common.delete')}?`,
      t('common.confirm'),
      { type: 'warning' },
    )
  } catch {
    return
  }

  const ok = await menuStore.deleteDish(dish.id)
  if (ok) ElMessage.success(t('common.success'))
  else ElMessage.error(t('common.fail'))
}
</script>

<template>
  <div class="menu-list">
    <div class="toolbar">
      <div class="title">{{ t('routes.menuManagement') }}</div>
      <el-button type="primary" :icon="Plus" @click="goCreate">
        {{ t('menu.addItem') }}
      </el-button>
    </div>

    <el-card shadow="never" class="filters">
      <div class="filter-row">
        <el-select v-model="selectedCategoryId" style="width: 220px">
          <el-option v-for="c in categoryOptions" :key="c.id" :label="c.label" :value="c.id" />
        </el-select>

        <el-select v-model="selectedStatus" style="width: 180px">
          <el-option :label="t('common.all')" value="all" />
          <el-option :label="t('menu.onSale')" value="on_sale" />
          <el-option :label="t('menu.offShelf')" value="off_sale" />
          <el-option :label="t('menu.soldOut')" value="sold_out" />
        </el-select>

        <el-input v-model="keyword" :placeholder="t('common.search')" clearable style="max-width: 360px" />
      </div>
    </el-card>

    <el-card shadow="never" class="table-shell" style="--el-card-padding: 0px">
      <el-table :data="filteredDishes" v-loading="menuStore.dishesLoading" stripe>
        <el-table-column label="Dish" min-width="240">
          <template #default="{ row }">
            <div class="dish-cell">
              <img class="dish-thumb" :src="row.image" alt="" />
              <div class="dish-info">
                <div class="dish-name">{{ row.name.zh }} / {{ row.name.en }}</div>
                <div class="dish-desc">{{ row.description.zh }}</div>
              </div>
            </div>
          </template>
        </el-table-column>

        <el-table-column :label="t('menu.category')" width="160">
          <template #default="{ row }">
            {{ categoryNameMap.get(row.categoryId) || row.categoryId }}
          </template>
        </el-table-column>

        <el-table-column :label="t('menu.price')" width="120">
          <template #default="{ row }">¥{{ row.price.toFixed(2) }}</template>
        </el-table-column>

        <el-table-column :label="t('common.status')" width="140">
          <template #default="{ row }">
            <el-tag :type="row.status === 'on_sale' ? 'success' : row.status === 'sold_out' ? 'danger' : 'info'">
              {{ row.status }}
            </el-tag>
          </template>
        </el-table-column>

        <el-table-column :label="t('common.action')" width="260" fixed="right">
          <template #default="{ row }">
            <el-button-group>
              <el-button type="primary" :icon="Edit" @click="goEdit(row)">{{ t('common.edit') }}</el-button>
              <el-button :icon="SwitchFilled" @click="toggleStatus(row)">
                {{ row.status === 'on_sale' ? t('menu.offShelf') : t('menu.onSale') }}
              </el-button>
              <el-button type="danger" :icon="Delete" @click="removeDish(row)">{{ t('common.delete') }}</el-button>
            </el-button-group>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
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

.filters {
  margin-bottom: 16px;
}

.filter-row {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
  align-items: center;
}

.dish-cell {
  display: flex;
  gap: 12px;
  align-items: center;
}

.dish-thumb {
  width: 54px;
  height: 54px;
  object-fit: cover;
  border-radius: 10px;
  background: #f2f2f2;
  flex: none;
}

.dish-info {
  min-width: 0;
}

.dish-name {
  font-weight: 700;
  color: #303133;
  line-height: 1.2;
}

.dish-desc {
  font-size: 12px;
  color: #909399;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  max-width: 360px;
}

.table-shell {
  overflow: hidden;
}
</style>

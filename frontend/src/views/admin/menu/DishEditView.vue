<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { useMenuStore } from '@/stores/menu'
import DishForm, { type DishFormValue } from '@/components/admin/DishForm.vue'
import AdminPageHeader from '@/components/admin/AdminPageHeader.vue'
import type { Dish } from '@/types'

const { t } = useI18n()
const route = useRoute()
const router = useRouter()
const menuStore = useMenuStore()

const dishId = computed(() => (typeof route.params.id === 'string' ? route.params.id : ''))
const isEdit = computed(() => !!dishId.value)
const pageTitle = computed(() => (isEdit.value ? t('routes.editDish') : t('routes.addDish')))

const loading = ref(false)
const saving = ref(false)

function emptyDish(categoryId: string): DishFormValue {
  return {
    categoryId,
    name: { zh: '', en: '' },
    description: { zh: '', en: '' },
    price: 0,
    originalPrice: undefined,
    image: '',
    images: [],
    status: 'on_sale',
    options: [],
    tags: [],
    preparationTime: 10,
  }
}

const formModel = ref<DishFormValue>(emptyDish(''))

onMounted(async () => {
  loading.value = true
  try {
    if (!menuStore.categories.length) {
      await menuStore.fetchCategories()
    }

    if (!menuStore.dishes.length) {
      await menuStore.fetchDishes()
    }

    const defaultCategory = menuStore.categories[0]?.id || ''
    formModel.value = emptyDish(defaultCategory)

    if (isEdit.value) {
      const dish = await menuStore.fetchDishById(dishId.value)
      if (!dish) {
        ElMessage.error('Dish not found')
        router.replace({ name: 'AdminMenuList' })
        return
      }
      formModel.value = {
        categoryId: dish.categoryId,
        name: { ...dish.name },
        description: { ...dish.description },
        price: dish.price,
        originalPrice: dish.originalPrice,
        image: dish.image,
        images: dish.images ? [...dish.images] : [],
        status: dish.status,
        options: JSON.parse(JSON.stringify(dish.options)),
        tags: [...dish.tags],
        preparationTime: dish.preparationTime,
      }
    }
  } catch (e: any) {
    ElMessage.error(e?.message || t('common.fail'))
  } finally {
    loading.value = false
  }
})

function validate(model: DishFormValue): string | null {
  if (!model.name.zh.trim() || !model.name.en.trim()) return 'Please input ZH/EN name'
  if (!model.categoryId) return 'Please select category'
  if (!model.image.trim()) return 'Please input image URL'
  if (model.price < 0) return 'Invalid price'
  return null
}

async function save() {
  const err = validate(formModel.value)
  if (err) {
    ElMessage.warning(err)
    return
  }

  saving.value = true
  try {
    if (!isEdit.value) {
      await menuStore.createDish(formModel.value)
      ElMessage.success(t('common.success'))
      router.replace({ name: 'AdminMenuList' })
      return
    }

    const updated = await menuStore.updateDish(dishId.value, formModel.value as Partial<Omit<Dish, 'id' | 'createdAt'>>)
    if (!updated) {
      ElMessage.error(t('common.fail'))
      return
    }
    ElMessage.success(t('common.success'))
    router.replace({ name: 'AdminMenuList' })
  } catch (e: any) {
    ElMessage.error(e?.message || t('common.fail'))
  } finally {
    saving.value = false
  }
}
</script>

<template>
  <div class="dish-edit admin-page">
    <AdminPageHeader :title="pageTitle">
      <template #actions>
        <el-button @click="router.back()">{{ t('common.back') }}</el-button>
        <el-button type="primary" :loading="saving" @click="save">
          {{ t('common.save') }}
        </el-button>
      </template>
    </AdminPageHeader>

    <el-skeleton v-if="loading" :rows="8" animated />
    <DishForm
      v-else
      v-model="formModel"
      :categories="menuStore.categories"
      :disabled="saving"
    />
  </div>
</template>

<style scoped lang="scss">
.dish-edit {
  width: 100%;
}
</style>

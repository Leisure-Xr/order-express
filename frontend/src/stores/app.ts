import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useAppStore = defineStore('app', () => {
  const locale = ref<'zh-CN' | 'en'>(
    (localStorage.getItem('locale') as 'zh-CN' | 'en') || 'zh-CN',
  )
  const isLoading = ref(false)
  const sidebarCollapsed = ref(localStorage.getItem('sidebarCollapsed') === '1')
  const currentTableId = ref<string | null>(
    localStorage.getItem('tableId') || null,
  )

  function setLocale(newLocale: 'zh-CN' | 'en') {
    locale.value = newLocale
    localStorage.setItem('locale', newLocale)
  }

  function setTableFromQuery(tableId: string | null) {
    currentTableId.value = tableId
    if (tableId) {
      localStorage.setItem('tableId', tableId)
    } else {
      localStorage.removeItem('tableId')
    }
  }

  function setSidebarCollapsed(collapsed: boolean) {
    sidebarCollapsed.value = collapsed
    localStorage.setItem('sidebarCollapsed', collapsed ? '1' : '0')
  }

  function toggleSidebar() {
    setSidebarCollapsed(!sidebarCollapsed.value)
  }

  return {
    locale,
    isLoading,
    sidebarCollapsed,
    currentTableId,
    setLocale,
    setTableFromQuery,
    setSidebarCollapsed,
    toggleSidebar,
  }
})

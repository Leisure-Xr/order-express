import { ref, computed, onMounted, onUnmounted } from 'vue'

export function useResponsive() {
  const windowWidth = ref(typeof window !== 'undefined' ? window.innerWidth : 1024)

  let mobileQuery: MediaQueryList | null = null
  let tabletQuery: MediaQueryList | null = null
  let desktopQuery: MediaQueryList | null = null

  const isMobile = computed(() => windowWidth.value < 768)
  const isTablet = computed(() => windowWidth.value >= 768 && windowWidth.value < 1024)
  const isDesktop = computed(() => windowWidth.value >= 1024)

  function onMobileChange(e: MediaQueryListEvent | MediaQueryList) {
    if ('matches' in e) {
      updateWidth()
    }
  }

  function onTabletChange(e: MediaQueryListEvent | MediaQueryList) {
    if ('matches' in e) {
      updateWidth()
    }
  }

  function onDesktopChange(e: MediaQueryListEvent | MediaQueryList) {
    if ('matches' in e) {
      updateWidth()
    }
  }

  function updateWidth() {
    windowWidth.value = window.innerWidth
  }

  onMounted(() => {
    mobileQuery = window.matchMedia('(max-width: 767px)')
    tabletQuery = window.matchMedia('(min-width: 768px) and (max-width: 1023px)')
    desktopQuery = window.matchMedia('(min-width: 1024px)')

    mobileQuery.addEventListener('change', onMobileChange)
    tabletQuery.addEventListener('change', onTabletChange)
    desktopQuery.addEventListener('change', onDesktopChange)

    // Also listen to resize for more granular updates
    window.addEventListener('resize', updateWidth)
  })

  onUnmounted(() => {
    if (mobileQuery) mobileQuery.removeEventListener('change', onMobileChange)
    if (tabletQuery) tabletQuery.removeEventListener('change', onTabletChange)
    if (desktopQuery) desktopQuery.removeEventListener('change', onDesktopChange)
    window.removeEventListener('resize', updateWidth)
  })

  return {
    windowWidth,
    isMobile,
    isTablet,
    isDesktop,
  }
}

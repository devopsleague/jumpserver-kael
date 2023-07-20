import { computed, watch } from 'vue'
import { darkTheme, useOsTheme } from 'naive-ui'
import { useAppStore } from '@/store'

export function useTheme() {
  const appStore = useAppStore()
  const OsTheme = useOsTheme()

  const isDark = computed(() => {
    return appStore.theme === 'dark'
  })

  const theme = computed(() => {
    return isDark.value ? darkTheme : OsTheme
  })

  const themeOverrides = computed(() => {
    if (isDark.value) {
      return {
        common: {},
        Button: {
          borderHover: '1px solid rgba(255, 255, 255, 0.3)'
        }
      }
    }
    return {}
  })

  watch(
    () => isDark.value,
    (dark) => {
      if (dark) {
        document.documentElement.classList.add('dark')
      } else {
        document.documentElement.classList.remove('dark')
      }
    },
    { immediate: true }
  )

  return { theme, themeOverrides }
}

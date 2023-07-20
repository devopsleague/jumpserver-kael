<script setup>
import NaiveProvider from '@/components/NaiveProvider/index.vue'
import { useTheme } from '@/hooks/useTheme'

const { theme, themeOverrides } = useTheme()
</script>

<template>
  <n-config-provider
    :theme="theme"
    :theme-overrides="themeOverrides"
  >
    <NaiveProvider>
      <RouterView />
    </NaiveProvider>
  </n-config-provider>
</template>

<style scoped>
</style>

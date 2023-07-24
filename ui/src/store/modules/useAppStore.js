import { defineStore } from 'pinia'

export const useAppStore = defineStore('app-store', {
  state: () => (JSON.parse(localStorage.getItem('appSetting')) || {
    theme: 'dark',
    sidebarWidth: 240
  }),
  actions: {
    setTheme(theme) {
      this.theme = theme
      this.setLocalSetting()
    },
    setSidebarWidth(width) {
      this.sidebarWidth = width
    },
    setLocalSetting() {
      localStorage.setItem('appSetting', JSON.stringify(this.$state))
    },
    switchSidebar() {
      this.sidebarWidth = this.sidebarWidth === 0 ? 240 : 0
    }
  }
})
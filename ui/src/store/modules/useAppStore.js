import { defineStore } from 'pinia'

export const useAppStore = defineStore('app-store', {
  state: () => (JSON.parse(localStorage.getItem('appSetting')) || {
    theme: 'dark'
  }),
  actions: {
    setTheme(theme) {
      this.theme = theme
      this.setLocalSetting()
    },
    setLocalSetting() {
      localStorage.setItem('appSetting', JSON.stringify(this.$state))
    }
  }
})
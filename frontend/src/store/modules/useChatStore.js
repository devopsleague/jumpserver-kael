import { defineStore } from 'pinia'

export const useChatStore = defineStore('chat-store', {
  state: () => ({
    activeTab: 1,
  }),
  actions: {
    setActiveTab(id) {
      this.activeTab++
    }
  }
})
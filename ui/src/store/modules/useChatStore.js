import { defineStore } from 'pinia'

export const useChatStore = defineStore('chat-store', {
  state: () => ({
    tabNum: 0,
    activeTab: 0,
    sessionsStore: []
  }),
  actions: {
    setTabNum() {
      this.tabNum++
    },
    setActiveNum(id) {
      if (id === this.activeTab) return

      this.activeTab = id
    },
    addSessionsStore(data) {
      this.sessionsStore.unshift(data)
      if (data?.id) {
        this.setActiveNum(data.id)
      }
    },
    removeSessionsStore(id) {
      this.sessionsStore = this.sessionsStore.filter(item => item.id !== id)
    },
    currentActiveTab() {
      return this.sessionsStore[this.activeTab]
    }
  }
})
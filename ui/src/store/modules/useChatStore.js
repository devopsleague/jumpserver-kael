import { defineStore } from 'pinia'

export const useChatStore = defineStore('chat-store', {
  state: () => ({
    loading: false,
    tabNum: 0,
    activeTab: 0,
    sessionsStore: [],
    filterChat: {}
  }),
  actions: {
    setLoading(loading) {
      this.loading = loading
    },

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
      const sessionsStore = this.sessionsStore.filter(item => item.id !== id)
      this.sessionsStore = sessionsStore
      const hasActiveTab = sessionsStore.find(item => item.id === this.activeTab)
      if (!hasActiveTab) {
        this.activeTab = sessionsStore?.[0]?.id
      }
    },
  
    currentActiveTab() {
      return this.sessionsStore[this.activeTab]
    },

    // 过滤当前的聊天
    filterCurrentChat () {
      this.filterChat = this.sessionsStore.filter((chat) => chat.id === this.activeTab)?.[0] || {}
    },

    addChatsById(chat) {
      this.filterChat.chats?.push(chat)
    },

    updateChatConversationContentById(id, content) {
      const chats = this.filterChat.chats || []
      const filterChat = chats.filter((chat) => chat.message.id === id)?.[0] || {}
      filterChat.message.content = content
    }
  }
})
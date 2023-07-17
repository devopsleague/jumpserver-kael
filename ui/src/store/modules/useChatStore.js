import { defineStore } from 'pinia'

export const useChatStore = defineStore('chat-store', {
  state: () => ({
    loading: false,
    tabNum: 0,
    activeTab: 0,
    sessionsStore: []
  }),
  getters: {
    filterChat(state) {
      const currentChat = state.sessionsStore.filter((chat) => chat.id === state.activeTab)?.[0] || {}
      return currentChat
    }
  },
  actions: {
    setLoading(loading) {
      this.loading = loading
    },

    setTabNum() {
      this.tabNum++
    },

    setFilterChatDisabled(disabled) {
      this.filterChat.disabled = disabled
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

    addChatsById(chat) {
      this.filterChat.chats?.push(chat)
    },

    removeLastChat() {
      const lastChat = this.filterChat.chats[this.filterChat.chats.length - 1]
      if (lastChat?.message?.content === 'loading') {
        this.filterChat.chats.pop()
      }
    },

    updateChatConversationContentById(id, content) {
      const chats = this.filterChat.chats || []
      const filterChat = chats.filter((chat) => chat.message.id === id)?.[0] || {}
      filterChat.message.content = content
    },

    updateChatConversationDisabledById(index, disabled) {
      const chats = this.filterChat.chats || []
      const chat = chats[index]
      chat.disabled = disabled
    }
  }
})
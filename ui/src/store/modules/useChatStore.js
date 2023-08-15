import { defineStore } from 'pinia'

export const useChatStore = defineStore('chat-store', {
  state: () => ({
    loading: false,
    tabNum: 0,
    activeTab: 0,
    globalDisabled: false,
    chatsStore: []
  }),
  getters: {
    activeChat(state) {
      const chat = state.chatsStore.filter((chat) => chat.id === state.activeTab)?.[0] || {}
      return chat
    }
  },
  actions: {
    setLoading(loading) {
      this.loading = loading
    },

    setTabNum() {
      this.tabNum++
    },

    setActiveChatDisabled(disabled) {
      this.activeChat.disabled = disabled
    },

    setGlobalDisabled(disabled) {
      this.globalDisabled = disabled
    },

    setActiveNum(id) {
      if (id === this.activeTab) return

      this.activeTab = id
    },

    setActiveChatConversationId(data) {
      this.activeChat.conversation_id = data
    },

    addChatToStore(data) {
      this.chatsStore.unshift(data)
      if (data?.id) {
        this.setActiveNum(data.id)
      }
    },

    removeChatInStore(id) {
      const chatsStore = this.chatsStore.filter(item => item.id !== id)
      this.chatsStore = chatsStore
      const hasActiveTab = chatsStore.find(item => item.id === this.activeTab)
      if (!hasActiveTab) {
        this.activeTab = chatsStore?.[0]?.id
      }
    },

    addConversationToActiveChat(chat) {
      this.activeChat.chats?.push(chat)
    },

    removeLastChat() {
      const length = this.activeChat.chats?.length
      if (length > 0) {
        const lastChat = this.activeChat.chats[length - 1]
        if (lastChat?.message?.content === 'loading') {
          this.activeChat.chats.pop()
        }
      }
    },

    updateChatConversationContentById(id, content) {
      const chats = this.activeChat.chats || []
      const filterChat = chats.filter((chat) => chat.message.id === id)?.[0] || {}
      filterChat.message.content = content
    },

    updateChatConversationDisabledByIndex(index, disabled) {
      const chats = this.activeChat.chats || []
      const chat = chats[index]
      chat.disabled = disabled
    }
  }
})
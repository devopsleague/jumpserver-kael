import { useChatStore } from '@/store'
import { pageScroll } from '@/utils/common'

export function useChat() {
  const chatStore = useChatStore()

  const setLoading = (loading) => {
    chatStore.setLoading(loading)
  }

  const addChatConversationById = (chat) => {
    chatStore.filterCurrentChat()
    chatStore.addChatsById(chat)
    pageScroll('scrollRef')
  }

  const updateChatConversationContentById = (id, content) => {
    chatStore.updateChatConversationContentById(id, content)
    pageScroll('scrollRef')
  }

  const hasChat = (id) => {
    const chats = chatStore.filterChat.chats
    const filterChat = chats.filter((chat) => chat.message.id === id)
    if (filterChat.length > 0) {
      return false
    }
    return true
  }

  return {
    chatStore,
    hasChat,
    setLoading,
    addChatConversationById,
    updateChatConversationContentById
  }
}

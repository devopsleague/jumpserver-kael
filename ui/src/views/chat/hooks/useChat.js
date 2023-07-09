import { useChatStore } from '@/store'


export function useChat() {
  const chatStore = useChatStore()

  const addChatConversationById = (chat) => {
    chatStore.filterChatId()
    chatStore.addChatsById(chat)
  }

  const addChatConversationContentById = (id, content) => {
    chatStore.addChatConversationContentById(id, content)
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
    addChatConversationById,
    addChatConversationContentById
  }
}

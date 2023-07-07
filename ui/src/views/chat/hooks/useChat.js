import { ref, reactive } from 'vue'
import { useChatStore } from '@/store'


export function useChat() {
  const chatStore1 = useChatStore()
  const chatStore = reactive({})
  let activeId = ref(1)

  const addChatConversationById = (chat) => {
    chatStore1.filterChatId()
    chatStore1.addChatsById(chat)
  }

  const addChatConversationContentById = (id, content) => {
    chatStore1.addChatConversationContentById(id, content)
  }

  const hasChat = (id) => {
    const chats = chatStore1.filterChat.chats
    const filterChat = chats.filter((chat) => chat.message.id === id)
    if (filterChat.length > 0) {
      return false
    }
    return true
  }

  return {
    chatStore,
    activeId,
    hasChat,
    addChatConversationById,
    addChatConversationContentById
  }
}

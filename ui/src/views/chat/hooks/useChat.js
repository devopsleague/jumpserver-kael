import { ref, reactive } from 'vue'

export function useChat() {
  const chatStore = reactive({})
  let activeId = ref(1)

  const addChatConversationById = (chat) => {
    if (!chatStore[activeId.value]) {
      chatStore[activeId.value] = [chat]
    } else {
      chatStore[activeId.value].push(chat)
    }
  }

  const addChatConversationContentById = (id, content) => {
    const filterChat = filterChatId(id)?.[0]
    filterChat.message.content = content
  }

  const filterChatId = (id) => {
    const filterChat = chatStore[activeId.value].filter((chat) => chat.message.id === id)
    return filterChat
  }

  return {
    chatStore,
    activeId,
    filterChatId,
    addChatConversationById,
    addChatConversationContentById
  }
}

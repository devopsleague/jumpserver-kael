import { reactive } from 'vue'

export function useChat() {
  const chatStore = reactive(JSON.parse(localStorage.getItem('chatStorage')) || {})


  const addChat = (uuid, chat) => {
    // chatStore.addChatByUuid(uuid, chat)
  }

  const addChatConversationById = (id, chat) => {
    chatStore[id] = chat
  }

  const addChatConversationContentById = (id, content) => {
    if (chatStore[id]) {
      chatStore[id].message.content = content
    }
  }

  const updateChatStorage = (data) => {
    localStorage.setItem('chatStorage', JSON.stringify(data))
  }

  return {
    chatStore,
    addChat,
    addChatConversationById,
    addChatConversationContentById,
    updateChatStorage
  }
}

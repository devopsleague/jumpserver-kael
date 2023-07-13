import { useChatStore } from '@/store'
import { pageScroll } from '@/utils/common'

export function useChat() {
  const chatStore = useChatStore()

  const setLoading = (loading) => {
    chatStore.setLoading(loading)
  }

  const onNewChat = (name) => {
    chatStore.setTabNum()
    console.log('chatStore: ', chatStore);
    const data = {
      name: name || `new chat ${chatStore.tabNum}`,
      id: chatStore.tabNum,
      chats: []
    }
    chatStore.addSessionsStore(data)
  } 

  const addChatConversationById = (chat) => {
    chatStore.filterCurrentChat()
    chatStore.addChatsById(chat)
    pageScroll('scrollRef')
  }

  const onNewChatOrAddChatConversationById = (chat) => {
    debugger
    onNewChat(chat.message.content)
    addChatConversationById(chat)
    console.log(chatStore.sessionsStore)
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
    onNewChat,
    onNewChatOrAddChatConversationById,
    hasChat,
    setLoading,
    addChatConversationById,
    updateChatConversationContentById
  }
}

import { useChatStore } from '@/store'
import { pageScroll } from '@/utils/common'

export function useChat() {
  const chatStore = useChatStore()

  const setLoading = (loading) => {
    chatStore.setLoading(loading)
  }

  const getInputFocus = () => {
    const dom = document.getElementsByClassName('n-input__textarea-el')[0]
    dom?.focus()
  }

  const onNewChat = (name) => {
    chatStore.setTabNum()
    const data = {
      name: name || `new chat ${chatStore.tabNum}`,
      id: chatStore.tabNum,
      conversation_id: '',
      chats: []
    }
    chatStore.addSessionsStore(data)
  } 

  const addChatConversationById = (chat) => {
    chatStore.addChatsById(chat)
    chatStore.setActiveChatConversationId(chat?.conversation_id)
    pageScroll('scrollRef')
  }

  const addSystemMessageToCurrentChat = (data) => {
    const currentChat = chatStore.filterChat
    if (currentChat.conversation_id === data.conversation_id) {
      chatStore.addChatsById(data)
    } else {
      const sessionsStore = chatStore.sessionsStore.filter(item => item.conversation_id === data.conversation_id)
      if (sessionsStore.length > 0) {
        sessionsStore[0].chats?.push(data)
      }
    }
    pageScroll('scrollRef')
  }

  const addTemporaryLoadingChat = () => {
    const temporaryChat = {
      message: {
        content: 'loading',
        role: 'assistant',
        create_time: new Date()
      }
    }
    addChatConversationById(temporaryChat)
  }

  const onNewChatOrAddChatConversationById = (chat) => {
    onNewChat(chat.message.content)
    addChatConversationById(chat)
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
    getInputFocus,
    setLoading,
    addChatConversationById,
    addTemporaryLoadingChat,
    addSystemMessageToCurrentChat,
    updateChatConversationContentById
  }
}

import { defineStore } from 'pinia'

const getLocalState = () => {
  const localState = localStorage.getItem('chatStorage')
  if (localState) {
    return JSON.parse(localState)
  }
  return {}
}

export const useChatStore = defineStore('chat-store', {
  state: ()=> getLocalState(),
  action: {
    addChatById(uuid, chat) {}
  }
})
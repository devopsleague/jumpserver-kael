<script setup>
import { ref, onMounted, computed, inject, onUnmounted, reactive } from 'vue'
import Message from '../Message/index.vue'
import { useChat } from '../../hooks/useChat.js'
import { createWebSocket, onSend, closeWs } from '@/utils/socket'
import { useChatStore } from '@/store'
import { pageScroll } from '@/utils/common'

const { hasChat, setLoading, addChatConversationById, updateChatConversationContentById } = useChat()
const chatStore = useChatStore()
const value = ref('')
const $axios = inject("$axios")
const currentConversationId = ref('')

const loading = computed(() => {
  return chatStore.loading
})
const currentSessionStore = computed(() => {
  return chatStore.filterChat
})

const onWebSocketMessage = (data) => {
  setLoading(true)
  if (data.type === 'message') {
    currentConversationId.value = data.conversation_id
    if (hasChat(data.message.id)) {
      addChatConversationById(data)
    } else {
      updateChatConversationContentById(data.message.id, data.message.content)

    }
  } else if (data.type === 'finish') {
    setLoading(false)
  }
}

const onSendHandle = () => {
  const chat = {
    message: {
      content: value.value,
      role: 'user',
      create_time: new Date()
    }
  }
  addChatConversationById(chat)
  const message = {
    content: value.value,
    conversation_id: currentConversationId.value || null
  }
  onSend(message)
  value.value = ''
}

const initWebSocket = () => {
  const protocol = window.location.protocol === 'https:' ? 'wss' : 'ws'
  const path = `${protocol}://127.0.0.1:8083/kael/chat`
  createWebSocket(path, onWebSocketMessage)
}

const onStopHandle = () => {
  $axios.post(
    '/kael/interrupt_current_ask',
    { id: currentConversationId.value || '' }
  ).then(res => {
    console.log('res:----------------', res)
  })
  setLoading(false)
}

const onKeyUpEnter = () => {
  onSendHandle()
}

onMounted(() => {
  initWebSocket()
  pageScroll('scrollRef')
})

onUnmounted(() => {
  closeWs()
})
</script>

<template>
  <div class="content">
    <main class="flex-1 overflow-y-auto">
      <div id="scrollRef" class="overflow-hidden p-4">
        <div>
          <div class="overflow-y-auto">
            <Message
              v-for="(item, index) of currentSessionStore.chats"
              :key="index"
              :item="item"
              :message="item.message"
              @delete="handleDelete(index)"
            />
            <div v-if="loading" class="sticky bottom-0 left-0 flex justify-center">
              <n-button type="warning" @click="onStopHandle()">
                <i class="fa fa-stop-circle-o"></i> 停止
              </n-button>
            </div>
          </div>
        </div>
      </div>
    </main>
    <footer class="footer p-4">
      <div class="flex">
        <n-input
          v-model:value="value"
          type="text"
          placeholder="来说点什么吧..."
          :disabled="loading"
          @keyup.enter="onKeyUpEnter"
        />
        <n-button
          type="primary"
          class="ml-10px"
          :disabled="loading"
          @click="onSendHandle"
        >
          <i class="fa fa-send"></i>
        </n-button>
      </div>
    </footer>
  </div>
</template>

<style lang="scss" scoped>
.content {
  display: flex;
  flex-direction: column;
  width: 100%;
  height: 100%;
}
</style>
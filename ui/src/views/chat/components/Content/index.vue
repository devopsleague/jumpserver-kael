<script setup>
import { ref, onMounted, computed, onUnmounted } from 'vue'
import Message from '../Message/index.vue'
import { useChat } from '../../hooks/useChat.js'
import { pageScroll } from '@/utils/common'
import { createWebSocket, onSend, onClose, closeWs } from '@/utils/socket'
import { useChatStore } from '@/store'

const { hasChat, addChatConversationById, addChatConversationContentById } = useChat()
const chatStore = useChatStore()
const value = ref('')
const loading = ref(false)

const currentSessionStore = computed(() => {
  return chatStore.filterChat
})

const onWebSocketMessage = (data) => {
  console.log('data: -----------------', data)
  if (data.type === 'message') {
    loading.value = true
    
    if (hasChat(data.message.id)) {
      addChatConversationById(data)
    } else {
      addChatConversationContentById(data.message.id, data.message.content)
      pageScroll('scrollRef')
    }
  } else if (data.type === 'finish') {
    loading.value = false
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
  pageScroll('scrollRef')
  const message = {
    content: value.value
  }
  onSend(message)
  value.value = ''
}

const initWebSocket = () => {
  const protocol = window.location.protocol === 'https:' ? 'wss' : 'ws'
  const path = `${protocol}://127.0.0.1:8880/chat`
  createWebSocket(path, onWebSocketMessage)
}

const handleStop = () => {
  loading.value = false
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
              <n-button type="warning" @click="handleStop">
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
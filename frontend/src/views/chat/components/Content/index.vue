<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import Message from '../Message/index.vue'
import { useChat } from '../../hooks/useChat.js'
import dayjs from 'dayjs'
import { pageScroll } from '@/utils/common'

const { chatStore, activeId, filterChatId, addChatConversationById, addChatConversationContentById, updateChatStorage } = useChat()
const value = ref('')
const loading = ref(false)
let webSocket = reactive({})

const currentChatStore = computed(() => {
  return chatStore[activeId.value]
})

const onWebsocketOpen = (msg) => {
  console.log('msg: -----------------open', msg)
}

const onWebSocketMessage = (msg) => {
  const data = JSON.parse(msg.data)
  if (data.type === 'message') {
    loading.value = true
    if (filterChatId(data.message.id).length < 1) {
      addChatConversationById(data)
    } else {
      addChatConversationContentById(data.message.id, data.message.content)
      pageScroll('scrollRef')
    }
  } else if (data.type === 'finish') {
    loading.value = false
    updateChatStorage(chatStore)
  }
}
const onWebSocketError = (msg) => {
  console.log('msg:=> onWebSocketError ', msg)
}
const onWebSocketClose = (msg) => {
  console.log('msg:=> onWebSocketClose ', msg)
}

const onSend = () => {
  const time = dayjs().format('YYYY-MM-DD HH:mm:ss');
  const chat = {
    message: {
      content: value.value,
      role: "user",
      create_time: time
    }
  }
  addChatConversationById(chat)
  pageScroll('scrollRef')
  const message = {
    content: value.value,
    sender: "user",
    new_conversation: true,
    model: 'gpt_3_5',
  }
  webSocket.send(JSON.stringify(message))
  value.value = ''
}

const initWebSocket = () => {
  const path = 'ws://127.0.0.1:8800/chat'
  webSocket = new WebSocket(path)
  webSocket.onopen = onWebsocketOpen
  webSocket.onmessage = onWebSocketMessage
  webSocket.onerror = onWebSocketError
  webSocket.onclose = onWebSocketClose
}

const handleStop = () => {
  loading.value = false
}

onMounted(() => {
  initWebSocket()
})
</script>

<template>
  <div class="content">
    <main class="flex-1 overflow-y-auto">
      <div id="scrollRef" class="overflow-hidden p-4">
        <div>
          <div class="overflow-y-auto">
            <Message
              v-for="(item, index) of currentChatStore"
              :key="index"
              :loading="loading"
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
        <n-input v-model:value="value" type="text" placeholder="来说点什么吧..." />
        <n-button
          type="primary"
          class="ml-10px"
          :disabled="loading"
          @click="onSend"
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
main {
  &::-webkit-scrollbar {
    background-color: transparent;
    width: 8px;
  }
  &::-webkit-scrollbar-thumb {
    border-radius: 8px;
    box-shadow: inset 8px 10px 10px #c6c6c6;
    border: 1px solid rgba(0,0,0,0);
  }
}
</style>
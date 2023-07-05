<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import Message from '../Message/index.vue'
import { useChat } from '../../hooks/useChat.js'

const { chatStore, addChatConversationById, addChatConversationContentById, updateChatStorage } = useChat()
const value = ref('')
const loading = ref(false)

let webSocket = reactive({})

const onWebsocketOpen = (msg) => {
  console.log('msg: -----------------open', msg)
}

const onWebSocketMessage = (msg) => {
  const data = JSON.parse(msg.data)
  if (data.type === 'message') {
    loading.value = true
    if (!chatStore[data.conversation_id]) {
      addChatConversationById(data.conversation_id, data)
    } else {
      addChatConversationContentById(data.conversation_id, data.message.content)
    }
  } else if (data.type === 'finish') {
    loading.value = false
    updateChatStorage(chatStore)
  }
}
const onWebSocketError = (msg) => {
  console.log('msg:===================onWebSocketError ', msg)
}
const onWebSocketClose = (msg) => {
  console.log('msg:===================onWebSocketClose ', msg)
}

const onSend = () => {
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
    <main class="flex-1 overflow-hidden">
      <div id="scrollRef" class="h-1/1 overflow-hidden overflow-y-auto p-4">
        <div class="h-1/1">
          <div v-if="Object.keys(chatStore).length < 1">
          </div>
          <template v-else>
            <div class="overflow-y-auto">
              <Message
                v-for="(item, index) of chatStore"
                :key="index"
                :error="item.error_detail"
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
          </template>
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
#scrollRef::-webkit-scrollbar {
  background-color: rgba(0, 0, 0, 0.25);
  width: 5px;
}
</style>
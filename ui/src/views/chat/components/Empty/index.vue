<script setup>
import { reactive } from 'vue'
import { useChat } from '../../hooks/useChat.js'
import Footer from '../Footer/index.vue'
import { onSend } from '@/utils/socket'

const { onNewChat, addChatConversationById, addTemporaryLoadingChat } = useChat()

const lists = reactive([
  {
    title: 'Examples',
    icon: 'fa fa-life-bouy',
    children: [
      'Explain quantum computing in simple terms',
      'Got any creative ideas for a 10 year old’s birthday?',
      'How do I make an HTTP request in Javascript?'
    ]
  },
  {
    title: 'Capabilities',
    icon: 'fa fa-flash',
    children: [
      'Remembers what user said earlier in the conversation',
      'Allows user to provide follow-up corrections',
      'Trained to decline inappropriate requests'
    ]
  },
  {
    title: 'Limitations',
    icon: 'fa fa-handshake-o',
    children: [
      'May occasionally generate incorrect information',
      'May occasionally produce harmful instructions or biased content',
      'Limited knowledge of world and events after 2021'
    ]
  }
])

const onSendHandle = (value) => {
  const chat = {
    message: {
      content: value,
      role: 'user',
      create_time: new Date()
    }
  }
  onNewChat(value)
  addChatConversationById(chat)
  addTemporaryLoadingChat()
  const message = {
    content: value,
    conversation_id: null
  }
  onSend(message)
}

</script>

<template>
  <div class="empty column-alignment dark:bg-[#343540]">
    <div class="header column-alignment">
      <span class="title">ChatGPT</span>
      <span class="sub-title column-alignment">via JumpServer</span>
    </div>
    <div class="content">
      <div v-for="(item) in lists" class="layout">
        <i :class="item.icon" class="text-center"></i>
        <p class="text-center font-normal text-lg">{{ item.title }}</p>
        <ul class="layout">
          <li v-for="(child) in item.children" class="box">
            {{ child }}
          </li>
        </ul>
      </div>
    </div>
    <div class="footer">
      <Footer @send="onSendHandle" />
    </div>
  </div>
</template>

<style lang="scss" scoped>
.layout {
  display: flex;
  flex-direction: column;
  gap: 14px;
  flex: 1;
}
.column-alignment {
  display: flex;
  flex-direction: column;
  align-items: center;
}
.empty {
  justify-content: space-between;
  height: 100vh;
  .header {
    margin-top: 12vh;
    .title {
      font-weight: 600;
      font-size: 2.25rem;
    }
    .sub-title {
      gap: 12px;
      margin-top: 6px;
      font-weight: 300;
      color: #959598;
      &::before {
        position: relative;
        width: 66%;
        content: '';
        border-top: 1px solid #646466;
      }
    }
  }
  .content {
    display: flex;
    justify-content: center;
    gap: 14px;
    max-width: 690px;
    .box {
      padding: 12px;
      background-color: #40424c;
      border-radius: 6px;
    }
  }
}
</style>
<script setup>
import { inject, ref, toRefs, watch,computed, onMounted } from 'vue'
import { useChatStore } from '@/store'

const chatStore = useChatStore()

const props = defineProps({
  item: Object,
  index: Number
})
const $axios = inject("$axios")
const { item, index } = toRefs(props)
const isDisabled = ref(false)

const currentChat = computed(() => {
  return chatStore.filterChat.chats[index.value] || {}
})

watch(isDisabled, (value) => {
  if (isDisabled.value) {
    chatStore.setFilterChatDisabled(false)
  } else {
    chatStore.setFilterChatDisabled(true)
  }
}, { immediate: true })


const onClick = (value) => {
  $axios.post(
    '/kael/jms_state/',
    { 
      id: item.value.conversation_id,
      activate_review: value
    }
    ).finally(() => {
      isDisabled.value = true
      chatStore.updateChatConversationDisabledById(index.value, true)
    })
  }

onMounted(() => {
  const currentChatValue = { ...currentChat.value }
  if (currentChatValue.hasOwnProperty('disabled')) {
    isDisabled.value = currentChatValue.disabled
  } else {
    isDisabled.value = false
  }
})
</script>

<template>
  <div>
    <n-card v-if="item.meta?.activate_review">
      <span>{{ item.system_message }}</span>
      <template #footer>
        <n-button
          secondary
          size="small"
          :disabled="isDisabled"
          @click="onClick(true)"
        >
          是
        </n-button>
        <n-button
          secondary
          class="ml-6px"
          size="small"
          :disabled="isDisabled"
          @click="onClick(false)"
        >
          否
        </n-button>
      </template>
    </n-card>
    <span v-else>{{ item.system_message }}</span>
  </div>
</template>
<script setup>
import { ref, computed, reactive } from 'vue'
import { useChatStore } from '@/store'

const chatStore = useChatStore()

const value = ref('')
const emit = defineEmits()
const loading = computed(() => {
  return chatStore.loading
})
const disabled = computed(() => {
  return chatStore.filterChat.disabled || false
})

const onSendHandle = () => {
  if (!value.value) return

  emit('send', value.value)
  value.value = ''
}

const onStopHandle = () => {
  emit('stop')
}

const onKeyUpEnter = () => {
  onSendHandle()
}
</script>

<template>
  <footer class="footer dark:bg-[#343540]">
    <div v-if="loading" class="sticky bottom-0 left-0 flex justify-center">
      <n-button type="warning" @click="onStopHandle()">
        <i class="fa fa-stop-circle-o mr-4px"></i> 停止
      </n-button>
    </div>
    <div class="flex w-800px mx-auto  p-4">
      <n-input
        v-model:value="value"
        type="text"
        placeholder="来说点什么吧..."
        class="dark:bg-[#40414f] hover:border-transparent"
        style="--n-border-hover: 1px solid transparent; --n-color-focus: transparent; --n-border-focus: 1px solid transparent; --n-box-shadow-focus: 0 0 8px 0 rgba(193, 194, 198, 0.3);"
        :disabled="loading || disabled"
        @keyup.enter="onKeyUpEnter"
      >
        <template #suffix>
          <n-button
            quaternary
            class="ml-10px"
            :disabled="loading || disabled"
            @click="onSendHandle"
          >
            <i class="fa fa-send"></i>
          </n-button>
        </template>
      </n-input>
    </div>
  </footer>
</template>

<style lang="scss" scoped>
.footer {
  .n-input {
    height: 58px;
    line-height: 58px;
    border-radius: 12px;
  }
}
</style>
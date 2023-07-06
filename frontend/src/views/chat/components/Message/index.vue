<script setup>
import { ref, toRefs, computed } from 'vue'
import Text from './Text.vue'
import { useMessage, useDialog } from 'naive-ui'
import { copy } from '@/utils/common'
import defaultAvatar from '@/assets/avatar.jpg'
import robot from '@/assets/pwa-192x192.png'
import dayjs from 'dayjs'

const props = defineProps({
  id: Number,
  loading: Boolean,
  message: Object,
  error: String
})

const { message = {} } = toRefs(props)
const NMessage = useMessage()
const NDialog = useDialog()
console.log('message: ', message.value)
const asRawText = ref(message.role !== 'assistant')

const options = computed(() => {
  const common = [
    {
      label: '复制',
      key: 'copyText',
    },
    {
      label: '删除',
      key: 'delete',
    },
  ]

  return common
})

const handleSelect = (value) => {
  console.log('value: ', value)
  switch(value) {
    case 'copyText':
      copy(message.value.content)
      NMessage.success('复制成功')
      break
    case 'delete':
      NDialog.warning({
        title: '删除',
        content: '是否删除此消息？',
        positiveText: '是',
        negativeText: '否',
        onPositiveClick: () => {
          NMessage.success('确定')
        }
      })
      break
  }
}

</script>
<template>
  <div ref="messageRef" class="mb-30px">
    <div class="flex" :class="[message?.role === 'assistant' ? 'flex-row': 'flex-row-reverse']">
      <div class="avatar mr-6px ml-6px">
        <n-avatar :src="message?.role === 'assistant' ? robot : defaultAvatar" />
      </div>
      <div class="overflow-hidden flex-1 text-sm flex flex-col" :class="[message?.role === 'assistant' ? 'items-start': 'items-end']">
        <p style="color: #b6bdc6">
          {{ dayjs(message?.create_time).format('YYYY-MM-DD HH:mm:ss') }}
        </p>
        <div class="message flex">
          <Text :message="message" :as-raw-text="asRawText" :error="error" />
          <n-dropdown trigger="hover" :options="options" @select="handleSelect">
            <div style="display: inline-block; color: #b6bdc6" class="hover:cursor-pointer">
              <i class="fa fa-ellipsis-v caret-"></i>
            </div>
          </n-dropdown>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped lang="scss">
.avatar {
  width: 34px;
  height: 34px;
  .n-avatar {
    width: 100%;
    height: 100%;
    border-radius: 50% !important;
  }
}
</style>
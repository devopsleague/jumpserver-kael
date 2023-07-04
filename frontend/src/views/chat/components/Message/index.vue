<script setup>
import { ref, toRefs, computed } from 'vue'
import defaultAvatar from '@/assets/avatar.jpg'
import Text from './Text.vue'
import { useMessage, useDialog } from 'naive-ui'
import { copy } from '@/utils/common'

const props = defineProps({
  id: Number,
  loading: Boolean,
  message: Object,
  error: String
})

const { message, loading } = toRefs(props)
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
  <div ref="messageRef" class="mb-30px inline">
    <div class="flex" :class="[message.role === 'assistant' ? 'flex-row': 'flex-row-reverse']">
      <div class="avatar mr-6px ml-6px">
        <n-avatar :src="defaultAvatar" />
      </div>
      <div class="overflow-hidden text-sm items-start">
        <p>{{ message.create_time }}</p>
        <div class="message">
          <Text :message="message" :as-raw-text="asRawText" :error="error" />
          <div style="display: inline-block;">
            <n-dropdown trigger="hover" :options="options" @select="handleSelect">
              <i class="iconfont fa-align-justify"></i>
            </n-dropdown>
          </div>
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
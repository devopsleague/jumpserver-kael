<script setup>
import { ref, toRefs, computed } from 'vue'
import Text from './Text.vue'
import { useMessage, useDialog } from 'naive-ui'
import { copy } from '@/utils/common'
import defaultAvatar from '@/assets/avatar.jpg'
import robot from '@/assets/pwa-192x192.png'
import dayjs from 'dayjs'

const props = defineProps({
  item: Object,
})

const { item = {} } = toRefs(props)
const NMessage = useMessage()
const NDialog = useDialog()
const asRawText = ref(props.item.message.role === 'assistant')

const options = computed(() => {
  const common = [
    {
      label: '复制',
      key: 'copyText',
      props: {
        onClick: () => {
          console.log('item: ', item)
          copy(props.item.message.content)
          NMessage.success('复制成功')
        }
      }
    },
    {
      label: '删除',
      key: 'delete',
      props: {
        onClick: () => {
          NDialog.warning({
            title: '删除',
            content: '是否删除此消息？',
            positiveText: '是',
            negativeText: '否',
            onPositiveClick: () => {
              NMessage.success('确定')
            }
          })
        }
      }
    },
  ]

  return common
})

</script>
<template>
  <div ref="messageRef" class="mb-20px">
    <div class="flex" :class="[asRawText ? 'flex-row': 'flex-row-reverse']">
      <div class="avatar mr-6px ml-6px">
        <n-avatar :src="asRawText ? robot : defaultAvatar" />
      </div>
      <div class="overflow-hidden flex-1 text-sm flex flex-col" :class="[asRawText ? 'items-start': 'items-end']">
        <p style="color: #b6bdc6">
          {{ dayjs(item.message?.create_time).format('YYYY-MM-DD HH:mm:ss') }}
        </p>
        <div class="message flex">
          <Text :message="item.message" :as-raw-text="asRawText" :error="item?.error" />
          <n-dropdown trigger="hover" :options="options">
            <div style="display: inline-block; color: #b6bdc6; align-self: end;" class="hover:cursor-pointer">
              <i class="fa fa-ellipsis-v"></i>
            </div>
          </n-dropdown>
        </div>
      </div>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.avatar {
  width: 34px;
  height: 34px;
  .n-avatar {
    width: 100%;
    height: 100%;
    border-radius: 50% !important;
  }
}
.message {
  & > div {
    display: inline-block;
    padding: 6px 10px;
    border-radius: 6px;
  }
}
</style>
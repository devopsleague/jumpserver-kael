<script setup>
import { ref, toRefs, computed } from 'vue'
import Text from './Text.vue'
import System from './System.vue'
import { useMessage } from 'naive-ui'
import { copy } from '@/utils/common'
import robot from '@/assets/pwa-192x192.png'
import dayjs from 'dayjs'

const props = defineProps({
  item: Object,
  index: Number,
})

const { item = {} } = toRefs(props)
const NMessage = useMessage()
const asRawText = ref(props.item.message.role === 'assistant')
const userAvatar = computed(() => {
  return '/api/v1/settings/logo/'
})

const options = computed(() => {
  const common = [
    {
      label: '复制',
      key: 'copyText',
      icons: 'copy',
      props: {
        onClick: () => {
          console.log('item: ', item)
          copy(props.item.message.content)
          NMessage.success('复制成功')
        }
      }
    }
  ]

  return common
})

</script>
<template>
  <div
    ref="messageRef" :class="{'dark:bg-[#444654]': asRawText}">
    <div class="flex w-full max-w-800px mx-auto p-4">
      <div class="avatar mr-6px ml-6px">
        <n-avatar round :src="asRawText ? robot : userAvatar" />
      </div>
      <div class="overflow-hidden flex-1 text-sm flex flex-col">
        <p style="color: #b6bdc6" class="flex justify-between">
          <span>
            {{ dayjs(item.message?.create_time).format('YYYY-MM-DD HH:mm:ss') }}
          </span>
          <div class="inline-block">
            <span v-if="options.length < 3">
              <span v-for="(item) in options" class="cursor-pointer hover:text-light-100">
                <i v-if="item.icons.startsWith('fa')" :class="item.icons" class="ml-4px" @click="item.props.onClick"></i>
                <SvgIcon v-else :name="item.icons" class="ml-4px" @click="item.props.onClick" />
              </span>
            </span>
            <n-dropdown v-else trigger="hover" :options="options">
              <div style="display: inline-block; color: #b6bdc6; align-self: end;" class="hover:cursor-pointer">
                <i class="fa fa-ellipsis-v"></i>
              </div>
            </n-dropdown>
          </div>
        </p>
        <div class="message flex">
          <template v-if="item.type && item.type === 'waiting' && item.meta?.activate_review">
            <System :item="item" :index="index" />
          </template>
          <template v-else>
            <Text :message="item.message" :as-raw-text="asRawText" :error="item?.error" />
          </template>
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
    background-color: #444654;
  }
}
.message {
  & > div {
    display: inline-block;
    padding: 6px 0;
    background-color: transparent;
  }
}
</style>
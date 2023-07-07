<script setup>
import { reactive, computed, onMounted } from 'vue'
import { useChatStore } from '@/store'

const chatStore = useChatStore()
const tabNum = computed(() => chatStore.tabNum)
const activeTab = computed(() => chatStore.activeTab)
const sessions = computed(() => chatStore.sessionsStore)

const onNewChat = () => {
  chatStore.setTabNum()
  const data = {
    id: tabNum.value,
    name: 'new chat ' + tabNum.value,
    chats: []
  }
  chatStore.addSessionsStore(data)
  console.log('sessions: ===============', sessions.value)
}

const switchTab = (id) => {
  chatStore.setActiveNum(id)
  chatStore.filterChatId()
}

const onDelete = (id) => {
  chatStore.removeSessionsStore(id)
}

onMounted(() => {
  if (sessions.value.length < 1) {
    onNewChat()
  }
})

</script>
<template>
  <div class="box-border">
    <n-button type="primary" dashed class="mb-16px w-1/1" @click="onNewChat">
      新建聊天
    </n-button>
    <div 
      v-for="(item, index) in sessions"
      :key="index"
      class="card"
      :class="[activeTab === item.id ? 'active-tab' : '']"
      @click="switchTab(item.id)"
    >
      <span class="title">
        <i class="fa fa-commenting-o mr-8px"></i>
        <span style="user-select: none;">{{ item.name }}</span>
      </span>
      <span v-if="activeTab === item.id" class="action">
        <i class="fa fa-trash-o cursor-pointer" @click.stop="onDelete(item.id)"></i>
      </span>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.card {
  display: flex;
  width: 100%;
  height: 46px;
  padding: 12px;
  margin-bottom: 8px;
  border: 1px solid #e5e7eb;
  border-radius: 6px;
  &:hover {
    border-color: #d2d6dc;
    background-color: rgb(245 245 245 / 1);
  }
  .title {
    flex: 1;
    overflow: hidden;
    white-space: nowrap;
    text-overflow: ellipsis;
  }
}
.active-tab {
  border-color: #36ad6a;
  background-color: rgb(245 245 245 / 1);
  &:hover {
    border-color: #36ad6a;
  }
}
</style>
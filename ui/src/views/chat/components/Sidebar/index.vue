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
}

const switchTab = (id) => {
  chatStore.setActiveNum(id)
  chatStore.filterCurrentChat()
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
  <n-layout-sider
    collapse-mode="width"
    :collapsed-width="0"
    :width="240"
    show-trigger="arrow-circle"
    content-style="padding: 16px;"
    bordered
  >
    <div class="box-border">
      <n-button type="primary" dashed class="mb-16px w-1/1" @click="onNewChat">
        新建聊天
      </n-button>
      <div 
        v-for="(item, index) in sessions"
        :key="index"
        class="card border hover:bg-neutral-100 dark:hover:bg-[#24272e] border-[#e5e7eb] dark:border-neutral-800"
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
  </n-layout-sider>
</template>

<style lang="scss" scoped>
.card {
  display: flex;
  width: 100%;
  height: 46px;
  padding: 12px;
  margin-bottom: 8px;
  border-radius: 6px;
  cursor: pointer;
  .title {
    flex: 1;
    overflow: hidden;
    white-space: nowrap;
    text-overflow: ellipsis;
  }
}
.active-tab {
  border-color: #36ad6a;
  background-color: rgba(36, 39, 46, 1);
}
</style>
<script setup>
import { ref, reactive, computed } from 'vue'

const active = ref(false)
const sessions = reactive([
  {
    id: 1,
    message: {
      content: '你好',
      role: 'user',
      create_time: '2021-08-01 12:00:00'
    }
  },
  {
    id: 2,
    message: {
      content: '阿哈阿',
      role: 'user',
      create_time: '2021-08-01 13:00:00'
    }
  }
])

</script>
<template>
  <div class="box-border">
    <n-button type="primary" dashed class="mb-16px w-1/1">
      新建聊天
    </n-button>
    <div  v-for="(item) in sessions" class="card">
      <span class="title">
        <i class="fa fa-commenting-o mr-8px"></i>
        <span>{{ item.message.content }}</span>
      </span>
      <span class="action">
        <i class="fa fa-trash-o"></i>
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
</style>
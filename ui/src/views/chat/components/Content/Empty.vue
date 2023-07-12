<script setup>
import { reactive } from 'vue';
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
</script>

<template>
  <div class="empty">
    <h1 class="title">chatGPT</h1>
    <div class="content">
      <div v-for="(item, index) in lists" class="layout">
        <i :class="item.icon" class="text-center"></i>
        <p class="text-center font-normal text-lg">{{ item.title }}</p>
        <ul class="layout">
          <li v-for="(child, i) in item.children" class="box">
            {{ child }}
          </li>
        </ul>
      </div>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.layout {
  display: flex;
  flex-direction: column;
  gap: 14px;
}
.empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  .title {
    margin-top: 16vh;
    margin-bottom: 34px;
    font-weight: 600;
    font-size: 2.25rem;
  }
  .content {
    display: flex;
    justify-content: center;
    gap: 14px;
    max-width: 730px;
    .box {
      padding: 12px;
      background-color: #40424c;
      border-radius: 6px;
    }
  }
}
</style>
<script setup>
import { onMounted } from 'vue'
import Sidebar from './components/Sidebar/index.vue'
import Content from './components/Content/index.vue'
import { LunaEvent } from '@/utils/luna'

const lunaEvent = new LunaEvent()

onMounted(() => {
  lunaEvent.init()
})

</script>

<template>
  <div class="root">
    <n-layout
      has-sider
      class="root-layout"
    >
      <Sidebar />
      <n-layout-content class="dark:bg-[#343540]">
        <Content />
      </n-layout-content>
    </n-layout>
  </div>
</template>

<style>
::v-deep(.n-layout) {
  background-color: transparent;
}
</style>

<style lang="scss" scoped>
.root {
  width: 100%;
  height: 100vh;
  box-sizing: border-box;
}
.root-layout {
  height: 100%;
}
</style>

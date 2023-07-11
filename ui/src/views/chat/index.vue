<script setup>
import Sidebar from './components/Sidebar/index.vue'
import Content from './components/Content/index.vue'
import { useMessage } from 'naive-ui'

window.$message = useMessage()

</script>

<template>
  <div class="root">
    <n-layout
      has-sider
      class="root-layout dark:bg-gray-800"
    >
      <Sidebar />
      <n-layout-content>
        <Content />
      </n-layout-content>
    </n-layout>
  </div>
</template>

<style>
::v-deep .n-layout {
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

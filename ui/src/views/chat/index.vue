<script setup>
import Sider from './components/Sidebar/index.vue'
import Content from './components/Content/index.vue'
import { useMessage } from 'naive-ui'

window.$message = useMessage()

</script>

<template>
  <div class="root">
    <n-layout has-sider class="root-layout">
      <n-layout-sider
        collapse-mode="width"
        :collapsed-width="0"
        :width="240"
        show-trigger="arrow-circle"
        content-style="padding: 16px;"
        bordered
      >
        <Sider />
      </n-layout-sider>
      <n-layout-content>
        <Content />
      </n-layout-content>
    </n-layout>
  </div>
</template>

<style lang="scss" scoped>
.root {
  width: 100%;
  height: 100%;
  padding: 20px;
  box-sizing: border-box;
}
.root-layout {
  height: 100%;
  border: 1px solid #e5e7eb;
  border-radius: 6px;
} 
</style>

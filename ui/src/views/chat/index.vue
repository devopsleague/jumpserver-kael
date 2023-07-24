<script setup>
import { onMounted } from 'vue'
import Sidebar from './components/Sidebar/index.vue'
import Content from './components/Content/index.vue'
import { LunaEvent } from '@/utils/luna'
import { useAppStore } from '@/store'

const appStore = useAppStore()
const lunaEvent = new LunaEvent()

const onSwitchSidebar = () => {
  appStore.switchSidebar()
}

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
        <button
          v-if="!appStore.sidebarWidth"
          secondary
          class="switch border border-solid border-[#545557] h-44px rounded-6px p-13px text-[0px] hover:bg-[#40414f]"
          @click="onSwitchSidebar"
        >
          <SvgIcon name="switch" />
        </button>
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
.switch {
  position: absolute;
  top: 16px;
  left: 16px;
}
</style>

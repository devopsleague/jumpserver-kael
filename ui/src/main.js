import { createApp } from 'vue'
import "virtual:windi.css"
import App from './App.vue'
import { setupRouter } from './router'
import naive from 'naive-ui'
import 'font-awesome/css/font-awesome.css'
import lodash from 'lodash'
import { setupStore } from './store'
import '@/styles/index.scss'
import axios from '@/utils/axios'
import 'virtual:svg-icons-register'
import SvgIcon from '@/components/SvgIcon/index.vue'

window._ = lodash
const app = createApp(App)
app.provide('$axios', axios)
app.component('SvgIcon', SvgIcon)

setupStore(app)
setupRouter(app)

app.use(naive)
app.mount('#app')

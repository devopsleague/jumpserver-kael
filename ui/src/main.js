import { createApp } from 'vue'
import "virtual:windi.css"
import App from './App.vue'
import { setupRouter } from './router'
import naive from 'naive-ui'
import 'font-awesome/css/font-awesome.css'
import lodash from 'lodash'
import { setupStore } from './store'
import '@/styles/index.scss'

window._ = lodash
const app = createApp(App)
setupStore(app)
setupRouter(app)

app.use(naive)
app.mount('#app')

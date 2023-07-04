import { createApp } from 'vue'
import './style.css'
import "virtual:windi.css"
import App from './App.vue'
import { setupRouter } from './router'
import naive from 'naive-ui'
import 'font-awesome/css/font-awesome.min.css'
import lodash from 'lodash'

window._ = lodash
const app = createApp(App)

setupRouter(app)
app.use(naive)
app.mount('#app')

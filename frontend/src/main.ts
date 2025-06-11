import './assets/main.css'

import { createApp } from 'vue'
import App from './App.vue'
import router from './router'

// 导入Element Plus
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'

const app = createApp(App)

// 注册Element Plus
app.use(ElementPlus)

app.use(router)

app.mount('#app')

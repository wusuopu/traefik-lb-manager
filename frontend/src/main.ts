import './assets/main.css'

import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import router from './router'

// 导入Element Plus
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import 'element-plus/theme-chalk/dark/css-vars.css'
import './assets/css-vars.css'

import localizeEn from 'ajv-i18n/localize/en';
import { i18n } from '@lljj/vue3-form-element';
i18n.useLocal(localizeEn);

const app = createApp(App)
app.use(createPinia())

// 注册Element Plus
app.use(ElementPlus)

app.use(router)

app.mount('#app')

// @ticket 00-project：应用入口——挂载 Vue、Pinia、Router
import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import router from './router'
// 执行 import './styles.css'，实现对应的前端逻辑

import './styles.css'

// 调用 createApp(App).use(createPinia()).use(router).mount('#app')，触发对应功能

createApp(App).use(createPinia()).use(router).mount('#app')

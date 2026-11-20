// @ticket B端-前端：应用入口
import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import router from './router'
// 执行 import './styles/global.css'，实现对应的前端逻辑

import './styles/global.css'

// 调用 createApp(App).use(createPinia()).use(router).mount('#app')，触发对应功能

createApp(App).use(createPinia()).use(router).mount('#app')

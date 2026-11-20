// @ticket B端-前端：应用入口
import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import router from './router'
import './styles/global.css'

createApp(App).use(createPinia()).use(router).mount('#app')

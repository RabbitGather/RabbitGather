import { createApp } from 'vue'
import '@/style.css'
import App from '@/App.vue'
import router from '@/router' // 引入router
import store from '@/store' // 引入store
import 'virtual:windi.css'

createApp(App)
    .use(store) // 使用 store
    .use(router) // 使用 router
    .mount('#app')

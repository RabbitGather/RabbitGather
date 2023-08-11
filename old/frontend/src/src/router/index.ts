import { createRouter, createWebHistory, RouteRecordRaw} from 'vue-router'
import Home from '../components/Home.vue'

const routes: Array<RouteRecordRaw> = [
    {
        path: '/',
        name: 'Home',
        component: Home
    }
]


// 創建一個router，
const router = createRouter({
    history: createWebHistory(),
    routes
})


// 將 router 物件導出成一個 module
export default router
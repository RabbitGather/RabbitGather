import { createRouter, createWebHistory, RouteRecordRaw} from 'vue-router'
import HelloWorld from '@/components/HelloWorld.vue'

// 創建一個 routes 陣列，定義所有 route 路徑
const routes: Array<RouteRecordRaw> = [
    {
        path: '/hello',
        name: 'HelloWorld',
        component: HelloWorld
    }
]


// 創建一個router，
const router = createRouter({
    history: createWebHistory(), // 這邊我使用 # 路徑模式
    routes // 導入使用上方 routes 所定義的路徑
})


// 將 router 物件導出成一個 module
export default router

import { createApp } from 'vue'
import { createPinia } from 'pinia'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'

import App from './App.vue'
import router from './router'
import { useUserStore } from './stores/user'

const app = createApp(App)
const pinia = createPinia()

app.use(pinia) // 先 use pinia
app.use(router)
app.use(ElementPlus)

// --- Token 持久化逻辑 ---
// 尝试从 localStorage 中恢复 token
const userStore = useUserStore()
const storedToken = localStorage.getItem('questflow_token')
const storedUserInfo = localStorage.getItem('questflow_userinfo')

if (storedToken && storedUserInfo) {
  userStore.token = storedToken
  userStore.userInfo = JSON.parse(storedUserInfo)
}

// 监听 Pinia store 的变化
userStore.$subscribe((mutation, state) => {
  // 每当 token 或 userInfo 变化时，就存入 localStorage
  localStorage.setItem('questflow_token', state.token)
  localStorage.setItem('questflow_userinfo', JSON.stringify(state.userInfo))
})

app.mount('#app')

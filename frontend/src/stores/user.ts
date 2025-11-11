import { ref, computed } from 'vue'
import { defineStore } from 'pinia'
import { loginAPI } from '@/api/user'
import type { LoginRequest } from '@/api/user'

// 定义用户信息的接口类型
interface UserInfo {
  id: number
  username: string
  nickname: string
  role: number
}

export const useUserStore = defineStore(
  'user',
  () => {
    // --- State ---
    // 用户认证凭证
    const token = ref('')
    // 用户基本信息
    const userInfo = ref<UserInfo>({
      id: 0,
      username: '',
      nickname: '',
      role: 0
    })

    // --- Getters ---
    // 计算属性，判断用户是否已登录
    const isAuthenticated = computed(() => !!token.value)

    // --- Actions ---
    /**
     * 处理用户登录逻辑
     * @param loginData - 包含用户名和密码的对象
     */
    async function login(loginData: LoginRequest) {
      // 调用登录 API
      const data = await loginAPI(loginData)
      // 保存 token 和用户信息到 state
      token.value = data.token
      userInfo.value = data.user_info
    }

    /**
     * 处理用户登出逻辑
     */
    function logout() {
      // 清除 Pinia state
      token.value = ''
      userInfo.value = { id: 0, username: '', nickname: '', role: 0 }

      // 清除 localStorage (这一步是为了配合 main.ts 中的持久化逻辑)
      localStorage.removeItem('questflow_token')
      localStorage.removeItem('questflow_userinfo')
    }

    return { token, userInfo, isAuthenticated, login, logout }
  }
)

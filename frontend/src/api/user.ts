import request from './request'

// --- 类型定义 ---
export interface LoginRequest {
  username?: string
  password?: string
}

interface UserInfo {
  id: number
  username: string
  nickname: string
  role: number
}

interface LoginResponse {
  token: string
  user_info: UserInfo
}

// 新增：更新用户信息请求类型
export interface UpdateProfileRequest {
  username?: string;
  password?: string;
}

// --- API 函数 ---

/**
 * 用户登录
 * @param data - 包含用户名和密码的对象
 */
export const loginAPI = (data: LoginRequest) => {
  return request<any, LoginResponse>({
    url: '/users/login',
    method: 'POST',
    data
  })
}

export const updateProfileAPI = (data: UpdateProfileRequest) => {
  return request<any, null>({ // 成功时后端不返回 data
    url: '/users/profile',
    method: 'PUT',
    data
  })
}

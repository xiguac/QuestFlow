<template>
  <div class="profile-container">
    <el-card>
      <template #header>
        <span>个人中心</span>
      </template>
      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="100px"
        style="max-width: 500px"
      >
        <el-form-item label="当前用户名">
          <el-input :value="userStore.userInfo.username" disabled />
        </el-form-item>
        <el-form-item label="新用户名" prop="username">
          <el-input v-model="form.username" placeholder="留空则不修改" />
        </el-form-item>
        <el-form-item label="新密码" prop="password">
          <el-input v-model="form.password" type="password" show-password placeholder="留空则不修改" />
        </el-form-item>
        <el-form-item label="确认新密码" prop="confirmPassword">
          <el-input v-model="form.confirmPassword" type="password" show-password placeholder="请再次输入新密码" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSubmit" :loading="loading">
            保存修改
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { updateProfileAPI } from '@/api/user'
import type { FormInstance, FormRules } from 'element-plus'
import { ElMessage, ElMessageBox } from 'element-plus'

const userStore = useUserStore()
const router = useRouter()

const formRef = ref<FormInstance>()
const loading = ref(false)
const form = reactive({
  username: '',
  password: '',
  confirmPassword: ''
})

const validateConfirmPassword = (rule: any, value: any, callback: any) => {
  if (form.password && value !== form.password) {
    callback(new Error('两次输入的密码不一致!'))
  } else {
    callback()
  }
}

const rules = reactive<FormRules>({
  username: [
    { min: 3, max: 20, message: '用户名长度应为 3 到 20 个字符', trigger: 'blur' }
  ],
  password: [
    { min: 6, max: 30, message: '密码长度应为 6 到 30 个字符', trigger: 'blur' }
  ],
  confirmPassword: [
    { validator: validateConfirmPassword, trigger: 'blur' },
    // 只有在输入了新密码时，才要求确认密码
    {
      required: true,
      message: '请再次输入新密码',
      trigger: 'blur',
      validator: (rule, value, callback) => {
        if (form.password && !value) {
          callback(new Error('请再次输入新密码'));
        } else {
          callback();
        }
      },
    },
  ]
})

const handleSubmit = () => {
  formRef.value?.validate(async (valid) => {
    if (valid) {
      if (!form.username && !form.password) {
        ElMessage.warning('请输入要修改的内容')
        return
      }
      loading.value = true
      try {
        await updateProfileAPI({
          username: form.username || undefined,
          password: form.password || undefined
        })
        ElMessageBox.alert(
          '用户信息已更新，请重新登录。',
          '操作成功',
          {
            confirmButtonText: '好的',
            callback: () => {
              userStore.logout()
              router.push('/login')
            },
          }
        )
      } catch (error) {
        console.error('Update profile failed:', error)
      } finally {
        loading.value = false
      }
    }
  })
}
</script>

<style scoped>
.profile-container {
  padding: 10px;
}
</style>

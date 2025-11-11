<template>
  <div class="form-fill-container">
    <el-card v-if="!loading && form" class="form-card">
      <template #header>
        <div class="card-header">
          <h1>{{ form.title }}</h1>
          <p>{{ form.description }}</p>
        </div>
      </template>

      <el-form :model="answers" label-position="top">
        <el-form-item
          v-for="question in form.definition.questions"
          :key="question.id"
          :label="question.title"
        >
          <!-- 动态渲染不同类型的问题 -->
          <div v-if="question.type === 'single_choice'">
            <el-radio-group v-model="answers[question.id]">
              <el-radio
                v-for="option in question.options"
                :key="option.id"
                :label="option.id"
              >
                {{ option.text }}
              </el-radio>
            </el-radio-group>
          </div>

          <div v-if="question.type === 'text_input'">
            <el-input
              v-model="answers[question.id]"
              type="textarea"
              :rows="3"
              placeholder="请输入您的回答"
            />
          </div>

          <!-- 未来可以添加更多问题类型 -->

        </el-form-item>

        <el-form-item>
          <el-button type="primary" @click="handleSubmit" :loading="submitting">
            提 交
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <div v-if="loading" class="loading-container">
      <el-skeleton :rows="5" animated />
    </div>

    <el-result
      v-if="!loading && !form"
      status="error"
      title="表单加载失败"
      sub-title="您访问的表单不存在或已关闭。"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, reactive } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getPublicFormAPI, submitFormAPI, type PublicForm } from '@/api/form'
import { ElMessage } from 'element-plus'

const route = useRoute()
const router = useRouter()

const formKey = route.params.formKey as string

const form = ref<PublicForm | null>(null)
const loading = ref(true)
const submitting = ref(false)

// 使用 reactive 来存储所有问题的答案
// key 是 question.id, value 是答案
const answers = reactive<Record<string, any>>({})

const fetchFormDefinition = async () => {
  try {
    loading.value = true
    const res = await getPublicFormAPI(formKey)
    form.value = res
  } catch (error) {
    console.error('Failed to fetch form definition:', error)
    form.value = null
  } finally {
    loading.value = false
  }
}

const handleSubmit = async () => {
  try {
    submitting.value = true
    await submitFormAPI(formKey, { data: answers })

    // 提交成功后不再是 ElMessage 提示，而是直接跳转
    router.push({ name: 'success' })

  } catch (error) {
    console.error('Submission failed:', error)
    // 错误消息已由 axios 拦截器处理
  } finally {
    submitting.value = false
  }
}

onMounted(() => {
  fetchFormDefinition()
})
</script>

<style scoped>
.form-fill-container {
  display: flex;
  justify-content: center;
  padding: 40px 20px;
  min-height: 100vh;
  background-color: #f0f2f5;
  box-sizing: border-box;
}
.form-card {
  width: 100%;
  max-width: 800px;
}
.card-header h1 {
  margin: 0;
  font-size: 24px;
}
.card-header p {
  margin: 10px 0 0;
  color: #606266;
}
.loading-container {
  width: 100%;
  max-width: 800px;
}
</style>

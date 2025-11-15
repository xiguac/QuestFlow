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
          v-for="(question, index) in form.definition.questions"
          :key="question.id"
          :label="`${index + 1}. ${question.title}`"
        >
          <!-- 单选题 -->
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

          <!-- 【新增】多选题 -->
          <div v-if="question.type === 'multi_choice'">
            <el-checkbox-group v-model="answers[question.id]">
              <el-checkbox
                v-for="option in question.options"
                :key="option.id"
                :label="option.id"
              >
                {{ option.text }}
              </el-checkbox>
            </el-checkbox-group>
          </div>

          <!-- 【新增】判断题 -->
          <div v-if="question.type === 'judgment'">
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

          <!-- 文本题 -->
          <div v-if="question.type === 'text_input'">
            <el-input
              v-model="answers[question.id]"
              type="textarea"
              :rows="3"
              placeholder="请输入您的回答"
            />
          </div>

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
import { ref, onMounted, reactive, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getPublicFormAPI, submitFormAPI, type PublicForm } from '@/api/form'

const route = useRoute()
const router = useRouter()

const formKey = route.params.formKey as string

const form = ref<PublicForm | null>(null)
const loading = ref(true)
const submitting = ref(false)

const answers = reactive<Record<string, any>>({})

// 【新增】监听 form 数据的变化，为多选题初始化答案数组
watch(form, (newForm) => {
  if (newForm) {
    newForm.definition.questions.forEach(question => {
      if (question.type === 'multi_choice') {
        answers[question.id] = [] // 多选题的答案必须是数组
      }
    })
  }
})

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

    router.push({ name: 'success' })

  } catch (error) {
    console.error('Submission failed:', error)
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

/* 样式调整，让选项垂直排列 */
:deep(.el-radio-group),
:deep(.el-checkbox-group) {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
}
:deep(.el-radio),
:deep(.el-checkbox) {
  margin-bottom: 8px;
}

</style>

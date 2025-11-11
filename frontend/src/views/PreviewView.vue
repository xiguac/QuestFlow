<template>
  <div class="form-fill-container">
    <el-card v-if="form" class="form-card">
      <template #header>
        <div class="card-header">
          <h1>{{ form.title }}</h1>
          <p>{{ form.description }}</p>
        </div>
      </template>

      <el-form :model="answers" label-position="top">
        <el-form-item
          v-for="question in form.questions"
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
        </el-form-item>

        <el-form-item>
          <el-button type="primary" disabled>
            提交 (预览模式)
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <el-result
      v-else
      status="warning"
      title="预览加载失败"
      sub-title="没有找到可供预览的表单数据，请返回编辑器重试。"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, reactive } from 'vue'
import type { FormDefinition } from '@/stores/editor'

const form = ref<FormDefinition | null>(null)
const answers = reactive<Record<string, any>>({})

onMounted(() => {
  const previewData = sessionStorage.getItem('questflow_form_preview')
  if (previewData) {
    try {
      form.value = JSON.parse(previewData)
    } catch (e) {
      console.error('Failed to parse preview data:', e)
      form.value = null
    }
  }
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
</style>

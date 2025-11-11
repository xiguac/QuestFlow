<template>
  <div class="stats-container">
    <el-card v-if="!loading && stats" v-loading="loading">
      <template #header>
        <div class="card-header">
          <span>表单统计</span>
          <el-tag type="success">总提交数: {{ stats.total_submissions }}</el-tag>
        </div>
      </template>

      <div v-for="(question, index) in stats.question_stats" :key="question.question_id" class="question-stat">
        <h4>{{ index + 1 }}. {{ question.title }}</h4>

        <!-- 单选题统计 -->
        <div v-if="question.question_type === 'single_choice' && question.option_stats">
          <div v-for="option in question.option_stats" :key="option.text" class="option-stat">
            <span class="option-text">{{ option.text }}</span>
            <el-progress :percentage="getPercentage(option.count, stats.total_submissions)" class="option-progress" />
            <span class="option-count">{{ option.count }} 票</span>
          </div>
        </div>

        <!-- 文本题统计 -->
        <div v-if="question.question_type === 'text_input' && question.text_answers">
          <el-table :data="formatTextAnswers(question.text_answers)" stripe border size="small">
            <el-table-column type="index" label="#" width="50" />
            <el-table-column prop="answer" label="用户回答" />
          </el-table>
        </div>
        <el-divider />
      </div>

      <el-empty v-if="stats.total_submissions === 0" description="暂无提交数据" />

    </el-card>

    <el-skeleton v-if="loading" :rows="10" animated />

    <el-result
      v-if="!loading && error"
      status="error"
      title="加载失败"
      :sub-title="error"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { getFormStatsAPI, type FormStats } from '@/api/form'

const route = useRoute()
const formId = Number(route.params.formId)

const stats = ref<FormStats | null>(null)
const loading = ref(true)
const error = ref<string | null>(null)

const fetchStats = async () => {
  try {
    loading.value = true
    error.value = null
    const res = await getFormStatsAPI(formId)
    stats.value = res
  } catch (err: any) {
    console.error("Failed to fetch stats:", err)
    error.value = err.message || '获取统计数据失败，请稍后重试。'
  } finally {
    loading.value = false
  }
}

const getPercentage = (count: number, total: number) => {
  if (total === 0) return 0
  return parseFloat(((count / total) * 100).toFixed(2))
}

const formatTextAnswers = (answers: string[] | undefined) => {
  if (!answers) return []
  return answers.map(answer => ({ answer }))
}

onMounted(() => {
  if (formId) {
    fetchStats()
  } else {
    error.value = "无效的表单ID"
    loading.value = false
  }
})
</script>

<style scoped>
.stats-container {
  padding: 10px;
}
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.question-stat {
  margin-bottom: 20px;
}
.question-stat h4 {
  margin-bottom: 15px;
}
.option-stat {
  display: flex;
  align-items: center;
  margin-bottom: 8px;
}
.option-text {
  width: 150px;
  text-align: right;
  margin-right: 15px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
.option-progress {
  flex-grow: 1;
}
.option-count {
  width: 80px;
  margin-left: 15px;
}
</style>

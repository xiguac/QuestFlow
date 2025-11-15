<template>
  <div class="stats-container">
    <el-card v-if="!loading && stats" v-loading="loading">
      <template #header>
        <div class="card-header">
          <span>表单统计</span>
          <div class="header-actions">
            <el-tag type="success">总提交数: {{ stats.total_submissions }}</el-tag>
            <el-button
              type="primary"
              :icon="Download"
              @click="openExportDialog"
              :disabled="stats.total_submissions === 0"
            >
              导出数据
            </el-button>
          </div>
        </div>
      </template>

      <div v-for="(question, index) in stats.question_stats" :key="question.question_id" class="question-stat">
        <h4>{{ index + 1 }}. {{ question.title }}</h4>
        <div v-if="['single_choice', 'multi_choice', 'judgment'].includes(question.question_type) && question.option_stats">
          <div v-for="option in question.option_stats" :key="option.text" class="option-stat">
            <span class="option-text">{{ option.text }}</span>
            <el-progress :percentage="getPercentage(option.count, getTotalVotes(question))" class="option-progress" />
            <span class="option-count">{{ option.count }} 票</span>
          </div>
        </div>
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

    <el-result v-if="!loading && error" status="error" title="加载失败" :sub-title="error" />

    <el-dialog v-model="exportDialogVisible" title="导出提交数据" width="700px">
      <el-form label-width="100px" class="export-form">
        <el-form-item label="提交时间">
          <el-date-picker
            v-model="dateRange"
            type="datetimerange"
            range-separator="至"
            start-placeholder="开始时间"
            end-placeholder="结束时间"
            :shortcuts="shortcuts"
            clearable
          />
        </el-form-item>

        <el-divider>筛选条件</el-divider>

        <div v-for="(condition, index) in exportConditions" :key="index" class="condition-row">
          <el-select
            v-model="condition.questionId"
            placeholder="选择题目"
            class="condition-item"
            @change="onQuestionChange(condition)"
          >
            <el-option
              v-for="q in availableQuestions"
              :key="q.id"
              :label="q.title"
              :value="q.id"
            />
          </el-select>

          <el-select v-model="condition.operator" placeholder="操作" class="condition-item short">
            <el-option v-for="op in getOperators(condition.questionType)" :key="op.value" :label="op.label" :value="op.value" />
          </el-select>

          <el-select
            v-if="['single_choice', 'multi_choice', 'judgment'].includes(condition.questionType)"
            v-model="condition.value"
            placeholder="选择选项"
            class="condition-item"
            :multiple="condition.questionType === 'multi_choice'"
            clearable
          >
            <el-option v-for="opt in getOptions(condition.questionId)" :key="opt.id" :label="opt.text" :value="opt.id" />
          </el-select>

          <el-input v-if="condition.questionType === 'text_input'" v-model="condition.value" placeholder="输入文本" class="condition-item" />

          <el-button type="danger" :icon="Delete" circle plain @click="removeCondition(index)"></el-button>
        </div>

        <el-button type="primary" link @click="addCondition">
          <el-icon><Plus /></el-icon>添加筛选条件
        </el-button>
      </el-form>

      <template #footer>
        <span class="dialog-footer">
          <el-button @click="exportDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="handleExport" :loading="isExporting">
            确认导出
          </el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed, reactive } from 'vue'
import { useRoute } from 'vue-router'
import { getFormStatsAPI, getFormDetailsAPI, exportSubmissionsAPI, type FormStats, type Question, type QuestionType, type FilterCondition, type ExportRequestPayload } from '@/api/form'
import { Download, Delete, Plus } from '@element-plus/icons-vue'
import { downloadBlob } from '@/utils/download'
import { ElMessage } from 'element-plus'

interface QuestionStat {
  question_id: string;
  question_type: string;
  title: string;
  option_stats?: { text: string; count: number; }[];
  text_answers?: string[];
}
type Operator = { label: string; value: string };
interface LocalFilterCondition {
  questionId: string;
  questionType: QuestionType;
  operator: string;
  value: string | string[];
}

const route = useRoute()
const formId = Number(route.params.formId)
const stats = ref<FormStats | null>(null)
const loading = ref(true)
const error = ref<string | null>(null)
const formDefinition = ref<Question[]>([])

const exportDialogVisible = ref(false)
const isExporting = ref(false)
const dateRange = ref<[Date, Date] | null>(null)
const exportConditions = ref<LocalFilterCondition[]>([])

const availableQuestions = computed(() => formDefinition.value.filter(q => ['single_choice', 'multi_choice', 'judgment', 'text_input'].includes(q.type)))

const fetchStats = async () => {
  try {
    loading.value = true
    error.value = null
    const [statsRes, detailsRes] = await Promise.all([
      getFormStatsAPI(formId),
      getFormDetailsAPI(formId)
    ]);
    stats.value = statsRes;
    formDefinition.value = detailsRes.Definition.questions || [];
  } catch (err: any) {
    console.error("Failed to fetch data:", err)
    error.value = err.message || '获取页面数据失败，请稍后重试。'
  } finally {
    loading.value = false
  }
}
const getTotalVotes = (question: QuestionStat): number => {
  if (!stats.value) return 0;
  if (question.question_type === 'multi_choice') {
    return question.option_stats?.reduce((sum, option) => sum + option.count, 0) || 0;
  }
  return stats.value.total_submissions;
}
const getPercentage = (count: number, total: number) => {
  if (total === 0) return 0
  return parseFloat(((count / total) * 100).toFixed(2))
}
const formatTextAnswers = (answers: string[] | undefined) => {
  if (!answers) return []
  return answers.map(answer => ({ answer }))
}

const openExportDialog = () => {
  dateRange.value = null;
  exportConditions.value = [];
  addCondition();
  exportDialogVisible.value = true;
}

const addCondition = () => {
  exportConditions.value.push(reactive({
    questionId: '',
    questionType: 'single_choice',
    operator: 'equals',
    value: ''
  }));
}

const removeCondition = (index: number) => {
  exportConditions.value.splice(index, 1);
}

const onQuestionChange = (condition: LocalFilterCondition) => {
  const selectedQuestion = formDefinition.value.find(q => q.id === condition.questionId);
  if (selectedQuestion) {
    condition.questionType = selectedQuestion.type;
    condition.operator = getOperators(selectedQuestion.type)[0].value;
    condition.value = selectedQuestion.type === 'multi_choice' ? [] : '';
  }
}

const getOperators = (type: QuestionType): Operator[] => {
  switch (type) {
    case 'single_choice':
    case 'judgment':
      return [{ label: '等于', value: 'equals' }, { label: '不等于', value: 'not_equals' }];
    case 'multi_choice':
      return [
        { label: '包含任意', value: 'contains' },
        { label: '不包含', value: 'not_contains' },
        { label: '完全匹配', value: 'equals' }
      ];
    case 'text_input':
      return [{ label: '等于', value: 'equals' }];
    default:
      return [];
  }
}

const getOptions = (questionId: string) => {
  const question = formDefinition.value.find(q => q.id === questionId);
  return question?.options || [];
}

const handleExport = async () => {
  isExporting.value = true;
  try {
    const payload: ExportRequestPayload = {};
    if (dateRange.value && dateRange.value.length === 2) {
      payload.startTime = dateRange.value[0].toISOString();
      payload.endTime = dateRange.value[1].toISOString();
    }

    const validConditions: FilterCondition[] = exportConditions.value
      .filter(c => c.questionId && c.value && (Array.isArray(c.value) ? c.value.length > 0 : c.value !== ''))
      .map(c => {
        const normalizedValue = Array.isArray(c.value) ? c.value : [c.value];
        return {
          questionId: c.questionId,
          questionType: c.questionType,
          operator: c.operator as FilterCondition['operator'],
          value: normalizedValue
        };
      });

    if (validConditions.length > 0) {
      payload.conditions = validConditions;
    }

    const res = await exportSubmissionsAPI(formId, payload);

    downloadBlob(res.blob, res.fileName);

    exportDialogVisible.value = false;
    ElMessage.success('导出任务已开始，请注意浏览器下载');
  } catch (error) {
    console.error('Export failed:', error);
    // 错误现在由 request.ts 拦截器统一处理并弹出 ElMessage
  } finally {
    isExporting.value = false;
  }
}

const shortcuts = [
  { text: '最近一周', value: () => { const end = new Date(); const start = new Date(); start.setTime(start.getTime() - 3600 * 1000 * 24 * 7); return [start, end]; }},
  { text: '最近一个月', value: () => { const end = new Date(); const start = new Date(); start.setMonth(start.getMonth() - 1); return [start, end]; }},
  { text: '最近三个月', value: () => { const end = new Date(); const start = new Date(); start.setMonth(start.getMonth() - 3); return [start, end]; }},
]

onMounted(() => {
  if (formId) {
    fetchStats();
  } else {
    error.value = "无效的表单ID";
    loading.value = false;
  }
});
</script>

<style scoped>
.stats-container { padding: 10px; }
.card-header { display: flex; justify-content: space-between; align-items: center; }
.header-actions { display: flex; align-items: center; gap: 16px; }
.question-stat { margin-bottom: 20px; }
.question-stat h4 { margin-bottom: 15px; }
.option-stat { display: flex; align-items: center; margin-bottom: 8px; }
.option-text { width: 150px; text-align: right; margin-right: 15px; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
.option-progress { flex-grow: 1; }
.option-count { width: 80px; margin-left: 15px; }
.export-form .condition-row { display: flex; align-items: center; gap: 10px; margin-bottom: 15px; }
.export-form .condition-item { flex: 1; }
.export-form .condition-item.short { flex: 0 0 120px; }
</style>

import { ref, computed } from 'vue'
import { defineStore } from 'pinia'
import type { Question } from '@/api/form' // 复用我们之前定义的 Question 类型
import { ElMessage } from 'element-plus'

// 导出 FormDefinition 接口，以便其他文件（如 api/form.ts）可以引用
export interface FormDefinition {
  title: string
  description: string
  questions: Question[]
}

export const useEditorStore = defineStore('editor', () => {
  // --- State ---
  // 当前正在编辑的表单 ID，null 表示新建
  const formId = ref<number | null>(null)

  // 存储整个表单的定义
  const formDefinition = ref<FormDefinition>({
    title: '',
    description: '',
    questions: []
  })

  // 存储当前被选中的问题ID
  const selectedQuestionId = ref<string | null>(null)

  // --- Getters ---
  // 计算属性，返回当前被选中的问题对象
  const selectedQuestion = computed(() => {
    if (selectedQuestionId.value === null) {
      return null
    }
    return formDefinition.value.questions.find((q) => q.id === selectedQuestionId.value) || null
  })

  // --- Actions ---
  /**
   * 初始化/重置为一个新的空白表单
   */
  function initNewForm() {
    formId.value = null
    formDefinition.value = {
      title: '请输入表单标题',
      description: '请输入表单描述',
      questions: []
    }
    selectedQuestionId.value = null
  }

  /**
   * 加载一个已有的表单数据到 Store
   * @param id 表单ID
   * @param definition 表单的定义
   */
  function setForm(id: number, definition: FormDefinition) {
    formId.value = id
    formDefinition.value = definition
    selectedQuestionId.value = null // 加载后不选中任何问题
  }


  /**
   * 添加一个新问题到表单
   * @param question - 要添加的问题对象
   */
  function addQuestion(question: Question) {
    formDefinition.value.questions.push(question)
    // 自动选中新添加的问题
    selectQuestion(question.id)
    ElMessage.success('已添加新题目')
  }

  /**
   * 选中一个问题，用于在右侧设置面板中显示其属性
   * @param id - 问题的ID，如果传入 null，则取消选中（选中整个表单）
   */
  function selectQuestion(id: string | null) {
    selectedQuestionId.value = id
  }

  /**
   * 删除一个问题
   * @param id - 要删除的问题的ID
   */
  function deleteQuestion(id: string) {
    const index = formDefinition.value.questions.findIndex((q) => q.id === id)
    if (index !== -1) {
      formDefinition.value.questions.splice(index, 1)
      ElMessage.success('题目已删除')

      // 如果删除的是当前选中的问题，则取消选中
      if (selectedQuestionId.value === id) {
        selectedQuestionId.value = null
      }
    }
  }


  return {
    formId,
    formDefinition,
    selectedQuestionId,
    selectedQuestion,
    initNewForm,
    setForm,
    addQuestion,
    selectQuestion,
    deleteQuestion, // <-- 导出新 action
  }
})

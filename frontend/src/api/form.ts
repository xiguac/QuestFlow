import request from './request'
import type { FormDefinition } from '@/stores/editor' // 从 editor store 导入类型

// --- 类型定义 ---

// 后端返回的完整表单信息
export interface FormInfo {
  ID: number
  FormKey: string
  CreatorID: number
  Title: string
  Description: string
  Status: number // 1-草稿, 2-已发布, 3-已关闭
  CreatedAt: string
  UpdatedAt: string
}

// 用于编辑时获取的表单详情，包含 Definition
export interface FormDetails extends FormInfo {
  Definition: {
    questions: Question[]
  }
}

// 表单定义中单个问题的类型
export interface Question {
  id: string
  type: 'single_choice' | 'text_input' // 随着功能增加，这里会扩展
  title: string
  options?: { id: string; text: string }[]
}

// 公开表单的完整定义类型
export interface PublicForm {
  form_key: string
  title: string
  description: string
  definition: {
    questions: Question[]
  }
}

// 表单统计数据的类型 (与后端对应)
export interface FormStats {
  total_submissions: number
  question_stats: {
    question_id: string
    question_type: string
    title: string
    option_stats?: {
      text: string
      count: number
    }[]
    text_answers?: string[]
  }[]
}

// 创建或更新表单时发送的数据
export interface UpsertFormData {
  title: string
  description: string
  definition: FormDefinition['questions']
}

// --- API 函数 ---

/**
 * 创建一个新表单
 * @param formData
 */
export const createFormAPI = (formData: UpsertFormData) => {
  const requestData = {
    title: formData.title,
    description: formData.description,
    definition: { questions: formData.definition }
  }
  return request<any, FormInfo>({
    url: '/forms',
    method: 'POST',
    data: requestData
  })
}

/**
 * 更新一个已有表单
 * @param formId
 * @param formData
 */
export const updateFormAPI = (formId: number, formData: UpsertFormData) => {
  const requestData = {
    title: formData.title,
    description: formData.description,
    definition: { questions: formData.definition }
  }
  return request<any, FormInfo>({
    url: `/forms/${formId}`,
    method: 'PUT',
    data: requestData
  })
}

/**
 * 获取当前用户创建的表单列表
 */
export const getMyFormsAPI = () => {
  return request<any, FormInfo[]>({
    url: '/forms/my',
    method: 'GET'
  })
}

/**

 * 获取单个表单的详细信息以供编辑
 * @param formId
 */
export const getFormDetailsAPI = (formId: number) => {
  return request<any, FormDetails>({
    url: `/forms/${formId}/details`,
    method: 'GET'
  })
}

/**
 * 获取公开的表单定义
 * @param formKey
 */
export const getPublicFormAPI = (formKey: string) => {
  return request<any, PublicForm>({
    url: `/public/forms/${formKey}`,
    method: 'GET'
  })
}

/**
 * 提交表单数据
 * @param formKey
 * @param submissionData
 */
export const submitFormAPI = (formKey: string, submissionData: { data: Record<string, any> }) => {
  return request<{ message_id: string }>({
    url: `/public/forms/${formKey}/submissions`,
    method: 'POST',
    data: submissionData
  })
}

/**
 * 获取表单统计数据
 * @param formId
 */
export const getFormStatsAPI = (formId: number) => {
  return request<any, FormStats>({
    url: `/forms/${formId}/stats`,
    method: 'GET'
  })
}

/**
 * 删除一个表单
 * @param formId
 */
export const deleteFormAPI = (formId: number) => {
  return request<any, null>({
    url: `/forms/${formId}`,
    method: 'DELETE'
  })
}

/**
 * 更新表单状态
 * @param formId
 * @param status 1-草稿, 2-发布, 3-关闭
 */
export const updateFormStatusAPI = (formId: number, status: 1 | 2 | 3) => {
  return request<any, null>({
    url: `/forms/${formId}/status`,
    method: 'PUT',
    data: { status }
  })
}

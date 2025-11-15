import request from './request'
import type { FormDefinition } from '@/stores/editor'

// --- 类型定义 ---
export interface FormInfo {
  ID: number
  FormKey: string
  CreatorID: number
  Title: string
  Description: string
  Status: number
  CreatedAt: string
  UpdatedAt: string
}

export interface FormDetails extends FormInfo {
  Definition: {
    questions: Question[]
  }
}

export type QuestionType = 'single_choice' | 'multi_choice' | 'judgment' | 'text_input';

export interface Question {
  id: string
  type: QuestionType
  title: string
  options?: { id: string; text: string }[]
}

export interface PublicForm {
  form_key: string
  title: string
  description: string
  definition: {
    questions: Question[]
  }
}

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

export interface UpsertFormData {
  title: string
  description: string
  definition: FormDefinition['questions']
}

export interface FilterCondition {
  questionId: string
  questionType: QuestionType
  operator: 'equals' | 'not_equals' | 'contains' | 'not_contains'
  value: string[]
}
export interface ExportRequestPayload {
  startTime?: string
  endTime?: string
  conditions?: FilterCondition[]
}

export interface ExportResponse {
  blob: Blob;
  fileName: string;
}

// --- API 函数 ---

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

export const getMyFormsAPI = () => {
  return request<any, FormInfo[]>({
    url: '/forms/my',
    method: 'GET'
  })
}

export const getFormDetailsAPI = (formId: number) => {
  return request<any, FormDetails>({
    url: `/forms/${formId}/details`,
    method: 'GET'
  })
}

export const getPublicFormAPI = (formKey: string) => {
  return request<any, PublicForm>({
    url: `/public/forms/${formKey}`,
    method: 'GET'
  })
}

export const submitFormAPI = (formKey: string, submissionData: { data: Record<string, any> }) => {
  return request<{ message_id: string }>({
    url: `/public/forms/${formKey}/submissions`,
    method: 'POST',
    data: submissionData
  })
}

export const getFormStatsAPI = (formId: number) => {
  return request<any, FormStats>({
    url: `/forms/${formId}/stats`,
    method: 'GET'
  })
}

export const deleteFormAPI = (formId: number) => {
  return request<any, null>({
    url: `/forms/${formId}`,
    method: 'DELETE'
  })
}

export const updateFormStatusAPI = (formId: number, status: 1 | 2 | 3) => {
  return request<any, null>({
    url: `/forms/${formId}/status`,
    method: 'PUT',
    data: { status }
  })
}

export const exportSubmissionsAPI = (formId: number, payload: ExportRequestPayload) => {
  return request<any, ExportResponse>({
    url: `/forms/${formId}/export`,
    method: 'POST',
    data: payload,
    responseType: 'blob'
  })
}

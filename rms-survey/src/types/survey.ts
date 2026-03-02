export interface SurveyLoadRequest {
  code: string
  password?: string
}

export interface SurveyView {
  id: string
  name: string
  survey: object
  setting: object
  status: number
  mode: string
}

export interface PublicAnswerView {
  id: string
  setting: object
}

// Survey schema types for editor
export interface SurveySchema {
  id: string
  title: string
  pages: SurveyPage[]
}

export interface SurveyPage {
  id: string
  title: string
  elements: SurveyElement[]
}

export interface SurveyElement {
  id: string
  type: 'radio' | 'checkbox' | 'fillBlank' | 'dropdown' | 'rating'
  title: string
  required: boolean
  options?: SurveyOption[]
  placeholder?: string
  minLength?: number
  maxLength?: number
  min?: number
  max?: number
  attachment?: AttachmentConfig
}

export interface SurveyOption {
  id: string
  text: string
  value?: string
}

// Attachment configuration for survey elements
export interface AttachmentConfig {
  enabled: boolean
  maxFiles: number
  maxSize: number       // bytes
  allowedTypes: string[]  // ['.pdf', '.doc', ...]
}

// Attachment info returned after upload
export interface AttachmentInfo {
  fileId: string
  fileName: string
  fileSize: number
  fileType: string
}

// Answer value structure supporting nested attachments
export interface AnswerValue {
  value: string | string[] | number | null
  attachments?: AttachmentInfo[]
}

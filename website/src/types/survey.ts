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

// Survey statistics
export interface SurveyStatistics {
  totalAnswers: number
  completedAnswers: number
  tempSavedAnswers: number
  questionStats: QuestionStat[]
}

export interface QuestionStat {
  questionId: string
  questionTitle: string
  answerCount: number
  options?: OptionStat[]
}

export interface OptionStat {
  optionId: string
  optionText: string
  count: number
  percentage: number
}

// Survey settings
export interface SurveySetting {
  projectId: string
  startTime?: string
  endTime?: string
  answerLimit?: number
  showProgressBar?: boolean
  showQuestionNumber?: boolean
  shuffleQuestions?: boolean
  allowBack?: boolean
  allowSave?: boolean
  captchaRequired?: boolean
  [key: string]: any
}

export interface SurveySettingRequest {
  projectId: string
  setting: SurveySetting
}

// Survey logic
export interface SurveyLogic {
  projectId: string
  rules: LogicRule[]
}

export interface LogicRule {
  id: string
  condition: LogicCondition
  action: LogicAction
}

export interface LogicCondition {
  questionId: string
  operator: string
  value: any
}

export interface LogicAction {
  type: string
  targetId: string
}

export interface SurveyLogicRequest {
  projectId: string
  logic: SurveyLogic
}

// Query/Dict request types
export interface QueryRequest {
  projectId: string
  questionId: string
  keyword?: string
}

export interface DictRequest {
  dictCode: string
  keyword?: string
}

// Public answer view for validation
export interface ProjectValidation {
  valid: boolean
  message?: string
  project?: SurveyView
}

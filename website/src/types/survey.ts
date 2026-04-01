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
  logic?: SurveyLogic
}

export interface SurveyPage {
  id: string
  title: string
  description?: string
  elements: SurveyElement[]
}

// Cloze blank definition
export interface ClozeBlank {
  id: string
  placeholder?: string
  maxLength?: number
}

export interface SurveyElement {
  id: string
  type: 'radio' | 'checkbox' | 'fillBlank' | 'dropdown' | 'rating' | 'cloze'
  title: string
  required: boolean
  options?: SurveyOption[]
  placeholder?: string
  minLength?: number
  maxLength?: number
  min?: number
  max?: number
  attachment?: AttachmentConfig
  questionAttachments?: QuestionAttachment[]  // Files uploaded by creator for display
  blanks?: ClozeBlank[]  // For cloze type: blank definitions
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

// Question attachment (uploaded by survey creator for display)
export interface QuestionAttachment {
  fileId: string
  fileName: string
  fileType: string
}

// Answer value structure supporting nested attachments
export interface AnswerValue {
  value: string | string[] | number | Record<string, string> | null
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
  ipLimitEnabled?: boolean      // IP restriction master switch
  ipMaxSubmissions?: number     // Max submissions per IP (0 = unlimited)
  ipInterval?: number           // Min seconds between submissions from same IP (0 = no limit)
  [key: string]: any
}

export interface SurveySettingRequest {
  projectId: string
  setting: SurveySetting
}

// Logic operator types
export type LogicOperator =
  | 'eq'        // equals (radio, dropdown)
  | 'neq'       // not equals
  | 'in'        // contains (checkbox)
  | 'not_in'    // not contains
  | 'gt'        // greater than (rating)
  | 'gte'       // greater than or equal
  | 'lt'        // less than
  | 'lte'       // less than or equal
  | 'empty'     // is empty
  | 'not_empty' // is not empty

// Logic action type
export type LogicActionType = 'show' | 'hide'

// Single condition
export interface LogicCondition {
  questionId: string
  operator: LogicOperator
  value: string | number | string[]
}

// Condition group (supports AND/OR combination)
export interface LogicConditionGroup {
  id: string
  type: 'and' | 'or'
  conditions: LogicCondition[]
}

// Logic rule with condition group and action
export interface LogicRule {
  id: string
  name?: string
  condition: LogicConditionGroup
  action: {
    type: LogicActionType
    targetIds: string[]
  }
  enabled: boolean
}

// Survey logic configuration
export interface SurveyLogic {
  rules: LogicRule[]
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

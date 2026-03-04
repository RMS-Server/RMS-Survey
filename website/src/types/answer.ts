export interface AnswerQuery {
  pageIndex?: number
  pageSize?: number
  projectId?: string
  id?: string
  tempSave?: number
  deleted?: boolean
}

export interface AnswerRequest {
  id?: string
  projectId: string
  answer: object
  metaInfo?: object
  tempSave?: number
  ids?: string[]
}

export interface AnswerView {
  id: string
  projectId: string
  answer: object
  metaInfo: object
  tempSave: number | null
  examScore: number | null
  createBy: string
  createAt: string
  updateAt: string
  isRead: boolean
  readAt?: string
  readBy?: string
  ipAddress: string
}

export interface DownloadQuery {
  projectId: string
  locale?: string
}

// Answer upload result
export interface AnswerUploadResult {
  projectId: string
  schema: object
  importedCount: number
  errors?: string[]
}

// Answer create/update request
export interface AnswerCreateRequest {
  id?: string
  projectId: string
  answer: object
  metaInfo?: object
}

// Timing info for a single question
export interface QuestionTiming {
  questionId: string      // Question ID
  startTime: number       // Timestamp when starting this question (ms)
  firstAnswerTime: number // Timestamp of first answer (ms)
  duration: number        // Time spent on this question (ms)
}

// Timing info for entire survey
export interface TimingInfo {
  surveyStartTime: number        // Survey start timestamp (ms)
  submitTime: number             // Submit timestamp (ms)
  totalDuration: number          // Total time spent (ms)
  questionTimings: QuestionTiming[]
}

// Draft data stored in localStorage
export interface DraftData {
  projectId: string
  answers: Record<string, unknown>
  timing: TimingInfo
  savedAt: number  // Draft save timestamp
}

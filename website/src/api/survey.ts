import request from './index'
import type {
  SurveyLoadRequest,
  SurveyView,
  PublicAnswerView,
  AttachmentInfo,
  SurveyStatistics,
  SurveySetting,
  SurveySettingRequest,
  QueryRequest,
  DictRequest,
  ProjectValidation
} from '@/types/survey'
import type { AnswerRequest } from '@/types/answer'

export const surveyApi = {
  // Load survey for public answering
  loadProject(data: SurveyLoadRequest): Promise<SurveyView> {
    return request.post('/public/loadProject', data)
  },

  // Submit survey answer
  saveAnswer(data: AnswerRequest): Promise<PublicAnswerView> {
    return request.post('/public/saveAnswer', data)
  },

  // Temp save survey answer
  tempSaveAnswer(data: AnswerRequest): Promise<PublicAnswerView> {
    return request.post('/public/tempSaveAnswer', data)
  },

  // Get captcha image
  getCaptcha(): Promise<{ captchaId: string; captchaImg: string }> {
    return request.get('/captcha/get')
  },

  // Verify captcha
  checkCaptcha(captchaId: string, captchaCode: string): Promise<boolean> {
    return request.post('/captcha/check', { captchaId, captchaCode })
  },

  // Upload attachment for a question
  uploadAttachment(
    projectId: string,
    questionId: string,
    file: File
  ): Promise<AttachmentInfo> {
    const formData = new FormData()
    formData.append('projectId', projectId)
    formData.append('questionId', questionId)
    formData.append('file', file)
    return request.post('/public/uploadAttachment', formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    })
  },

  // Validate survey access
  validateProject(data: SurveyLoadRequest): Promise<ProjectValidation> {
    return request.post('/public/validateProject', data)
  },

  // Get survey statistics
  getStatistics(data: SurveyLoadRequest): Promise<SurveyStatistics> {
    return request.post('/public/statistics', data)
  },

  // Preview attachment
  getPreviewUrl(attachmentId: string): string {
    return `/api/public/preview/${attachmentId}`
  },

  // Load query data
  loadQuery(data: QueryRequest): Promise<any[]> {
    return request.post('/public/loadQuery', data)
  },

  // Get query result
  getQueryResult(data: QueryRequest): Promise<any> {
    return request.post('/public/getQueryResult', data)
  },

  // Load dictionary data
  loadDict(data: DictRequest): Promise<any[]> {
    return request.post('/public/loadDict', data)
  },

  // Load exam result
  loadExamResult(data: { projectId: string; answerId: string }): Promise<any> {
    return request.post('/public/loadExamResult', data)
  },

  // Load link result
  loadLinkResult(data: { projectId: string; answerId: string }): Promise<any> {
    return request.post('/public/loadLinkResult', data)
  },

  // Get survey settings
  getSetting(projectId: string): Promise<SurveySetting> {
    return request.get('/survey/setting', { params: { projectId } })
  },

  // Update survey settings
  updateSetting(data: SurveySettingRequest): Promise<void> {
    return request.post('/survey/setting', data)
  }
}

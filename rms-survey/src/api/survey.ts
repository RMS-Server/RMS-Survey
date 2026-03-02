import request from './index'
import type { SurveyLoadRequest, SurveyView, PublicAnswerView, AttachmentInfo } from '@/types/survey'
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
  }
}

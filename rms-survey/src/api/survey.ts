import request from './index'
import type { SurveyLoadRequest, SurveyView, PublicAnswerView } from '@/types/survey'
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
  }
}

<template>
  <div class="survey-fill-container">
    <a-spin :spinning="loading">
      <SurveyEndPage v-if="error" status="error" :title="error">
        <a-button type="primary" @click="loadSurvey">重试</a-button>
      </SurveyEndPage>

      <SurveyEndPage
        v-else-if="deviceBlocked"
        status="warning"
        title="该设备已提交过"
        sub-title="本设备已达到该问卷的最大提交次数，不允许再次提交。"
      />

      <SurveyEndPage v-else-if="submitted" />

      <template v-else-if="survey">
        <h1 class="survey-title">{{ survey.name }}</h1>
        <SurveyRenderer
          :survey="surveySchema"
          :answers="answers"
          :project-id="surveyId"
          :paged="surveySchema.pages.length > 1"
          @update:answers="answers = $event"
          @answer-change="handleAnswerChange"
          @submit="handleSubmit"
        />
      </template>
    </a-spin>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { message } from 'ant-design-vue'
import { surveyApi } from '@/api/survey'
import SurveyRenderer from '@/components/survey/SurveyRenderer.vue'
import SurveyEndPage from '@/components/survey/SurveyEndPage.vue'
import { useLogicEvaluator } from '@/composables/useLogicEvaluator'
import { useTimingTracker } from '@/composables/useTimingTracker'
import { useDraftManager } from '@/composables/useDraftManager'
import { useDeviceFingerprint } from '@/composables/useDeviceFingerprint'
import type { SurveySchema, SurveyElement, AnswerValue } from '@/types/survey'

const route = useRoute()

const code = computed(() => route.params.code as string)

const loading = ref(true)
const error = ref('')
const deviceBlocked = ref(false)
const submitted = ref(false)
const survey = ref<{ id: string; name: string; survey: SurveySchema; setting: Record<string, unknown> } | null>(null)
const answers = ref<Record<string, unknown>>({})

const surveyId = computed(() => survey.value?.id || '')

const surveySchema = computed(() => {
  if (!survey.value?.survey) return { id: '', title: '', pages: [] }
  return survey.value.survey
})

const surveySchemaComputed = computed(() => surveySchema.value)
const answersComputed = computed(() => answers.value)
const { isQuestionVisible } = useLogicEvaluator(surveySchemaComputed, answersComputed)

const { initTiming, recordAnswer, getTimingInfo } = useTimingTracker()
const { saveDraft, loadDraft, clearDraft } = useDraftManager()
const { getFingerprint } = useDeviceFingerprint()

const fingerprint = ref('')

onMounted(() => {
  loadSurvey()
})

async function loadSurvey() {
  loading.value = true
  error.value = ''
  deviceBlocked.value = false
  submitted.value = false
  try {
    // Load survey schema first — we need the project id to run the device check.
    const result = await surveyApi.loadProject({ code: code.value })
    survey.value = result as unknown as typeof survey.value

    // Generate fingerprint and ask the server whether this device is still allowed.
    // If the server says no, short-circuit to the block screen without initializing
    // the form — the draft stays untouched in case the user unblocks later.
    fingerprint.value = await getFingerprint()
    try {
      const check = await surveyApi.checkDevice({
        projectId: surveyId.value,
        deviceFingerprint: fingerprint.value
      })
      if (!check.allowed) {
        deviceBlocked.value = true
        return
      }
    } catch {
      // Check endpoint failing should not block submission attempts — let the
      // server's final check under advisory lock be authoritative.
    }

    if (result.survey) {
      const surveyData = result.survey as SurveySchema
      if (surveyData.pages) {
        surveyData.pages.forEach(page => {
          page.elements?.forEach((el: SurveyElement) => {
            if (el.type === 'checkbox') {
              answers.value[el.id] = []
            } else {
              answers.value[el.id] = ''
            }
          })
        })
      }
    }

    if (surveyId.value) {
      const draft = loadDraft(surveyId.value)
      if (draft) {
        for (const [key, value] of Object.entries(draft.answers)) {
          if (key in answers.value) {
            answers.value[key] = value
          }
        }
        initTiming(draft.timing)
      } else {
        initTiming()
      }
    } else {
      initTiming()
    }
  } catch (e: unknown) {
    const err = e as { message?: string }
    error.value = err.message || '加载问卷失败'
  } finally {
    loading.value = false
  }
}

function handleAnswerChange(questionId: string) {
  recordAnswer(questionId)
  if (surveyId.value) {
    saveDraft(surveyId.value, answers.value, getTimingInfo())
  }
}

async function handleSubmit(submittedAnswers: Record<string, AnswerValue>) {
  if (survey.value?.survey) {
    const surveyData = survey.value.survey
    if (surveyData.pages) {
      for (const page of surveyData.pages) {
        for (const el of page.elements || []) {
          if (!isQuestionVisible(el.id)) continue

          if (el.required) {
            const answer = submittedAnswers[el.id]
            const value = answer?.value
            let isEmpty = false

            if (value === '' || value === undefined || value === null) {
              isEmpty = true
            } else if (Array.isArray(value) && value.length === 0) {
              isEmpty = true
            } else if (el.type === 'cloze' && typeof value === 'object') {
              const clozeValue = value as Record<string, string>
              const blanks = el.blanks || []
              isEmpty = blanks.some(blank => !clozeValue[blank.id]?.trim())
            }

            if (isEmpty) {
              message.warning(`请回答: ${el.title || `题目 ${el.id}`}`)
              return
            }
          }
        }
      }
    }
  }

  try {
    const timing = getTimingInfo()
    await surveyApi.saveAnswer({
      projectId: surveyId.value,
      answer: submittedAnswers,
      metaInfo: { timing },
      deviceFingerprint: fingerprint.value
    })
    clearDraft(surveyId.value)
    submitted.value = true
  } catch (error: any) {
    const errorMsg = error?.response?.data?.message || error?.message || ''
    if (errorMsg.includes('device already submitted')) {
      deviceBlocked.value = true
    } else if (errorMsg.includes('ip already submitted')) {
      message.error('该IP已达到最大提交次数')
    } else if (errorMsg.includes('submission interval too short')) {
      message.error('提交过于频繁，请稍后再试')
    } else {
      message.error('提交失败，请重试')
    }
  }
}

</script>

<style scoped>
.survey-fill-container {
  max-width: 800px;
  margin: 0 auto;
  padding: 20px 16px;
  min-height: 100vh;
  background: transparent;
}

@media (min-width: 768px) {
  .survey-fill-container {
    padding: 40px 20px;
  }
}

.survey-title {
  text-align: center;
  font-size: 20px;
  font-weight: 600;
  margin-bottom: 24px;
  color: var(--color-text-main);
  padding: 0 8px;
}

@media (min-width: 768px) {
  .survey-title {
    font-size: 24px;
    margin-bottom: 32px;
    padding: 0;
  }
}

</style>

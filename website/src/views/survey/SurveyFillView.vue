<template>
  <div class="survey-fill-container">
    <a-spin :spinning="loading">
      <div v-if="error" class="error-state">
        <a-result status="error" :title="error">
          <template #extra>
            <a-button type="primary" @click="loadSurvey">重试</a-button>
          </template>
        </a-result>
      </div>

      <template v-else-if="survey">
        <h1 class="survey-title">{{ survey.name }}</h1>
        <SurveyRenderer
          :survey="surveySchema"
          :answers="answers"
          :project-id="surveyId"
          @update:answers="answers = $event"
          @answer-change="handleAnswerChange"
          @submit="handleSubmit"
          @temp-save="handleTempSave"
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
import { useLogicEvaluator } from '@/composables/useLogicEvaluator'
import { useTimingTracker } from '@/composables/useTimingTracker'
import { useDraftManager } from '@/composables/useDraftManager'
import type { SurveySchema, SurveyElement, AnswerValue } from '@/types/survey'

const route = useRoute()

const code = computed(() => route.params.code as string)

const loading = ref(true)
const error = ref('')
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

// Timing tracker
const { initTiming, recordAnswer, getTimingInfo } = useTimingTracker()

// Draft manager
const { saveDraft, loadDraft, clearDraft } = useDraftManager()

onMounted(() => {
  loadSurvey()
})

async function loadSurvey() {
  loading.value = true
  error.value = ''
  try {
    const result = await surveyApi.loadProject({ code: code.value })
    survey.value = result as unknown as typeof survey.value
    // Initialize answers
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

    // Check for existing draft
    if (surveyId.value) {
      const draft = loadDraft(surveyId.value)
      if (draft) {
        // Merge draft answers with initialized answers (draft takes precedence)
        for (const [key, value] of Object.entries(draft.answers)) {
          if (key in answers.value) {
            answers.value[key] = value
          }
        }
        // Restore timing from draft
        initTiming(draft.timing)
      } else {
        // Fresh start
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

// Handle answer change event from SurveyRenderer
function handleAnswerChange(questionId: string) {
  recordAnswer(questionId)
}

async function handleSubmit(submittedAnswers: Record<string, AnswerValue>) {
  // Validate required fields
  if (survey.value?.survey) {
    const surveyData = survey.value.survey
    if (surveyData.pages) {
      for (const page of surveyData.pages) {
        for (const el of page.elements || []) {
          // Skip validation for hidden questions
          if (!isQuestionVisible(el.id)) continue

          if (el.required) {
            const answer = submittedAnswers[el.id]
            const value = answer?.value
            if (value === '' || value === undefined || value === null ||
                (Array.isArray(value) && value.length === 0)) {
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
      metaInfo: { timing }
    })
    // Clear draft after successful submission
    clearDraft(surveyId.value)
    message.success('提交成功')
  } catch {
    message.error('提交失败，请重试')
  }
}

async function handleTempSave(submittedAnswers: Record<string, AnswerValue>) {
  try {
    const timing = getTimingInfo()
    // Save to localStorage
    saveDraft(surveyId.value, submittedAnswers, timing)
    // Also save to backend
    await surveyApi.tempSaveAnswer({
      projectId: surveyId.value,
      answer: submittedAnswers,
      tempSave: 1
    })
    message.success('保存成功')
  } catch {
    message.error('保存失败，请重试')
  }
}
</script>

<style scoped>
.survey-fill-container {
  max-width: 800px;
  margin: 0 auto;
  padding: 40px 20px;
  min-height: 100vh;
  background: transparent;
}

.survey-title {
  text-align: center;
  font-size: 24px;
  font-weight: 600;
  margin-bottom: 32px;
  color: var(--color-text-main);
}

.error-state {
  padding: 100px 0;
}
</style>

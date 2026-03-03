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
  } catch (e: unknown) {
    const err = e as { message?: string }
    error.value = err.message || '加载问卷失败'
  } finally {
    loading.value = false
  }
}

async function handleSubmit(submittedAnswers: Record<string, AnswerValue>) {
  // Validate required fields
  if (survey.value?.survey) {
    const surveyData = survey.value.survey
    if (surveyData.pages) {
      for (const page of surveyData.pages) {
        for (const el of page.elements || []) {
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
    await surveyApi.saveAnswer({
      projectId: surveyId.value,
      answer: submittedAnswers
    })
    message.success('提交成功')
  } catch {
    message.error('提交失败，请重试')
  }
}

async function handleTempSave(submittedAnswers: Record<string, AnswerValue>) {
  try {
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
  background: #f5f5f5;
}

.survey-title {
  text-align: center;
  font-size: 24px;
  font-weight: 600;
  margin-bottom: 32px;
  color: #1f1f1f;
}

.error-state {
  padding: 100px 0;
}
</style>

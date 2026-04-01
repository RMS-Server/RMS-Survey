<template>
  <div class="survey-renderer">
    <!-- Paged mode: single page with navigation -->
    <template v-if="paged">
      <div class="answer-progress-track">
        <div
          class="answer-progress-fill"
          :style="{ width: totalVisibleQuestions > 0 ? `${Math.round(answeredQuestions / totalVisibleQuestions * 100)}%` : '0%' }"
        ></div>
      </div>
      <div class="page-progress">第 {{ visibleCursor + 1 }} 页 / 共 {{ visiblePageIndices.length }} 页</div>
      <div class="survey-page">
        <h2 v-if="currentPage.title && currentPage.title !== 'Page 1'" class="page-title">{{ currentPage.title }}</h2>
        <p v-if="currentPage.description" class="page-description">{{ currentPage.description }}</p>

        <template v-for="(element, elementIndex) in currentPage.elements" :key="element.id">
          <div v-if="isQuestionVisible(element.id)" class="question-item">
            <div class="question-title">
              <span class="question-number">{{ globalElementIndex(currentPageRawIndex, elementIndex) }}</span>
              <span v-if="element.type !== 'cloze'">{{ element.title }}</span>
              <span v-if="element.required" class="required-mark">*</span>
            </div>

            <div v-if="element.questionAttachments?.length" class="question-attachments-display">
              <div
                v-for="file in element.questionAttachments"
                :key="file.fileId"
                class="question-attachment-item"
              >
                <a-image
                  v-if="isImage(file.fileType)"
                  :src="getPreviewUrl(file.fileId)"
                  class="attachment-image"
                />
                <a-button
                  v-else
                  type="link"
                  :href="getPreviewUrl(file.fileId)"
                  target="_blank"
                  class="attachment-download-btn"
                >
                  <DownloadOutlined />
                  {{ file.fileName }}
                </a-button>
              </div>
            </div>

            <div class="question-content">
              <RadioElement
                v-if="element.type === 'radio'"
                :element="element"
                :value="internalAnswers[element.id]?.value as string"
                @update:value="updateAnswer(element.id, $event)"
              />
              <CheckboxElement
                v-else-if="element.type === 'checkbox'"
                :element="element"
                :value="internalAnswers[element.id]?.value as string[]"
                @update:value="updateAnswer(element.id, $event)"
              />
              <FillBlankElement
                v-else-if="element.type === 'fillBlank'"
                :element="element"
                :value="internalAnswers[element.id]?.value as string"
                @update:value="updateAnswer(element.id, $event)"
              />
              <DropdownElement
                v-else-if="element.type === 'dropdown'"
                :element="element"
                :value="internalAnswers[element.id]?.value as string"
                @update:value="updateAnswer(element.id, $event)"
              />
              <RatingElement
                v-else-if="element.type === 'rating'"
                :element="element"
                :value="internalAnswers[element.id]?.value as number"
                @update:value="updateAnswer(element.id, $event)"
              />
              <ClozeElement
                v-else-if="element.type === 'cloze'"
                :element="element"
                :value="(internalAnswers[element.id]?.value as Record<string, string>) || {}"
                @update:value="updateAnswer(element.id, $event)"
              />
            </div>

            <div v-if="getAttachmentConfig(element) && !preview" class="question-attachment">
              <AttachmentUpload
                :project-id="projectId || ''"
                :question-id="element.id"
                :config="element.attachment!"
                :model-value="internalAnswers[element.id]?.attachments || []"
                @update:model-value="updateAttachments(element.id, $event)"
              />
            </div>
          </div>
        </template>
      </div>

      <div v-if="!preview" class="survey-actions paged-actions">
        <a-button v-if="visibleCursor > 0" @click="prevPage">上一页</a-button>
        <a-button v-if="visibleCursor < visiblePageIndices.length - 1" type="primary" @click="nextPage">下一页</a-button>
        <a-button v-else type="primary" @click="handleSubmit">提交</a-button>
      </div>
    </template>

    <!-- Non-paged mode: all pages at once (preview / single page) -->
    <template v-else>
      <div v-for="(page, pageIndex) in survey.pages" :key="page.id" class="survey-page">
        <h2 v-if="page.title && page.title !== 'Page 1'" class="page-title">{{ page.title }}</h2>
        <p v-if="page.description" class="page-description">{{ page.description }}</p>

        <template v-for="(element, elementIndex) in page.elements" :key="element.id">
          <div v-if="isQuestionVisible(element.id)" class="question-item">
            <div class="question-title">
              <span class="question-number">{{ pageIndex + 1 }}.{{ elementIndex + 1 }}</span>
              <span v-if="element.type !== 'cloze'">{{ element.title }}</span>
              <span v-if="element.required" class="required-mark">*</span>
            </div>

            <div v-if="element.questionAttachments?.length" class="question-attachments-display">
              <div
                v-for="file in element.questionAttachments"
                :key="file.fileId"
                class="question-attachment-item"
              >
                <a-image
                  v-if="isImage(file.fileType)"
                  :src="getPreviewUrl(file.fileId)"
                  class="attachment-image"
                />
                <a-button
                  v-else
                  type="link"
                  :href="getPreviewUrl(file.fileId)"
                  target="_blank"
                  class="attachment-download-btn"
                >
                  <DownloadOutlined />
                  {{ file.fileName }}
                </a-button>
              </div>
            </div>

            <div class="question-content">
              <RadioElement
                v-if="element.type === 'radio'"
                :element="element"
                :value="internalAnswers[element.id]?.value as string"
                @update:value="updateAnswer(element.id, $event)"
              />
              <CheckboxElement
                v-else-if="element.type === 'checkbox'"
                :element="element"
                :value="internalAnswers[element.id]?.value as string[]"
                @update:value="updateAnswer(element.id, $event)"
              />
              <FillBlankElement
                v-else-if="element.type === 'fillBlank'"
                :element="element"
                :value="internalAnswers[element.id]?.value as string"
                @update:value="updateAnswer(element.id, $event)"
              />
              <DropdownElement
                v-else-if="element.type === 'dropdown'"
                :element="element"
                :value="internalAnswers[element.id]?.value as string"
                @update:value="updateAnswer(element.id, $event)"
              />
              <RatingElement
                v-else-if="element.type === 'rating'"
                :element="element"
                :value="internalAnswers[element.id]?.value as number"
                @update:value="updateAnswer(element.id, $event)"
              />
              <ClozeElement
                v-else-if="element.type === 'cloze'"
                :element="element"
                :value="(internalAnswers[element.id]?.value as Record<string, string>) || {}"
                @update:value="updateAnswer(element.id, $event)"
              />
            </div>

            <div v-if="getAttachmentConfig(element) && !preview" class="question-attachment">
              <AttachmentUpload
                :project-id="projectId || ''"
                :question-id="element.id"
                :config="element.attachment!"
                :model-value="internalAnswers[element.id]?.attachments || []"
                @update:model-value="updateAttachments(element.id, $event)"
              />
            </div>
          </div>
        </template>
      </div>

      <div v-if="!preview" class="survey-actions">
        <a-button type="primary" @click="handleSubmit">提交</a-button>
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, computed } from 'vue'
import { message } from 'ant-design-vue'
import { DownloadOutlined } from '@ant-design/icons-vue'
import type { SurveySchema, SurveyElement, AttachmentConfig, AnswerValue, AttachmentInfo } from '@/types/survey'
import { useLogicEvaluator } from '@/composables/useLogicEvaluator'
import RadioElement from './elements/RadioElement.vue'
import CheckboxElement from './elements/CheckboxElement.vue'
import FillBlankElement from './elements/FillBlankElement.vue'
import DropdownElement from './elements/DropdownElement.vue'
import RatingElement from './elements/RatingElement.vue'
import ClozeElement from './elements/ClozeElement.vue'
import AttachmentUpload from './AttachmentUpload.vue'

const props = defineProps<{
  survey: SurveySchema
  answers: Record<string, unknown>
  preview?: boolean
  projectId?: string
  paged?: boolean
}>()

const emit = defineEmits<{
  (e: 'submit', answers: Record<string, AnswerValue>): void
  (e: 'update:answers', answers: Record<string, unknown>): void
  (e: 'answerChange', questionId: string): void
}>()

// Internal answers with attachment support
const internalAnswers = ref<Record<string, AnswerValue>>({})

// Convert survey prop to computed for evaluator
const surveyComputed = computed(() => props.survey)

// Convert internalAnswers to plain format for logic evaluation
const answerVersion = ref(0)
const answersForLogic = computed<Record<string, unknown>>(() => {
  void answerVersion.value
  const plain: Record<string, unknown> = {}
  for (const [key, answer] of Object.entries(internalAnswers.value)) {
    plain[key] = answer.value
  }
  return plain
})

const { isQuestionVisible } = useLogicEvaluator(surveyComputed, answersForLogic)

// Paged navigation: track position in visible pages only
const visibleCursor = ref(0)

// Raw indices of pages that have at least one visible question
const visiblePageIndices = computed<number[]>(() => {
  void answerVersion.value
  return props.survey.pages
    .map((page, idx) => ({ idx, page }))
    .filter(({ page }) => page.elements?.some(el => isQuestionVisible(el.id)))
    .map(({ idx }) => idx)
})

// Clamp cursor when visible pages change (e.g. logic rules hide all questions on current page)
watch(visiblePageIndices, (newIndices) => {
  if (visibleCursor.value >= newIndices.length) {
    visibleCursor.value = Math.max(0, newIndices.length - 1)
  }
})

const currentPageRawIndex = computed(() => visiblePageIndices.value[visibleCursor.value] ?? 0)
const currentPage = computed(() => props.survey.pages[currentPageRawIndex.value] ?? props.survey.pages[0])

// Overall answer progress
const totalVisibleQuestions = computed(() => {
  void answerVersion.value
  let count = 0
  for (const page of props.survey.pages) {
    for (const el of page.elements || []) {
      if (isQuestionVisible(el.id)) count++
    }
  }
  return count
})

const answeredQuestions = computed(() => {
  void answerVersion.value
  let count = 0
  for (const page of props.survey.pages) {
    for (const el of page.elements || []) {
      if (!isQuestionVisible(el.id)) continue
      const answer = internalAnswers.value[el.id]
      const val = answer?.value
      if (val === null || val === undefined || val === '') continue
      if (Array.isArray(val) && val.length === 0) continue
      count++
    }
  }
  return count
})

// Initialize from props
watch(() => props.answers, (val) => {
  const converted: Record<string, AnswerValue> = {}
  for (const [key, value] of Object.entries(val)) {
    if (typeof value === 'object' && value !== null && 'value' in value) {
      converted[key] = value as AnswerValue
    } else {
      converted[key] = { value: value as string | string[] | number | Record<string, string> | null, attachments: [] }
    }
  }
  internalAnswers.value = converted
  answerVersion.value++
}, { immediate: true })

// Reset cursor when survey changes
watch(() => props.survey.id, () => {
  visibleCursor.value = 0
})

function updateAnswer(questionId: string, value: string | string[] | number | Record<string, string> | null) {
  if (!internalAnswers.value[questionId]) {
    internalAnswers.value[questionId] = { value, attachments: [] }
  } else {
    internalAnswers.value[questionId].value = value
  }
  answerVersion.value++
  emitUpdate()
  emit('answerChange', questionId)
}

function updateAttachments(questionId: string, attachments: AttachmentInfo[]) {
  if (!internalAnswers.value[questionId]) {
    internalAnswers.value[questionId] = { value: null, attachments }
  } else {
    internalAnswers.value[questionId].attachments = attachments
  }
  emitUpdate()
}

function emitUpdate() {
  const plain: Record<string, unknown> = {}
  for (const [key, answer] of Object.entries(internalAnswers.value)) {
    if (answer.attachments && answer.attachments.length > 0) {
      plain[key] = { value: answer.value, attachments: answer.attachments }
    } else {
      plain[key] = answer.value
    }
  }
  emit('update:answers', plain)
}

// Validate required fields on the given page, returns true if valid
function validatePage(pageIndex: number): boolean {
  const page = props.survey.pages[pageIndex]
  if (!page) return true
  for (const element of page.elements) {
    if (!element.required) continue
    if (!isQuestionVisible(element.id)) continue
    const answer = internalAnswers.value[element.id]
    const val = answer?.value
    if (element.type === 'checkbox') {
      if (!Array.isArray(val) || val.length === 0) {
        message.warning('请完成当前页的必填题目')
        return false
      }
    } else if (element.type === 'cloze') {
      const blanks = element.blanks || []
      const record = (val as Record<string, string>) || {}
      const incomplete = blanks.some(b => !record[b.id] || record[b.id].trim() === '')
      if (incomplete) {
        message.warning('请完成当前页的必填题目')
        return false
      }
    } else {
      if (val === '' || val === null || val === undefined) {
        message.warning('请完成当前页的必填题目')
        return false
      }
    }
  }
  return true
}

function nextPage() {
  if (!validatePage(currentPageRawIndex.value)) return
  if (visibleCursor.value < visiblePageIndices.value.length - 1) {
    visibleCursor.value++
    window.scrollTo({ top: 0, behavior: 'smooth' })
  }
}

function prevPage() {
  if (visibleCursor.value > 0) {
    visibleCursor.value--
    window.scrollTo({ top: 0, behavior: 'smooth' })
  }
}

// Global question number for paged mode (sequential, counting only visible questions)
function globalElementIndex(pageIndex: number, elementIndex: number): string {
  let count = 0
  for (let pi = 0; pi <= pageIndex; pi++) {
    const page = props.survey.pages[pi]
    for (let ei = 0; ei < page.elements.length; ei++) {
      const el = page.elements[ei]
      if (!isQuestionVisible(el.id)) continue
      count++
      if (pi === pageIndex && ei === elementIndex) return String(count)
    }
  }
  return String(count)
}

function getAttachmentConfig(element: SurveyElement): AttachmentConfig | undefined {
  return element.attachment?.enabled ? element.attachment : undefined
}

function isImage(fileType: string): boolean {
  return ['.jpg', '.jpeg', '.png', '.gif'].includes(fileType.toLowerCase())
}

function getPreviewUrl(fileId: string): string {
  return `/api/public/preview/${fileId}`
}

function handleSubmit() {
  if (props.paged && !validatePage(currentPageRawIndex.value)) return
  emit('submit', internalAnswers.value)
}
</script>

<style scoped>
.survey-renderer {
  background: var(--surface-glass);
  backdrop-filter: blur(var(--blur-strength));
  -webkit-backdrop-filter: blur(var(--blur-strength));
  border-radius: var(--radius-lg);
  padding: 16px;
  border: 1px solid var(--border-glass);
  box-shadow: var(--shadow-raised);
  overflow: hidden;
  word-break: break-word;
  overflow-wrap: break-word;
}

@media (min-width: 768px) {
  .survey-renderer {
    padding: 24px;
  }
}

.answer-progress-track {
  height: 4px;
  background: var(--border-glass);
  border-radius: 2px;
  overflow: hidden;
  margin-bottom: 12px;
}

.answer-progress-fill {
  height: 100%;
  background: var(--color-primary);
  border-radius: 2px;
  transition: width 0.3s ease;
}

.page-progress {
  font-size: 13px;
  color: var(--color-text-muted);
  text-align: center;
  margin-bottom: 16px;
  padding: 6px 12px;
  background: var(--surface-glass-input);
  border-radius: var(--radius-sm);
  border: 1px solid var(--border-glass);
}

.survey-page {
  margin-bottom: 16px;
}

@media (min-width: 768px) {
  .survey-page {
    margin-bottom: 24px;
  }
}

.page-title {
  font-size: 16px;
  font-weight: 600;
  margin-bottom: 8px;
  padding-bottom: 8px;
  border-bottom: 1px solid var(--border-glass);
  color: var(--color-text-main);
}

@media (min-width: 768px) {
  .page-title {
    font-size: 18px;
    margin-bottom: 8px;
  }
}

.page-description {
  font-size: 13px;
  color: var(--color-text-muted);
  margin-bottom: 16px;
  line-height: 1.6;
  white-space: pre-wrap;
  word-break: break-word;
}

@media (min-width: 768px) {
  .page-description {
    font-size: 14px;
    margin-bottom: 20px;
  }
}

.question-item {
  margin-bottom: 16px;
  padding: 12px;
  background: var(--surface-glass-input);
  border-radius: var(--radius-md);
  border: 1px solid var(--border-glass);
  box-shadow: var(--shadow-inset);
  overflow: hidden;
  word-break: break-word;
  overflow-wrap: break-word;
}

@media (min-width: 768px) {
  .question-item {
    margin-bottom: 24px;
    padding: 16px;
  }
}

.question-title {
  font-size: 14px;
  font-weight: 500;
  margin-bottom: 12px;
  color: var(--color-text-main);
  line-height: 1.5;
  word-break: break-word;
  overflow-wrap: break-word;
  white-space: pre-wrap;
}

@media (min-width: 768px) {
  .question-title {
    font-size: 15px;
  }
}

.question-number {
  color: var(--color-primary);
  margin-right: 8px;
}

.required-mark {
  color: #ff4d4f;
  margin-left: 4px;
}

.question-attachments-display {
  margin-bottom: 12px;
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
}

.question-attachment-item {
  display: inline-block;
}

.attachment-image {
  max-width: 100%;
  max-height: 200px;
  border-radius: var(--radius-sm);
}

@media (min-width: 768px) {
  .attachment-image {
    max-width: 300px;
    max-height: 300px;
  }
}

.attachment-download-btn {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 4px 8px;
  height: auto;
}

.question-content {
  margin-top: 8px;
  word-break: break-word;
  overflow-wrap: break-word;
}

.question-content :deep(.ant-radio-wrapper),
.question-content :deep(.ant-checkbox-wrapper) {
  word-break: break-word;
  overflow-wrap: break-word;
}

.question-attachment {
  margin-top: 12px;
  padding-top: 12px;
  border-top: 1px dashed var(--border-glass);
}

.survey-actions {
  display: flex;
  flex-direction: column;
  gap: 12px;
  margin-top: 24px;
  padding-top: 16px;
  border-top: 1px solid var(--border-glass);
}

@media (min-width: 576px) {
  .survey-actions {
    flex-direction: row;
    justify-content: center;
    gap: 16px;
    margin-top: 32px;
    padding-top: 24px;
  }
}

.survey-actions :deep(.ant-btn) {
  min-height: 40px;
}

@media (max-width: 575px) {
  .survey-actions :deep(.ant-btn) {
    width: 100%;
  }
}

.paged-actions {
  justify-content: space-between;
}

@media (min-width: 576px) {
  .paged-actions {
    justify-content: center;
  }
}
</style>

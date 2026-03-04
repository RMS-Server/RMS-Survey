<template>
  <div class="survey-renderer">
    <div v-for="(page, pageIndex) in survey.pages" :key="page.id" class="survey-page">
      <h2 v-if="page.title" class="page-title">{{ page.title }}</h2>

      <template v-for="(element, elementIndex) in page.elements" :key="element.id">
        <div v-if="isQuestionVisible(element.id)" class="question-item">
        <div class="question-title">
          <span class="question-number">{{ pageIndex + 1 }}.{{ elementIndex + 1 }}</span>
          {{ element.title }}
          <span v-if="element.required" class="required-mark">*</span>
        </div>

        <!-- Question attachments (uploaded by creator, visible to respondents) -->
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
        </div>

        <!-- Attachment upload area -->
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
      <a-button @click="handleTempSave">保存草稿</a-button>
      <a-button type="primary" @click="handleSubmit">提交</a-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, computed } from 'vue'
import { DownloadOutlined } from '@ant-design/icons-vue'
import type { SurveySchema, SurveyElement, AttachmentConfig, AnswerValue, AttachmentInfo } from '@/types/survey'
import { useLogicEvaluator } from '@/composables/useLogicEvaluator'
import RadioElement from './elements/RadioElement.vue'
import CheckboxElement from './elements/CheckboxElement.vue'
import FillBlankElement from './elements/FillBlankElement.vue'
import DropdownElement from './elements/DropdownElement.vue'
import RatingElement from './elements/RatingElement.vue'
import AttachmentUpload from './AttachmentUpload.vue'

const props = defineProps<{
  survey: SurveySchema
  answers: Record<string, unknown>
  preview?: boolean
  projectId?: string
}>()

const emit = defineEmits<{
  (e: 'submit', answers: Record<string, AnswerValue>): void
  (e: 'tempSave', answers: Record<string, AnswerValue>): void
  (e: 'update:answers', answers: Record<string, unknown>): void
  (e: 'answerChange', questionId: string): void
}>()

// Internal answers with attachment support - this is the source of truth for answers
const internalAnswers = ref<Record<string, AnswerValue>>({})

// Convert survey prop to computed for evaluator
const surveyComputed = computed(() => props.survey)

// Convert internalAnswers to plain format for logic evaluation
// Use a reactive counter to force re-evaluation
const answerVersion = ref(0)
const answersForLogic = computed<Record<string, unknown>>(() => {
  // Access answerVersion to create dependency
  void answerVersion.value
  const plain: Record<string, unknown> = {}
  for (const [key, answer] of Object.entries(internalAnswers.value)) {
    plain[key] = answer.value
  }
  return plain
})

// Logic evaluator for conditional rendering
const { isQuestionVisible } = useLogicEvaluator(surveyComputed, answersForLogic)

// Initialize from props
watch(() => props.answers, (val) => {
  const converted: Record<string, AnswerValue> = {}
  for (const [key, value] of Object.entries(val)) {
    if (typeof value === 'object' && value !== null && 'value' in value) {
      converted[key] = value as AnswerValue
    } else {
      converted[key] = { value: value as string | string[] | number | null, attachments: [] }
    }
  }
  internalAnswers.value = converted
  answerVersion.value++
}, { immediate: true })

// Emit updates when internal answers change
function updateAnswer(questionId: string, value: string | string[] | number | null) {
  if (!internalAnswers.value[questionId]) {
    internalAnswers.value[questionId] = { value, attachments: [] }
  } else {
    internalAnswers.value[questionId].value = value
  }
  answerVersion.value++  // Trigger logic re-evaluation
  emit('answerChange', questionId)
  emitUpdate()
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

// Get attachment config for an element
function getAttachmentConfig(element: SurveyElement): AttachmentConfig | undefined {
  return element.attachment?.enabled ? element.attachment : undefined
}

// Check if file type is an image
function isImage(fileType: string): boolean {
  return ['.jpg', '.jpeg', '.png', '.gif'].includes(fileType.toLowerCase())
}

// Get preview URL for a file
function getPreviewUrl(fileId: string): string {
  return `/api/public/preview/${fileId}`
}

function handleSubmit() {
  emit('submit', internalAnswers.value)
}

function handleTempSave() {
  emit('tempSave', internalAnswers.value)
}
</script>

<style scoped>
.survey-renderer {
  background: var(--surface-glass);
  backdrop-filter: blur(var(--blur-strength));
  -webkit-backdrop-filter: blur(var(--blur-strength));
  border-radius: var(--radius-lg);
  padding: 24px;
  border: 1px solid var(--border-glass);
  box-shadow: var(--shadow-raised);
}

.survey-page {
  margin-bottom: 24px;
}

.page-title {
  font-size: 18px;
  font-weight: 600;
  margin-bottom: 16px;
  padding-bottom: 8px;
  border-bottom: 1px solid var(--border-glass);
  color: var(--color-text-main);
}

.question-item {
  margin-bottom: 24px;
  padding: 16px;
  background: var(--surface-glass-input);
  border-radius: var(--radius-md);
  border: 1px solid var(--border-glass);
  box-shadow: var(--shadow-inset);
}

.question-title {
  font-size: 15px;
  font-weight: 500;
  margin-bottom: 12px;
  color: var(--color-text-main);
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
  max-width: 300px;
  max-height: 300px;
  border-radius: var(--radius-sm);
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
}

.question-attachment {
  margin-top: 12px;
  padding-top: 12px;
  border-top: 1px dashed var(--border-glass);
}

.survey-actions {
  display: flex;
  justify-content: center;
  gap: 16px;
  margin-top: 32px;
  padding-top: 24px;
  border-top: 1px solid var(--border-glass);
}
</style>

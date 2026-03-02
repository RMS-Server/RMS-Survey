<template>
  <div class="survey-renderer">
    <div v-for="(page, pageIndex) in survey.pages" :key="page.id" class="survey-page">
      <h2 v-if="page.title" class="page-title">{{ page.title }}</h2>

      <div
        v-for="(element, elementIndex) in page.elements"
        :key="element.id"
        class="question-item"
      >
        <div class="question-title">
          <span class="question-number">{{ pageIndex + 1 }}.{{ elementIndex + 1 }}</span>
          {{ element.title }}
          <span v-if="element.required" class="required-mark">*</span>
        </div>

        <div class="question-content">
          <RadioElement
            v-if="element.type === 'radio'"
            :element="element"
            :value="internalAnswers[element.id]?.value as string"
            @update:value="internalAnswers[element.id]!.value = $event"
          />
          <CheckboxElement
            v-else-if="element.type === 'checkbox'"
            :element="element"
            :value="internalAnswers[element.id]?.value as string[]"
            @update:value="internalAnswers[element.id]!.value = $event"
          />
          <FillBlankElement
            v-else-if="element.type === 'fillBlank'"
            :element="element"
            :value="internalAnswers[element.id]?.value as string"
            @update:value="internalAnswers[element.id]!.value = $event"
          />
          <DropdownElement
            v-else-if="element.type === 'dropdown'"
            :element="element"
            :value="internalAnswers[element.id]?.value as string"
            @update:value="internalAnswers[element.id]!.value = $event"
          />
          <RatingElement
            v-else-if="element.type === 'rating'"
            :element="element"
            :value="internalAnswers[element.id]?.value as number"
            @update:value="internalAnswers[element.id]!.value = $event"
          />
        </div>

        <!-- Attachment upload area -->
        <div v-if="getAttachmentConfig(element) && !preview" class="question-attachment">
          <AttachmentUpload
            :project-id="projectId || ''"
            :question-id="element.id"
            :config="element.attachment!"
            :model-value="internalAnswers[element.id]?.attachments || []"
            @update:model-value="internalAnswers[element.id]!.attachments = $event"
          />
        </div>
      </div>
    </div>

    <div v-if="!preview" class="survey-actions">
      <a-button @click="handleTempSave">保存草稿</a-button>
      <a-button type="primary" @click="handleSubmit">提交</a-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import type { SurveySchema, SurveyElement, AttachmentConfig, AnswerValue } from '@/types/survey'
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
}>()

// Internal answers with attachment support
const internalAnswers = ref<Record<string, AnswerValue>>({})

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
}, { immediate: true, deep: true })

// Get attachment config for an element
function getAttachmentConfig(element: SurveyElement): AttachmentConfig | undefined {
  return element.attachment?.enabled ? element.attachment : undefined
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
  background: #fff;
  border-radius: 8px;
  padding: 24px;
}

.survey-page {
  margin-bottom: 24px;
}

.page-title {
  font-size: 18px;
  font-weight: 600;
  margin-bottom: 16px;
  padding-bottom: 8px;
  border-bottom: 1px solid #e8e8e8;
}

.question-item {
  margin-bottom: 24px;
  padding: 16px;
  background: #fafafa;
  border-radius: 4px;
}

.question-title {
  font-size: 15px;
  font-weight: 500;
  margin-bottom: 12px;
  color: #1f1f1f;
}

.question-number {
  color: #1890ff;
  margin-right: 8px;
}

.required-mark {
  color: #ff4d4f;
  margin-left: 4px;
}

.question-content {
  margin-top: 8px;
}

.question-attachment {
  margin-top: 12px;
  padding-top: 12px;
  border-top: 1px dashed #e8e8e8;
}

.survey-actions {
  display: flex;
  justify-content: center;
  gap: 16px;
  margin-top: 32px;
  padding-top: 24px;
  border-top: 1px solid #e8e8e8;
}
</style>

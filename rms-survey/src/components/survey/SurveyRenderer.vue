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
            :value="answers[element.id] as string"
            @update:value="answers[element.id] = $event"
          />
          <CheckboxElement
            v-else-if="element.type === 'checkbox'"
            :element="element"
            :value="answers[element.id] as string[]"
            @update:value="answers[element.id] = $event"
          />
          <FillBlankElement
            v-else-if="element.type === 'fillBlank'"
            :element="element"
            :value="answers[element.id] as string"
            @update:value="answers[element.id] = $event"
          />
          <DropdownElement
            v-else-if="element.type === 'dropdown'"
            :element="element"
            :value="answers[element.id] as string"
            @update:value="answers[element.id] = $event"
          />
          <RatingElement
            v-else-if="element.type === 'rating'"
            :element="element"
            :value="answers[element.id] as number"
            @update:value="answers[element.id] = $event"
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
import { defineProps, defineEmits } from 'vue'
import type { SurveySchema } from '@/types/survey'
import RadioElement from './elements/RadioElement.vue'
import CheckboxElement from './elements/CheckboxElement.vue'
import FillBlankElement from './elements/FillBlankElement.vue'
import DropdownElement from './elements/DropdownElement.vue'
import RatingElement from './elements/RatingElement.vue'

defineProps<{
  survey: SurveySchema
  answers: Record<string, unknown>
  preview?: boolean
}>()

const emit = defineEmits<{
  (e: 'submit'): void
  (e: 'tempSave'): void
}>()

function handleSubmit() {
  emit('submit')
}

function handleTempSave() {
  emit('tempSave')
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

.survey-actions {
  display: flex;
  justify-content: center;
  gap: 16px;
  margin-top: 32px;
  padding-top: 24px;
  border-top: 1px solid #e8e8e8;
}
</style>

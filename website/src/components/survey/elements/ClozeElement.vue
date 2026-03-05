<template>
  <div class="cloze-element">
    <div class="cloze-content">
      <template v-for="(part, index) in parts" :key="index">
        <span v-if="part.type === 'text'" class="cloze-text">{{ part.content }}</span>
        <span v-else-if="part.blankId" class="cloze-blank">
          <a-input
            v-model:value="blankValues[part.blankId]"
            :placeholder="getPlaceholder(part.blankId)"
            :maxlength="getMaxLength(part.blankId)"
            class="cloze-input"
            @change="handleChange"
          />
        </span>
      </template>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, computed } from 'vue'
import type { SurveyElement } from '@/types/survey'

const props = defineProps<{
  element: SurveyElement
  value: Record<string, string>
}>()

const emit = defineEmits<{
  (e: 'update:value', value: Record<string, string>): void
}>()

// Initialize blank values
const blankValues = ref<Record<string, string>>({})

// Parse title into parts (text and blanks)
interface ClozePart {
  type: 'text' | 'blank'
  content?: string
  blankId?: string
}

const parts = computed<ClozePart[]>(() => {
  const blanks = props.element.blanks || []
  const title = props.element.title || ''

  // Match ____ patterns
  const regex = /____/g
  const result: ClozePart[] = []
  let lastIndex = 0
  let match
  let blankIndex = 0

  while ((match = regex.exec(title)) !== null) {
    // Add text before match
    if (match.index > lastIndex) {
      result.push({ type: 'text', content: title.slice(lastIndex, match.index) })
    }

    // Use sequential blanks
    const blankId = blanks[blankIndex]?.id || `blank_${blankIndex}`
    blankIndex++

    result.push({ type: 'blank', blankId })
    lastIndex = match.index + match[0].length
  }

  // Add remaining text
  if (lastIndex < title.length) {
    result.push({ type: 'text', content: title.slice(lastIndex) })
  }

  return result
})

function getPlaceholder(blankId: string): string {
  const blank = props.element.blanks?.find(b => b.id === blankId)
  return blank?.placeholder || '填写'
}

function getMaxLength(blankId: string): number | undefined {
  const blank = props.element.blanks?.find(b => b.id === blankId)
  return blank?.maxLength
}

function handleChange() {
  emit('update:value', { ...blankValues.value })
}

// Watch for external value changes
watch(() => props.value, (val) => {
  if (val && typeof val === 'object') {
    blankValues.value = { ...val }
  }
}, { immediate: true })

// Initialize blank values from blanks array
watch(() => props.element.blanks, (blanks) => {
  if (blanks) {
    for (const blank of blanks) {
      if (!(blank.id in blankValues.value)) {
        blankValues.value[blank.id] = ''
      }
    }
  }
}, { immediate: true })
</script>

<style scoped>
.cloze-element {
  width: 100%;
  font-size: 14px;
  font-weight: 500;
  color: var(--color-text-main);
  line-height: 1.5;
  word-break: break-word;
  overflow-wrap: break-word;
}

@media (min-width: 768px) {
  .cloze-element {
    font-size: 15px;
  }
}

.cloze-content {
  display: inline;
}

.cloze-text {
  white-space: pre-wrap;
}

.cloze-blank {
  display: inline-block;
  margin: 0 4px;
  vertical-align: baseline;
}

.cloze-input {
  width: 120px;
  min-width: 80px;
  max-width: 200px;
}

@media (max-width: 575px) {
  .cloze-input {
    width: 100px;
    min-width: 60px;
    max-width: 150px;
  }
}

.cloze-input :deep(.ant-input) {
  text-align: center;
  border-bottom: 2px solid var(--color-primary);
  border-radius: 0;
  background: transparent;
}

.cloze-input :deep(.ant-input:focus) {
  box-shadow: none;
  border-bottom-color: var(--color-primary);
}
</style>

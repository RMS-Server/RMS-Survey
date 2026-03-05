<template>
  <a-checkbox-group v-model:value="localValue" @change="handleChange">
    <a-checkbox
      v-for="option in element.options"
      :key="option.id"
      :value="option.value || option.text"
    >
      {{ option.text }}
    </a-checkbox>
  </a-checkbox-group>
</template>

<script setup lang="ts">
import { ref, watch, defineProps, defineEmits } from 'vue'
import type { SurveyElement } from '@/types/survey'

const props = defineProps<{
  element: SurveyElement
  value: string[]
}>()

const emit = defineEmits<{
  (e: 'update:value', value: string[]): void
}>()

const localValue = ref<string[]>(props.value || [])

watch(() => props.value, (val) => {
  localValue.value = val || []
})

function handleChange() {
  emit('update:value', localValue.value)
}
</script>

<style scoped>
.ant-checkbox-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.ant-checkbox-group :deep(.ant-checkbox-wrapper) {
  display: flex;
  align-items: flex-start;
  gap: 8px;
  word-break: break-word;
  overflow-wrap: break-word;
  white-space: pre-wrap;
}

.ant-checkbox-group :deep(.ant-checkbox) {
  flex-shrink: 0;
  margin-top: 2px;
}
</style>

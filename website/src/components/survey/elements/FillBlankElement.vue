<template>
  <a-textarea
    v-model:value="localValue"
    :placeholder="element.placeholder || '请输入您的答案'"
    :minlength="element.minLength"
    :maxlength="element.maxLength"
    :auto-size="{ minRows: 1, maxRows: 6 }"
    @change="handleChange"
  />
  <div v-if="element.maxLength" class="char-count">
    {{ localValue.length }} / {{ element.maxLength }}
  </div>
</template>

<script setup lang="ts">
import { ref, watch, defineProps, defineEmits } from 'vue'
import type { SurveyElement } from '@/types/survey'

const props = defineProps<{
  element: SurveyElement
  value: string
}>()

const emit = defineEmits<{
  (e: 'update:value', value: string): void
}>()

const localValue = ref(props.value || '')

watch(() => props.value, (val) => {
  localValue.value = val || ''
})

function handleChange() {
  emit('update:value', localValue.value)
}
</script>

<style scoped>
.char-count {
  text-align: right;
  font-size: 12px;
  color: var(--color-text-muted);
  margin-top: 4px;
}
</style>

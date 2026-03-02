<template>
  <a-select
    v-model:value="localValue"
    :placeholder="element.placeholder || '请选择'"
    style="width: 100%"
    @change="handleChange"
  >
    <a-select-option
      v-for="option in element.options"
      :key="option.id"
      :value="option.value || option.text"
    >
      {{ option.text }}
    </a-select-option>
  </a-select>
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

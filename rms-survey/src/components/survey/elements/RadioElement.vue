<template>
  <a-radio-group v-model:value="localValue" @change="handleChange">
    <a-radio
      v-for="option in element.options"
      :key="option.id"
      :value="option.value || option.text"
    >
      {{ option.text }}
    </a-radio>
  </a-radio-group>
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
.ant-radio-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}
</style>

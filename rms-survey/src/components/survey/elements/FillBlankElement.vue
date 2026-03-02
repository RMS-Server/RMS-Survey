<template>
  <a-textarea
    v-if="isLongText"
    v-model:value="localValue"
    :placeholder="element.placeholder || '请输入您的答案'"
    :minlength="element.minLength"
    :maxlength="element.maxLength"
    :rows="4"
    @change="handleChange"
  />
  <a-input
    v-else
    v-model:value="localValue"
    :placeholder="element.placeholder || '请输入您的答案'"
    :minlength="element.minLength"
    :maxlength="element.maxLength"
    @change="handleChange"
  />
  <div v-if="element.maxLength" class="char-count">
    {{ localValue.length }} / {{ element.maxLength }}
  </div>
</template>

<script setup lang="ts">
import { ref, watch, computed, defineProps, defineEmits } from 'vue'
import type { SurveyElement } from '@/types/survey'

const props = defineProps<{
  element: SurveyElement
  value: string
}>()

const emit = defineEmits<{
  (e: 'update:value', value: string): void
}>()

const localValue = ref(props.value || '')

const isLongText = computed(() => (element.maxLength || 0) > 200)

watch(() => props.value, (val) => {
  localValue.value = val || ''
})

function handleChange() {
  emit('update:value', localValue.value)
}

const { element } = props
</script>

<style scoped>
.char-count {
  text-align: right;
  font-size: 12px;
  color: #999;
  margin-top: 4px;
}
</style>

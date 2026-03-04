<template>
  <div class="rating-element">
    <a-rate
      v-model:value="localValue"
      :count="element.max || 5"
      :allow-half="false"
      @change="handleChange"
    />
    <span v-if="localValue" class="rating-value">{{ localValue }} / {{ element.max || 5 }}</span>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, defineProps, defineEmits } from 'vue'
import type { SurveyElement } from '@/types/survey'

const props = defineProps<{
  element: SurveyElement
  value: number
}>()

const emit = defineEmits<{
  (e: 'update:value', value: number): void
}>()

const localValue = ref(props.value || 0)

watch(() => props.value, (val) => {
  localValue.value = val || 0
})

function handleChange() {
  emit('update:value', localValue.value)
}
</script>

<style scoped>
.rating-element {
  display: flex;
  align-items: center;
  gap: 12px;
}

.rating-value {
  font-size: 14px;
  color: var(--color-text-muted);
}
</style>

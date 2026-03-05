<template>
  <div class="options-editor">
    <div v-for="(option, index) in localOptions" :key="option.id" class="option-item">
      <a-textarea
        v-model:value="option.text"
        :placeholder="`选项 ${index + 1}`"
        :auto-size="{ minRows: 1, maxRows: 4 }"
        @change="handleUpdate"
      />
      <a-button
        type="text"
        danger
        :disabled="localOptions.length <= 1"
        @click="removeOption(index)"
      >
        <DeleteOutlined />
      </a-button>
    </div>
    <a-button type="dashed" block @click="addOption">
      <PlusOutlined />
      添加选项
    </a-button>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { v4 as uuidv4 } from 'uuid'
import { DeleteOutlined, PlusOutlined } from '@ant-design/icons-vue'
import type { SurveyOption } from '@/types/survey'

const props = defineProps<{
  options?: SurveyOption[]
}>()

const emit = defineEmits<{
  (e: 'update:options', options: SurveyOption[]): void
}>()

const localOptions = ref<SurveyOption[]>([])

watch(() => props.options, (val) => {
  localOptions.value = val ? [...val] : []
}, { immediate: true, deep: true })

function handleUpdate() {
  emit('update:options', [...localOptions.value])
}

function addOption() {
  localOptions.value.push({
    id: uuidv4(),
    text: ''
  })
  handleUpdate()
}

function removeOption(index: number) {
  localOptions.value.splice(index, 1)
  handleUpdate()
}
</script>

<style scoped>
.options-editor {
  margin-top: 8px;
}

.option-item {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 8px;
}

.option-item .ant-input {
  flex: 1;
}
</style>

<template>
  <div class="cloze-editor">
    <div class="cloze-toolbar">
      <a-button type="dashed" size="small" @click="insertBlank">
        <template #icon><PlusOutlined /></template>
        插入空格
      </a-button>
      <span class="blank-count">已插入 {{ blanks.length }} 个空格</span>
    </div>
    <div class="blanks-config" v-if="blanks.length > 0">
      <div class="blanks-header">空格设置</div>
      <div v-for="(blank, index) in blanks" :key="blank.id" class="blank-item">
        <span class="blank-label">空格 {{ index + 1 }}</span>
        <a-input
          v-model:value="blank.placeholder"
          placeholder="占位提示（可选）"
          class="blank-placeholder"
          @blur="emitUpdate"
        />
        <a-input-number
          v-model:value="blank.maxLength"
          :min="1"
          :max="500"
          placeholder="最大长度"
          class="blank-max-length"
          @blur="emitUpdate"
        />
        <a-button type="text" danger size="small" @click="removeBlank(index)">
          <DeleteOutlined />
        </a-button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { v4 as uuidv4 } from 'uuid'
import { PlusOutlined, DeleteOutlined } from '@ant-design/icons-vue'
import type { SurveyElement, ClozeBlank } from '@/types/survey'

const props = defineProps<{
  element: SurveyElement
}>()

const emit = defineEmits<{
  (e: 'update:element', element: SurveyElement): void
}>()

const blanks = ref<ClozeBlank[]>([])

// Initialize once on mount
onMounted(() => {
  if (props.element.blanks && props.element.blanks.length > 0) {
    blanks.value = props.element.blanks.map(b => ({ ...b }))
  }
})

// Insert a blank
function insertBlank() {
  const currentTitle = props.element.title || ''
  const newTitle = currentTitle + '____'

  const newBlank: ClozeBlank = {
    id: uuidv4(),
    placeholder: '',
    maxLength: undefined
  }

  blanks.value.push(newBlank)

  emit('update:element', {
    ...props.element,
    title: newTitle,
    blanks: blanks.value.map(b => ({ ...b }))
  })
}

// Remove a blank
function removeBlank(index: number) {
  if (index < 0 || index >= blanks.value.length) return

  // Remove the (index+1)th ____ from title
  let count = 0
  let newTitle = props.element.title || ''
  newTitle = newTitle.replace(/____/g, (match) => {
    count++
    if (count === index + 1) {
      return ''
    }
    return match
  })

  // Clean up extra spaces
  newTitle = newTitle.replace(/\s+/g, ' ').trim()

  blanks.value.splice(index, 1)

  emit('update:element', {
    ...props.element,
    title: newTitle,
    blanks: blanks.value.map(b => ({ ...b }))
  })
}

// Emit update when blank properties change
function emitUpdate() {
  emit('update:element', {
    ...props.element,
    blanks: blanks.value.map(b => ({ ...b }))
  })
}
</script>

<style scoped>
.cloze-editor {
  margin-top: 8px;
}

.cloze-toolbar {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 8px;
}

.blank-count {
  font-size: 12px;
  color: var(--color-text-muted);
}

.blanks-config {
  margin-top: 12px;
  padding: 12px;
  background: var(--surface-glass-input);
  border: 1px solid var(--border-glass);
  border-radius: var(--radius-sm);
}

.blanks-header {
  font-size: 13px;
  font-weight: 500;
  color: var(--color-text-main);
  margin-bottom: 8px;
}

.blank-item {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 8px;
  flex-wrap: wrap;
}

.blank-item:last-child {
  margin-bottom: 0;
}

.blank-label {
  font-size: 12px;
  color: var(--color-text-muted);
  min-width: 50px;
}

.blank-placeholder {
  flex: 1;
  min-width: 100px;
}

.blank-max-length {
  width: 100px;
}
</style>

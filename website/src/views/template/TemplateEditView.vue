<template>
  <div class="template-edit">
    <div class="edit-header">
      <div class="header-left">
        <a-button @click="handleBack">
          <template #icon><ArrowLeftOutlined /></template>
          返回
        </a-button>
        <a-input
          v-model:value="templateName"
          placeholder="模板名称"
          style="width: 300px; margin-left: 16px"
        />
      </div>
      <div class="header-right">
        <a-button @click="handlePreview">
          <template #icon><EyeOutlined /></template>
          预览
        </a-button>
        <a-button type="primary" :loading="saving" @click="handleSave">
          <template #icon><SaveOutlined /></template>
          保存
        </a-button>
      </div>
    </div>

    <div class="edit-body">
      <div class="toolbar">
        <h4>模板信息</h4>
        <a-form layout="vertical" class="template-info-form">
          <a-form-item label="模式" required>
            <a-select v-model:value="templateMode" placeholder="选择模式">
              <a-select-option value="survey">问卷</a-select-option>
              <a-select-option value="exam">考试</a-select-option>
              <a-select-option value="vote">投票</a-select-option>
            </a-select>
          </a-form-item>
          <a-form-item label="分类">
            <a-input
              v-model:value="templateCategory"
              placeholder="输入分类"
            />
          </a-form-item>
          <a-form-item label="标签">
            <a-input
              v-model:value="templateTag"
              placeholder="输入标签"
            />
          </a-form-item>
          <a-form-item label="共享">
            <a-switch v-model:checked="templateShared" />
          </a-form-item>
        </a-form>

        <a-divider />

        <h4>添加题目</h4>
        <div class="question-types">
          <div
            v-for="type in questionTypes"
            :key="type.value"
            class="question-type-item"
            @click="addQuestion(type.value)"
          >
            <component :is="type.icon" />
            <span>{{ type.label }}</span>
          </div>
        </div>
      </div>

      <div class="canvas">
        <div v-if="elements.length === 0" class="empty-canvas">
          <p>点击左侧添加题目</p>
        </div>
        <draggable
          v-else
          v-model="elements"
          item-key="id"
          handle=".drag-handle"
          class="question-list"
        >
          <template #item="{ element, index }">
            <div
              class="question-item"
              :class="{ selected: selectedIndex === index }"
              @click="selectQuestion(index)"
            >
              <div class="question-header">
                <span class="drag-handle">
                  <HolderOutlined />
                </span>
                <span class="question-number">Q{{ index + 1 }}</span>
                <span class="question-type-tag">{{ getTypeLabel(element.type) }}</span>
                <div class="question-actions">
                  <a-button type="text" size="small" @click.stop="moveUp(index)" :disabled="index === 0">
                    <UpOutlined />
                  </a-button>
                  <a-button type="text" size="small" @click.stop="moveDown(index)" :disabled="index === elements.length - 1">
                    <DownOutlined />
                  </a-button>
                  <a-button type="text" size="small" danger @click.stop="removeQuestion(index)">
                    <DeleteOutlined />
                  </a-button>
                </div>
              </div>
              <div class="question-content">
                <a-input
                  v-model:value="element.title"
                  placeholder="题目内容"
                  @click.stop
                />
                <component
                  :is="getEditorComponent(element.type)"
                  v-model:options="element.options"
                />
              </div>
              <div class="question-footer">
                <a-checkbox v-model:checked="element.required">必填</a-checkbox>
              </div>
            </div>
          </template>
        </draggable>
      </div>

      <div v-if="selectedElement" class="settings-panel">
        <h4>题目设置</h4>
        <a-form layout="vertical">
          <a-form-item label="标题">
            <a-input v-model:value="selectedElement.title" />
          </a-form-item>
          <a-form-item label="必填">
            <a-switch v-model:checked="selectedElement.required" />
          </a-form-item>
          <template v-if="selectedElement.type === 'fillBlank'">
            <a-form-item label="最小长度">
              <a-input-number v-model:value="selectedElement.minLength" :min="0" />
            </a-form-item>
            <a-form-item label="最大长度">
              <a-input-number v-model:value="selectedElement.maxLength" :min="1" />
            </a-form-item>
          </template>
          <template v-if="selectedElement.type === 'rating'">
            <a-form-item label="最小值">
              <a-input-number v-model:value="selectedElement.min" :min="1" />
            </a-form-item>
            <a-form-item label="最大值">
              <a-input-number v-model:value="selectedElement.max" :min="2" />
            </a-form-item>
          </template>
        </a-form>
      </div>
    </div>

    <a-modal
      v-model:open="previewVisible"
      title="预览"
      width="800px"
      :footer="null"
    >
      <SurveyRenderer :survey="surveyData" :answers="previewAnswers" :preview="true" />
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { message } from 'ant-design-vue'
import draggable from 'vuedraggable'
import { v4 as uuidv4 } from 'uuid'
import { templateApi } from '@/api/template'
import {
  ArrowLeftOutlined,
  SaveOutlined,
  EyeOutlined,
  HolderOutlined,
  UpOutlined,
  DownOutlined,
  DeleteOutlined,
  CheckSquareOutlined,
  BorderOutlined,
  FormOutlined,
  DownCircleOutlined,
  StarOutlined
} from '@ant-design/icons-vue'
import SurveyRenderer from '@/components/survey/SurveyRenderer.vue'
import OptionsEditor from '@/components/survey/OptionsEditor.vue'
import type { SurveyElement, SurveySchema } from '@/types/survey'

const router = useRouter()
const route = useRoute()

const templateId = computed(() => route.params.id as string | undefined)
const isEdit = computed(() => !!templateId.value)

const templateName = ref('')
const templateMode = ref('survey')
const templateCategory = ref('')
const templateTag = ref('')
const templateShared = ref(false)
const elements = ref<SurveyElement[]>([])
const selectedIndex = ref(-1)
const saving = ref(false)
const previewVisible = ref(false)

const selectedElement = computed(() => {
  if (selectedIndex.value >= 0 && selectedIndex.value < elements.value.length) {
    return elements.value[selectedIndex.value]
  }
  return null
})

const surveyData = computed(() => ({
  id: templateId.value || 'preview',
  title: templateName.value,
  pages: [{ id: 'page1', title: 'Page 1', elements: elements.value }]
}))

const previewAnswers = computed(() => {
  const answers: Record<string, unknown> = {}
  elements.value.forEach(el => {
    if (el.type === 'checkbox') {
      answers[el.id] = []
    } else {
      answers[el.id] = ''
    }
  })
  return answers
})

const questionTypes = [
  { value: 'radio', label: '单选题', icon: CheckSquareOutlined },
  { value: 'checkbox', label: '多选题', icon: BorderOutlined },
  { value: 'fillBlank', label: '填空题', icon: FormOutlined },
  { value: 'dropdown', label: '下拉题', icon: DownCircleOutlined },
  { value: 'rating', label: '评分题', icon: StarOutlined }
]

const editorComponents: Record<string, unknown> = {
  radio: OptionsEditor,
  checkbox: OptionsEditor,
  dropdown: OptionsEditor
}

function getEditorComponent(type: string) {
  return editorComponents[type] || null
}

function getTypeLabel(type: string) {
  const item = questionTypes.find(t => t.value === type)
  return item?.label || type
}

function addQuestion(type: string) {
  const element: SurveyElement = {
    id: uuidv4(),
    type: type as SurveyElement['type'],
    title: '',
    required: false,
    options: type === 'radio' || type === 'checkbox' || type === 'dropdown'
      ? [{ id: uuidv4(), text: '选项1' }, { id: uuidv4(), text: '选项2' }]
      : undefined,
    min: type === 'rating' ? 1 : undefined,
    max: type === 'rating' ? 5 : undefined,
    attachment: { enabled: false, maxFiles: 1, maxSize: 10485760, allowedTypes: ['.pdf', '.doc', '.docx', '.jpg', '.png'] }
  }
  elements.value.push(element)
  selectedIndex.value = elements.value.length - 1
}

function selectQuestion(index: number) {
  selectedIndex.value = index
}

function removeQuestion(index: number) {
  elements.value.splice(index, 1)
  if (selectedIndex.value === index) {
    selectedIndex.value = -1
  } else if (selectedIndex.value > index) {
    selectedIndex.value--
  }
}

function moveUp(index: number) {
  if (index > 0) {
    const temp = elements.value[index]
    elements.value[index] = elements.value[index - 1]
    elements.value[index - 1] = temp
    if (selectedIndex.value === index) selectedIndex.value--
    else if (selectedIndex.value === index - 1) selectedIndex.value++
  }
}

function moveDown(index: number) {
  if (index < elements.value.length - 1) {
    const temp = elements.value[index]
    elements.value[index] = elements.value[index + 1]
    elements.value[index + 1] = temp
    if (selectedIndex.value === index) selectedIndex.value++
    else if (selectedIndex.value === index + 1) selectedIndex.value--
  }
}

function handleBack() {
  router.push('/template')
}

function handlePreview() {
  previewVisible.value = true
}

async function handleSave() {
  if (!templateName.value) {
    message.warning('请输入模板名称')
    return
  }
  if (!templateMode.value) {
    message.warning('请选择模式')
    return
  }

  saving.value = true
  try {
    const survey: SurveySchema = {
      id: templateId.value || uuidv4(),
      title: templateName.value,
      pages: [{ id: 'page1', title: 'Page 1', elements: elements.value }]
    }

    if (isEdit.value) {
      await templateApi.update({
        id: templateId.value,
        name: templateName.value,
        mode: templateMode.value,
        category: templateCategory.value,
        tag: templateTag.value,
        shared: templateShared.value,
        template: survey as unknown as Record<string, unknown>
      })
      message.success('保存成功')
    } else {
      const id = await templateApi.create({
        name: templateName.value,
        mode: templateMode.value,
        category: templateCategory.value,
        tag: templateTag.value,
        shared: templateShared.value,
        template: survey as unknown as Record<string, unknown>
      })
      message.success('创建成功')
      router.push(`/template/${id}/edit`)
    }
  } catch {
    message.error('保存失败')
  } finally {
    saving.value = false
  }
}

onMounted(async () => {
  if (templateId.value) {
    try {
      const template = await templateApi.get(templateId.value)
      templateName.value = template.name
      templateMode.value = template.mode || 'survey'
      templateCategory.value = template.category || ''
      templateTag.value = template.tag || ''
      templateShared.value = template.shared || false
      if (template.template) {
        const survey = template.template as SurveySchema
        if (survey.pages && survey.pages[0]) {
          elements.value = survey.pages[0].elements || []
        }
      }
    } catch {
      message.error('加载模板失败')
      router.push('/template')
    }
  }
})
</script>

<style scoped>
.template-edit {
  height: 100vh;
  display: flex;
  flex-direction: column;
}

.edit-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 24px;
  background: var(--surface-glass);
  backdrop-filter: blur(var(--blur-strength));
  -webkit-backdrop-filter: blur(var(--blur-strength));
  border-bottom: 1px solid var(--border-glass);
  box-shadow: var(--shadow-inset);
}

.header-left {
  display: flex;
  align-items: center;
}

.header-right {
  display: flex;
  gap: 8px;
}

.edit-body {
  flex: 1;
  display: flex;
  overflow: hidden;
}

.toolbar {
  width: 240px;
  background: var(--surface-glass);
  backdrop-filter: blur(var(--blur-strength));
  -webkit-backdrop-filter: blur(var(--blur-strength));
  border-right: 1px solid var(--border-glass);
  box-shadow: var(--shadow-inset);
  padding: 16px;
  overflow-y: auto;
}

.toolbar h4 {
  margin-bottom: 12px;
  color: var(--color-text-main);
}

.template-info-form {
  margin-bottom: 16px;
}

.question-types {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.question-type-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 12px;
  background: var(--surface-glass-input);
  border: 1px solid var(--border-glass);
  border-radius: var(--radius-sm);
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  color: var(--color-text-main);
  box-shadow: var(--shadow-inset);
}

.question-type-item:hover {
  border-color: var(--color-primary);
  color: var(--color-primary);
  transform: translateY(-2px);
  box-shadow: var(--shadow-raised);
}

.canvas {
  flex: 1;
  padding: 24px;
  overflow-y: auto;
  background: transparent;
}

.empty-canvas {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 300px;
  color: var(--color-text-muted);
}

.question-list {
  max-width: 800px;
  margin: 0 auto;
}

.question-item {
  background: var(--surface-glass);
  backdrop-filter: blur(var(--blur-strength));
  -webkit-backdrop-filter: blur(var(--blur-strength));
  border: 1px solid var(--border-glass);
  border-radius: var(--radius-md);
  padding: 16px;
  margin-bottom: 12px;
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  box-shadow: var(--shadow-raised);
}

.question-item:hover {
  border-color: var(--color-primary);
  box-shadow: var(--shadow-floating);
  transform: translateY(-2px);
}

.question-item.selected {
  border-color: var(--color-primary);
  box-shadow: 0 0 0 2px rgba(252, 121, 97, 0.3), var(--shadow-floating);
}

.question-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 12px;
}

.drag-handle {
  cursor: move;
  color: var(--color-text-muted);
}

.question-number {
  font-weight: 600;
  color: var(--color-primary);
}

.question-type-tag {
  font-size: 12px;
  color: var(--color-text-muted);
  background: var(--surface-glass-input);
  padding: 2px 8px;
  border-radius: var(--radius-sm);
  border: 1px solid var(--border-glass);
}

.question-actions {
  margin-left: auto;
  display: flex;
  gap: 4px;
}

.question-content {
  margin-bottom: 12px;
}

.question-footer {
  display: flex;
  justify-content: flex-end;
}

.settings-panel {
  width: 280px;
  background: var(--surface-glass);
  backdrop-filter: blur(var(--blur-strength));
  -webkit-backdrop-filter: blur(var(--blur-strength));
  border-left: 1px solid var(--border-glass);
  box-shadow: var(--shadow-inset);
  padding: 16px;
  overflow-y: auto;
}

.settings-panel h4 {
  margin-bottom: 16px;
  color: var(--color-text-main);
}
</style>

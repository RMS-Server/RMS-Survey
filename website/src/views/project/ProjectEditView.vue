<template>
  <div class="project-edit">
    <div class="edit-header">
      <div class="header-left">
        <a-button @click="handleBack">
          <template #icon><ArrowLeftOutlined /></template>
          返回
        </a-button>
        <a-input
          v-model:value="projectName"
          placeholder="问卷名称"
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
        <a-button type="primary" :loading="publishing" @click="handlePublish">
          <template #icon><SendOutlined /></template>
          发布
        </a-button>
      </div>
    </div>

    <div class="edit-body">
      <div class="toolbar">
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
          <!-- 题目附件（答题者可见） -->
          <a-divider>题目附件</a-divider>
          <QuestionAttachmentUpload v-model="questionAttachments" />
          <!-- 答题附件设置 -->
          <a-divider>答题附件</a-divider>
          <a-form-item label="允许上传附件">
            <a-switch v-model:checked="attachmentEnabled" />
          </a-form-item>
          <template v-if="attachmentEnabled">
            <a-form-item label="最大文件数">
              <a-input-number v-model:value="attachmentMaxFiles" :min="1" :max="10" />
            </a-form-item>
            <a-form-item label="最大文件大小 (MB)">
              <a-input-number v-model:value="attachmentMaxSizeMB" :min="1" :max="50" />
            </a-form-item>
            <a-form-item label="允许的文件类型">
              <a-select
                v-model:value="attachmentTypes"
                mode="tags"
                placeholder="输入或选择文件类型，如 .pdf, .custom"
                :options="fileTypeOptions"
              />
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
import { useProjectStore } from '@/stores/project'
import {
  ArrowLeftOutlined,
  SaveOutlined,
  SendOutlined,
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
import QuestionAttachmentUpload from '@/components/survey/QuestionAttachmentUpload.vue'
import type { SurveyElement, SurveySchema, QuestionAttachment } from '@/types/survey'

const router = useRouter()
const route = useRoute()
const projectStore = useProjectStore()

const projectId = computed(() => route.params.id as string | undefined)
const isEdit = computed(() => !!projectId.value)

const projectName = ref('')
const elements = ref<SurveyElement[]>([])
const selectedIndex = ref(-1)
const saving = ref(false)
const publishing = ref(false)
const previewVisible = ref(false)

const selectedElement = computed(() => {
  if (selectedIndex.value >= 0 && selectedIndex.value < elements.value.length) {
    return elements.value[selectedIndex.value]
  }
  return null
})

// File type options for attachment config
const fileTypeOptions = [
  { value: '.pdf', label: 'PDF (.pdf)' },
  { value: '.doc', label: 'Word (.doc)' },
  { value: '.docx', label: 'Word (.docx)' },
  { value: '.xls', label: 'Excel (.xls)' },
  { value: '.xlsx', label: 'Excel (.xlsx)' },
  { value: '.jpg', label: '图片 (.jpg)' },
  { value: '.jpeg', label: '图片 (.jpeg)' },
  { value: '.png', label: '图片 (.png)' },
  { value: '.gif', label: '图片 (.gif)' },
]

// Attachment computed properties
const attachmentEnabled = computed({
  get: () => selectedElement.value?.attachment?.enabled ?? false,
  set: (val) => {
    if (selectedElement.value) {
      if (!selectedElement.value.attachment) {
        selectedElement.value.attachment = { enabled: false, maxFiles: 1, maxSize: 10485760, allowedTypes: ['.pdf', '.doc', '.docx', '.jpg', '.png'] }
      }
      selectedElement.value.attachment.enabled = val
    }
  }
})

const attachmentMaxFiles = computed({
  get: () => selectedElement.value?.attachment?.maxFiles ?? 1,
  set: (val) => {
    if (selectedElement.value?.attachment) {
      selectedElement.value.attachment.maxFiles = val
    }
  }
})

const attachmentMaxSizeMB = computed({
  get: () => Math.floor((selectedElement.value?.attachment?.maxSize ?? 10485760) / 1048576),
  set: (val) => {
    if (selectedElement.value?.attachment) {
      selectedElement.value.attachment.maxSize = val * 1048576
    }
  }
})

const attachmentTypes = computed({
  get: () => selectedElement.value?.attachment?.allowedTypes ?? [],
  set: (val) => {
    if (selectedElement.value?.attachment) {
      selectedElement.value.attachment.allowedTypes = val
    }
  }
})

const questionAttachments = computed<QuestionAttachment[]>({
  get: () => selectedElement.value?.questionAttachments ?? [],
  set: (val) => {
    if (selectedElement.value) {
      selectedElement.value.questionAttachments = val
    }
  }
})

const surveyData = computed(() => ({
  id: projectId.value || 'preview',
  title: projectName.value,
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

const editorComponents: Record<string, any> = {
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
  router.push('/project')
}

function handlePreview() {
  previewVisible.value = true
}

async function handleSave() {
  if (!projectName.value) {
    message.warning('请输入问卷名称')
    return
  }

  saving.value = true
  try {
    const survey: SurveySchema = {
      id: projectId.value || uuidv4(),
      title: projectName.value,
      pages: [{ id: 'page1', title: 'Page 1', elements: elements.value }]
    }

    if (isEdit.value) {
      await projectStore.updateProject({
        id: projectId.value,
        name: projectName.value,
        survey: survey as unknown as Record<string, unknown>,
        setting: {}
      })
      message.success('保存成功')
    } else {
      const result = await projectStore.createProject({
        name: projectName.value,
        survey: survey as unknown as Record<string, unknown>,
        setting: {}
      })
      message.success('创建成功')
      router.push(`/project/${result.id}/edit`)
    }
  } catch {
    message.error('保存失败')
  } finally {
    saving.value = false
  }
}

async function handlePublish() {
  await handleSave()
  publishing.value = true
  try {
    await projectStore.updateProject({
      id: projectId.value!,
      name: projectName.value,
      status: 1
    })
    message.success('发布成功')
  } catch {
    message.error('发布失败')
  } finally {
    publishing.value = false
  }
}

onMounted(async () => {
  if (projectId.value) {
    try {
      const project = await projectStore.fetchProject(projectId.value)
      projectName.value = project.name
      if (project.survey) {
        const survey = project.survey as SurveySchema
        if (survey.pages && survey.pages[0]) {
          elements.value = survey.pages[0].elements || []
        }
      }
    } catch {
      message.error('加载问卷失败')
      router.push('/project')
    }
  }
})
</script>

<style scoped>
.project-edit {
  height: 100vh;
  display: flex;
  flex-direction: column;
}

.edit-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 24px;
  background: #fff;
  border-bottom: 1px solid #e8e8e8;
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
  width: 200px;
  background: #fafafa;
  border-right: 1px solid #e8e8e8;
  padding: 16px;
}

.toolbar h4 {
  margin-bottom: 12px;
  color: #666;
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
  padding: 8px 12px;
  background: #fff;
  border: 1px solid #e8e8e8;
  border-radius: 4px;
  cursor: pointer;
  transition: all 0.2s;
}

.question-type-item:hover {
  border-color: #1890ff;
  color: #1890ff;
}

.canvas {
  flex: 1;
  padding: 24px;
  overflow-y: auto;
  background: #f5f5f5;
}

.empty-canvas {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 300px;
  color: #999;
}

.question-list {
  max-width: 800px;
  margin: 0 auto;
}

.question-item {
  background: #fff;
  border: 1px solid #e8e8e8;
  border-radius: 4px;
  padding: 16px;
  margin-bottom: 12px;
  cursor: pointer;
  transition: all 0.2s;
}

.question-item:hover {
  border-color: #1890ff;
}

.question-item.selected {
  border-color: #1890ff;
  box-shadow: 0 0 0 2px rgba(24, 144, 255, 0.2);
}

.question-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 12px;
}

.drag-handle {
  cursor: move;
  color: #999;
}

.question-number {
  font-weight: 600;
  color: #1890ff;
}

.question-type-tag {
  font-size: 12px;
  color: #666;
  background: #f5f5f5;
  padding: 2px 8px;
  border-radius: 4px;
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
  background: #fff;
  border-left: 1px solid #e8e8e8;
  padding: 16px;
  overflow-y: auto;
}

.settings-panel h4 {
  margin-bottom: 16px;
}
</style>

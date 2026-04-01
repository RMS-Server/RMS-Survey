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
        <!-- Page tab bar -->
        <div class="page-tabs-bar">
          <div class="page-tabs-scroll">
            <draggable
              v-model="pages"
              item-key="id"
              handle=".page-drag-handle"
              class="page-tab-list"
              @start="onPageDragStart"
              @end="onPageReorder"
            >
              <template #item="{ element: page, index }">
                <div
                  class="page-tab"
                  :class="{ active: activePageIndex === index }"
                  @click="activePageIndex = index"
                >
                  <HolderOutlined class="page-drag-handle" />
                  <span class="page-tab-name">{{ page.title || `第 ${index + 1} 页` }}</span>
                  <a-button
                    v-if="pages.length > 1"
                    type="text"
                    size="small"
                    class="page-tab-close"
                    @click.stop="confirmRemovePage(index)"
                  >
                    <CloseOutlined />
                  </a-button>
                </div>
              </template>
            </draggable>
          </div>
          <a-button type="dashed" size="small" class="add-page-btn" @click="addPage">
            <PlusOutlined /> 添加页面
          </a-button>
        </div>

        <!-- Page header editor -->
        <div class="page-header-editor">
          <a-input
            v-model:value="activePage.title"
            placeholder="页面标题（可选）"
            class="page-title-input"
          />
          <a-textarea
            v-model:value="activePage.description"
            placeholder="页面说明文字（可选）"
            :auto-size="{ minRows: 1, maxRows: 4 }"
            class="page-desc-input"
          />
        </div>

        <div v-if="activePage.elements.length === 0" class="empty-canvas">
          <p>点击左侧添加题目</p>
        </div>
        <draggable
          v-else
          v-model="activePage.elements"
          item-key="id"
          handle=".drag-handle"
          class="question-list"
        >
          <template #item="{ element, index }">
            <div
              class="question-item"
              :class="{ selected: selectedElementId === element.id }"
              @click="selectQuestion(element.id)"
            >
              <div class="question-header">
                <span class="drag-handle">
                  <HolderOutlined />
                </span>
                <span class="question-number">Q{{ globalQuestionIndex(element.id) }}</span>
                <span class="question-type-tag">{{ getTypeLabel(element.type) }}</span>
                <div class="question-actions">
                  <a-button type="text" size="small" @click.stop="moveUp(index)" :disabled="index === 0">
                    <UpOutlined />
                  </a-button>
                  <a-button type="text" size="small" @click.stop="moveDown(index)" :disabled="index === activePage.elements.length - 1">
                    <DownOutlined />
                  </a-button>
                  <a-dropdown v-if="pages.length > 1" :trigger="['click']">
                    <a-button type="text" size="small" @click.stop>
                      <SwapOutlined />
                    </a-button>
                    <template #overlay>
                      <a-menu @click="handleMoveToPage(element.id, $event)">
                        <a-menu-item
                          v-for="(p, pi) in pages"
                          :key="p.id"
                          :disabled="pi === activePageIndex"
                        >
                          {{ p.title || `第 ${pi + 1} 页` }}
                        </a-menu-item>
                      </a-menu>
                    </template>
                  </a-dropdown>
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

    <!-- 删除页面确认模态框 -->
    <a-modal
      v-model:open="removePageModalVisible"
      title="删除页面"
      @ok="confirmRemovePageAction"
      ok-text="确认"
      cancel-text="取消"
    >
      <p>该页面包含 <strong>{{ removePageTargetQuestionCount }}</strong> 道题目，请选择处理方式：</p>
      <a-radio-group v-model:value="removePageAction" style="margin-top: 12px">
        <a-radio value="move" style="display: block; margin-bottom: 8px">将题目移动到相邻页面</a-radio>
        <a-radio value="delete" style="display: block">删除题目</a-radio>
      </a-radio-group>
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
import { usePageEditor } from '@/composables/usePageEditor'
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
  StarOutlined,
  PlusOutlined,
  CloseOutlined,
  SwapOutlined
} from '@ant-design/icons-vue'
import SurveyRenderer from '@/components/survey/SurveyRenderer.vue'
import OptionsEditor from '@/components/survey/OptionsEditor.vue'
import type { SurveyElement, SurveySchema, SurveyLogic } from '@/types/survey'

const router = useRouter()
const route = useRoute()

const templateId = computed(() => route.params.id as string | undefined)
const isEdit = computed(() => !!templateId.value)

const templateName = ref('')
const templateMode = ref('survey')
const templateCategory = ref('')
const templateTag = ref('')
const templateShared = ref(false)
const selectedElementId = ref<string | null>(null)
const saving = ref(false)
const previewVisible = ref(false)

// Template editor doesn't use logic rules, pass a stub
const surveyLogic = ref<SurveyLogic>({ rules: [] })

// Remove page confirmation state
const removePageModalVisible = ref(false)
const removePageTargetIndex = ref(-1)
const removePageTargetQuestionCount = ref(0)
const removePageAction = ref<'move' | 'delete'>('move')

const {
  pages,
  activePageIndex,
  activePage,
  allElements,
  addPage,
  removePage,
  deletePageWithQuestions,
  moveQuestionToPage,
  cleanupLogicRules,
  loadPages
} = usePageEditor(surveyLogic)

const selectedElement = computed(() => {
  if (!selectedElementId.value) return null
  for (const page of pages.value) {
    const el = page.elements.find(e => e.id === selectedElementId.value)
    if (el) return el
  }
  return null
})

function globalQuestionIndex(elementId: string): number {
  let count = 1
  for (const page of pages.value) {
    for (const el of page.elements) {
      if (el.id === elementId) return count
      count++
    }
  }
  return count
}

let _activePageIdBeforeReorder: string | null = null

function onPageDragStart() {
  _activePageIdBeforeReorder = pages.value[activePageIndex.value]?.id ?? null
}

function onPageReorder() {
  if (_activePageIdBeforeReorder) {
    const newIdx = pages.value.findIndex(p => p.id === _activePageIdBeforeReorder)
    activePageIndex.value = newIdx >= 0 ? newIdx : Math.min(activePageIndex.value, pages.value.length - 1)
  }
  _activePageIdBeforeReorder = null
}

function confirmRemovePage(index: number) {
  const count = pages.value[index].elements.length
  if (count === 0) {
    removePage(index)
    return
  }
  removePageTargetIndex.value = index
  removePageTargetQuestionCount.value = count
  removePageAction.value = 'move'
  removePageModalVisible.value = true
}

function confirmRemovePageAction() {
  const index = removePageTargetIndex.value
  if (removePageAction.value === 'move') {
    removePage(index)
  } else {
    const deletedIds = deletePageWithQuestions(index)
    if (selectedElementId.value && deletedIds.includes(selectedElementId.value)) {
      selectedElementId.value = null
    }
  }
  removePageModalVisible.value = false
}

const surveyData = computed(() => ({
  id: templateId.value || 'preview',
  title: templateName.value,
  pages: pages.value
}))

const previewAnswers = computed(() => {
  const answers: Record<string, unknown> = {}
  allElements.value.forEach(el => {
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
  activePage.value.elements.push(element)
  selectedElementId.value = element.id
}

function selectQuestion(id: string) {
  selectedElementId.value = id
}

function removeQuestion(index: number) {
  const deletedId = activePage.value.elements[index]?.id
  activePage.value.elements.splice(index, 1)
  if (selectedElementId.value === deletedId) {
    selectedElementId.value = null
  }
  if (deletedId) {
    cleanupLogicRules(deletedId)
  }
}

function handleMoveToPage(elementId: string, info: unknown) {
  const key = (info as { key: string | number }).key
  const targetPageIndex = moveQuestionToPage(elementId, String(key))
  activePageIndex.value = targetPageIndex
  selectedElementId.value = elementId
}

function moveUp(index: number) {
  const els = activePage.value.elements
  if (index > 0) {
    const temp = els[index]
    els[index] = els[index - 1]
    els[index - 1] = temp
  }
}

function moveDown(index: number) {
  const els = activePage.value.elements
  if (index < els.length - 1) {
    const temp = els[index]
    els[index] = els[index + 1]
    els[index + 1] = temp
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
      pages: pages.value
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
        loadPages(survey.pages)
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

/* Page tabs bar */
.page-tabs-bar {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 12px;
}

.page-tabs-scroll {
  flex: 1;
  overflow-x: auto;
  scrollbar-width: none;
}

.page-tabs-scroll::-webkit-scrollbar {
  display: none;
}

.page-tab-list {
  display: flex;
  gap: 6px;
  white-space: nowrap;
  min-width: min-content;
}

.page-tab {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 4px 10px;
  background: var(--surface-glass-input);
  border: 1px solid var(--border-glass);
  border-radius: var(--radius-sm);
  cursor: pointer;
  font-size: 13px;
  color: var(--color-text-muted);
  transition: all 0.2s;
  user-select: none;
  white-space: nowrap;
}

.page-tab:hover {
  border-color: var(--color-primary);
  color: var(--color-text-main);
}

.page-tab.active {
  background: var(--color-primary);
  border-color: var(--color-primary);
  color: #fff;
}

.page-drag-handle {
  cursor: move;
  font-size: 12px;
  opacity: 0.6;
}

.page-tab-name {
  max-width: 100px;
  overflow: hidden;
  text-overflow: ellipsis;
}

.page-tab-close {
  width: 16px !important;
  height: 16px !important;
  min-width: 16px !important;
  padding: 0 !important;
  font-size: 10px;
  color: inherit !important;
  opacity: 0.7;
}

.add-page-btn {
  flex-shrink: 0;
}

/* Page header editor */
.page-header-editor {
  display: flex;
  flex-direction: column;
  gap: 8px;
  margin-bottom: 16px;
  padding: 12px;
  background: var(--surface-glass);
  border: 1px solid var(--border-glass);
  border-radius: var(--radius-md);
}

.page-title-input {
  font-size: 15px;
  font-weight: 500;
}

.page-desc-input {
  font-size: 13px;
  color: var(--color-text-muted);
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

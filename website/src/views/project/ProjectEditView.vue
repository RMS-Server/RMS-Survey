<template>
  <div class="project-edit">
    <div class="edit-header">
      <div class="header-left">
        <a-button @click="handleBack">
          <template #icon><ArrowLeftOutlined /></template>
          <span class="btn-text">返回</span>
        </a-button>
        <a-input
          v-model:value="projectName"
          placeholder="问卷名称"
          class="project-name-input"
        />
      </div>
      <div class="header-right">
        <a-button @click="settingModalVisible = true">
          <template #icon><SettingOutlined /></template>
          <span class="btn-text">设置</span>
        </a-button>
        <a-button @click="logicModalVisible = true">
          <template #icon><BranchesOutlined /></template>
          <span class="btn-text">逻辑</span>
          <a-badge v-if="surveyLogic.rules.length > 0" :count="surveyLogic.rules.length" :offset="[5, -5]" />
        </a-button>
        <a-button @click="handlePreview">
          <template #icon><EyeOutlined /></template>
          <span class="btn-text">预览</span>
        </a-button>
        <a-button type="primary" :loading="saving" @click="handleSave">
          <template #icon><SaveOutlined /></template>
          <span class="btn-text">保存</span>
        </a-button>
        <a-button type="primary" :loading="publishing" @click="handlePublish">
          <template #icon><SendOutlined /></template>
          <span class="btn-text">发布</span>
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
            <PlusOutlined />
            <span class="btn-text">添加页面</span>
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
                  <!-- Move to page dropdown -->
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
                <a-textarea
                  v-model:value="element.title"
                  placeholder="题目内容"
                  :auto-size="{ minRows: 1, maxRows: 6 }"
                  @click.stop
                />
                <component
                  :is="getEditorComponent(element.type)"
                  v-model:options="element.options"
                  :element="element"
                  @update:element="updateElement(index, $event)"
                />
              </div>
              <div class="question-footer">
                <a-checkbox v-model:checked="element.required">必填</a-checkbox>
                <a-button
                  type="link"
                  size="small"
                  class="settings-toggle-btn"
                  @click.stop="showMobileSettings(element.id)"
                >
                  更多设置
                </a-button>
              </div>
            </div>
          </template>
        </draggable>
      </div>

      <!-- Desktop Settings Panel -->
      <div v-if="selectedElement" class="settings-panel">
        <h4>题目设置</h4>
        <a-form layout="vertical">
          <a-form-item label="标题">
            <a-textarea v-model:value="selectedElement.title" :auto-size="{ minRows: 1, maxRows: 6 }" />
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

    <!-- Mobile Settings Drawer -->
    <a-drawer
      v-model:open="mobileSettingsVisible"
      title="题目设置"
      placement="bottom"
      :height="'70%'"
      class="mobile-settings-drawer"
    >
      <div v-if="selectedElement" class="mobile-settings-content">
        <a-form layout="vertical">
          <a-form-item label="标题">
            <a-textarea v-model:value="selectedElement.title" :auto-size="{ minRows: 1, maxRows: 6 }" />
          </a-form-item>
          <a-form-item label="必填">
            <a-switch v-model:checked="selectedElement.required" />
          </a-form-item>
          <template v-if="selectedElement.type === 'fillBlank'">
            <a-form-item label="最小长度">
              <a-input-number v-model:value="selectedElement.minLength" :min="0" style="width: 100%" />
            </a-form-item>
            <a-form-item label="最大长度">
              <a-input-number v-model:value="selectedElement.maxLength" :min="1" style="width: 100%" />
            </a-form-item>
          </template>
          <template v-if="selectedElement.type === 'rating'">
            <a-form-item label="最小值">
              <a-input-number v-model:value="selectedElement.min" :min="1" style="width: 100%" />
            </a-form-item>
            <a-form-item label="最大值">
              <a-input-number v-model:value="selectedElement.max" :min="2" style="width: 100%" />
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
              <a-input-number v-model:value="attachmentMaxFiles" :min="1" :max="10" style="width: 100%" />
            </a-form-item>
            <a-form-item label="最大文件大小 (MB)">
              <a-input-number v-model:value="attachmentMaxSizeMB" :min="1" :max="50" style="width: 100%" />
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
    </a-drawer>

    <!-- 逻辑跳转模态框 -->
    <a-modal
      v-model:open="logicModalVisible"
      title="逻辑跳转规则"
      :width="700"
      :footer="null"
    >
      <LogicRuleEditor
        :logic="surveyLogic"
        :elements="allElements"
        @update:logic="surveyLogic = $event"
      />
    </a-modal>

    <!-- 设置模态框 -->
    <a-modal
      v-model:open="settingModalVisible"
      title="问卷设置"
      :width="600"
      :footer="null"
    >
      <SurveySettingEditor
        :project-id="projectId || ''"
        :setting="surveySetting"
        @update:setting="handleUpdateSetting"
        @cancel="settingModalVisible = false"
      />
    </a-modal>

    <!-- 预览模态框 -->
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
import { useProjectStore } from '@/stores/project'
import { surveyApi } from '@/api/survey'
import { usePageEditor } from '@/composables/usePageEditor'
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
  StarOutlined,
  BranchesOutlined,
  SettingOutlined,
  EditOutlined,
  PlusOutlined,
  CloseOutlined,
  SwapOutlined
} from '@ant-design/icons-vue'
import SurveyRenderer from '@/components/survey/SurveyRenderer.vue'
import LogicRuleEditor from '@/components/survey/LogicRuleEditor.vue'
import OptionsEditor from '@/components/survey/OptionsEditor.vue'
import ClozeEditor from '@/components/survey/ClozeEditor.vue'
import QuestionAttachmentUpload from '@/components/survey/QuestionAttachmentUpload.vue'
import SurveySettingEditor from '@/components/survey/SurveySettingEditor.vue'
import type { SurveyElement, SurveySchema, QuestionAttachment, SurveyLogic, SurveySetting } from '@/types/survey'

const router = useRouter()
const route = useRoute()
const projectStore = useProjectStore()

const projectId = computed(() => route.params.id as string | undefined)
const isEdit = computed(() => !!projectId.value)

const projectName = ref('')
const surveyLogic = ref<SurveyLogic>({ rules: [] })
const surveySetting = ref<SurveySetting>({ projectId: '' })
const selectedElementId = ref<string | null>(null)
const saving = ref(false)
const publishing = ref(false)
const previewVisible = ref(false)
const logicModalVisible = ref(false)
const settingModalVisible = ref(false)
const mobileSettingsVisible = ref(false)

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

// Track active page by ID before drag, then restore after reorder
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
    // Empty page, just remove directly
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
    // Element moved to another page; selectedElement computed still finds it
  } else {
    const deletedIds = deletePageWithQuestions(index)
    if (selectedElementId.value && deletedIds.includes(selectedElementId.value)) {
      selectedElementId.value = null
    }
  }
  removePageModalVisible.value = false
}

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
  pages: pages.value,
  logic: surveyLogic.value
}))

const previewAnswers = computed(() => {
  const answers: Record<string, unknown> = {}
  allElements.value.forEach(el => {
    if (el.type === 'checkbox') {
      answers[el.id] = []
    } else if (el.type === 'cloze') {
      answers[el.id] = {}
    } else {
      answers[el.id] = ''
    }
  })
  return answers
})

const questionTypes = [
  { value: 'radio', label: '单选题', icon: CheckSquareOutlined },
  { value: 'checkbox', label: '多选题', icon: BorderOutlined },
  { value: 'fillBlank', label: '简答题', icon: FormOutlined },
  { value: 'cloze', label: '填空题', icon: EditOutlined },
  { value: 'dropdown', label: '下拉题', icon: DownCircleOutlined },
  { value: 'rating', label: '评分题', icon: StarOutlined }
]

const editorComponents: Record<string, any> = {
  radio: OptionsEditor,
  checkbox: OptionsEditor,
  dropdown: OptionsEditor,
  cloze: ClozeEditor
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
    blanks: type === 'cloze' ? [] : undefined,
    attachment: { enabled: false, maxFiles: 1, maxSize: 10485760, allowedTypes: ['.pdf', '.doc', '.docx', '.jpg', '.png'] }
  }
  activePage.value.elements.push(element)
  selectedElementId.value = element.id
}

function updateElement(index: number, updatedElement: SurveyElement) {
  if (index >= 0 && index < activePage.value.elements.length) {
    activePage.value.elements[index] = updatedElement
  }
}

function selectQuestion(id: string) {
  selectedElementId.value = id
}

function showMobileSettings(id: string) {
  selectedElementId.value = id
  mobileSettingsVisible.value = true
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
  // Switch to the target page so the user can see the moved element
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
      pages: pages.value,
      logic: surveyLogic.value
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

async function handleUpdateSetting(setting: SurveySetting) {
  if (!projectId.value) return
  try {
    await surveyApi.updateSetting({ projectId: projectId.value, setting })
    surveySetting.value = setting
    settingModalVisible.value = false
  } catch {
    message.error('保存设置失败')
  }
}

async function loadSetting() {
  if (!projectId.value) return
  try {
    const setting = await surveyApi.getSetting(projectId.value)
    surveySetting.value = setting
  } catch {
    // Ignore - setting may not exist yet
  }
}

onMounted(async () => {
  if (projectId.value) {
    try {
      const project = await projectStore.fetchProject(projectId.value)
      projectName.value = project.name
      if (project.survey) {
        const survey = project.survey as SurveySchema
        loadPages(survey.pages)
        if (survey.logic) {
          surveyLogic.value = survey.logic
        }
      }
      await loadSetting()
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
  padding: 12px 16px;
  background: var(--surface-glass);
  backdrop-filter: blur(var(--blur-strength));
  -webkit-backdrop-filter: blur(var(--blur-strength));
  border-bottom: 1px solid var(--border-glass);
  box-shadow: var(--shadow-inset);
  flex-wrap: wrap;
  gap: 8px;
}

@media (min-width: 768px) {
  .edit-header {
    padding: 12px 24px;
    flex-wrap: nowrap;
  }
}

.header-left {
  display: flex;
  align-items: center;
  gap: 12px;
  flex-wrap: wrap;
}

.project-name-input {
  width: 100%;
  max-width: 200px;
}

@media (min-width: 576px) {
  .project-name-input {
    width: 300px;
    max-width: none;
  }
}

.header-right {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

@media (max-width: 767px) {
  .header-right {
    width: 100%;
    justify-content: flex-end;
  }
}

/* Hide button text on small screens */
@media (max-width: 575px) {
  .btn-text {
    display: none;
  }
}

.edit-body {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

@media (min-width: 992px) {
  .edit-body {
    flex-direction: row;
  }
}

.toolbar {
  width: 100%;
  background: var(--surface-glass);
  backdrop-filter: blur(var(--blur-strength));
  -webkit-backdrop-filter: blur(var(--blur-strength));
  border-right: none;
  border-bottom: 1px solid var(--border-glass);
  box-shadow: var(--shadow-inset);
  padding: 12px 16px;
  flex-shrink: 0;
}

@media (min-width: 992px) {
  .toolbar {
    width: 200px;
    padding: 16px;
    border-right: 1px solid var(--border-glass);
    border-bottom: none;
  }
}

.toolbar h4 {
  display: none;
}

@media (min-width: 992px) {
  .toolbar h4 {
    display: block;
    margin-bottom: 12px;
    color: var(--color-text-main);
  }
}

.question-types {
  display: flex;
  flex-direction: row;
  flex-wrap: wrap;
  gap: 8px;
}

@media (min-width: 992px) {
  .question-types {
    flex-direction: column;
    flex-wrap: nowrap;
  }
}

.question-type-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 12px;
  background: var(--surface-glass-input);
  border: 1px solid var(--border-glass);
  border-radius: var(--radius-sm);
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  color: var(--color-text-main);
  box-shadow: var(--shadow-inset);
  flex: 1;
  min-width: calc(50% - 4px);
  font-size: 13px;
}

@media (min-width: 576px) {
  .question-type-item {
    min-width: auto;
    flex: 0 0 auto;
  }
}

@media (min-width: 992px) {
  .question-type-item {
    flex: 0 0 auto;
    width: 100%;
    padding: 10px 12px;
    font-size: 14px;
  }
}

.question-type-item:hover {
  border-color: var(--color-primary);
  color: var(--color-primary);
  transform: translateY(-2px);
  box-shadow: var(--shadow-raised);
}

.canvas {
  flex: 1;
  padding: 16px;
  overflow-y: auto;
  background: transparent;
}

@media (min-width: 768px) {
  .canvas {
    padding: 24px;
  }
}

/* Page tabs bar */
.page-tabs-bar {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 12px;
  flex-wrap: nowrap;
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

.page-tab.active .page-drag-handle {
  opacity: 0.8;
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

.page-tab-close:hover {
  opacity: 1;
}

.add-page-btn {
  flex-shrink: 0;
  white-space: nowrap;
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
  max-width: 100%;
  margin: 0 auto;
}

@media (min-width: 768px) {
  .question-list {
    max-width: 800px;
  }
}

.question-item {
  background: var(--surface-glass);
  backdrop-filter: blur(var(--blur-strength));
  -webkit-backdrop-filter: blur(var(--blur-strength));
  border: 1px solid var(--border-glass);
  border-radius: var(--radius-md);
  padding: 12px;
  margin-bottom: 8px;
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  box-shadow: var(--shadow-raised);
}

@media (min-width: 768px) {
  .question-item {
    padding: 16px;
    margin-bottom: 12px;
  }
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
  flex-wrap: wrap;
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
  display: flex;
  gap: 4px;
  margin-left: auto;
}

@media (max-width: 575px) {
  .question-actions {
    width: 100%;
    justify-content: flex-end;
    margin-top: 8px;
    margin-left: 0;
  }
}

.question-content {
  margin-bottom: 12px;
}

.question-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.settings-toggle-btn {
  display: inline-flex;
}

@media (min-width: 992px) {
  .settings-toggle-btn {
    display: none;
  }
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
  flex-shrink: 0;
}

@media (max-width: 991px) {
  .settings-panel {
    display: none;
  }
}

.settings-panel h4 {
  margin-bottom: 16px;
  color: var(--color-text-main);
}

/* Mobile settings drawer */
.mobile-settings-content {
  padding: 0 4px;
}

:deep(.mobile-settings-drawer .ant-drawer-content) {
  background: var(--surface-glass-strong);
  backdrop-filter: blur(var(--blur-strength));
  -webkit-backdrop-filter: blur(var(--blur-strength));
}

:deep(.mobile-settings-drawer .ant-drawer-header) {
  background: transparent;
  border-bottom: 1px solid var(--border-glass);
}
</style>

<template>
  <div class="page-container">
    <div class="page-header">
      <div class="header-top">
        <h1 class="page-title">答卷详情</h1>
        <div class="header-actions">
          <a-button
            v-if="answer && !answer.isRead"
            type="primary"
            :loading="markingRead"
            @click="handleMarkRead"
          >
            <template #icon><CheckOutlined /></template>
            标记为已读
          </a-button>
          <a-button
            v-if="answer && answer.isRead"
            :loading="markingRead"
            @click="handleMarkUnread"
          >
            <template #icon><CloseOutlined /></template>
            标记为未读
          </a-button>
          <a-button @click="handleBack">
            <template #icon><ArrowLeftOutlined /></template>
            返回
          </a-button>
        </div>
      </div>
    </div>

    <a-spin :spinning="loading">
      <a-card v-if="answer" :bordered="false">
        <a-descriptions title="答卷信息" bordered>
          <a-descriptions-item label="ID">{{ answer.id }}</a-descriptions-item>
          <a-descriptions-item label="问卷">{{ projectName }}</a-descriptions-item>
          <a-descriptions-item label="状态">
            <a-tag :color="answer.tempSave ? 'orange' : 'green'">
              {{ answer.tempSave ? '草稿' : '已提交' }}
            </a-tag>
          </a-descriptions-item>
          <a-descriptions-item label="提交时间">
            {{ formatDate(answer.createAt) }}
          </a-descriptions-item>
          <a-descriptions-item label="提交人">
            {{ answer.createBy || '匿名' }}
          </a-descriptions-item>
          <a-descriptions-item label="IP地址">
            {{ answer.ipAddress || '-' }}
          </a-descriptions-item>
          <a-descriptions-item v-if="answer.examScore" label="得分">
            {{ answer.examScore }}
          </a-descriptions-item>
          <a-descriptions-item label="阅读状态">
            <a-tag :color="answer.isRead ? 'blue' : 'default'">
              {{ answer.isRead ? '已读' : '未读' }}
            </a-tag>
            <span v-if="answer.isRead && answer.readAt" class="read-time">
              ({{ formatDate(answer.readAt) }})
            </span>
          </a-descriptions-item>
          <a-descriptions-item v-if="hasTimingInfo" label="答题时长">
            <span class="timing-badge">
              <ClockCircleOutlined />
              {{ formatDuration(timingInfo!.totalDuration) }}
            </span>
          </a-descriptions-item>
        </a-descriptions>

        <a-divider />

        <h3>答题内容</h3>
        <div class="answer-content">
          <div
            v-for="(value, key) in answerData"
            :key="key"
            class="answer-item"
          >
            <div class="answer-label">{{ getElementTitle(key as string) }}</div>
            <div class="answer-value">
              {{ formatValue(value) }}
              <div v-if="getAttachments(key as string).length > 0" class="attachment-list">
                <div class="attachment-label">
                  <PaperClipOutlined /> 附件：
                </div>
                <div class="attachment-items">
                  <a
                    v-for="att in getAttachments(key as string)"
                    :key="att.fileId"
                    :href="getAttachmentUrl(att.fileId)"
                    target="_blank"
                    class="attachment-link"
                  >
                    {{ att.fileName }} ({{ formatFileSize(att.fileSize) }})
                  </a>
                </div>
              </div>
              <div v-if="getQuestionTiming(key as string)?.duration" class="question-timing">
                <ClockCircleOutlined />
                耗时: {{ formatDuration(getQuestionTiming(key as string)!.duration) }}
              </div>
            </div>
          </div>
        </div>
      </a-card>
    </a-spin>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { message } from 'ant-design-vue'
import { answerApi } from '@/api/answer'
import { projectApi } from '@/api/project'
import { ArrowLeftOutlined, PaperClipOutlined, ClockCircleOutlined, CheckOutlined, CloseOutlined } from '@ant-design/icons-vue'
import type { AnswerView } from '@/types/answer'
import type { TimingInfo, QuestionTiming } from '@/types/answer'
import type { SurveySchema, SurveyElement, AttachmentInfo } from '@/types/survey'

const route = useRoute()
const router = useRouter()

const answerId = computed(() => route.params.id as string)

const loading = ref(false)
const markingRead = ref(false)
const answer = ref<AnswerView | null>(null)
const projectName = ref('')
const elementMap = ref<Map<string, SurveyElement>>(new Map())
const timingInfo = ref<TimingInfo | null>(null)

interface AnswerValueWithAttachments {
  value: unknown
  attachments?: AttachmentInfo[]
}

const answerData = computed(() => {
  if (!answer.value?.answer) return {}
  return answer.value.answer as Record<string, unknown>
})

onMounted(() => {
  fetchAnswer()
})

async function fetchAnswer() {
  loading.value = true
  try {
    answer.value = await answerApi.get(answerId.value)
    if (answer.value?.projectId) {
      const project = await projectApi.get(answer.value.projectId)
      projectName.value = project.name
      buildElementMap(project.survey as SurveySchema)
    }
    // Extract timing info from metaInfo
    if (answer.value?.metaInfo && typeof answer.value.metaInfo === 'object') {
      const meta = answer.value.metaInfo as { timing?: TimingInfo }
      timingInfo.value = meta.timing || null
    }
  } catch {
    message.error('加载答卷失败')
  } finally {
    loading.value = false
  }
}

function buildElementMap(survey: SurveySchema) {
  elementMap.value.clear()
  if (!survey?.pages) return
  for (const page of survey.pages) {
    if (!page.elements) continue
    for (const element of page.elements) {
      elementMap.value.set(element.id, element)
    }
  }
}

function getElementTitle(elementId: string): string {
  const element = elementMap.value.get(elementId)
  return element?.title || elementId
}

function handleBack() {
  router.push('/answer')
}

async function handleMarkRead() {
  if (!answer.value) return
  markingRead.value = true
  try {
    await answerApi.markRead(answer.value.id)
    answer.value = { ...answer.value, isRead: true }
    message.success('已标记为已读')
  } catch {
    message.error('操作失败')
  } finally {
    markingRead.value = false
  }
}

async function handleMarkUnread() {
  if (!answer.value) return
  markingRead.value = true
  try {
    await answerApi.markUnread(answer.value.id)
    answer.value = { ...answer.value, isRead: false }
    message.success('已标记为未读')
  } catch {
    message.error('操作失败')
  } finally {
    markingRead.value = false
  }
}

function formatDate(dateStr: string) {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  return date.toLocaleString()
}

function formatValue(value: unknown): string {
  if (value === null || value === undefined) return ''
  // Check if this is an answer with attachments structure
  if (typeof value === 'object' && value !== null && 'value' in value) {
    const answerVal = value as AnswerValueWithAttachments
    return formatValue(answerVal.value)
  }
  if (Array.isArray(value)) return value.join(', ')
  if (typeof value === 'object') {
    // For cloze answers (Record<string, string>), format nicely
    const obj = value as Record<string, string>
    const values = Object.values(obj)
    if (values.length > 0 && typeof values[0] === 'string') {
      return values.join(' / ')
    }
    return JSON.stringify(value, null, 2)
  }
  return String(value)
}

function getAttachments(key: string): AttachmentInfo[] {
  const value = answerData.value[key]
  if (typeof value === 'object' && value !== null && 'attachments' in value) {
    const answerVal = value as AnswerValueWithAttachments
    return answerVal.attachments || []
  }
  return []
}

function getAttachmentUrl(fileId: string): string {
  return `/api/public/preview/${fileId}`
}

function formatFileSize(bytes: number): string {
  if (bytes < 1024) return bytes + ' B'
  if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + ' KB'
  return (bytes / (1024 * 1024)).toFixed(1) + ' MB'
}

// Format duration in ms to human readable
function formatDuration(ms: number): string {
  if (ms < 1000) return `${ms}ms`
  if (ms < 60000) return `${(ms / 1000).toFixed(1)}s`
  const minutes = Math.floor(ms / 60000)
  const seconds = Math.round((ms % 60000) / 1000)
  return `${minutes}m ${seconds}s`
}

// Get timing for a specific question
function getQuestionTiming(questionId: string): QuestionTiming | undefined {
  if (!timingInfo.value?.questionTimings) return undefined
  return timingInfo.value.questionTimings.find(t => t.questionId === questionId)
}

// Check if timing info is available
const hasTimingInfo = computed(() => {
  return timingInfo.value && timingInfo.value.totalDuration > 0
})
</script>

<style scoped>
.page-header {
  background: var(--surface-glass);
  backdrop-filter: blur(var(--blur-strength));
  -webkit-backdrop-filter: blur(var(--blur-strength));
  padding: 16px 24px;
  margin-bottom: 16px;
  border-radius: var(--radius-md);
  border: 1px solid var(--border-glass);
  box-shadow: var(--shadow-raised);
}

@media (max-width: 767px) {
  .page-header {
    padding: 12px 16px;
    margin-bottom: 12px;
  }
}

.header-top {
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: 12px;
}

.header-actions {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

@media (max-width: 575px) {
  .header-actions .ant-btn {
    padding: 4px 8px;
    font-size: 13px;
  }
}

.answer-content {
  background: var(--surface-glass-input);
  border-radius: var(--radius-md);
  padding: 16px;
  border: 1px solid var(--border-glass);
}

@media (max-width: 767px) {
  .answer-content {
    padding: 12px;
  }
}

.answer-item {
  display: flex;
  flex-direction: column;
  padding: 12px 0;
  border-bottom: 1px solid var(--border-glass);
}

@media (min-width: 576px) {
  .answer-item {
    flex-direction: row;
  }
}

.answer-item:last-child {
  border-bottom: none;
}

.answer-label {
  font-weight: 500;
  color: var(--color-text-muted);
  margin-bottom: 4px;
}

@media (min-width: 576px) {
  .answer-label {
    width: 200px;
    flex-shrink: 0;
    margin-bottom: 0;
  }
}

.answer-value {
  flex: 1;
  color: var(--color-text-main);
  word-break: break-word;
  overflow-wrap: break-word;
  white-space: pre-wrap;
}

.attachment-list {
  margin-top: 8px;
  padding: 10px 14px;
  background: var(--surface-glass);
  border-radius: var(--radius-sm);
  border: 1px solid var(--border-glass);
}

.attachment-label {
  font-size: 12px;
  color: var(--color-text-muted);
  margin-bottom: 4px;
}

.attachment-label .anticon {
  margin-right: 4px;
}

.attachment-items {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.attachment-link {
  color: var(--color-primary);
  text-decoration: none;
  font-size: 13px;
  word-break: break-all;
}

.attachment-link:hover {
  text-decoration: underline;
}

.timing-badge {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  color: var(--color-primary);
}

.question-timing {
  margin-top: 6px;
  font-size: 12px;
  color: var(--color-text-muted);
  display: flex;
  align-items: center;
  gap: 4px;
}

.read-time {
  margin-left: 8px;
  color: var(--color-text-muted);
  font-size: 12px;
}

/* Responsive descriptions */
:deep(.ant-descriptions-bordered .ant-descriptions-item-label),
:deep(.ant-descriptions-bordered .ant-descriptions-item-content) {
  border-right: 1px solid var(--border-glass) !important;
}

@media (max-width: 575px) {
  :deep(.ant-descriptions-bordered .ant-descriptions-row) {
    display: flex;
    flex-direction: column;
  }

  :deep(.ant-descriptions-bordered .ant-descriptions-item-label),
  :deep(.ant-descriptions-bordered .ant-descriptions-item-content) {
    display: block;
    width: 100%;
    padding: 8px 12px;
    border-right: none !important;
    border-bottom: 1px solid var(--border-glass);
  }

  :deep(.ant-descriptions-bordered .ant-descriptions-item-label) {
    background: rgba(255, 255, 255, 0.15);
    font-weight: 600;
  }
}
</style>

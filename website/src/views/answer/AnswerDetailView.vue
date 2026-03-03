<template>
  <div class="page-container">
    <div class="page-header">
      <div class="header-top">
        <h1 class="page-title">答卷详情</h1>
        <a-button @click="handleBack">
          <template #icon><ArrowLeftOutlined /></template>
          返回
        </a-button>
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
          <a-descriptions-item v-if="answer.examScore" label="得分">
            {{ answer.examScore }}
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
import { ArrowLeftOutlined, PaperClipOutlined } from '@ant-design/icons-vue'
import type { AnswerView } from '@/types/answer'
import type { SurveySchema, SurveyElement, AttachmentInfo } from '@/types/survey'

const route = useRoute()
const router = useRouter()

const answerId = computed(() => route.params.id as string)

const loading = ref(false)
const answer = ref<AnswerView | null>(null)
const projectName = ref('')
const elementMap = ref<Map<string, SurveyElement>>(new Map())

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
  if (typeof value === 'object') return JSON.stringify(value, null, 2)
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
</script>

<style scoped>
.page-header {
  margin-bottom: 16px;
}

.header-top {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.answer-content {
  background: #fafafa;
  border-radius: 4px;
  padding: 16px;
}

.answer-item {
  display: flex;
  padding: 8px 0;
  border-bottom: 1px solid #e8e8e8;
}

.answer-item:last-child {
  border-bottom: none;
}

.answer-label {
  width: 200px;
  font-weight: 500;
  color: #666;
}

.answer-value {
  flex: 1;
}

.attachment-list {
  margin-top: 8px;
  padding: 8px 12px;
  background: #f5f5f5;
  border-radius: 4px;
}

.attachment-label {
  font-size: 12px;
  color: #666;
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
  color: #1890ff;
  text-decoration: none;
  font-size: 13px;
}

.attachment-link:hover {
  text-decoration: underline;
}
</style>

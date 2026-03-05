<template>
  <div class="page-container">
    <div class="page-header">
      <h1 class="page-title">回收站</h1>
    </div>

    <a-card :bordered="false">
      <a-tabs v-model:activeKey="activeTab" @change="handleTabChange">
        <a-tab-pane key="projects" tab="问卷项目">
          <div class="table-responsive">
            <a-table
              :columns="projectColumns"
              :data-source="projects"
              :loading="projectLoading"
              :pagination="projectPagination"
              :scroll="{ x: 500 }"
              row-key="id"
              @change="handleProjectTableChange"
            >
            <template #bodyCell="{ column, record }">
              <template v-if="column.key === 'createAt'">
                {{ formatDate(record.createAt) }}
              </template>
              <template v-if="column.key === 'status'">
                <a-tag :color="getStatusColor(record.status)">
                  {{ getStatusText(record.status) }}
                </a-tag>
              </template>
              <template v-if="column.key === 'actions'">
                <a-space>
                  <a-button type="link" size="small" @click="handleRestoreProject(record)">
                    恢复
                  </a-button>
                  <a-popconfirm
                    title="确定要永久删除此问卷吗？此操作不可恢复。"
                    @confirm="handleDestroyProject(record)"
                  >
                    <a-button type="link" size="small" danger>
                      永久删除
                    </a-button>
                  </a-popconfirm>
                </a-space>
              </template>
            </template>
          </a-table>
          </div>
        </a-tab-pane>

        <a-tab-pane key="answers" tab="答卷">
          <div class="table-responsive">
            <a-table
              :columns="answerColumns"
              :data-source="answers"
              :loading="answerLoading"
              :pagination="answerPagination"
              :scroll="{ x: 500 }"
              row-key="id"
              @change="handleAnswerTableChange"
            >
            <template #bodyCell="{ column, record }">
              <template v-if="column.key === 'createAt'">
                {{ formatDate(record.createAt) }}
              </template>
              <template v-if="column.key === 'tempSave'">
                <a-tag :color="record.tempSave ? 'orange' : 'green'">
                  {{ record.tempSave ? '草稿' : '已提交' }}
                </a-tag>
              </template>
              <template v-if="column.key === 'actions'">
                <a-space>
                  <a-button type="link" size="small" @click="handleRestoreAnswer(record)">
                    恢复
                  </a-button>
                  <a-popconfirm
                    title="确定要永久删除此答卷吗？此操作不可恢复。"
                    @confirm="handleDestroyAnswer(record)"
                  >
                    <a-button type="link" size="small" danger>
                      永久删除
                    </a-button>
                  </a-popconfirm>
                </a-space>
              </template>
            </template>
          </a-table>
          </div>
        </a-tab-pane>
      </a-tabs>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { message } from 'ant-design-vue'
import { projectApi } from '@/api/project'
import { answerApi } from '@/api/answer'
import type { ProjectView } from '@/types/project'
import type { AnswerView } from '@/types/answer'

const activeTab = ref('projects')

const projectLoading = ref(false)
const projects = ref<ProjectView[]>([])
const projectPagination = reactive({
  current: 1,
  pageSize: 10,
  total: 0,
  showSizeChanger: true,
  showTotal: (total: number) => `共 ${total} 条`
})

const answerLoading = ref(false)
const answers = ref<AnswerView[]>([])
const answerPagination = reactive({
  current: 1,
  pageSize: 10,
  total: 0,
  showSizeChanger: true,
  showTotal: (total: number) => `共 ${total} 条`
})

const projectColumns = [
  { title: '问卷名称', dataIndex: 'name', ellipsis: true },
  { title: '状态', key: 'status', width: 100 },
  { title: '删除时间', key: 'createAt', width: 180 },
  { title: '操作', key: 'actions', width: 150 }
]

const answerColumns = [
  { title: 'ID', dataIndex: 'id', width: 200, ellipsis: true },
  { title: '状态', key: 'tempSave', width: 100 },
  { title: '删除时间', key: 'createAt', width: 180 },
  { title: '操作', key: 'actions', width: 150 }
]

onMounted(async () => {
  await fetchProjects()
})

async function fetchProjects() {
  projectLoading.value = true
  try {
    const result = await projectApi.getTrash()
    projects.value = result
    projectPagination.total = result.length
  } finally {
    projectLoading.value = false
  }
}

async function fetchAnswers() {
  answerLoading.value = true
  try {
    const result = await answerApi.getTrash({
      pageIndex: answerPagination.current,
      pageSize: answerPagination.pageSize
    })
    answers.value = result.list
    answerPagination.total = result.total
  } finally {
    answerLoading.value = false
  }
}

function handleTabChange(key: string) {
  if (key === 'answers' && answers.value.length === 0) {
    fetchAnswers()
  }
}

function handleProjectTableChange(pag: { current: number; pageSize: number }) {
  projectPagination.current = pag.current
  projectPagination.pageSize = pag.pageSize
}

function handleAnswerTableChange(pag: { current: number; pageSize: number }) {
  answerPagination.current = pag.current
  answerPagination.pageSize = pag.pageSize
  fetchAnswers()
}

async function handleRestoreProject(record: ProjectView) {
  try {
    await projectApi.restore(record.id)
    message.success('恢复成功')
    fetchProjects()
  } catch {
    message.error('恢复失败')
  }
}

async function handleDestroyProject(record: ProjectView) {
  try {
    await projectApi.destroy(record.id)
    message.success('永久删除成功')
    fetchProjects()
  } catch {
    message.error('永久删除失败')
  }
}

async function handleRestoreAnswer(record: AnswerView) {
  try {
    await answerApi.restore(record.id)
    message.success('恢复成功')
    fetchAnswers()
  } catch {
    message.error('恢复失败')
  }
}

async function handleDestroyAnswer(record: AnswerView) {
  try {
    await answerApi.destroy(record.id)
    message.success('永久删除成功')
    fetchAnswers()
  } catch {
    message.error('永久删除失败')
  }
}

function getStatusColor(status: number): string {
  switch (status) {
    case 0:
      return 'default'
    case 1:
      return 'processing'
    case 2:
      return 'success'
    default:
      return 'default'
  }
}

function getStatusText(status: number): string {
  switch (status) {
    case 0:
      return '草稿'
    case 1:
      return '进行中'
    case 2:
      return '已结束'
    default:
      return '未知'
  }
}

function formatDate(dateStr: string) {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  return date.toLocaleString()
}
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

.table-responsive {
  overflow-x: auto;
  -webkit-overflow-scrolling: touch;
}

.table-responsive :deep(.ant-table) {
  min-width: 500px;
}

@media (max-width: 767px) {
  .table-responsive :deep(.ant-table-thead > tr > th) {
    padding: 12px 8px;
    font-size: 13px;
    white-space: nowrap;
  }

  .table-responsive :deep(.ant-table-tbody > tr > td) {
    padding: 12px 8px;
    font-size: 13px;
  }

  .table-responsive :deep(.ant-space) {
    flex-wrap: wrap;
    gap: 4px !important;
  }
}
</style>

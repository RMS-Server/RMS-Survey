<template>
  <div class="page-container">
    <div class="page-header">
      <h1 class="page-title">答题管理</h1>
    </div>

    <a-card :bordered="false">
      <div class="filter-bar">
        <a-select
          v-model:value="filterProjectId"
          placeholder="选择问卷"
          style="width: 200px"
          allowClear
          show-search
          :filter-option="filterOption"
          @change="handleSearch"
        >
          <a-select-option
            v-for="project in projects"
            :key="project.id"
            :value="project.id"
          >
            {{ project.name }}
          </a-select-option>
        </a-select>
        <a-checkbox v-model:checked="showDeleted" @change="handleSearch">
          显示已删除
        </a-checkbox>
        <a-button
          type="primary"
          :disabled="!filterProjectId"
          :loading="exporting"
          @click="handleExport"
        >
          <template #icon><DownloadOutlined /></template>
          导出Excel
        </a-button>
      </div>

      <a-table
        :columns="columns"
        :data-source="answers"
        :loading="loading"
        :pagination="pagination"
        row-key="id"
        @change="handleTableChange"
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
              <a-button type="link" size="small" @click="handleView(record)">
                查看
              </a-button>
              <a-popconfirm
                v-if="!record.deleted"
                title="确定要删除此答卷吗？"
                @confirm="handleDelete(record)"
              >
                <a-button type="link" size="small" danger>
                  删除
                </a-button>
              </a-popconfirm>
              <a-button
                v-else
                type="link"
                size="small"
                @click="handleRestore(record)"
              >
                恢复
              </a-button>
            </a-space>
          </template>
        </template>
      </a-table>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { message } from 'ant-design-vue'
import { answerApi } from '@/api/answer'
import { projectApi } from '@/api/project'
import { DownloadOutlined } from '@ant-design/icons-vue'
import type { AnswerView } from '@/types/answer'
import type { ProjectView } from '@/types/project'

const router = useRouter()

const loading = ref(false)
const exporting = ref(false)
const showDeleted = ref(false)
const filterProjectId = ref<string | undefined>()
const projects = ref<ProjectView[]>([])
const answers = ref<AnswerView[]>([])

const pagination = reactive({
  current: 1,
  pageSize: 10,
  total: 0,
  showSizeChanger: true,
  showTotal: (total: number) => `共 ${total} 条`
})

const columns = [
  { title: 'ID', dataIndex: 'id', width: 200, ellipsis: true },
  { title: '状态', key: 'tempSave', width: 100 },
  { title: '提交时间', key: 'createAt', width: 180 },
  { title: '操作', key: 'actions', width: 150 }
]

onMounted(async () => {
  await fetchProjects()
  await fetchAnswers()
})

async function fetchProjects() {
  try {
    const result = await projectApi.list({ pageSize: 1000 })
    projects.value = result.list
  } catch {
    // Ignore
  }
}

async function fetchAnswers() {
  loading.value = true
  try {
    const result = await answerApi.list({
      pageIndex: pagination.current,
      pageSize: pagination.pageSize,
      projectId: filterProjectId.value,
      deleted: showDeleted.value
    })
    answers.value = result.list
    pagination.total = result.total
  } finally {
    loading.value = false
  }
}

function handleSearch() {
  pagination.current = 1
  fetchAnswers()
}

function handleTableChange(pag: { current: number; pageSize: number }) {
  pagination.current = pag.current
  pagination.pageSize = pag.pageSize
  fetchAnswers()
}

function filterOption(input: string, option: { children: string }) {
  return option.children.toLowerCase().indexOf(input.toLowerCase()) >= 0
}

function handleView(record: AnswerView) {
  router.push(`/answer/${record.id}?project=${record.projectId}`)
}

async function handleDelete(record: AnswerView) {
  try {
    await answerApi.delete(record.id)
    message.success('删除成功')
    fetchAnswers()
  } catch {
    message.error('删除失败')
  }
}

async function handleRestore(record: AnswerView) {
  try {
    await answerApi.restore(record.id)
    message.success('恢复成功')
    fetchAnswers()
  } catch {
    message.error('恢复失败')
  }
}

async function handleExport() {
  if (!filterProjectId.value) return

  exporting.value = true
  try {
    const blob = await answerApi.download({ projectId: filterProjectId.value })
    const url = window.URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = `answers_${filterProjectId.value}.xlsx`
    link.click()
    window.URL.revokeObjectURL(url)
    message.success('导出成功')
  } catch {
    message.error('导出失败')
  } finally {
    exporting.value = false
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

.filter-bar {
  display: flex;
  gap: 16px;
  margin-bottom: 16px;
}
</style>

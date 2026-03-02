<template>
  <div class="page-container">
    <div class="page-header">
      <div class="header-top">
        <h1 class="page-title">问卷管理</h1>
        <a-button type="primary" @click="handleCreate">
          <template #icon><PlusOutlined /></template>
          新建
        </a-button>
      </div>
    </div>

    <a-card :bordered="false">
      <div class="filter-bar">
        <a-input-search
          v-model:value="searchName"
          placeholder="按名称搜索"
          style="width: 200px"
          @search="handleSearch"
        />
        <a-select
          v-model:value="filterStatus"
          placeholder="状态"
          style="width: 120px"
          allowClear
          @change="handleSearch"
        >
          <a-select-option :value="0">草稿</a-select-option>
          <a-select-option :value="1">已发布</a-select-option>
        </a-select>
        <a-checkbox v-model:checked="showDeleted" @change="handleSearch">
          显示已删除
        </a-checkbox>
      </div>

      <a-table
        :columns="columns"
        :data-source="projects"
        :loading="loading"
        :pagination="pagination"
        row-key="id"
        @change="handleTableChange"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'name'">
            <a @click="handleEdit(record)">{{ record.name }}</a>
          </template>
          <template v-if="column.key === 'status'">
            <a-tag :color="record.status === 1 ? 'green' : 'default'">
              {{ record.status === 1 ? '已发布' : '草稿' }}
            </a-tag>
          </template>
          <template v-if="column.key === 'mode'">
            <a-tag>{{ record.mode || 'survey' }}</a-tag>
          </template>
          <template v-if="column.key === 'createAt'">
            {{ formatDate(record.createAt) }}
          </template>
          <template v-if="column.key === 'actions'">
            <a-space>
              <a-button type="link" size="small" @click="handleEdit(record)">
                编辑
              </a-button>
              <a-button type="link" size="small" @click="handleCopyLink(record)">
                复制链接
              </a-button>
              <a-popconfirm
                v-if="!record.deleted"
                title="确定要删除此问卷吗？"
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
import { useProjectStore } from '@/stores/project'
import { PlusOutlined } from '@ant-design/icons-vue'
import type { ProjectView } from '@/types/project'

const router = useRouter()
const projectStore = useProjectStore()

const searchName = ref('')
const filterStatus = ref<number | undefined>()
const showDeleted = ref(false)
const loading = ref(false)
const projects = ref<ProjectView[]>([])

const pagination = reactive({
  current: 1,
  pageSize: 10,
  total: 0,
  showSizeChanger: true,
  showTotal: (total: number) => `共 ${total} 条`
})

const columns = [
  { title: '名称', key: 'name', dataIndex: 'name' },
  { title: '状态', key: 'status', width: 100 },
  { title: '模式', key: 'mode', width: 100 },
  { title: '创建时间', key: 'createAt', width: 180 },
  { title: '操作', key: 'actions', width: 200 }
]

onMounted(() => {
  fetchProjects()
})

async function fetchProjects() {
  loading.value = true
  try {
    const result = await projectStore.fetchProjects({
      pageIndex: pagination.current,
      pageSize: pagination.pageSize,
      name: searchName.value || undefined,
      status: filterStatus.value,
      deleted: showDeleted.value
    })
    projects.value = result.list
    pagination.total = result.total
  } finally {
    loading.value = false
  }
}

function handleSearch() {
  pagination.current = 1
  fetchProjects()
}

function handleTableChange(pag: { current: number; pageSize: number }) {
  pagination.current = pag.current
  pagination.pageSize = pag.pageSize
  fetchProjects()
}

function handleCreate() {
  router.push('/project/create')
}

function handleEdit(record: ProjectView) {
  router.push(`/project/${record.id}/edit`)
}

async function handleDelete(record: ProjectView) {
  try {
    await projectStore.deleteProject(record.id)
    message.success('删除成功')
    fetchProjects()
  } catch {
    message.error('删除失败')
  }
}

async function handleRestore(record: ProjectView) {
  try {
    await projectStore.restoreProject(record.id)
    message.success('恢复成功')
    fetchProjects()
  } catch {
    message.error('恢复失败')
  }
}

function handleCopyLink(record: ProjectView) {
  const link = `${window.location.origin}/s/${record.id}`
  navigator.clipboard.writeText(link)
  message.success('链接已复制到剪贴板')
}

function formatDate(dateStr: string) {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  return date.toLocaleString()
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

.filter-bar {
  display: flex;
  gap: 16px;
  margin-bottom: 16px;
}
</style>

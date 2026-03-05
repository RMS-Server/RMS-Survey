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
          class="filter-search"
          @search="handleSearch"
        />
        <a-select
          v-model:value="filterStatus"
          placeholder="状态"
          class="filter-select"
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

      <div class="table-responsive">
        <a-table
          :columns="columns"
          :data-source="projects"
          :loading="loading"
          :pagination="pagination"
          :scroll="{ x: 600 }"
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
                <a-button v-if="record.isOwner" type="link" size="small" @click="handlePartners(record)">
                  参与者
                </a-button>
                <a-button type="link" size="small" @click="handleCopyLink(record)">
                  复制链接
                </a-button>
                <a-popconfirm
                  v-if="!record.deleted && record.isOwner"
                  title="确定要删除此问卷吗？"
                  @confirm="handleDelete(record)"
                >
                  <a-button type="link" size="small" danger>
                    删除
                  </a-button>
                </a-popconfirm>
                <a-button
                  v-else-if="record.deleted && record.isOwner"
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
      </div>
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
  { title: '操作', key: 'actions', width: 260 }
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

function handlePartners(record: ProjectView) {
  router.push(`/project/${record.id}/partners`)
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

.filter-bar {
  display: flex;
  gap: 12px;
  margin-bottom: 16px;
  flex-wrap: wrap;
  align-items: center;
}

.filter-search {
  width: 100%;
  max-width: 200px;
}

.filter-select {
  width: 120px;
}

@media (max-width: 575px) {
  .filter-search {
    max-width: 100%;
  }

  .filter-select {
    flex: 1;
    min-width: 100px;
  }
}

.table-responsive {
  overflow-x: auto;
  -webkit-overflow-scrolling: touch;
}

.table-responsive :deep(.ant-table) {
  min-width: 600px;
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

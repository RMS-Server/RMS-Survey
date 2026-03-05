<template>
  <div class="partner-manager">
    <div class="partner-header">
      <h3>参与者管理</h3>
      <div class="header-actions">
        <a-input-search
          v-model:value="searchName"
          placeholder="搜索参与者"
          style="width: 200px"
          @search="loadPartners"
          allow-clear
        />
        <a-button v-if="!readOnly" type="primary" @click="showAddModal = true">
          <template #icon><PlusOutlined /></template>
          添加参与者
        </a-button>
        <a-dropdown v-if="!readOnly">
          <a-button>
            <template #icon><MoreOutlined /></template>
            更多
          </a-button>
          <template #overlay>
            <a-menu>
              <a-menu-item @click="handleImport">
                <UploadOutlined /> 导入名单
              </a-menu-item>
              <a-menu-item @click="handleExport">
                <DownloadOutlined /> 导出名单
              </a-menu-item>
              <a-menu-divider />
              <a-menu-item @click="showBatchDelete = true" :disabled="selectedRowKeys.length === 0">
                <DeleteOutlined /> 批量删除 ({{ selectedRowKeys.length }})
              </a-menu-item>
            </a-menu>
          </template>
        </a-dropdown>
        <a-button v-else @click="handleExport">
          <template #icon><DownloadOutlined /></template>
          导出名单
        </a-button>
      </div>
    </div>

    <a-table
      :columns="columns"
      :data-source="partners"
      :loading="loading"
      :pagination="pagination"
      :row-selection="readOnly ? undefined : { selectedRowKeys, onChange: onSelectionChange }"
      row-key="id"
      @change="handleTableChange"
    >
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'type'">
          <a-tag :color="record.type === 1 ? 'blue' : 'green'">
            {{ record.type === 1 ? '系统用户' : '导入用户' }}
          </a-tag>
        </template>
        <template v-else-if="column.key === 'status'">
          <a-tag :color="statusColors[record.status] || 'default'">
            {{ statusTexts[record.status] || '未知' }}
          </a-tag>
        </template>
        <template v-else-if="column.key === 'action'">
          <a-popconfirm
            title="确定要删除此参与者吗？"
            @confirm="handleDelete(record.id)"
          >
            <a-button type="link" danger size="small">删除</a-button>
          </a-popconfirm>
        </template>
      </template>
    </a-table>

    <!-- Add Partner Modal -->
    <a-modal
      v-model:open="showAddModal"
      title="添加参与者"
      :confirm-loading="addLoading"
      @ok="handleAdd"
      @cancel="resetAddForm"
    >
      <a-form :model="addForm" layout="vertical">
        <a-form-item label="添加方式" required>
          <a-radio-group v-model:value="addForm.mode">
            <a-radio value="user">从系统用户选择</a-radio>
            <a-radio value="manual">手动输入</a-radio>
          </a-radio-group>
        </a-form-item>

        <template v-if="addForm.mode === 'user'">
          <a-form-item label="选择用户" required>
            <a-select
              v-model:value="addForm.userId"
              show-search
              placeholder="输入姓名搜索用户"
              :filter-option="false"
              :loading="userSearchLoading"
              @search="handleUserSearch"
              allow-clear
            >
              <a-select-option
                v-for="user in userOptions"
                :key="user.id"
                :value="user.id"
              >
                {{ user.name }} ({{ user.email || user.phone || user.id }})
              </a-select-option>
            </a-select>
          </a-form-item>
        </template>

        <template v-else>
          <a-form-item label="参与者姓名" required>
            <a-input v-model:value="addForm.userName" placeholder="输入参与者姓名" />
          </a-form-item>
        </template>
      </a-form>
    </a-modal>

    <!-- Import Modal -->
    <a-modal
      v-model:open="showImportModal"
      title="导入参与者名单"
      :confirm-loading="importLoading"
      @ok="confirmImport"
      @cancel="resetImport"
    >
      <a-alert
        message="导入说明"
        description="上传 Excel 文件，第一行为表头，第一列为参与者姓名。系统将自动匹配系统用户，未匹配的将作为导入用户添加。"
        type="info"
        show-icon
        style="margin-bottom: 16px"
      />
      <a-upload-dragger
        v-model:file-list="importFile"
        :before-upload="() => false"
        accept=".xlsx,.xls"
        :max-count="1"
      >
        <p class="ant-upload-drag-icon">
          <InboxOutlined />
        </p>
        <p class="ant-upload-text">点击或拖拽文件到此区域</p>
        <p class="ant-upload-hint">支持 .xlsx, .xls 格式</p>
      </a-upload-dragger>
    </a-modal>

    <!-- Batch Delete Confirm -->
    <a-modal
      v-model:open="showBatchDelete"
      title="批量删除"
      :confirm-loading="batchDeleteLoading"
      @ok="handleBatchDelete"
      @cancel="showBatchDelete = false"
    >
      <p>确定要删除选中的 {{ selectedRowKeys.length }} 个参与者吗？</p>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, onUnmounted } from 'vue'
import { message } from 'ant-design-vue'
import type { UploadFile } from 'ant-design-vue'
import {
  PlusOutlined,
  MoreOutlined,
  UploadOutlined,
  DownloadOutlined,
  DeleteOutlined,
  InboxOutlined
} from '@ant-design/icons-vue'
import { projectApi } from '@/api/project'
import type { ProjectPartnerView, SelectUserView } from '@/types/project'

const props = defineProps<{
  projectId: string
  readOnly?: boolean
}>()

const loading = ref(false)
const partners = ref<ProjectPartnerView[]>([])
const selectedRowKeys = ref<string[]>([])
const searchName = ref('')

const pagination = reactive({
  current: 1,
  pageSize: 10,
  total: 0,
  showSizeChanger: true,
  showTotal: (total: number) => `共 ${total} 条`
})

const columns = computed(() => {
  const baseColumns: Array<{ title: string; dataIndex?: string; key: string; width?: number }> = [
    { title: '姓名', dataIndex: 'userName', key: 'userName' },
    { title: '类型', dataIndex: 'type', key: 'type' },
    { title: '状态', dataIndex: 'status', key: 'status' }
  ]
  if (!props.readOnly) {
    baseColumns.push({ title: '操作', key: 'action', width: 100 })
  }
  return baseColumns
})

const statusColors: Record<number, string> = {
  0: 'default',
  1: 'processing',
  2: 'success',
  3: 'error'
}

const statusTexts: Record<number, string> = {
  0: '未答题',
  1: '已答题',
  2: '已通过',
  3: '未通过'
}

// Add modal state
const showAddModal = ref(false)
const addLoading = ref(false)
const addForm = reactive({
  mode: 'user' as 'user' | 'manual',
  userId: '',
  userName: ''
})
const userSearchLoading = ref(false)
const userOptions = ref<SelectUserView[]>([])

// Import state
const showImportModal = ref(false)
const importLoading = ref(false)
const importFile = ref<UploadFile[]>([])

// Batch delete state
const showBatchDelete = ref(false)
const batchDeleteLoading = ref(false)

// Debounce timer
let userSearchTimer: ReturnType<typeof setTimeout> | null = null

async function loadPartners() {
  loading.value = true
  try {
    const res = await projectApi.listPartners({
      projectId: props.projectId,
      pageIndex: pagination.current,
      pageSize: pagination.pageSize,
      userName: searchName.value || undefined
    })
    partners.value = res.list
    pagination.total = res.total
  } catch {
    message.error('加载参与者列表失败')
  } finally {
    loading.value = false
  }
}

function handleTableChange(pag: { current: number; pageSize: number }) {
  pagination.current = pag.current
  pagination.pageSize = pag.pageSize
  loadPartners()
}

function onSelectionChange(keys: string[]) {
  selectedRowKeys.value = keys
}

async function handleUserSearch(value: string) {
  if (userSearchTimer) clearTimeout(userSearchTimer)
  if (!value.trim()) {
    userOptions.value = []
    return
  }
  userSearchTimer = setTimeout(async () => {
    userSearchLoading.value = true
    try {
      userOptions.value = await projectApi.selectUser(value)
    } catch {
      userOptions.value = []
    } finally {
      userSearchLoading.value = false
    }
  }, 300)
}

async function handleAdd() {
  if (addForm.mode === 'user' && !addForm.userId) {
    message.warning('请选择用户')
    return
  }
  if (addForm.mode === 'manual' && !addForm.userName.trim()) {
    message.warning('请输入参与者姓名')
    return
  }

  addLoading.value = true
  try {
    if (addForm.mode === 'user') {
      const selectedUser = userOptions.value.find(u => u.id === addForm.userId)
      await projectApi.addPartner({
        projectId: props.projectId,
        userId: addForm.userId,
        userName: selectedUser?.name,
        type: 1
      })
    } else {
      await projectApi.addPartner({
        projectId: props.projectId,
        userName: addForm.userName,
        type: 2
      })
    }
    message.success('添加成功')
    showAddModal.value = false
    resetAddForm()
    loadPartners()
  } catch {
    message.error('添加失败')
  } finally {
    addLoading.value = false
  }
}

function resetAddForm() {
  addForm.mode = 'user'
  addForm.userId = ''
  addForm.userName = ''
  userOptions.value = []
}

async function handleDelete(id: string) {
  try {
    await projectApi.removePartner({ id, projectId: props.projectId })
    message.success('删除成功')
    loadPartners()
  } catch {
    message.error('删除失败')
  }
}

function handleImport() {
  importFile.value = []
  showImportModal.value = true
}

async function confirmImport() {
  if (importFile.value.length === 0) {
    message.warning('请选择文件')
    return
  }
  const file = importFile.value[0].originFileObj
  if (!file) {
    message.warning('请选择有效的文件')
    return
  }
  importLoading.value = true
  try {
    await projectApi.importPartners(props.projectId, file)
    message.success('导入成功')
    showImportModal.value = false
    loadPartners()
  } catch {
    message.error('导入失败')
  } finally {
    importLoading.value = false
  }
}

function resetImport() {
  importFile.value = []
}

async function handleExport() {
  try {
    const blob = await projectApi.downloadPartners(props.projectId)
    const url = window.URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `participants_${props.projectId}.xlsx`
    a.click()
    window.URL.revokeObjectURL(url)
  } catch {
    message.error('导出失败')
  }
}

async function handleBatchDelete() {
  batchDeleteLoading.value = true
  let successCount = 0
  let failCount = 0

  for (const id of selectedRowKeys.value) {
    try {
      await projectApi.removePartner({ id, projectId: props.projectId })
      successCount++
    } catch {
      failCount++
    }
  }

  if (failCount === 0) {
    message.success(`成功删除 ${successCount} 个参与者`)
  } else if (successCount === 0) {
    message.error('批量删除失败')
  } else {
    message.warning(`成功删除 ${successCount} 个，失败 ${failCount} 个`)
  }

  showBatchDelete.value = false
  selectedRowKeys.value = []
  loadPartners()
  batchDeleteLoading.value = false
}

onMounted(() => {
  loadPartners()
})

onUnmounted(() => {
  if (userSearchTimer) {
    clearTimeout(userSearchTimer)
  }
})

defineExpose({
  loadPartners
})
</script>

<style scoped>
.partner-manager {
  padding: 16px;
}

.partner-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.partner-header h3 {
  margin: 0;
}

.header-actions {
  display: flex;
  gap: 8px;
}
</style>

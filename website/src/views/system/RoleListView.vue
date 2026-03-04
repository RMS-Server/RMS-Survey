<template>
  <div class="page-container">
    <div class="page-header">
      <h1 class="page-title">角色管理</h1>
      <a-button type="primary" @click="handleCreate">
        <template #icon><PlusOutlined /></template>
        添加角色
      </a-button>
    </div>

    <a-card :bordered="false">
      <a-table
        :columns="columns"
        :data-source="roles"
        :loading="loading"
        :pagination="pagination"
        row-key="id"
        @change="handleTableChange"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'status'">
            <a-tag :color="record.status ? 'green' : 'red'">
              {{ record.status ? '启用' : '禁用' }}
            </a-tag>
          </template>
          <template v-if="column.key === 'actions'">
            <a-space>
              <a-button type="link" size="small" @click="handleEdit(record)">
                编辑
              </a-button>
              <a-popconfirm
                title="确定要删除此角色吗？"
                @confirm="handleDelete(record)"
              >
                <a-button type="link" size="small" danger>
                  删除
                </a-button>
              </a-popconfirm>
            </a-space>
          </template>
        </template>
      </a-table>
    </a-card>

    <a-modal
      v-model:open="modalVisible"
      :title="isEdit ? '编辑角色' : '添加角色'"
      @ok="handleSubmit"
      @cancel="handleCancel"
    >
      <a-form :model="formState" :rules="rules" layout="vertical" ref="formRef">
        <a-form-item name="name" label="名称">
          <a-input v-model:value="formState.name" />
        </a-form-item>
        <a-form-item name="code" label="编码">
          <a-input v-model:value="formState.code" :disabled="isEdit" />
        </a-form-item>
        <a-form-item name="remark" label="备注">
          <a-textarea v-model:value="formState.remark" :rows="3" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { message } from 'ant-design-vue'
import { roleApi } from '@/api/system'
import { PlusOutlined } from '@ant-design/icons-vue'
import type { RoleView } from '@/types/user'

const loading = ref(false)
const roles = ref<RoleView[]>([])
const modalVisible = ref(false)
const isEdit = ref(false)
const formRef = ref()

const formState = reactive({
  id: '',
  name: '',
  code: '',
  remark: ''
})

const rules = {
  name: [{ required: true, message: '请输入名称' }],
  code: [{ required: true, message: '请输入编码' }]
}

const pagination = reactive({
  current: 1,
  pageSize: 10,
  total: 0,
  showSizeChanger: true,
  showTotal: (total: number) => `共 ${total} 条`
})

const columns = [
  { title: '名称', dataIndex: 'name' },
  { title: '编码', dataIndex: 'code' },
  { title: '操作', key: 'actions', width: 150 }
]

onMounted(() => {
  fetchRoles()
})

async function fetchRoles() {
  loading.value = true
  try {
    const result = await roleApi.list({
      pageIndex: pagination.current,
      pageSize: pagination.pageSize
    })
    roles.value = result.list
    pagination.total = result.total
  } finally {
    loading.value = false
  }
}

function handleTableChange(pag: { current: number; pageSize: number }) {
  pagination.current = pag.current
  pagination.pageSize = pag.pageSize
  fetchRoles()
}

function handleCreate() {
  isEdit.value = false
  Object.assign(formState, {
    id: '',
    name: '',
    code: '',
    remark: ''
  })
  modalVisible.value = true
}

function handleEdit(record: RoleView) {
  isEdit.value = true
  Object.assign(formState, {
    id: record.id,
    name: record.name,
    code: record.code,
    remark: ''
  })
  modalVisible.value = true
}

async function handleSubmit() {
  try {
    await formRef.value?.validate()
    if (isEdit.value) {
      await roleApi.update({
        id: formState.id,
        name: formState.name,
        remark: formState.remark
      })
      message.success('更新成功')
    } else {
      await roleApi.create({
        name: formState.name,
        code: formState.code,
        remark: formState.remark
      })
      message.success('创建成功')
    }
    modalVisible.value = false
    fetchRoles()
  } catch {
    // Form validation failed
  }
}

function handleCancel() {
  modalVisible.value = false
}

async function handleDelete(record: RoleView) {
  try {
    await roleApi.delete(record.id)
    message.success('删除成功')
    fetchRoles()
  } catch {
    message.error('删除失败')
  }
}
</script>

<style scoped>
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  background: var(--surface-glass);
  backdrop-filter: blur(var(--blur-strength));
  -webkit-backdrop-filter: blur(var(--blur-strength));
  padding: 16px 24px;
  margin-bottom: 16px;
  border-radius: var(--radius-md);
  border: 1px solid var(--border-glass);
  box-shadow: var(--shadow-raised);
}
</style>

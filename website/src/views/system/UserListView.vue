<template>
  <div class="page-container">
    <div class="page-header">
      <h1 class="page-title">用户管理</h1>
      <a-button type="primary" @click="handleCreate">
        <template #icon><PlusOutlined /></template>
        添加用户
      </a-button>
    </div>

    <a-card :bordered="false">
      <div class="filter-bar">
        <a-input-search
          v-model:value="searchName"
          placeholder="按名称搜索"
          style="width: 200px"
          @search="handleSearch"
        />
      </div>

      <a-table
        :columns="columns"
        :data-source="users"
        :loading="loading"
        :pagination="pagination"
        row-key="id"
        @change="handleTableChange"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'status'">
            <a-tag :color="record.status === 1 ? 'green' : 'red'">
              {{ record.status === 1 ? '启用' : '禁用' }}
            </a-tag>
          </template>
          <template v-if="column.key === 'roles'">
            <a-tag v-for="role in record.roles" :key="role.id" color="blue">
              {{ role.name }}
            </a-tag>
          </template>
          <template v-if="column.key === 'actions'">
            <a-space>
              <a-button type="link" size="small" @click="handleEdit(record)">
                编辑
              </a-button>
              <a-popconfirm
                title="确定要删除此用户吗？"
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
      :title="isEdit ? '编辑用户' : '添加用户'"
      @ok="handleSubmit"
      @cancel="handleCancel"
    >
      <a-form :model="formState" :rules="rules" layout="vertical" ref="formRef">
        <a-form-item name="username" label="用户名">
          <a-input v-model:value="formState.username" :disabled="isEdit" />
        </a-form-item>
        <a-form-item name="name" label="姓名">
          <a-input v-model:value="formState.name" />
        </a-form-item>
        <a-form-item v-if="!isEdit" name="password" label="密码">
          <a-input-password v-model:value="formState.password" />
        </a-form-item>
        <a-form-item name="email" label="邮箱">
          <a-input v-model:value="formState.email" />
        </a-form-item>
        <a-form-item name="phone" label="手机号">
          <a-input v-model:value="formState.phone" />
        </a-form-item>
        <a-form-item name="status" label="状态">
          <a-select v-model:value="formState.status">
            <a-select-option :value="1">启用</a-select-option>
            <a-select-option :value="0">禁用</a-select-option>
          </a-select>
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { message } from 'ant-design-vue'
import { userApi } from '@/api/auth'
import { PlusOutlined } from '@ant-design/icons-vue'
import type { UserView, CreateUserRequest } from '@/types/user'

const loading = ref(false)
const searchName = ref('')
const users = ref<UserView[]>([])
const modalVisible = ref(false)
const isEdit = ref(false)
const formRef = ref()

const formState = reactive({
  id: '',
  username: '',
  name: '',
  password: '',
  email: '',
  phone: '',
  status: 1
})

const rules = {
  username: [{ required: true, message: '请输入用户名' }],
  name: [{ required: true, message: '请输入姓名' }],
  password: [{ required: true, message: '请输入密码' }]
}

const pagination = reactive({
  current: 1,
  pageSize: 10,
  total: 0,
  showSizeChanger: true,
  showTotal: (total: number) => `共 ${total} 条`
})

const columns = [
  { title: '用户名', dataIndex: 'username' },
  { title: '姓名', dataIndex: 'name' },
  { title: '邮箱', dataIndex: 'email' },
  { title: '手机号', dataIndex: 'phone' },
  { title: '状态', key: 'status', width: 100 },
  { title: '角色', key: 'roles', width: 150 },
  { title: '操作', key: 'actions', width: 150 }
]

onMounted(() => {
  fetchUsers()
})

async function fetchUsers() {
  loading.value = true
  try {
    const result = await userApi.listUsers({
      pageIndex: pagination.current,
      pageSize: pagination.pageSize,
      name: searchName.value || undefined
    })
    users.value = result.list
    pagination.total = result.total
  } finally {
    loading.value = false
  }
}

function handleSearch() {
  pagination.current = 1
  fetchUsers()
}

function handleTableChange(pag: { current: number; pageSize: number }) {
  pagination.current = pag.current
  pagination.pageSize = pag.pageSize
  fetchUsers()
}

function handleCreate() {
  isEdit.value = false
  Object.assign(formState, {
    id: '',
    username: '',
    name: '',
    password: '',
    email: '',
    phone: '',
    status: 1
  })
  modalVisible.value = true
}

function handleEdit(record: UserView) {
  isEdit.value = true
  Object.assign(formState, {
    id: record.id,
    username: record.username,
    name: record.name,
    password: '',
    email: record.email,
    phone: record.phone,
    status: record.status
  })
  modalVisible.value = true
}

async function handleSubmit() {
  try {
    await formRef.value?.validate()
    if (isEdit.value) {
      await userApi.updateUser({
        id: formState.id,
        name: formState.name,
        email: formState.email,
        phone: formState.phone,
        status: formState.status
      })
      message.success('更新成功')
    } else {
      await userApi.createUser(formState as CreateUserRequest)
      message.success('创建成功')
    }
    modalVisible.value = false
    fetchUsers()
  } catch {
    // Form validation failed
  }
}

function handleCancel() {
  modalVisible.value = false
}

async function handleDelete(record: UserView) {
  try {
    await userApi.deleteUser(record.id)
    message.success('删除成功')
    fetchUsers()
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

.filter-bar {
  margin-bottom: 16px;
}
</style>

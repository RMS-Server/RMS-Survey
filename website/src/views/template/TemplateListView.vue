<template>
  <div class="page-container">
    <div class="page-header">
      <div class="header-top">
        <h1 class="page-title">模板管理</h1>
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
          v-model:value="filterMode"
          placeholder="模式"
          style="width: 120px"
          allowClear
          @change="handleSearch"
        >
          <a-select-option value="survey">问卷</a-select-option>
          <a-select-option value="exam">考试</a-select-option>
          <a-select-option value="vote">投票</a-select-option>
        </a-select>
        <a-select
          v-model:value="filterCategory"
          placeholder="分类"
          style="width: 150px"
          allowClear
          @change="handleSearch"
        >
          <a-select-option v-for="cat in categories" :key="cat.id" :value="cat.name">
            {{ cat.name }}
          </a-select-option>
        </a-select>
      </div>

      <a-table
        :columns="columns"
        :data-source="templates"
        :loading="loading"
        :pagination="pagination"
        row-key="id"
        @change="handleTableChange"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'name'">
            <a @click="handleEdit(record)">{{ record.name }}</a>
          </template>
          <template v-if="column.key === 'mode'">
            <a-tag>{{ getModeLabel(record.mode) }}</a-tag>
          </template>
          <template v-if="column.key === 'category'">
            <a-tag v-if="record.category" color="blue">{{ record.category }}</a-tag>
            <span v-else class="text-muted">-</span>
          </template>
          <template v-if="column.key === 'shared'">
            <a-tag v-if="record.shared" color="green">已共享</a-tag>
            <span v-else class="text-muted">-</span>
          </template>
          <template v-if="column.key === 'actions'">
            <a-space>
              <a-button type="link" size="small" @click="handleEdit(record)">
                编辑
              </a-button>
              <a-button type="link" size="small" @click="handleUseTemplate(record)">
                使用
              </a-button>
              <a-popconfirm
                title="确定要删除此模板吗？"
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
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { message } from 'ant-design-vue'
import { PlusOutlined } from '@ant-design/icons-vue'
import { templateApi } from '@/api/template'
import type { TemplateView, CategoryView } from '@/types/template'

const router = useRouter()

const searchName = ref('')
const filterMode = ref<string | undefined>()
const filterCategory = ref<string | undefined>()
const loading = ref(false)
const templates = ref<TemplateView[]>([])
const categories = ref<CategoryView[]>([])

const pagination = reactive({
  current: 1,
  pageSize: 10,
  total: 0,
  showSizeChanger: true,
  showTotal: (total: number) => `共 ${total} 条`
})

const columns = [
  { title: '名称', key: 'name', dataIndex: 'name' },
  { title: '模式', key: 'mode', width: 100 },
  { title: '分类', key: 'category', width: 120 },
  { title: '共享', key: 'shared', width: 100 },
  { title: '操作', key: 'actions', width: 200 }
]

onMounted(() => {
  fetchCategories()
  fetchTemplates()
})

async function fetchCategories() {
  try {
    categories.value = await templateApi.listCategories({})
  } catch {
    // Ignore error for categories
  }
}

async function fetchTemplates() {
  loading.value = true
  try {
    const result = await templateApi.list({
      pageIndex: pagination.current,
      pageSize: pagination.pageSize,
      name: searchName.value || undefined,
      mode: filterMode.value,
      category: filterCategory.value
    })
    templates.value = result.list
    pagination.total = result.total
  } catch {
    message.error('加载模板失败')
  } finally {
    loading.value = false
  }
}

function handleSearch() {
  pagination.current = 1
  fetchTemplates()
}

function handleTableChange(pag: { current: number; pageSize: number }) {
  pagination.current = pag.current
  pagination.pageSize = pag.pageSize
  fetchTemplates()
}

function handleCreate() {
  router.push('/template/create')
}

function handleEdit(record: TemplateView) {
  router.push(`/template/${record.id}/edit`)
}

function handleUseTemplate(record: TemplateView) {
  router.push({ path: '/project/create', query: { templateId: record.id } })
}

async function handleDelete(record: TemplateView) {
  try {
    await templateApi.delete(record.id)
    message.success('删除成功')
    fetchTemplates()
  } catch {
    message.error('删除失败')
  }
}

function getModeLabel(mode: string) {
  const labels: Record<string, string> = {
    survey: '问卷',
    exam: '考试',
    vote: '投票'
  }
  return labels[mode] || mode || '问卷'
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

.text-muted {
  color: var(--color-text-muted);
}
</style>

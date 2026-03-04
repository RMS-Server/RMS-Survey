<template>
  <div class="partner-page">
    <div class="page-header">
      <div class="header-left">
        <a-button @click="handleBack">
          <template #icon><ArrowLeftOutlined /></template>
          返回
        </a-button>
        <h2>{{ projectName }} - 参与者管理</h2>
      </div>
    </div>
    <a-card :bordered="false">
      <PartnerManager v-if="projectId" :project-id="projectId" />
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { message } from 'ant-design-vue'
import { ArrowLeftOutlined } from '@ant-design/icons-vue'
import { useProjectStore } from '@/stores/project'
import PartnerManager from './PartnerManager.vue'

const router = useRouter()
const route = useRoute()
const projectStore = useProjectStore()

const projectId = computed(() => route.params.id as string)
const projectName = ref('')

function handleBack() {
  router.push('/project')
}

onMounted(async () => {
  if (!projectId.value) {
    message.error('项目ID不存在')
    router.push('/project')
    return
  }
  try {
    const project = await projectStore.fetchProject(projectId.value)
    projectName.value = project.name
  } catch {
    message.error('加载项目信息失败')
  }
})
</script>

<style scoped>
.partner-page {
  padding: 16px;
}

.page-header {
  margin-bottom: 16px;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 16px;
}

.header-left h2 {
  margin: 0;
}
</style>

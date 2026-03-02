<template>
  <div class="survey-end-page">
    <a-result
      status="success"
      title="感谢参与！"
      sub-title="您的答卷已提交成功。"
    >
      <template #extra>
        <a-button type="primary" @click="handleGoHome">
          返回首页
        </a-button>
      </template>
    </a-result>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'

const router = useRouter()

function handleGoHome() {
  router.push('/')
}
</script>

<style scoped>
.survey-end-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #f5f5f5;
}
</style>

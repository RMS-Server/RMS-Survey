<template>
  <div class="survey-end-page">
    <div class="glass-card survey-end-card">
      <Icon
        :icon="iconName"
        class="survey-end-icon"
        :class="`survey-end-icon--${status}`"
      />

      <h2 class="survey-end-title">{{ title }}</h2>
      <p v-if="subTitle" class="survey-end-subtitle">{{ subTitle }}</p>

      <div v-if="$slots.default" class="survey-end-actions">
        <slot />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { Icon } from '@iconify/vue'

interface Props {
  status?: 'success' | 'warning' | 'error'
  title?: string
  subTitle?: string
}

const props = withDefaults(defineProps<Props>(), {
  status: 'success',
  title: '感谢参与！',
  subTitle: '您的答卷已提交成功。',
})

// Phosphor regular — clean, warm, Anthropic-adjacent. One icon per status:
//   success → check inside a soft circle
//   warning → key-lock, reading as "this device is locked from further submissions"
//   error   → exclamation circle
const iconName = computed(() => {
  switch (props.status) {
    case 'warning':
      return 'ph:lock-key'
    case 'error':
      return 'ph:warning-circle'
    default:
      return 'ph:check-circle'
  }
})
</script>

<style scoped>
.survey-end-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: var(--spacing-lg);
  background: transparent;
}

.survey-end-card {
  width: 100%;
  max-width: 520px;
  padding: var(--spacing-xxl) var(--spacing-xl);
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
  gap: var(--spacing-md);
}

.survey-end-icon {
  font-size: 96px;
  margin-bottom: var(--spacing-sm);
  filter: drop-shadow(0 6px 16px rgba(0, 0, 0, 0.08));
}

.survey-end-icon--success {
  color: #2eb582;
}

.survey-end-icon--warning {
  color: #d89614;
}

.survey-end-icon--error {
  color: #d63031;
}

.survey-end-title {
  font-size: 1.5rem;
  font-weight: 600;
  color: var(--color-text-main);
  margin: 0;
  letter-spacing: -0.3px;
}

.survey-end-subtitle {
  font-size: 0.95rem;
  color: var(--color-text-muted);
  margin: 0;
  line-height: 1.6;
}

.survey-end-actions {
  margin-top: var(--spacing-md);
}
</style>

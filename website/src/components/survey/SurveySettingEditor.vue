<template>
  <div class="survey-setting-editor">
    <a-form layout="vertical">
      <a-divider>IP 限制</a-divider>
      <a-form-item label="启用 IP 限制">
        <a-switch
          v-model:checked="localSetting.ipLimitEnabled"
          checked-children="开启"
          un-checked-children="关闭"
        />
      </a-form-item>

      <template v-if="localSetting.ipLimitEnabled">
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item label="单IP最大提交次数">
              <a-input-number
                v-model:value="localSetting.ipMaxSubmissions"
                :min="0"
                placeholder="0表示不限制"
                style="width: 100%"
              />
              <span class="setting-hint">同一IP最多可提交次数</span>
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="提交间隔(秒)">
              <a-input-number
                v-model:value="localSetting.ipInterval"
                :min="0"
                placeholder="0表示不限制"
                style="width: 100%"
              />
              <span class="setting-hint">同一IP两次提交之间的最小间隔</span>
            </a-form-item>
          </a-col>
        </a-row>
      </template>

      <a-divider>设备限制</a-divider>
      <a-form-item label="启用设备限制">
        <a-switch
          v-model:checked="localSetting.deviceLimitEnabled"
          checked-children="开启"
          un-checked-children="关闭"
        />
        <span class="setting-hint">基于浏览器指纹识别同一设备，打开链接时即预检</span>
      </a-form-item>

      <template v-if="localSetting.deviceLimitEnabled">
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item label="单设备最大提交次数">
              <a-input-number
                v-model:value="localSetting.deviceMaxSubmissions"
                :min="0"
                placeholder="0表示不限制"
                style="width: 100%"
              />
              <span class="setting-hint">同一设备最多可提交次数，达到后再次打开链接将直接拦截</span>
            </a-form-item>
          </a-col>
        </a-row>
      </template>
    </a-form>

    <div class="setting-actions">
      <a-button @click="handleCancel">取消</a-button>
      <a-button type="primary" :loading="saving" @click="handleSave">保存设置</a-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { message } from 'ant-design-vue'
import type { SurveySetting } from '@/types/survey'

interface Props {
  projectId: string
  setting: SurveySetting
}

const props = defineProps<Props>()
const emit = defineEmits<{
  (e: 'update:setting', setting: SurveySetting): void
  (e: 'cancel'): void
}>()

const saving = ref(false)

interface LocalSetting {
  ipLimitEnabled: boolean
  ipMaxSubmissions: number
  ipInterval: number
  deviceLimitEnabled: boolean
  deviceMaxSubmissions: number
}

const localSetting = ref<LocalSetting>({
  ipLimitEnabled: false,
  ipMaxSubmissions: 0,
  ipInterval: 0,
  deviceLimitEnabled: false,
  deviceMaxSubmissions: 0
})

watch(() => props.setting, (newSetting) => {
  localSetting.value = {
    ipLimitEnabled: newSetting.ipLimitEnabled || false,
    ipMaxSubmissions: newSetting.ipMaxSubmissions || 0,
    ipInterval: newSetting.ipInterval || 0,
    deviceLimitEnabled: newSetting.deviceLimitEnabled || false,
    deviceMaxSubmissions: newSetting.deviceMaxSubmissions || 0
  }
}, { immediate: true, deep: true })

function handleCancel() {
  emit('cancel')
}

async function handleSave() {
  saving.value = true
  try {
    const setting: SurveySetting = {
      projectId: props.projectId,
      ipLimitEnabled: localSetting.value.ipLimitEnabled,
      ipMaxSubmissions: localSetting.value.ipMaxSubmissions || undefined,
      ipInterval: localSetting.value.ipInterval || undefined,
      deviceLimitEnabled: localSetting.value.deviceLimitEnabled,
      deviceMaxSubmissions: localSetting.value.deviceMaxSubmissions || undefined
    }
    emit('update:setting', setting)
    message.success('设置已保存')
  } finally {
    saving.value = false
  }
}
</script>

<style scoped>
.survey-setting-editor {
  padding: 16px 0;
}

.setting-hint {
  display: block;
  font-size: 12px;
  color: var(--color-text-muted);
  margin-top: 4px;
}

.setting-actions {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
  margin-top: 24px;
  padding-top: 16px;
  border-top: 1px solid var(--border-glass);
}
</style>

<template>
  <div class="attachment-upload">
    <a-upload
      v-model:file-list="fileList"
      :action="uploadUrl"
      :data="uploadData"
      :before-upload="beforeUpload"
      :accept="acceptTypes"
      :max-count="maxFiles"
      @change="handleChange"
    >
      <a-button>
        <UploadOutlined />
        上传附件
      </a-button>
    </a-upload>
    <div v-if="tip" class="upload-tip">{{ tip }}</div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { message } from 'ant-design-vue'
import { UploadOutlined } from '@ant-design/icons-vue'
import type { UploadFile, UploadChangeParam } from 'ant-design-vue'
import type { AttachmentConfig, AttachmentInfo } from '@/types/survey'

const props = defineProps<{
  projectId: string
  questionId: string
  config: AttachmentConfig
  modelValue: AttachmentInfo[]
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: AttachmentInfo[]): void
}>()

const uploadUrl = '/api/public/uploadAttachment'
const fileList = ref<UploadFile[]>([])
const attachments = ref<AttachmentInfo[]>([...props.modelValue])
const uploadingCount = ref(0)

// Compute accept types for file input
const acceptTypes = computed(() => {
  if (!props.config.allowedTypes?.length) {
    return ''
  }
  return props.config.allowedTypes.join(',')
})

const maxFiles = computed(() => props.config.maxFiles || 1)

const tip = computed(() => {
  const maxSizeMB = Math.floor((props.config.maxSize || 10485760) / 1048576)
  const types = props.config.allowedTypes?.join(', ') || '所有类型'
  return `最多 ${maxFiles.value} 个文件，每个最大 ${maxSizeMB}MB，允许: ${types}`
})

const uploadData = computed(() => ({
  projectId: props.projectId,
  questionId: props.questionId
}))

function beforeUpload(file: File) {
  // Check file size
  const maxSize = props.config.maxSize || 10485760
  if (file.size > maxSize) {
    message.error(`文件 ${file.name} 超过最大限制 ${Math.floor(maxSize / 1048576)}MB`)
    return false
  }

  // Check file type
  if (props.config.allowedTypes?.length) {
    const ext = '.' + file.name.split('.').pop()?.toLowerCase()
    if (!props.config.allowedTypes.some(t => t.toLowerCase() === ext)) {
      message.error(`文件类型 ${ext} 不允许上传`)
      return false
    }
  }

  // Check max files (including files being uploaded)
  const currentCount = attachments.value.length + uploadingCount.value
  if (currentCount >= maxFiles.value) {
    message.error(`最多只能上传 ${maxFiles.value} 个文件`)
    return false
  }

  uploadingCount.value++
  return true
}

function handleChange(info: UploadChangeParam) {
  if (info.file.status === 'uploading') {
    return
  }
  if (info.file.status === 'done') {
    uploadingCount.value--
    const response = info.file.response
    if (response && response.data) {
      const attachment: AttachmentInfo = {
        fileId: response.data.fileId || response.data.id,
        fileName: response.data.fileName || info.file.name,
        fileSize: response.data.fileSize || info.file.size || 0,
        fileType: response.data.fileType || '.' + (info.file.name.split('.').pop()?.toLowerCase() || '')
      }
      attachments.value.push(attachment)
      emit('update:modelValue', attachments.value)
      message.success(`${info.file.name} 上传成功`)
    }
  } else if (info.file.status === 'error') {
    uploadingCount.value--
    message.error(`${info.file.name} 上传失败`)
  }
}

// Sync from parent
watch(() => props.modelValue, (val) => {
  attachments.value = [...val]
}, { deep: true })
</script>

<style scoped>
.attachment-upload {
  margin-top: 12px;
}

.upload-tip {
  margin-top: 8px;
  font-size: 12px;
  color: #999;
}
</style>

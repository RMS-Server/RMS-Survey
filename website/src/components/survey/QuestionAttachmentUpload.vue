<template>
  <div class="question-attachment-upload">
    <div class="attachment-list">
      <div
        v-for="(file, index) in modelValue"
        :key="file.fileId"
        class="attachment-item"
      >
        <a-image
          v-if="isImage(file.fileType)"
          :src="getPreviewUrl(file.fileId)"
          :width="80"
          :height="80"
          class="attachment-thumb"
        />
        <div v-else class="attachment-file-icon">
          <FileOutlined />
        </div>
        <div class="attachment-info">
          <span class="attachment-name">{{ file.fileName }}</span>
          <a-button type="link" size="small" danger @click="removeFile(index)">
            删除
          </a-button>
        </div>
      </div>
    </div>
    <a-upload
      :show-upload-list="false"
      :before-upload="beforeUpload"
      :custom-request="handleUpload"
    >
      <a-button>
        <UploadOutlined />
        上传题目附件
      </a-button>
    </a-upload>
    <div class="upload-tip">支持图片、PDF等文件，答题者可查看</div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { message } from 'ant-design-vue'
import { UploadOutlined, FileOutlined } from '@ant-design/icons-vue'
import { fileApi } from '@/api/file'
import type { QuestionAttachment } from '@/types/survey'

const props = defineProps<{
  modelValue: QuestionAttachment[]
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: QuestionAttachment[]): void
}>()

const uploading = ref(false)

const MAX_SIZE = 10 * 1024 * 1024 // 10MB
const ALLOWED_TYPES = ['.jpg', '.jpeg', '.png', '.gif', '.pdf', '.doc', '.docx', '.xls', '.xlsx']

function isImage(fileType: string): boolean {
  return ['.jpg', '.jpeg', '.png', '.gif'].includes(fileType.toLowerCase())
}

function getPreviewUrl(fileId: string): string {
  return `/api/public/preview/${fileId}`
}

function beforeUpload(file: File): boolean {
  const ext = '.' + file.name.split('.').pop()?.toLowerCase()
  if (!ALLOWED_TYPES.includes(ext)) {
    message.error(`不支持的文件类型: ${ext}`)
    return false
  }
  if (file.size > MAX_SIZE) {
    message.error('文件大小不能超过 10MB')
    return false
  }
  return true
}

async function handleUpload(options: any) {
  const file = options.file as File
  uploading.value = true
  try {
    const result = await fileApi.upload(file)
    const attachment: QuestionAttachment = {
      fileId: result.id,
      fileName: file.name,
      fileType: '.' + file.name.split('.').pop()?.toLowerCase() || ''
    }
    emit('update:modelValue', [...props.modelValue, attachment])
    message.success('上传成功')
  } catch {
    message.error('上传失败')
  } finally {
    uploading.value = false
  }
}

function removeFile(index: number) {
  const newList = [...props.modelValue]
  newList.splice(index, 1)
  emit('update:modelValue', newList)
}
</script>

<style scoped>
.question-attachment-upload {
  margin-top: 8px;
}

.attachment-list {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-bottom: 8px;
}

.attachment-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px;
  background: var(--surface-glass-input);
  border: 1px solid var(--border-glass);
  border-radius: var(--radius-sm);
  box-shadow: var(--shadow-inset);
}

.attachment-thumb {
  object-fit: cover;
  border-radius: var(--radius-sm);
}

.attachment-file-icon {
  width: 80px;
  height: 80px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--surface-glass);
  border-radius: var(--radius-sm);
  font-size: 32px;
  color: var(--color-text-muted);
}

.attachment-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.attachment-name {
  max-width: 150px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  font-size: 12px;
  color: var(--color-text-main);
}

.upload-tip {
  margin-top: 4px;
  font-size: 12px;
  color: var(--color-text-muted);
}
</style>

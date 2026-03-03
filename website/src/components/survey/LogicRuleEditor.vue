<template>
  <div class="logic-rule-editor">
    <a-card :bordered="false">
      <template #title>
        <div class="card-header">
          <span>逻辑跳转规则</span>
          <a-button type="primary" size="small" @click="addRule">
            <PlusOutlined />
            添加规则
          </a-button>
        </div>
      </template>

      <a-empty v-if="!logic.rules?.length" description="暂无逻辑跳转规则" />

      <div v-else class="rule-list">
        <div
          v-for="(rule, index) in logic.rules"
          :key="rule.id"
          class="rule-item"
        >
          <div class="rule-content">
            <a-typography-text>
              规则{{ index + 1 }}: {{ getRuleDescription(rule, elements) }}
            </a-typography-text>
          </div>
          <div class="rule-actions">
            <a-switch
              :checked="rule.enabled"
              size="small"
              @change="toggleRule(rule)"
            />
            <a-space>
              <a-button type="link" size="small" @click="editRule(rule)">
                编辑
              </a-button>
              <a-popconfirm
                title="确定要删除此规则吗？"
                @confirm="deleteRule(rule.id)"
              >
                <a-button type="link" size="small" danger>
                  删除
                </a-button>
              </a-popconfirm>
            </a-space>
          </div>
        </div>
      </div>
    </a-card>

    <ConditionBuilder
      :open="showModal"
      :rule="editingRule"
      :elements="elements"
      @update:open="showModal = $event"
      @save="handleSave"
    />
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { v4 as uuidv4 } from 'uuid'
import { PlusOutlined } from '@ant-design/icons-vue'
import ConditionBuilder from './ConditionBuilder.vue'
import type {
  SurveyLogic,
  SurveyElement,
  LogicRule,
  LogicCondition,
  LogicOperator
} from '@/types/survey'

const props = defineProps<{
  logic: SurveyLogic
  elements: SurveyElement[]
}>()

const emit = defineEmits<{
  (e: 'update:logic', value: SurveyLogic): void
}>()

const editingRule = ref<LogicRule | null>(null)
const showModal = ref(false)

// Operator display text mapping
const operatorTextMap: Record<LogicOperator, string> = {
  eq: '等于',
  neq: '不等于',
  in: '包含',
  not_in: '不包含',
  gt: '大于',
  gte: '大于等于',
  lt: '小于',
  lte: '小于等于',
  empty: '为空',
  not_empty: '不为空'
}

// Get question title by ID
function getQuestionTitle(questionId: string, elements: SurveyElement[]): string {
  const el = elements.find(e => e.id === questionId)
  return el?.title || questionId
}

// Get option text by ID for a question
function getOptionText(questionId: string, optionId: string, elements: SurveyElement[]): string {
  const el = elements.find(e => e.id === questionId)
  if (!el?.options) return optionId
  const option = el.options.find(o => o.id === optionId || o.value === optionId)
  return option?.text || String(optionId)
}

// Format condition value for display
function formatConditionValue(
  condition: LogicCondition,
  elements: SurveyElement[]
): string {
  const { operator, value, questionId } = condition

  // Operators with no value
  if (operator === 'empty' || operator === 'not_empty') {
    return ''
  }

  // Array value (for 'in', 'not_in')
  if (Array.isArray(value)) {
    const texts = value.map(v => getOptionText(questionId, v, elements))
    return texts.map(t => `"${t}"`).join('、')
  }

  // Check if this is a choice-based question (radio, checkbox, dropdown)
  const el = elements.find(e => e.id === questionId)
  if (el && (el.type === 'radio' || el.type === 'checkbox' || el.type === 'dropdown')) {
    return `"${getOptionText(questionId, String(value), elements)}"`
  }

  // Numeric or text value
  return String(value)
}

// Build condition description
function getConditionDescription(
  condition: LogicCondition,
  elements: SurveyElement[]
): string {
  const questionTitle = getQuestionTitle(condition.questionId, elements)
  const operatorText = operatorTextMap[condition.operator]
  const valueText = formatConditionValue(condition, elements)

  return `${questionTitle} ${operatorText}${valueText ? ' ' + valueText : ''}`
}

// Build condition group description
function getConditionGroupDescription(
  conditions: LogicCondition[],
  type: 'and' | 'or',
  elements: SurveyElement[]
): string {
  const descriptions = conditions.map(c => getConditionDescription(c, elements))
  const connector = type === 'and' ? '且' : '或'
  return descriptions.join(` ${connector} `)
}

// Get target element titles
function getTargetTitles(targetIds: string[], elements: SurveyElement[]): string {
  const titles = targetIds.map(id => getQuestionTitle(id, elements))
  return titles.join('、')
}

// Generate rule description from rule data
function getRuleDescription(rule: LogicRule, elements: SurveyElement[]): string {
  const conditions = rule.condition?.conditions || []
  const conditionType = rule.condition?.type || 'and'
  const conditionDesc = conditions.length > 0
    ? getConditionGroupDescription(conditions, conditionType, elements)
    : '(无条件)'
  const actionText = rule.action?.type === 'show' ? '显示' : '隐藏'
  const targetIds = rule.action?.targetIds || []
  const targetDesc = targetIds.length > 0
    ? getTargetTitles(targetIds, elements)
    : '(无目标)'

  return `当${conditionDesc}时, ${actionText}${targetDesc}`
}

// Add new rule
function addRule() {
  editingRule.value = {
    id: uuidv4(),
    condition: {
      id: uuidv4(),
      type: 'and',
      conditions: []
    },
    action: {
      type: 'show',
      targetIds: []
    },
    enabled: true
  }
  showModal.value = true
}

// Edit existing rule
function editRule(rule: LogicRule) {
  editingRule.value = JSON.parse(JSON.stringify(rule))
  showModal.value = true
}

// Delete rule
function deleteRule(ruleId: string) {
  const newRules = props.logic.rules.filter(r => r.id !== ruleId)
  emit('update:logic', { rules: newRules })
}

// Toggle rule enabled state
function toggleRule(rule: LogicRule) {
  const newRules = props.logic.rules.map(r =>
    r.id === rule.id ? { ...r, enabled: !r.enabled } : r
  )
  emit('update:logic', { rules: newRules })
}

// Handle save from ConditionBuilder
function handleSave(rule: LogicRule) {
  const existingIndex = props.logic.rules.findIndex(r => r.id === rule.id)
  let newRules: LogicRule[]

  if (existingIndex >= 0) {
    // Update existing rule
    newRules = props.logic.rules.map(r => r.id === rule.id ? rule : r)
  } else {
    // Add new rule
    newRules = [...props.logic.rules, rule]
  }

  emit('update:logic', { rules: newRules })
  showModal.value = false
  editingRule.value = null
}
</script>

<style scoped>
.logic-rule-editor {
  width: 100%;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  width: 100%;
}

.rule-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.rule-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  background: #fafafa;
  border-radius: 4px;
  border: 1px solid #f0f0f0;
}

.rule-content {
  flex: 1;
  min-width: 0;
}

.rule-actions {
  display: flex;
  align-items: center;
  gap: 12px;
  flex-shrink: 0;
}
</style>

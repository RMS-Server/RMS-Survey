<template>
  <a-modal
    :open="open"
    title="编辑规则"
    :width="600"
    @cancel="handleCancel"
    @ok="handleOk"
  >
    <a-form :label-col="{ span: 4 }" :wrapper-col="{ span: 20 }">
      <!-- Rule Name -->
      <a-form-item label="规则名称">
        <a-input
          v-model:value="localRule.name"
          placeholder="可选，便于识别规则"
        />
      </a-form-item>

      <!-- Condition Combination Type -->
      <a-form-item label="条件组合">
        <a-select v-model:value="localRule.condition.type" style="width: 120px">
          <a-select-option value="and">且</a-select-option>
          <a-select-option value="or">或</a-select-option>
        </a-select>
      </a-form-item>

      <!-- Conditions List -->
      <a-form-item label="条件列表">
        <div class="conditions-container">
          <div
            v-for="(condition, index) in localRule.condition.conditions"
            :key="index"
            class="condition-item"
          >
            <div class="condition-row">
              <label>题目:</label>
              <a-select
                v-model:value="condition.questionId"
                placeholder="选择题目"
                style="flex: 1"
                @change="(val: string) => onQuestionChange(condition, val)"
              >
                <a-select-option
                  v-for="el in conditionElements"
                  :key="el.id"
                  :value="el.id"
                >
                  {{ getElementLabel(el) }}
                </a-select-option>
              </a-select>
            </div>

            <div class="condition-row">
              <label>操作:</label>
              <a-select
                v-model:value="condition.operator"
                placeholder="选择操作"
                style="flex: 1"
                @change="() => onOperatorChange(condition)"
              >
                <a-select-option
                  v-for="op in getAvailableOperators(condition.questionId)"
                  :key="op"
                  :value="op"
                >
                  {{ operatorLabels[op] }}
                </a-select-option>
              </a-select>
            </div>

            <div class="condition-row" v-if="!isEmptyOperator(condition.operator)">
              <label>值:</label>
              <a-select
                v-if="isOptionBasedQuestion(condition.questionId) && !isMultiSelectQuestion(condition.questionId)"
                v-model:value="condition.value"
                placeholder="选择选项"
                style="flex: 1"
              >
                <a-select-option
                  v-for="opt in getQuestionOptions(condition.questionId)"
                  :key="opt.id"
                  :value="opt.value || opt.text"
                >
                  {{ opt.text }}
                </a-select-option>
              </a-select>
              <a-select
                v-else-if="isMultiSelectQuestion(condition.questionId)"
                v-model:value="condition.value"
                mode="multiple"
                placeholder="选择选项"
                style="flex: 1"
              >
                <a-select-option
                  v-for="opt in getQuestionOptions(condition.questionId)"
                  :key="opt.id"
                  :value="opt.value || opt.text"
                >
                  {{ opt.text }}
                </a-select-option>
              </a-select>
              <a-input-number
                v-else-if="isRatingQuestion(condition.questionId)"
                v-model:value="condition.value as number"
                :min="1"
                :max="5"
                style="flex: 1"
              />
            </div>

            <a-button
              type="text"
              danger
              :disabled="localRule.condition.conditions.length <= 1"
              @click="removeCondition(index)"
            >
              <DeleteOutlined />
            </a-button>
          </div>

          <a-button type="dashed" block @click="addCondition">
            <PlusOutlined />
            添加条件
          </a-button>
        </div>
      </a-form-item>

      <a-divider />

      <!-- Action Type -->
      <a-form-item label="动作类型">
        <a-select v-model:value="localRule.action.type" style="width: 120px">
          <a-select-option value="show">显示</a-select-option>
          <a-select-option value="hide">隐藏</a-select-option>
        </a-select>
      </a-form-item>

      <!-- Target Questions -->
      <a-form-item label="目标题目">
        <div class="target-questions">
          <a-checkbox
            v-for="el in targetElements"
            :key="el.id"
            :checked="localRule.action.targetIds.includes(el.id)"
            @change="(e: any) => toggleTarget(el.id, e.target.checked)"
          >
            {{ getElementLabel(el) }}
          </a-checkbox>
        </div>
      </a-form-item>
    </a-form>
  </a-modal>
</template>

<script setup lang="ts">
import { ref, watch, computed } from 'vue'
import { message } from 'ant-design-vue'
import { v4 as uuidv4 } from 'uuid'
import { DeleteOutlined, PlusOutlined } from '@ant-design/icons-vue'
import type {
  SurveyElement,
  LogicRule,
  LogicCondition,
  LogicOperator,
  SurveyOption
} from '@/types/survey'

const props = defineProps<{
  open: boolean
  rule: LogicRule | null
  elements: SurveyElement[]
}>()

const emit = defineEmits<{
  (e: 'update:open', value: boolean): void
  (e: 'save', rule: LogicRule): void
}>()

const operatorLabels: Record<LogicOperator, string> = {
  eq: '等于',
  neq: '不等于',
  in: '包含于',
  not_in: '不包含于',
  gt: '大于',
  gte: '大于等于',
  lt: '小于',
  lte: '小于等于',
  empty: '为空',
  not_empty: '不为空'
}

// Filter elements that can be used in conditions (have options or rating)
const conditionElements = computed(() => {
  return props.elements.filter(el =>
    el.type === 'radio' || el.type === 'checkbox' ||
    el.type === 'dropdown' || el.type === 'rating'
  )
})

// All elements can be targets
const targetElements = computed(() => props.elements)

// Create empty condition
function createEmptyCondition(): LogicCondition {
  return {
    questionId: '',
    operator: 'eq',
    value: ''
  }
}

// Create default rule
function createDefaultRule(): LogicRule {
  return {
    id: uuidv4(),
    name: '',
    condition: {
      id: uuidv4(),
      type: 'and',
      conditions: [createEmptyCondition()]
    },
    action: {
      type: 'show',
      targetIds: []
    },
    enabled: true
  }
}

const localRule = ref<LogicRule>(createDefaultRule())

// Watch for rule prop changes
watch(() => props.rule, (val) => {
  if (val) {
    localRule.value = JSON.parse(JSON.stringify(val))
  } else {
    localRule.value = createDefaultRule()
  }
}, { immediate: true })

// Get element label for display
function getElementLabel(el: SurveyElement): string {
  const index = props.elements.findIndex(e => e.id === el.id)
  return `Q${index + 1} ${el.title}`
}

// Get element by ID
function getElementById(id: string): SurveyElement | undefined {
  return props.elements.find(el => el.id === id)
}

// Get available operators based on question type
function getAvailableOperators(questionId: string): LogicOperator[] {
  const element = getElementById(questionId)
  if (!element) return ['eq', 'neq', 'empty', 'not_empty']

  switch (element.type) {
    case 'radio':
    case 'dropdown':
      return ['eq', 'neq', 'empty', 'not_empty']
    case 'checkbox':
      return ['in', 'not_in', 'empty', 'not_empty']
    case 'rating':
      return ['gt', 'gte', 'lt', 'lte', 'eq', 'neq', 'empty', 'not_empty']
    default:
      return ['empty', 'not_empty']
  }
}

// Check if operator is empty/not_empty (no value needed)
function isEmptyOperator(operator: LogicOperator): boolean {
  return operator === 'empty' || operator === 'not_empty'
}

// Check if question is option-based (radio, dropdown)
function isOptionBasedQuestion(questionId: string): boolean {
  const element = getElementById(questionId)
  return element?.type === 'radio' || element?.type === 'dropdown'
}

// Check if question is multi-select (checkbox)
function isMultiSelectQuestion(questionId: string): boolean {
  const element = getElementById(questionId)
  return element?.type === 'checkbox'
}

// Check if question is rating
function isRatingQuestion(questionId: string): boolean {
  const element = getElementById(questionId)
  return element?.type === 'rating'
}

// Get question options
function getQuestionOptions(questionId: string): SurveyOption[] {
  const element = getElementById(questionId)
  return element?.options || []
}

// Handle question change
function onQuestionChange(condition: LogicCondition, questionId: string) {
  condition.questionId = questionId
  // Reset operator and value when question changes
  const operators = getAvailableOperators(questionId)
  condition.operator = operators[0]
  condition.value = isMultiSelectQuestion(questionId) ? [] : ''
}

// Handle operator change
function onOperatorChange(condition: LogicCondition) {
  // Reset value when operator changes to/from empty operators
  if (isEmptyOperator(condition.operator)) {
    condition.value = ''
  } else if (isMultiSelectQuestion(condition.questionId)) {
    if (!Array.isArray(condition.value)) {
      condition.value = []
    }
  } else {
    if (Array.isArray(condition.value)) {
      condition.value = ''
    }
  }
}

// Add condition
function addCondition() {
  localRule.value.condition.conditions.push(createEmptyCondition())
}

// Remove condition
function removeCondition(index: number) {
  if (localRule.value.condition.conditions.length > 1) {
    localRule.value.condition.conditions.splice(index, 1)
  }
}

// Toggle target question
function toggleTarget(id: string, checked: boolean) {
  if (checked) {
    if (!localRule.value.action.targetIds.includes(id)) {
      localRule.value.action.targetIds.push(id)
    }
  } else {
    const index = localRule.value.action.targetIds.indexOf(id)
    if (index > -1) {
      localRule.value.action.targetIds.splice(index, 1)
    }
  }
}

// Handle cancel
function handleCancel() {
  emit('update:open', false)
}

// Handle OK
function handleOk() {
  // Validate conditions
  const validConditions = localRule.value.condition.conditions.filter(
    c => c.questionId && c.operator
  )

  if (validConditions.length === 0) {
    message.warning('请至少添加一个有效条件')
    return
  }

  // Validate target IDs
  if (localRule.value.action.targetIds.length === 0) {
    message.warning('请至少选择一个目标题目')
    return
  }

  const ruleToSave: LogicRule = {
    ...localRule.value,
    condition: {
      ...localRule.value.condition,
      conditions: validConditions
    }
  }

  emit('save', ruleToSave)
  emit('update:open', false)
}
</script>

<style scoped>
.conditions-container {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.condition-item {
  padding: 12px;
  border: 1px solid var(--border-glass);
  border-radius: var(--radius-sm);
  background: var(--surface-glass-input);
  box-shadow: var(--shadow-inset);
}

.condition-row {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 8px;
}

.condition-row:last-child {
  margin-bottom: 0;
}

.condition-row label {
  width: 40px;
  flex-shrink: 0;
  color: var(--color-text-main);
}

.target-questions {
  display: flex;
  flex-direction: column;
  gap: 8px;
}
</style>

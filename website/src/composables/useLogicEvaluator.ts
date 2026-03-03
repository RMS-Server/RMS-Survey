import { computed, type ComputedRef, type Ref } from 'vue'
import type { SurveySchema, LogicCondition, LogicConditionGroup, AnswerValue } from '@/types/survey'

export function useLogicEvaluator(
  schema: ComputedRef<SurveySchema> | Ref<SurveySchema>,
  answers: ComputedRef<Record<string, unknown>> | Ref<Record<string, unknown>>
) {
  // Get answer value for a question
  function getAnswerValue(questionId: string): unknown {
    const answer = answers.value[questionId]
    if (answer && typeof answer === 'object' && 'value' in answer) {
      return (answer as AnswerValue).value
    }
    return answer
  }

  // Check if value is empty
  function isEmpty(value: unknown): boolean {
    if (value === null || value === undefined || value === '') return true
    if (Array.isArray(value) && value.length === 0) return true
    return false
  }

  // Evaluate a single condition
  function evaluateCondition(condition: LogicCondition): boolean {
    const answerValue = getAnswerValue(condition.questionId)
    const conditionValue = condition.value

    switch (condition.operator) {
      case 'empty':
        return isEmpty(answerValue)
      case 'not_empty':
        return !isEmpty(answerValue)
      case 'eq':
        // Handle string/number comparison
        if (typeof answerValue === 'number' && typeof conditionValue === 'string') {
          return answerValue === Number(conditionValue)
        }
        if (typeof answerValue === 'string' && typeof conditionValue === 'number') {
          return answerValue === String(conditionValue)
        }
        return answerValue === conditionValue
      case 'neq':
        // Handle string/number comparison
        if (typeof answerValue === 'number' && typeof conditionValue === 'string') {
          return answerValue !== Number(conditionValue)
        }
        if (typeof answerValue === 'string' && typeof conditionValue === 'number') {
          return answerValue !== String(conditionValue)
        }
        return answerValue !== conditionValue
      case 'gt':
      case 'gte':
      case 'lt':
      case 'lte': {
        // Convert to numbers for comparison
        const numAnswer = typeof answerValue === 'number' ? answerValue : Number(answerValue)
        const numCond = typeof conditionValue === 'number' ? conditionValue : Number(conditionValue)
        if (isNaN(numAnswer) || isNaN(numCond)) return false
        switch (condition.operator) {
          case 'gt': return numAnswer > numCond
          case 'gte': return numAnswer >= numCond
          case 'lt': return numAnswer < numCond
          case 'lte': return numAnswer <= numCond
          default: return false
        }
      }
      case 'in': {
        if (!Array.isArray(conditionValue)) return false
        // For checkbox: answerValue is an array of selected values
        if (Array.isArray(answerValue)) {
          return answerValue.some(v => conditionValue.includes(v))
        }
        // For single value: check if it's in the conditionValue array
        return conditionValue.includes(answerValue as string)
      }
      case 'not_in': {
        if (!Array.isArray(conditionValue)) return true
        // For checkbox: answerValue is an array of selected values
        if (Array.isArray(answerValue)) {
          return !answerValue.some(v => conditionValue.includes(v))
        }
        // For single value: check if it's NOT in the conditionValue array
        return !conditionValue.includes(answerValue as string)
      }
      default:
        return false
    }
  }

  // Evaluate a condition group (AND/OR)
  function evaluateConditionGroup(group: LogicConditionGroup): boolean {
    if (!group.conditions || group.conditions.length === 0) {
      return true // Empty group is always true
    }

    const results = group.conditions.map(c => evaluateCondition(c))

    return group.type === 'and'
      ? results.every(Boolean)
      : results.some(Boolean)
  }

  // Get all question IDs from schema
  function getAllQuestionIds(): string[] {
    const ids: string[] = []
    if (schema.value?.pages) {
      for (const page of schema.value.pages) {
        for (const element of page.elements || []) {
          ids.push(element.id)
        }
      }
    }
    return ids
  }

  // Compute visible question IDs based on logic rules
  const visibleQuestionIds = computed<Set<string>>(() => {
    const allIds = getAllQuestionIds()
    const rules = schema.value?.logic?.rules

    if (!rules || rules.length === 0) {
      return new Set<string>(allIds) // No rules, all visible
    }

    // Track visibility state for each question
    // null = no explicit rule, true = should show, false = should hide
    const showStates = new Map<string, boolean | null>()

    // Track questions targeted by hide rules
    const hideStates = new Map<string, boolean>()

    for (const rule of rules) {
      if (!rule.enabled) continue

      const conditionMet = evaluateConditionGroup(rule.condition)

      for (const targetId of rule.action.targetIds) {
        if (!allIds.includes(targetId)) continue // Skip invalid target IDs

        if (rule.action.type === 'show') {
          // For show: any true condition wins
          const currentState = showStates.get(targetId)
          if (conditionMet) {
            showStates.set(targetId, true)
          } else if (currentState === null || currentState === undefined) {
            // No explicit show yet and condition not met
            showStates.set(targetId, false)
          }
          // If already true, keep it true
        } else if (rule.action.type === 'hide') {
          // For hide: any true condition wins
          if (conditionMet) {
            hideStates.set(targetId, true)
          }
        }
      }
    }

    // Build final visibility set
    const visible = new Set<string>()
    for (const id of allIds) {
      const showState = showStates.get(id)
      const hideState = hideStates.get(id)

      // Hide rule takes precedence if it evaluates to true
      if (hideState === true) {
        continue // Hidden by hide rule
      }

      // If there's a show rule for this question
      if (showState !== null && showState !== undefined) {
        if (showState === true) {
          visible.add(id)
        }
        // If showState is false, don't add to visible
      } else {
        // No show rule targeting this question, default visible
        visible.add(id)
      }
    }

    return visible
  })

  // Check if a question is visible
  function isQuestionVisible(questionId: string): boolean {
    return visibleQuestionIds.value.has(questionId)
  }

  return {
    evaluateCondition,
    evaluateConditionGroup,
    visibleQuestionIds,
    isQuestionVisible
  }
}

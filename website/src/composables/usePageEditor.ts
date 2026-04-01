import { ref, computed } from 'vue'
import { v4 as uuidv4 } from 'uuid'
import type { SurveyPage, SurveyElement, SurveyLogic } from '@/types/survey'

export function usePageEditor(surveyLogic: { value: SurveyLogic }) {
  const pages = ref<SurveyPage[]>([
    { id: uuidv4(), title: '', description: '', elements: [] }
  ])
  const activePageIndex = ref(0)

  const activePage = computed(() => pages.value[activePageIndex.value])
  const allElements = computed(() => pages.value.flatMap(p => p.elements))

  function addPage() {
    pages.value.push({ id: uuidv4(), title: '', description: '', elements: [] })
    activePageIndex.value = pages.value.length - 1
  }

  // Move elements to adjacent page, return the removed page's id
  function removePage(index: number): boolean {
    if (pages.value.length <= 1) return false
    const removed = pages.value[index]
    const targetIndex = index > 0 ? index - 1 : 1
    const target = pages.value[targetIndex]
    target.elements.push(...removed.elements)
    pages.value.splice(index, 1)
    // Adjust active page: if active was after the removed page, shift down
    if (activePageIndex.value > index) {
      activePageIndex.value--
    } else if (activePageIndex.value === index) {
      activePageIndex.value = Math.max(0, index - 1)
    }
    return true
  }

  // Delete page and its questions, returns the IDs of deleted elements
  function deletePageWithQuestions(index: number): string[] {
    if (pages.value.length <= 1) return []
    const removed = pages.value[index]
    const deletedIds = removed.elements.map(e => e.id)
    for (const id of deletedIds) {
      cleanupLogicRules(id)
    }
    pages.value.splice(index, 1)
    if (activePageIndex.value > index) {
      activePageIndex.value--
    } else if (activePageIndex.value === index) {
      activePageIndex.value = Math.max(0, index - 1)
    }
    return deletedIds
  }

  // Move element to target page; returns the target page index
  function moveQuestionToPage(elementId: string, targetPageId: string): number {
    let element: SurveyElement | undefined
    for (const page of pages.value) {
      const idx = page.elements.findIndex(e => e.id === elementId)
      if (idx !== -1) {
        element = page.elements.splice(idx, 1)[0]
        break
      }
    }
    if (!element) return activePageIndex.value
    const targetIdx = pages.value.findIndex(p => p.id === targetPageId)
    if (targetIdx !== -1) {
      pages.value[targetIdx].elements.push(element)
      return targetIdx
    }
    return activePageIndex.value
  }

  function cleanupLogicRules(deletedId: string) {
    const rulesToRemove: string[] = []
    for (const rule of surveyLogic.value.rules) {
      rule.action.targetIds = rule.action.targetIds.filter(id => id !== deletedId)
      rule.condition.conditions = rule.condition.conditions.filter(
        c => c.questionId !== deletedId
      )
      if (rule.action.targetIds.length === 0 || rule.condition.conditions.length === 0) {
        rulesToRemove.push(rule.id)
      }
    }
    surveyLogic.value.rules = surveyLogic.value.rules.filter(
      r => !rulesToRemove.includes(r.id)
    )
  }

  function loadPages(surveyPages: SurveyPage[]) {
    if (surveyPages && surveyPages.length > 0) {
      pages.value = surveyPages.map(p => ({
        ...p,
        // Clear the legacy default title 'Page 1' (hardcoded by old editor)
        title: p.title === 'Page 1' ? '' : p.title,
        description: p.description || ''
      }))
    } else {
      pages.value = [{ id: uuidv4(), title: '', description: '', elements: [] }]
    }
    activePageIndex.value = 0
  }

  return {
    pages,
    activePageIndex,
    activePage,
    allElements,
    addPage,
    removePage,
    deletePageWithQuestions,
    moveQuestionToPage,
    cleanupLogicRules,
    loadPages
  }
}

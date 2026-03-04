import type { TimingInfo, DraftData } from '@/types/answer'

const DRAFT_PREFIX = 'survey_draft_'

export function useDraftManager() {
  // Save draft to localStorage
  function saveDraft(projectId: string, answers: Record<string, unknown>, timing: TimingInfo) {
    const draft: DraftData = {
      projectId,
      answers,
      timing,
      savedAt: Date.now()
    }
    localStorage.setItem(DRAFT_PREFIX + projectId, JSON.stringify(draft))
  }

  // Load draft from localStorage
  function loadDraft(projectId: string): DraftData | null {
    const data = localStorage.getItem(DRAFT_PREFIX + projectId)
    if (!data) return null
    try {
      return JSON.parse(data) as DraftData
    } catch {
      return null
    }
  }

  // Clear draft from localStorage
  function clearDraft(projectId: string) {
    localStorage.removeItem(DRAFT_PREFIX + projectId)
  }

  return { saveDraft, loadDraft, clearDraft }
}

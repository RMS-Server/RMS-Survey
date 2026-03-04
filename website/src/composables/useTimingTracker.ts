import { ref } from 'vue'
import type { QuestionTiming, TimingInfo } from '@/types/answer'

export function useTimingTracker() {
  const surveyStartTime = ref<number>(0)
  const lastAnswerTime = ref<number>(0)
  const questionTimings = ref<Map<string, QuestionTiming>>(new Map())
  const firstAnswerRecorded = ref<Set<string>>(new Set())

  // Initialize timing when survey loads
  function initTiming(existingTiming?: TimingInfo) {
    if (existingTiming) {
      // Draft restore: preserve existing timing data, continue from now
      surveyStartTime.value = existingTiming.surveyStartTime
      lastAnswerTime.value = Date.now()
      questionTimings.value.clear()
      firstAnswerRecorded.value.clear()
      // Restore recorded question timings
      for (const qt of existingTiming.questionTimings) {
        questionTimings.value.set(qt.questionId, { ...qt })
        if (qt.firstAnswerTime > 0) {
          firstAnswerRecorded.value.add(qt.questionId)
        }
      }
    } else {
      // Fresh start
      surveyStartTime.value = Date.now()
      lastAnswerTime.value = surveyStartTime.value
      questionTimings.value.clear()
      firstAnswerRecorded.value.clear()
    }
  }

  // Record answer change (only record time on first answer)
  function recordAnswer(questionId: string) {
    if (firstAnswerRecorded.value.has(questionId)) return  // Already recorded

    const now = Date.now()
    const timing: QuestionTiming = {
      questionId,
      startTime: lastAnswerTime.value,
      firstAnswerTime: now,
      duration: now - lastAnswerTime.value
    }
    questionTimings.value.set(questionId, timing)
    firstAnswerRecorded.value.add(questionId)
    lastAnswerTime.value = now
  }

  // Generate timing data for submission or draft save
  function getTimingInfo(): TimingInfo {
    const submitTime = Date.now()
    const timings: QuestionTiming[] = []
    questionTimings.value.forEach((t) => timings.push(t))
    return {
      surveyStartTime: surveyStartTime.value,
      submitTime,
      totalDuration: submitTime - surveyStartTime.value,
      questionTimings: timings
    }
  }

  return { initTiming, recordAnswer, getTimingInfo }
}

// Device fingerprint generator for per-device submission limits.
//
// Collects a fixed set of stable browser signals (UA, screen, timezone,
// hardware hints, canvas/WebGL output) and hashes them with SHA-256.
// The resulting hex string is used server-side as the `device_fingerprint`
// column on t_answer. Cached in localStorage so we don't recompute on
// every page load, but any missing cache entry regenerates deterministically.

const CACHE_KEY = 'rms_device_fp_v1'

interface CachedFingerprint {
  fp: string
  ts: number
}

function canvasSignal(): string {
  try {
    const canvas = document.createElement('canvas')
    canvas.width = 240
    canvas.height = 60
    const ctx = canvas.getContext('2d')
    if (!ctx) return ''
    ctx.textBaseline = 'top'
    ctx.font = "14px 'Arial'"
    ctx.fillStyle = '#f60'
    ctx.fillRect(125, 1, 62, 20)
    ctx.fillStyle = '#069'
    ctx.fillText('rms-survey-fp-⚓🔒', 2, 15)
    ctx.strokeStyle = 'rgba(102,204,0,0.7)'
    ctx.beginPath()
    ctx.arc(50, 30, 20, 0, Math.PI * 2)
    ctx.stroke()
    return canvas.toDataURL()
  } catch {
    return ''
  }
}

function webglSignal(): string {
  try {
    const canvas = document.createElement('canvas')
    const gl = (canvas.getContext('webgl') || canvas.getContext('experimental-webgl')) as WebGLRenderingContext | null
    if (!gl) return ''
    const dbg = gl.getExtension('WEBGL_debug_renderer_info')
    const vendor = dbg ? gl.getParameter(dbg.UNMASKED_VENDOR_WEBGL) : gl.getParameter(gl.VENDOR)
    const renderer = dbg ? gl.getParameter(dbg.UNMASKED_RENDERER_WEBGL) : gl.getParameter(gl.RENDERER)
    return `${vendor}|${renderer}`
  } catch {
    return ''
  }
}

function collectSignals(): string {
  const nav = navigator as Navigator & { deviceMemory?: number; doNotTrack?: string | null }
  const parts = [
    nav.userAgent,
    nav.language,
    (nav.languages || []).join(','),
    nav.platform,
    String(nav.hardwareConcurrency ?? ''),
    String(nav.deviceMemory ?? ''),
    String(nav.maxTouchPoints ?? ''),
    `${screen.width}x${screen.height}x${screen.colorDepth}`,
    String(window.devicePixelRatio ?? ''),
    Intl.DateTimeFormat().resolvedOptions().timeZone || '',
    String(new Date().getTimezoneOffset()),
    canvasSignal(),
    webglSignal(),
  ]
  return parts.join('||')
}

async function sha256Hex(input: string): Promise<string> {
  const buf = new TextEncoder().encode(input)
  const digest = await crypto.subtle.digest('SHA-256', buf)
  return Array.from(new Uint8Array(digest))
    .map(b => b.toString(16).padStart(2, '0'))
    .join('')
}

let inflight: Promise<string> | null = null

export function useDeviceFingerprint() {
  async function getFingerprint(): Promise<string> {
    if (inflight) return inflight
    inflight = (async () => {
      try {
        const cached = localStorage.getItem(CACHE_KEY)
        if (cached) {
          const parsed = JSON.parse(cached) as CachedFingerprint
          if (parsed.fp && parsed.fp.length === 64) return parsed.fp
        }
      } catch { /* localStorage unavailable: fall through */ }

      const hash = await sha256Hex(collectSignals())
      try {
        localStorage.setItem(CACHE_KEY, JSON.stringify({ fp: hash, ts: Date.now() }))
      } catch { /* ignore quota / private-mode errors */ }
      return hash
    })()
    return inflight
  }

  return { getFingerprint }
}

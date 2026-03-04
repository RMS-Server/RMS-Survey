/**
 * OAuth PKCE utilities for secure OAuth authentication
 */

/**
 * Generate a random string for state parameter
 */
export function generateState(): string {
  const array = new Uint8Array(32)
  crypto.getRandomValues(array)
  return Array.from(array, byte => byte.toString(16).padStart(2, '0')).join('')
}

/**
 * Generate a code_verifier for PKCE
 * Must be between 43 and 128 characters
 */
export function generateCodeVerifier(): string {
  const array = new Uint8Array(32)
  crypto.getRandomValues(array)
  return base64URLEncode(array)
}

/**
 * Generate a code_challenge from code_verifier using S256 method
 */
export async function generateCodeChallenge(verifier: string): Promise<string> {
  const encoder = new TextEncoder()
  const data = encoder.encode(verifier)
  const hash = await crypto.subtle.digest('SHA-256', data)
  return base64URLEncode(new Uint8Array(hash))
}

/**
 * Base64 URL encode without padding
 */
function base64URLEncode(array: Uint8Array): string {
  let str = ''
  array.forEach(byte => {
    str += String.fromCharCode(byte)
  })
  return btoa(str)
    .replace(/\+/g, '-')
    .replace(/\//g, '_')
    .replace(/=+$/, '')
}

/**
 * Generate PKCE code_verifier and code_challenge pair
 */
export async function generatePKCE(): Promise<{ codeVerifier: string; codeChallenge: string }> {
  const codeVerifier = generateCodeVerifier()
  const codeChallenge = await generateCodeChallenge(codeVerifier)
  return { codeVerifier, codeChallenge }
}

/**
 * Store OAuth state and code_verifier in sessionStorage
 */
export function setOAuthState(state: string, codeVerifier: string): void {
  sessionStorage.setItem('oauth_state', state)
  sessionStorage.setItem('oauth_code_verifier', codeVerifier)
}

/**
 * Get OAuth state and code_verifier from sessionStorage
 */
export function getOAuthState(): { state: string; codeVerifier: string } | null {
  const state = sessionStorage.getItem('oauth_state')
  const codeVerifier = sessionStorage.getItem('oauth_code_verifier')
  if (!state || !codeVerifier) {
    return null
  }
  return { state, codeVerifier }
}

/**
 * Clear OAuth state from sessionStorage
 */
export function clearOAuthState(): void {
  sessionStorage.removeItem('oauth_state')
  sessionStorage.removeItem('oauth_code_verifier')
}

/**
 * Decode JWT token and get expiry time
 */
export function decodeJWT(token: string): { exp: number; [key: string]: unknown } | null {
  try {
    const parts = token.split('.')
    if (parts.length !== 3) {
      return null
    }
    // Convert URL-safe base64 to standard base64
    let payload = parts[1]
    payload = payload.replace(/-/g, '+').replace(/_/g, '/')
    // Add padding if needed
    while (payload.length % 4) {
      payload += '='
    }
    const decoded = atob(payload)
    return JSON.parse(decoded)
  } catch {
    return null
  }
}

/**
 * Check if JWT token is expired or will expire soon
 * @param token JWT token
 * @param bufferMinutes Minutes before expiry to consider as "expiring soon"
 */
export function isTokenExpiring(token: string, bufferMinutes: number = 5): boolean {
  const decoded = decodeJWT(token)
  if (!decoded || !decoded.exp) {
    return true
  }
  const expiresAt = decoded.exp * 1000 // Convert to milliseconds
  const bufferMs = bufferMinutes * 60 * 1000
  return Date.now() >= expiresAt - bufferMs
}

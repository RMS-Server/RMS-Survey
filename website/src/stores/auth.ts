import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { authApi, userApi } from '@/api/auth'
import { getToken, setToken, removeToken } from '@/utils/storage'
import { generatePKCE, generateState, setOAuthState, getOAuthState, clearOAuthState } from '@/utils/oauth'
import type { UserView } from '@/types/user'

interface OAuthConfig {
  authUrl: string
  clientId: string
  redirectUri: string
  scopes: string
}

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string | null>(getToken())
  const user = ref<UserView | null>(null)
  let refreshTimer: ReturnType<typeof setTimeout> | null = null
  let lastRefreshTime: number = 0
  const MIN_REFRESH_INTERVAL = 60000 // Minimum 1 minute between refreshes

  const isLoggedIn = computed(() => !!token.value && !!user.value)

  /**
   * Initiate OAuth login flow - redirects to SSO provider
   */
  async function initiateOAuthLogin() {
    try {
      // Get OAuth configuration from backend
      const config = await authApi.getOAuthAuthorize() as OAuthConfig

      // Generate PKCE challenge and state locally
      const { codeVerifier, codeChallenge } = await generatePKCE()
      const state = generateState()

      // Store state and verifier for callback verification
      setOAuthState(state, codeVerifier)

      // Build the authorization URL with PKCE
      const authUrl = new URL(config.authUrl)
      authUrl.searchParams.set('client_id', config.clientId)
      authUrl.searchParams.set('response_type', 'code')
      authUrl.searchParams.set('redirect_uri', config.redirectUri)
      authUrl.searchParams.set('scope', config.scopes)
      authUrl.searchParams.set('state', state)
      authUrl.searchParams.set('code_challenge', codeChallenge)
      authUrl.searchParams.set('code_challenge_method', 'S256')

      // Redirect to SSO provider
      window.location.href = authUrl.toString()
    } catch (error) {
      console.error('Failed to initiate OAuth login:', error)
      throw error
    }
  }

  /**
   * Handle OAuth callback - exchange code for JWT
   */
  async function handleOAuthCallback(code: string, state: string): Promise<void> {
    // Verify state
    const stored = getOAuthState()
    if (!stored || stored.state !== state) {
      throw new Error('OAuth state mismatch')
    }

    // Exchange code for token (POST with JSON body)
    const response = await authApi.oauthCallback({
      code,
      state,
      codeVerifier: stored.codeVerifier
    })

    // Clear stored state
    clearOAuthState()

    // Store token and user
    token.value = response.token
    user.value = response.user
    setToken(response.token)

    // Setup auto refresh
    setupAutoRefresh()
  }

  /**
   * Refresh the JWT token using stored refresh token
   */
  async function refreshToken(): Promise<void> {
    if (!token.value) {
      return
    }

    // Debounce: don't refresh more than once per minute
    const now = Date.now()
    if (now - lastRefreshTime < MIN_REFRESH_INTERVAL) {
      return
    }
    lastRefreshTime = now

    try {
      const response = await authApi.oauthRefresh()
      token.value = response.token
      user.value = response.user
      setToken(response.token)

      // Setup next refresh
      setupAutoRefresh()
    } catch (error) {
      console.error('Failed to refresh token:', error)
      // Token refresh failed, logout
      logout()
      throw error
    }
  }

  /**
   * Setup auto refresh timer to refresh token before expiry
   */
  function setupAutoRefresh(): void {
    // Clear existing timer
    if (refreshTimer) {
      clearTimeout(refreshTimer)
      refreshTimer = null
    }

    if (!token.value) {
      return
    }

    // Decode token to get expiry time
    const decoded = decodeJWT(token.value)
    if (!decoded || !decoded.exp) {
      console.warn('Failed to decode token for auto-refresh')
      return
    }

    // Schedule refresh 5 minutes before expiry
    const expiresAt = decoded.exp * 1000 // Convert to milliseconds
    const bufferMs = 5 * 60 * 1000 // 5 minutes buffer
    const refreshIn = expiresAt - Date.now() - bufferMs

    if (refreshIn > 0) {
      refreshTimer = setTimeout(() => {
        refreshToken().catch(() => {
          // Silent fail, user will be redirected to login on next API call
        })
      }, refreshIn)
    }
  }

  /**
   * Decode JWT token (simple base64 decode)
   */
  function decodeJWT(t: string): { exp: number; [key: string]: unknown } | null {
    try {
      const parts = t.split('.')
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
   * Logout - clear token and session
   */
  async function logout() {
    try {
      await authApi.logout()
    } catch {
      // Ignore logout API errors
    } finally {
      token.value = null
      user.value = null
      removeToken()
      if (refreshTimer) {
        clearTimeout(refreshTimer)
        refreshTimer = null
      }
    }
  }

  /**
   * Fetch current user info
   */
  async function fetchCurrentUser() {
    if (!token.value) return null
    try {
      const userData = await userApi.getCurrentUser()
      user.value = userData
      setupAutoRefresh()
      return userData
    } catch {
      logout()
      return null
    }
  }

  return {
    token,
    user,
    isLoggedIn,
    initiateOAuthLogin,
    handleOAuthCallback,
    refreshToken,
    logout,
    fetchCurrentUser
  }
})

import { getStoredToken, setStoredToken } from './authStorage.js'

export const API_BASE = import.meta.env.VITE_API_BASE_URL ?? 'http://localhost:8000/api/v1'

let onUnauthorized = () => {}

export function setUnauthorizedHandler(fn) {
  onUnauthorized = typeof fn === 'function' ? fn : () => {}
}

export const request = async (path, options = {}) => {
  const token = getStoredToken()
  const hadAuthHeader = Boolean(token)
  const response = await fetch(`${API_BASE}${path}`, {
    credentials: 'include',
    headers: {
      'Content-Type': 'application/json',
      ...(token ? { Authorization: `Bearer ${token}` } : {}),
      ...(options.headers ?? {}),
    },
    ...options,
  })

  const text = await response.text()
  let body = null
  if (text) {
    try {
      body = JSON.parse(text)
    } catch {
      body = null
    }
  }

  if (response.status === 401 && hadAuthHeader) {
    setStoredToken('')
    onUnauthorized()
  }

  return {
    ok: response.ok,
    status: response.status,
    body,
  }
}

export const extractError = (response, fallback) => response?.body?.response?.message || fallback

export const TOKEN_KEY = 'mustore_access_token'

export const getStoredToken = () => localStorage.getItem(TOKEN_KEY) ?? ''

export const setStoredToken = (token) => {
  if (token) {
    localStorage.setItem(TOKEN_KEY, token)
  } else {
    localStorage.removeItem(TOKEN_KEY)
  }
}

export const parseJwtPayload = (token) => {
  if (!token || typeof token !== 'string') {
    return null
  }
  const parts = token.split('.')
  if (parts.length < 2) {
    return null
  }
  try {
    const base64 = parts[1].replace(/-/g, '+').replace(/_/g, '/')
    const padded = base64.padEnd(base64.length + ((4 - (base64.length % 4)) % 4), '=')
    const json = atob(padded)
    return JSON.parse(json)
  } catch {
    return null
  }
}

export const isAccessTokenValid = (token) => {
  const payload = parseJwtPayload(token)
  if (!payload?.exp) {
    return false
  }
  return Date.now() < payload.exp * 1000
}

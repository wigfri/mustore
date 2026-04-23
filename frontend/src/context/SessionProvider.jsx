import { useCallback, useEffect, useMemo, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { getStoredToken, isAccessTokenValid, parseJwtPayload, setStoredToken } from '../lib/authStorage.js'
import { request, setUnauthorizedHandler } from '../lib/api.js'
import { SessionContext } from './sessionContext.js'

function readInitialSession() {
  const token = getStoredToken()
  if (!token || !isAccessTokenValid(token)) {
    if (token) {
      setStoredToken('')
    }
    return 'guest'
  }
  return 'user'
}

export function SessionProvider({ children }) {
  const navigate = useNavigate()
  const [session, setSession] = useState(readInitialSession)

  const markAuthenticated = useCallback(() => {
    setSession('user')
    navigate('/instruments', { replace: true })
  }, [navigate])

  const logout = useCallback(async () => {
    await request('/auth/logout', { method: 'POST' })
    setStoredToken('')
    setSession('guest')
    navigate('/login', { replace: true })
  }, [navigate])

  const unauthorized = useCallback(() => {
    setStoredToken('')
    setSession('guest')
    navigate('/login', { replace: true })
  }, [navigate])

  useEffect(() => {
    setUnauthorizedHandler(unauthorized)
    return () => setUnauthorizedHandler(() => {})
  }, [unauthorized])

  const value = useMemo(
    () => ({
      session,
      isAdmin: session === 'user' && parseJwtPayload(getStoredToken())?.role === 'admin',
      markAuthenticated,
      logout,
    }),
    [session, markAuthenticated, logout],
  )

  return <SessionContext.Provider value={value}>{children}</SessionContext.Provider>
}

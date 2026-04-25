import { Navigate, Outlet } from 'react-router-dom'
import { useSession } from '../hooks/useSession.js'

export function RequireAuth() {
  const { session } = useSession()
  if (session === 'guest') {
    return <Navigate to="/login" replace />
  }
  return <Outlet />
}

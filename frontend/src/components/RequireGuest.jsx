import { Navigate } from 'react-router-dom'
import { useSession } from '../hooks/useSession.js'

export function RequireGuest({ children }) {
  const { session } = useSession()
  if (session === 'user') {
    return <Navigate to="/instruments" replace />
  }
  return children
}

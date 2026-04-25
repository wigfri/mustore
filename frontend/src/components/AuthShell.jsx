import { API_BASE } from '../lib/api.js'
import { AuthTabs } from './AuthTabs.jsx'

export function AuthShell({ children }) {
  return (
    <div className="flex min-h-full flex-col items-center justify-center px-4 py-12">
      <div className="w-full max-w-md space-y-8 animate-page-in">
        <div className="text-center">
          <h1 className="text-2xl font-semibold tracking-tight text-neutral-900 transition hover:text-neutral-700">
            MuStore
          </h1>
          <p className="mt-1 text-sm text-neutral-500">Магазин музыкальных инструментов</p>
        </div>
        <AuthTabs />
        {children}
        <p className="text-center text-xs text-neutral-400">API: {API_BASE}</p>
      </div>
    </div>
  )
}

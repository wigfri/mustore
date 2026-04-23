import { NavLink, Outlet } from 'react-router-dom'
import { useSession } from '../hooks/useSession.js'
import { Button } from '../components/ui/Button.jsx'

const navClass = ({ isActive }) =>
  `rounded-md px-3 py-2 text-sm font-medium transition-all duration-200 ease-out ${
    isActive
      ? 'bg-neutral-900 text-white shadow-sm'
      : 'text-neutral-600 hover:bg-neutral-100 hover:text-neutral-900 hover:shadow-sm'
  }`

export function AppLayout() {
  const { logout, isAdmin } = useSession()

  return (
    <div className="min-h-full bg-neutral-50">
      <header className="border-b border-neutral-200 bg-white/95 shadow-sm backdrop-blur-sm transition-shadow duration-300 hover:shadow-md">
        <div className="mx-auto flex max-w-6xl flex-wrap items-center justify-between gap-4 px-4 py-4">
          <div className="flex items-center gap-6">
            <div>
              <p className="text-xs uppercase tracking-widest text-neutral-400">MuStore</p>
              <p className="text-sm text-neutral-500">Каталог</p>
            </div>
            <nav className="flex gap-1">
              <NavLink to="/instruments" className={navClass} end>
                Все товары
              </NavLink>
            </nav>
          </div>
          <div className="flex items-center gap-2">
            {isAdmin ? (
              <span className="rounded-md bg-neutral-100 px-2 py-1 text-xs text-neutral-600">Админ</span>
            ) : null}
            <Button variant="ghost" onClick={() => logout()}>
              Выйти
            </Button>
          </div>
        </div>
      </header>
      <Outlet />
    </div>
  )
}

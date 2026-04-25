import { NavLink } from 'react-router-dom'

const tabClass = ({ isActive }) =>
  `flex-1 rounded-md py-2 text-center text-sm font-medium transition-all duration-200 ease-out ${
    isActive
      ? 'bg-neutral-900 text-white shadow-sm'
      : 'text-neutral-600 hover:bg-neutral-50 hover:text-neutral-900'
  }`

export function AuthTabs() {
  return (
    <div className="flex rounded-lg border border-neutral-200 bg-white p-1 shadow-sm transition-shadow duration-300 hover:shadow-md">
      <NavLink to="/login" className={tabClass} end>
        Вход
      </NavLink>
      <NavLink to="/register" className={tabClass}>
        Регистрация
      </NavLink>
    </div>
  )
}

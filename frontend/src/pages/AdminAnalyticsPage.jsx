import { useEffect, useMemo, useState } from 'react'
import { Navigate } from 'react-router-dom'
import { useSession } from '../hooks/useSession.js'
import { request } from '../lib/api.js'

function MetricCard({ title, value, hint }) {
  return (
    <article className="rounded-xl border border-neutral-200 bg-white p-4 shadow-sm">
      <p className="text-xs uppercase tracking-wide text-neutral-500">{title}</p>
      <p className="mt-2 text-3xl font-semibold text-neutral-900">{value}</p>
      {hint ? <p className="mt-1 text-sm text-neutral-500">{hint}</p> : null}
    </article>
  )
}

export function AdminAnalyticsPage() {
  const { isAdmin } = useSession()
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [analytics, setAnalytics] = useState(null)

  useEffect(() => {
    if (!isAdmin) return
    let ignore = false
    const load = async () => {
      setLoading(true)
      setError('')
      const res = await request('/admin/analytics')
      if (ignore) return
      if (!res.ok) {
        setError(res?.body?.response?.message || 'Не удалось загрузить аналитику')
        setLoading(false)
        return
      }
      setAnalytics(res.body?.payload?.user ?? null)
      setLoading(false)
    }
    load()
    return () => {
      ignore = true
    }
  }, [isAdmin])

  const periodHint = useMemo(() => {
    if (!analytics?.period_start) return ''
    const start = new Date(analytics.period_start)
    return `Данные за ${start.toLocaleDateString()}`
  }, [analytics])

  if (!isAdmin) {
    return <Navigate to="/instruments" replace />
  }

  return (
    <main className="mx-auto max-w-6xl space-y-4 px-4 py-6">
      <section className="rounded-xl border border-neutral-200 bg-white p-5 shadow-sm">
        <h1 className="text-xl font-semibold text-neutral-900">Аналитика</h1>
        <p className="mt-1 text-sm text-neutral-500">Активность пользователей и ключевые метрики.</p>
      </section>

      {error ? <p className="rounded-md border border-red-200 bg-red-50 px-3 py-2 text-sm text-red-700">{error}</p> : null}

      {loading ? (
        <section className="rounded-xl border border-neutral-200 bg-white p-5 text-sm text-neutral-500 shadow-sm">Загрузка аналитики...</section>
      ) : null}

      {!loading && analytics ? (
        <section className="grid gap-4 md:grid-cols-2 xl:grid-cols-3">
          <MetricCard title="Всего пользователей" value={analytics.total_users} />
          <MetricCard title="Подтверждено email" value={analytics.verified_users} />
          <MetricCard title="Администраторы" value={analytics.admins_count} />
          <MetricCard title="Регистраций за день" value={analytics.registrations_day} hint={periodHint} />
          <MetricCard title="Входов за день" value={analytics.logins_day} hint={periodHint} />
          <MetricCard title="Активных пользователей за день" value={analytics.active_users_day} hint={periodHint} />
        </section>
      ) : null}
    </main>
  )
}


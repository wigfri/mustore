import { useCallback, useEffect, useMemo, useState } from 'react'
import { Link, useNavigate, useParams } from 'react-router-dom'
import { useSession } from '../hooks/useSession.js'
import { categoryLabel } from '../lib/constants.js'
import { formatPrice } from '../lib/format.js'
import { emptyInstrument, rowToForm } from '../lib/instrumentForm.js'
import { extractError, request } from '../lib/api.js'
import { InstrumentModal } from '../components/InstrumentModal.jsx'
import { Preloader } from '../components/Preloader.jsx'
import { Button } from '../components/ui/Button.jsx'

export function InstrumentDetailPage() {
  const { id } = useParams()
  const navigate = useNavigate()
  const { isAdmin } = useSession()
  const [item, setItem] = useState(null)
  const [busy, setBusy] = useState(true)
  const [error, setError] = useState('')
  const [editOpen, setEditOpen] = useState(false)

  const load = useCallback(async () => {
    if (!id) {
      return
    }
    setBusy(true)
    setError('')
    try {
      const res = await request(`/instruments/${id}`)
      if (!res.ok) {
        setItem(null)
        setError(extractError(res, 'Инструмент не найден'))
        return
      }
      setItem(res.body?.payload?.instrument ?? null)
    } catch {
      setError('Сеть недоступна')
      setItem(null)
    } finally {
      setBusy(false)
    }
  }, [id])

  useEffect(() => {
    load()
  }, [load])

  const editInitial = useMemo(() => rowToForm(item), [item])

  const handleUpdate = async (form) => {
    if (!item) {
      return
    }
    setBusy(true)
    const res = await request(`/instruments/${item.id}`, {
      method: 'PUT',
      body: JSON.stringify({ ...form, id: item.id }),
    })
    setBusy(false)
    if (!res.ok) {
      setError(extractError(res, 'Не удалось сохранить'))
      return
    }
    setEditOpen(false)
    load()
  }

  const handleDelete = async () => {
    if (!item || !window.confirm(`Удалить «${item.name}»?`)) {
      return
    }
    setBusy(true)
    const res = await request(`/instruments/${item.id}`, { method: 'DELETE' })
    setBusy(false)
    if (!res.ok) {
      setError(extractError(res, 'Не удалось удалить'))
      return
    }
    navigate('/instruments', { replace: true })
  }

  if (busy && !item) {
    return (
      <main className="mx-auto max-w-3xl px-4 py-8">
        <div className="rounded-xl border border-neutral-200 bg-white shadow-sm animate-fade-in">
          <Preloader label="Загружаем карточку…" />
        </div>
      </main>
    )
  }

  if (error && !item) {
    return (
      <main className="mx-auto max-w-3xl space-y-4 px-4 py-12 animate-page-in">
        <p className="text-sm text-red-600">{error}</p>
        <Link
          to="/instruments"
          className="inline-block rounded-md text-sm text-neutral-600 underline decoration-neutral-300 underline-offset-2 transition hover:bg-neutral-100 hover:text-neutral-900"
        >
          ← В каталог
        </Link>
      </main>
    )
  }

  if (!item) {
    return null
  }

  return (
    <main className="mx-auto max-w-3xl px-4 py-8 animate-page-in">
      <Link
        to="/instruments"
        className="mb-6 inline-block rounded-md px-1 py-1 text-sm text-neutral-500 transition-colors duration-200 hover:bg-neutral-100 hover:text-neutral-900"
      >
        ← Каталог
      </Link>

      <div className="overflow-hidden rounded-xl border border-neutral-200 bg-white shadow-sm transition-shadow duration-300 hover:shadow-md">
        <div className="aspect-[16/9] overflow-hidden bg-neutral-100 md:aspect-[21/9]">
          {item.image_url ? (
            <img
              src={item.image_url}
              alt=""
              className="h-full w-full object-cover transition duration-700 ease-out hover:scale-[1.02]"
            />
          ) : (
            <div className="flex h-full min-h-[200px] items-center justify-center text-neutral-400">Нет изображения</div>
          )}
        </div>
        <div className="space-y-4 p-6 md:p-8">
          <div className="flex flex-wrap items-start justify-between gap-4">
            <div>
              <p className="text-xs uppercase tracking-wide text-neutral-400">{categoryLabel(item.category)}</p>
              <h1 className="mt-1 text-2xl font-semibold tracking-tight text-neutral-900">{item.name}</h1>
              <p className="mt-1 text-neutral-600">{item.brand}</p>
              {item.type ? <p className="mt-2 text-sm text-neutral-500">{item.type}</p> : null}
            </div>
            <div className="text-right">
              <p className="text-xl font-semibold tabular-nums text-neutral-900">{formatPrice(item.price, item.currency)}</p>
              <p className="mt-1 text-sm text-neutral-500">
                В наличии: <span className="tabular-nums text-neutral-800">{item.stock}</span> шт.
              </p>
              <p className="mt-1 text-xs text-neutral-400">SKU: {item.sku}</p>
            </div>
          </div>

          {!item.is_active ? (
            <p className="rounded-md border border-amber-200 bg-amber-50 px-3 py-2 text-sm text-amber-800">
              Позиция скрыта в каталоге
            </p>
          ) : null}

          {item.description ? (
            <div>
              <h2 className="text-sm font-medium text-neutral-900">Описание</h2>
              <p className="mt-2 whitespace-pre-wrap text-sm leading-relaxed text-neutral-600">{item.description}</p>
            </div>
          ) : null}

          {isAdmin ? (
            <div className="flex flex-wrap gap-2 border-t border-neutral-100 pt-6">
              <Button variant="secondary" onClick={() => setEditOpen(true)} disabled={busy}>
                Редактировать
              </Button>
              <Button variant="ghost" className="text-red-600 hover:bg-red-50 hover:text-red-700" onClick={handleDelete} disabled={busy}>
                Удалить
              </Button>
            </div>
          ) : null}
        </div>
      </div>

      {error ? <p className="mt-4 text-sm text-red-600">{error}</p> : null}

      <InstrumentModal
        key={editOpen && item ? `detail-edit-${item.id}` : 'detail-edit-closed'}
        title="Редактирование"
        initial={item ? editInitial : emptyInstrument()}
        open={editOpen}
        busy={busy}
        onClose={() => setEditOpen(false)}
        onSave={handleUpdate}
      />
    </main>
  )
}

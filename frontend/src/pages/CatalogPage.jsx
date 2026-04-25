import { useCallback, useEffect, useMemo, useState } from 'react'
import { useSession } from '../hooks/useSession.js'
import { categories } from '../lib/constants.js'
import { emptyInstrument, rowToForm } from '../lib/instrumentForm.js'
import { extractError, request } from '../lib/api.js'
import { CatalogGridSkeleton, InlinePreloader } from '../components/Preloader.jsx'
import { InstrumentCard } from '../components/InstrumentCard.jsx'
import { InstrumentModal } from '../components/InstrumentModal.jsx'
import { Button } from '../components/ui/Button.jsx'
import { SelectField } from '../components/ui/SelectField.jsx'
import { TextField } from '../components/ui/TextField.jsx'

export function CatalogPage() {
  const { isAdmin } = useSession()
  const [items, setItems] = useState([])
  const [busy, setBusy] = useState(false)
  const [error, setError] = useState('')
  const [search, setSearch] = useState('')
  const [category, setCategory] = useState('')


  const [createOpen, setCreateOpen] = useState(false)
  const [editOpen, setEditOpen] = useState(false)
  const [editing, setEditing] = useState(null)

  const load = useCallback(async () => {
    setBusy(true)
    setError('')
    try {
      const params = new URLSearchParams()
      if (search) {
        params.set('search', search)
      }
      if (category) {
        params.set('category', category)
      }
      params.set('limit', '50')
      const res = await request(`/instruments/?${params.toString()}`)
      if (!res.ok) {
        setError(extractError(res, 'Не удалось загрузить каталог'))
        return
      }
      setItems(res.body?.payload?.instruments ?? [])
    } catch {
      setError('Сеть недоступна')
    } finally {
      setBusy(false)
    }
  }, [search, category])

  useEffect(() => {
    load()
  }, [load])

  const handleCreate = async (form) => {
    setBusy(true)
    const res = await request('/instruments/', {
      method: 'POST',
      body: JSON.stringify(form),
    })
    setBusy(false)
    if (!res.ok) {
      setError(extractError(res, 'Не удалось создать'))
      return
    }
    setCreateOpen(false)
    load()
  }

  const handleUpdate = async (form) => {
    if (!editing) {
      return
    }
    setBusy(true)
    const res = await request(`/instruments/${editing.id}`, {
      method: 'PUT',
      body: JSON.stringify({ ...form, id: editing.id }),
    })
    setBusy(false)
    if (!res.ok) {
      setError(extractError(res, 'Не удалось сохранить'))
      return
    }
    setEditOpen(false)
    setEditing(null)
    load()
  }

  const handleDelete = async (row) => {
    if (!window.confirm(`Удалить «${row.name}»?`)) {
      return
    }
    setBusy(true)
    const res = await request(`/instruments/${row.id}`, { method: 'DELETE' })
    setBusy(false)
    if (!res.ok) {
      setError(extractError(res, 'Не удалось удалить'))
      return
    }
    load()
  }

  const editInitial = useMemo(() => rowToForm(editing), [editing])

  return (
    <main className="mx-auto max-w-6xl space-y-6 px-4 py-8 animate-page-in">
      <div className="flex flex-wrap items-end gap-3 rounded-lg border border-neutral-200 bg-white p-4 shadow-sm transition-shadow duration-300 hover:shadow-md">
        <TextField
          label="Поиск"
          placeholder="Название или бренд"
          value={search}
          onChange={(e) => setSearch(e.target.value)}
          className="min-w-[200px] flex-1"
        />
        <SelectField label="Категория" value={category} onChange={(e) => setCategory(e.target.value)}>
          <option value="">Все</option>
          {categories.map((c) => (
            <option key={c.value} value={c.value}>
              {c.label}
            </option>
          ))}
        </SelectField>
        <Button variant="secondary" onClick={load} disabled={busy}>
          Обновить
        </Button>
        {isAdmin ? (
          <Button variant="primary" onClick={() => setCreateOpen(true)} disabled={busy}>
            Новый товар
          </Button>
        ) : null}
      </div>

      {error ? <p className="text-sm text-red-600">{error}</p> : null}

      {busy && items.length === 0 ? (
        <div className="rounded-xl border border-neutral-200 bg-white p-6 shadow-sm">
          <CatalogGridSkeleton count={6} />
        </div>
      ) : items.length === 0 ? (
        <p className="rounded-lg border border-dashed border-neutral-200 bg-white py-16 text-center text-sm text-neutral-500 transition-colors duration-300 hover:border-neutral-300 hover:bg-neutral-50/80">
          Пока нет позиций. Загрузите каталог или измените фильтры.
        </p>
      ) : (
        <div className="grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
          {items.map((row, index) => (
            <InstrumentCard
              key={row.id}
              item={row}
              animationDelay={Math.min(index, 14) * 40}
              isAdmin={isAdmin}
              onEdit={(item) => {
                setEditing(item)
                setEditOpen(true)
              }}
              onDelete={handleDelete}
            />
          ))}
        </div>
      )}

      {busy && items.length > 0 ? <InlinePreloader /> : null}

      <InstrumentModal
        key={createOpen ? 'create-open' : 'create-closed'}
        title="Новый инструмент"
        initial={emptyInstrument()}
        open={createOpen}
        busy={busy}
        onClose={() => setCreateOpen(false)}
        onSave={handleCreate}
      />
      <InstrumentModal
        key={editOpen && editing ? `edit-${editing.id}` : 'edit-closed'}
        title="Редактирование"
        initial={editInitial}
        open={editOpen}
        busy={busy}
        onClose={() => {
          setEditOpen(false)
          setEditing(null)
        }}
        onSave={handleUpdate}
      />
    </main>
  )
}

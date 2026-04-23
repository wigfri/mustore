import { useEffect, useState } from 'react'
import { createPortal } from 'react-dom'
import { categories, currencies } from '../lib/constants.js'
import { Button } from './ui/Button.jsx'
import { SelectField } from './ui/SelectField.jsx'
import { TextField } from './ui/TextField.jsx'

export function InstrumentModal({ title, initial, open, busy, onClose, onSave }) {
  const [form, setForm] = useState(initial)

  useEffect(() => {
    if (!open) {
      return undefined
    }
    const prev = document.body.style.overflow
    document.body.style.overflow = 'hidden'
    return () => {
      document.body.style.overflow = prev
    }
  }, [open])

  if (!open) {
    return null
  }

  const modal = (
    <div
      className="fixed inset-0 z-[200] flex animate-fade-in items-center justify-center bg-black/55 p-4 backdrop-blur-sm"
      role="presentation"
      onMouseDown={(e) => {
        if (e.target === e.currentTarget) {
          onClose()
        }
      }}
    >
      <div
        className="max-h-[min(90vh,calc(100dvh-2rem))] w-full max-w-lg animate-card-in overflow-y-auto rounded-xl border border-neutral-200 bg-white p-6 shadow-2xl"
        role="dialog"
        aria-modal="true"
        aria-labelledby="instrument-modal-title"
        onMouseDown={(e) => e.stopPropagation()}
      >
        <div className="mb-4 flex items-center justify-between">
          <h3 id="instrument-modal-title" className="text-lg font-medium">
            {title}
          </h3>
          <button
            type="button"
            className="flex h-9 w-9 items-center justify-center rounded-full text-lg text-neutral-400 transition hover:rotate-90 hover:bg-neutral-100 hover:text-neutral-700"
            onClick={onClose}
            aria-label="Закрыть"
          >
            ×
          </button>
        </div>
        <div className="grid gap-3">
          <TextField label="Название" value={form.name} onChange={(e) => setForm((f) => ({ ...f, name: e.target.value }))} />
          <TextField label="Бренд" value={form.brand} onChange={(e) => setForm((f) => ({ ...f, brand: e.target.value }))} />
          <SelectField
            label="Категория"
            value={form.category}
            onChange={(e) => setForm((f) => ({ ...f, category: e.target.value }))}
          >
            {categories.map((c) => (
              <option key={c.value} value={c.value}>
                {c.label}
              </option>
            ))}
          </SelectField>
          <TextField label="Тип" value={form.type} onChange={(e) => setForm((f) => ({ ...f, type: e.target.value }))} />
          <label className="flex flex-col gap-1 text-sm text-neutral-600">
            <span>Описание</span>
            <textarea
              className="min-h-[80px] rounded-md border border-neutral-200 bg-white px-3 py-2 text-neutral-900 outline-none focus:border-neutral-400 focus:ring-1 focus:ring-neutral-300"
              value={form.description}
              onChange={(e) => setForm((f) => ({ ...f, description: e.target.value }))}
            />
          </label>
          <div className="grid grid-cols-2 gap-3">
            <TextField
              label="Цена"
              type="number"
              min={0}
              value={form.price}
              onChange={(e) => setForm((f) => ({ ...f, price: Number(e.target.value) }))}
            />
            <SelectField
              label="Валюта"
              value={form.currency}
              onChange={(e) => setForm((f) => ({ ...f, currency: e.target.value }))}
            >
              {currencies.map((cur) => (
                <option key={cur} value={cur}>
                  {cur}
                </option>
              ))}
            </SelectField>
          </div>
          <TextField
            label="Остаток"
            type="number"
            min={0}
            value={form.stock}
            onChange={(e) => setForm((f) => ({ ...f, stock: Number(e.target.value) }))}
          />
          <TextField label="SKU" value={form.sku} onChange={(e) => setForm((f) => ({ ...f, sku: e.target.value }))} />
          <TextField
            label="Изображение (URL)"
            value={form.image_url}
            onChange={(e) => setForm((f) => ({ ...f, image_url: e.target.value }))}
          />
          <label className="flex items-center gap-2 text-sm text-neutral-600">
            <input
              type="checkbox"
              checked={form.is_active}
              onChange={(e) => setForm((f) => ({ ...f, is_active: e.target.checked }))}
            />
            Активен в каталоге
          </label>
        </div>
        <div className="mt-6 flex justify-end gap-2">
          <Button variant="secondary" onClick={onClose} disabled={busy}>
            Отмена
          </Button>
          <Button variant="primary" disabled={busy} onClick={() => onSave(form)}>
            Сохранить
          </Button>
        </div>
      </div>
    </div>
  )

  return createPortal(modal, document.body)
}

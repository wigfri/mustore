import { Link } from 'react-router-dom'
import { categoryLabel } from '../lib/constants.js'
import { formatPrice } from '../lib/format.js'

export function InstrumentCard({ item, isAdmin, onEdit, onDelete, animationDelay = 0 }) {
  return (
    <article
      className="group flex flex-col overflow-hidden rounded-xl border border-neutral-200 bg-white shadow-sm transition-all duration-300 ease-out animate-card-in hover:-translate-y-1 hover:border-neutral-300 hover:shadow-lg"
      style={{ animationDelay: `${animationDelay}ms` }}
    >
      <Link to={`/instruments/${item.id}`} className="flex flex-1 flex-col outline-none ring-neutral-900/10 transition focus-visible:ring-2">
        <div className="aspect-[4/3] overflow-hidden bg-neutral-100">
          {item.image_url ? (
            <img
              src={item.image_url}
              alt=""
              className="h-full w-full object-cover transition duration-500 ease-out group-hover:scale-[1.05]"
              loading="lazy"
            />
          ) : (
            <div className="flex h-full items-center justify-center text-xs text-neutral-400 transition group-hover:bg-neutral-200/50">
              Нет фото
            </div>
          )}
        </div>
        <div className="flex flex-1 flex-col gap-1 p-4">
          <p className="text-xs uppercase tracking-wide text-neutral-400 transition group-hover:text-neutral-500">
            {categoryLabel(item.category)}
          </p>
          <h2 className="font-medium leading-snug text-neutral-900 transition group-hover:text-neutral-800">{item.name}</h2>
          <p className="text-sm text-neutral-500 transition group-hover:text-neutral-600">{item.brand}</p>
          <p className="mt-auto pt-2 text-sm font-medium tabular-nums text-neutral-900">
            {formatPrice(item.price, item.currency)}
          </p>
          {!item.is_active ? <p className="text-xs text-amber-600">Не в продаже</p> : null}
        </div>
      </Link>
      {isAdmin ? (
        <div className="flex gap-2 border-t border-neutral-100 px-4 py-3">
          <button
            type="button"
            className="flex-1 rounded-md border border-neutral-200 py-1.5 text-xs text-neutral-700 transition hover:border-neutral-300 hover:bg-neutral-50 active:scale-[0.98]"
            onClick={(e) => {
              e.preventDefault()
              onEdit(item)
            }}
          >
            Изменить
          </button>
          <button
            type="button"
            className="flex-1 rounded-md border border-red-100 py-1.5 text-xs text-red-600 transition hover:border-red-200 hover:bg-red-50 active:scale-[0.98]"
            onClick={(e) => {
              e.preventDefault()
              onDelete(item)
            }}
          >
            Удалить
          </button>
        </div>
      ) : null}
    </article>
  )
}

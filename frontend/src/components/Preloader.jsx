/** Inline spinner + label. Use for catalog, detail page, and other async states. */
export function Preloader({ label = 'Загрузка…', size = 'default' }) {
  const dim = size === 'sm' ? 'h-9 w-9' : 'h-11 w-11'
  return (
    <div
      className="flex flex-col items-center justify-center gap-5 py-10"
      role="status"
      aria-live="polite"
      aria-busy="true"
    >
      <div className={`relative ${dim}`} aria-hidden>
        <div className="absolute inset-0 rounded-full border-2 border-neutral-200" />
        <div className="absolute inset-0 rounded-full border-2 border-transparent border-t-neutral-800 border-r-neutral-800/40 animate-spin-ring" />
      </div>
      <p className="text-sm text-neutral-500 animate-subtle-pulse">{label}</p>
    </div>
  )
}

/** Small row for background refetch (e.g. catalog already has items). */
export function InlinePreloader({ label = 'Обновление…' }) {
  return (
    <div className="flex items-center justify-center gap-3 py-3 text-neutral-500" role="status" aria-live="polite">
      <div className="relative h-5 w-5 shrink-0" aria-hidden>
        <div className="absolute inset-0 rounded-full border border-neutral-200" />
        <div className="absolute inset-0 rounded-full border border-transparent border-t-neutral-700 animate-spin-ring" />
      </div>
      <span className="text-sm animate-subtle-pulse">{label}</span>
    </div>
  )
}

/** Shimmer placeholders while the first catalog request is in flight. */
export function CatalogGridSkeleton({ count = 6 }) {
  return (
    <div className="grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
      {Array.from({ length: count }, (_, i) => (
        <div
          key={i}
          className="overflow-hidden rounded-xl border border-neutral-200 bg-white shadow-sm"
          style={{ animationDelay: `${i * 60}ms` }}
        >
          <div className="aspect-[4/3] skeleton-shimmer bg-neutral-100" />
          <div className="space-y-3 p-4">
            <div className="h-3 w-16 rounded bg-neutral-100 skeleton-shimmer" />
            <div className="h-4 w-[85%] rounded bg-neutral-100 skeleton-shimmer" />
            <div className="h-3 w-24 rounded bg-neutral-100 skeleton-shimmer" />
            <div className="h-4 w-28 rounded bg-neutral-100 skeleton-shimmer pt-1" />
          </div>
        </div>
      ))}
    </div>
  )
}

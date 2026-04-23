export function Button({ variant = 'primary', className = '', type = 'button', ...props }) {
  const base =
    'inline-flex items-center justify-center rounded-md px-4 py-2 text-sm font-medium transition duration-200 ease-out hover:-translate-y-0.5 active:translate-y-0 active:scale-[0.98] disabled:pointer-events-none disabled:opacity-50 disabled:hover:translate-y-0'
  const styles =
    variant === 'primary'
      ? 'bg-neutral-900 text-white shadow-sm hover:bg-neutral-800 hover:shadow-md'
      : variant === 'ghost'
        ? 'text-neutral-700 hover:bg-neutral-100 hover:shadow-sm'
        : 'border border-neutral-200 bg-white text-neutral-800 shadow-sm hover:border-neutral-300 hover:bg-neutral-50 hover:shadow-md'
  return <button type={type} className={`${base} ${styles} ${className}`} {...props} />
}

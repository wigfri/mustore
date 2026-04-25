export function SelectField({ label, children, ...props }) {
  return (
    <label className="flex flex-col gap-1 text-sm text-neutral-600">
      <span>{label}</span>
      <select
        className="rounded-md border border-neutral-200 bg-white px-3 py-2 text-neutral-900 outline-none transition duration-200 hover:border-neutral-300 focus:border-neutral-400 focus:ring-2 focus:ring-neutral-200/80"
        {...props}
      >
        {children}
      </select>
    </label>
  )
}

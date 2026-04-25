import { useMemo, useRef } from 'react'

const CODE_LENGTH = 6

export function CodeInput({ label, value, onChange, disabled = false, invalid = false }) {
  const inputRefs = useRef([])
  const chars = useMemo(() => {
    const normalized = String(value ?? '').replace(/\D/g, '').slice(0, CODE_LENGTH)
    return Array.from({ length: CODE_LENGTH }, (_, idx) => normalized[idx] ?? '')
  }, [value])

  const emit = (nextChars) => {
    onChange(nextChars.join('').replace(/\D/g, '').slice(0, CODE_LENGTH))
  }

  const focusAt = (idx) => {
    const input = inputRefs.current[idx]
    if (input) {
      input.focus()
      input.select()
    }
  }

  const handleChange = (idx, rawValue) => {
    const onlyDigits = rawValue.replace(/\D/g, '')
    if (!onlyDigits) {
      const nextChars = [...chars]
      nextChars[idx] = ''
      emit(nextChars)
      return
    }

    const nextChars = [...chars]
    const chunk = onlyDigits.slice(0, CODE_LENGTH - idx).split('')
    chunk.forEach((digit, chunkIdx) => {
      nextChars[idx + chunkIdx] = digit
    })
    emit(nextChars)

    const nextFocus = Math.min(idx + chunk.length, CODE_LENGTH - 1)
    focusAt(nextFocus)
  }

  const handlePaste = (idx, event) => {
    event.preventDefault()
    const pasted = event.clipboardData.getData('text').replace(/\D/g, '')
    if (!pasted) return

    const nextChars = [...chars]
    pasted
      .slice(0, CODE_LENGTH - idx)
      .split('')
      .forEach((digit, offset) => {
        nextChars[idx + offset] = digit
      })
    emit(nextChars)

    const nextFocus = Math.min(idx + pasted.length, CODE_LENGTH - 1)
    focusAt(nextFocus)
  }

  const handleKeyDown = (idx, event) => {
    if (event.key === 'Backspace' && !chars[idx] && idx > 0) {
      const nextChars = [...chars]
      nextChars[idx - 1] = ''
      emit(nextChars)
      focusAt(idx - 1)
      return
    }
    if (event.key === 'ArrowLeft' && idx > 0) {
      event.preventDefault()
      focusAt(idx - 1)
      return
    }
    if (event.key === 'ArrowRight' && idx < CODE_LENGTH - 1) {
      event.preventDefault()
      focusAt(idx + 1)
    }
  }

  return (
    <fieldset className="space-y-2" disabled={disabled}>
      <legend className="text-sm text-neutral-600">{label}</legend>
      <div className="flex items-center gap-2">
        {chars.map((char, idx) => {
          const isMiddleGap = idx === 3
          return (
            <input
              key={idx}
              ref={(el) => {
                inputRefs.current[idx] = el
              }}
              type="text"
              inputMode="numeric"
              autoComplete="one-time-code"
              maxLength={1}
              value={char}
              onChange={(event) => handleChange(idx, event.target.value)}
              onPaste={(event) => handlePaste(idx, event)}
              onKeyDown={(event) => handleKeyDown(idx, event)}
              className={`h-11 w-10 rounded-lg border bg-white text-center text-lg font-semibold text-neutral-900 outline-none transition focus:ring-2 ${invalid ? 'border-red-400 focus:border-red-500 focus:ring-red-100' : 'border-neutral-300 focus:border-neutral-700 focus:ring-neutral-200'} ${isMiddleGap ? 'ml-2' : ''}`}
              aria-label={`${label} ${idx + 1}`}
              disabled={disabled}
            />
          )
        })}
      </div>
    </fieldset>
  )
}

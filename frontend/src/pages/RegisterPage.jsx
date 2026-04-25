import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { extractError, request } from '../lib/api.js'
import { Button } from '../components/ui/Button.jsx'
import { CodeInput } from '../components/ui/CodeInput.jsx'
import { TextField } from '../components/ui/TextField.jsx'

export function RegisterPage() {
  const navigate = useNavigate()
  const [busy, setBusy] = useState(false)
  const [message, setMessage] = useState('')
  const [error, setError] = useState('')
  const [register, setRegister] = useState({
    username: '',
    email: '',
    password: '',
  })
  const [registerPhase, setRegisterPhase] = useState('form')
  const [verify, setVerify] = useState({ email: '', code: '' })

  const clearAlerts = () => {
    setMessage('')
    setError('')
  }

  const handleRegister = async (e) => {
    e.preventDefault()
    clearAlerts()
    setBusy(true)
    try {
      const res = await request('/user/', {
        method: 'POST',
        body: JSON.stringify({ ...register, role: 'base_user' }),
      })
      if (!res.ok) {
        setError(extractError(res, 'Регистрация не удалась'))
        return
      }
      setVerify((v) => ({ ...v, email: register.email }))
      setMessage('Аккаунт создан. Подтвердите email — код отправлен на почту.')
      setRegisterPhase('verify')
    } catch {
      setError('Сеть недоступна')
    } finally {
      setBusy(false)
    }
  }

  const handleVerify = async (e) => {
    e.preventDefault()
    clearAlerts()
    if (verify.code.length < 6) {
      setError('Введите код из 6 цифр')
      return
    }
    setBusy(true)
    try {
      const res = await request('/auth/verify-email', {
        method: 'POST',
        body: JSON.stringify(verify),
      })
      if (!res.ok) {
        setError(extractError(res, 'Неверный код или email'))
        return
      }
      setMessage('Email подтверждён. Теперь можно войти.')
      setRegisterPhase('form')
      navigate('/login', { replace: true, state: { email: verify.email } })
    } catch {
      setError('Сеть недоступна')
    } finally {
      setBusy(false)
    }
  }

  const handleResend = async () => {
    clearAlerts()
    setBusy(true)
    try {
      const res = await request('/auth/resend-verification', {
        method: 'POST',
        body: JSON.stringify({ email: verify.email }),
      })
      if (!res.ok) {
        setError(extractError(res, 'Не удалось отправить код'))
        return
      }
      setMessage('Код отправлен повторно.')
    } catch {
      setError('Сеть недоступна')
    } finally {
      setBusy(false)
    }
  }

  return (
    <>
      {message ? (
        <p className="animate-fade-in rounded-md border border-emerald-200 bg-emerald-50 px-3 py-2 text-sm text-emerald-800">
          {message}
        </p>
      ) : null}
      {error ? (
        <p className="animate-fade-in rounded-md border border-red-200 bg-red-50 px-3 py-2 text-sm text-red-800">{error}</p>
      ) : null}

      {registerPhase === 'form' ? (
        <form
          onSubmit={handleRegister}
          className="space-y-4 rounded-xl border border-neutral-200 bg-white p-6 shadow-sm transition-shadow duration-300 hover:shadow-md"
        >
          <TextField
            label="Имя пользователя"
            required
            value={register.username}
            onChange={(e) => setRegister((s) => ({ ...s, username: e.target.value }))}
          />
          <TextField
            label="Email"
            type="email"
            required
            value={register.email}
            onChange={(e) => setRegister((s) => ({ ...s, email: e.target.value }))}
          />
          <TextField
            label="Пароль"
            type="password"
            required
            value={register.password}
            onChange={(e) => setRegister((s) => ({ ...s, password: e.target.value }))}
          />
          <Button variant="primary" type="submit" className="w-full" disabled={busy}>
            Создать аккаунт
          </Button>
        </form>
      ) : null}

      {registerPhase === 'verify' ? (
        <form
          onSubmit={handleVerify}
          className="space-y-4 rounded-xl border border-neutral-200 bg-white p-6 shadow-sm transition-shadow duration-300 hover:shadow-md"
        >
          <p className="text-sm text-neutral-600">Подтверждение email</p>
          <TextField
            label="Email"
            type="email"
            required
            value={verify.email}
            onChange={(e) => setVerify((s) => ({ ...s, email: e.target.value }))}
          />
          <CodeInput
            label="Код"
            value={verify.code}
            onChange={(nextCode) => setVerify((s) => ({ ...s, code: nextCode }))}
            disabled={busy}
          />
          <div className="flex gap-2">
            <Button variant="primary" type="submit" className="flex-1" disabled={busy}>
              Подтвердить
            </Button>
            <Button variant="secondary" className="flex-1" disabled={busy} onClick={handleResend}>
              Ещё раз
            </Button>
          </div>
          <button
            type="button"
            className="text-sm text-neutral-500 underline decoration-neutral-300 underline-offset-2 hover:text-neutral-800"
            onClick={() => {
              clearAlerts()
              setRegisterPhase('form')
            }}
          >
            ← Назад к форме
          </button>
        </form>
      ) : null}
    </>
  )
}

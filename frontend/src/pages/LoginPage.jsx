import { useEffect, useState } from 'react'
import { useLocation } from 'react-router-dom'
import { useSession } from '../hooks/useSession.js'
import { extractError, request } from '../lib/api.js'
import { setStoredToken } from '../lib/authStorage.js'
import { Button } from '../components/ui/Button.jsx'
import { CodeInput } from '../components/ui/CodeInput.jsx'
import { TextField } from '../components/ui/TextField.jsx'

export function LoginPage() {
  const { markAuthenticated } = useSession()
  const location = useLocation()
  const [busy, setBusy] = useState(false)
  const [message, setMessage] = useState('')
  const [error, setError] = useState('')
  const [codeInvalid, setCodeInvalid] = useState(false)
  const [login, setLogin] = useState({ email: '', password: '', code: '' })
  const [loginStep, setLoginStep] = useState('password')

  useEffect(() => {
    const email = location.state?.email
    if (email && typeof email === 'string') {
      setLogin((s) => ({ ...s, email }))
    }
  }, [location.state])

  const clearAlerts = () => {
    setMessage('')
    setError('')
    setCodeInvalid(false)
  }

  const handleLoginSubmit = async (e) => {
    e.preventDefault()
    clearAlerts()
    if (loginStep === 'code' && login.code.length < 6) {
      setCodeInvalid(true)
      return
    }
    setBusy(true)
    try {
      const payload =
        loginStep === 'code'
          ? { email: login.email, password: login.password, code: login.code }
          : { email: login.email, password: login.password }
      const res = await request('/auth/login', {
        method: 'POST',
        body: JSON.stringify(payload),
      })
      if (!res.ok) {
        if (loginStep === 'code') {
          setCodeInvalid(true)
          return
        }
        setError(extractError(res, 'Не удалось войти'))
        return
      }
      const data = res.body?.payload ?? {}
      if (data.requires_code) {
        setLoginStep('code')
        setMessage('Код отправлен на почту. Введите его ниже.')
        return
      }
      const token = data.access_token
      if (token) {
        setStoredToken(token)
      }
      markAuthenticated()
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

      <form
        onSubmit={handleLoginSubmit}
        className="space-y-4 rounded-xl border border-neutral-200 bg-white p-6 shadow-sm transition-shadow duration-300 hover:shadow-md"
      >
        {loginStep !== 'code' ? (
          <>
            <TextField
              label="Email"
              type="email"
              autoComplete="email"
              required
              value={login.email}
              onChange={(e) => setLogin((s) => ({ ...s, email: e.target.value }))}
            />
            <TextField
              label="Пароль"
              type="password"
              autoComplete="current-password"
              required
              value={login.password}
              onChange={(e) => setLogin((s) => ({ ...s, password: e.target.value }))}
            />
          </>
        ) : null}
        {loginStep === 'code' ? (
          <CodeInput
            label="Код из письма"
            value={login.code}
            onChange={(nextCode) => {
              setCodeInvalid(false)
              setLogin((s) => ({ ...s, code: nextCode }))
            }}
            disabled={busy}
            invalid={codeInvalid}
          />
        ) : null}
        <Button variant="primary" type="submit" className="w-full" disabled={busy}>
          {loginStep === 'code' ? 'Войти' : 'Получить код'}
        </Button>
      </form>
    </>
  )
}

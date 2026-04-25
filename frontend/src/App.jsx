import { Route, Routes } from 'react-router-dom'
import { AuthShell } from './components/AuthShell.jsx'
import { RequireAuth } from './components/RequireAuth.jsx'
import { RequireGuest } from './components/RequireGuest.jsx'
import { RootRedirect } from './components/RootRedirect.jsx'
import { AppLayout } from './layouts/AppLayout.jsx'
import { AdminAnalyticsPage } from './pages/AdminAnalyticsPage.jsx'
import { CatalogPage } from './pages/CatalogPage.jsx'
import { InstrumentDetailPage } from './pages/InstrumentDetailPage.jsx'
import { LoginPage } from './pages/LoginPage.jsx'
import { RegisterPage } from './pages/RegisterPage.jsx'

export default function App() {
  return (
    <Routes>
      <Route path="/" element={<RootRedirect />} />
      <Route
        path="/login"
        element={
          <RequireGuest>
            <AuthShell>
              <LoginPage />
            </AuthShell>
          </RequireGuest>
        }
      />
      <Route
        path="/register"
        element={
          <RequireGuest>
            <AuthShell>
              <RegisterPage />
            </AuthShell>
          </RequireGuest>
        }
      />
      <Route element={<RequireAuth />}>
        <Route element={<AppLayout />}>
          <Route path="/instruments" element={<CatalogPage />} />
          <Route path="/instruments/:id" element={<InstrumentDetailPage />} />
          <Route path="/admin/analytics" element={<AdminAnalyticsPage />} />
        </Route>
      </Route>
      <Route path="*" element={<RootRedirect />} />
    </Routes>
  )
}

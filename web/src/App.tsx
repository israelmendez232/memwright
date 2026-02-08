import { Routes, Route, Navigate } from 'react-router-dom'
import { Toaster } from '@/components/ui/sonner'
import { MainLayout } from '@/components/layout'
import { ProtectedRoute } from '@/components/auth'
import { DashboardPage, LoginPage, RegisterPage } from '@/pages'

function App() {
  return (
    <>
      <Routes>
        <Route path="/login" element={<LoginPage />} />
        <Route path="/register" element={<RegisterPage />} />

        <Route
          path="/"
          element={
            <ProtectedRoute>
              <MainLayout>
                <DashboardPage />
              </MainLayout>
            </ProtectedRoute>
          }
        />

        <Route
          path="/decks"
          element={
            <ProtectedRoute>
              <MainLayout>
                <div className="text-muted-foreground">Decks page coming soon</div>
              </MainLayout>
            </ProtectedRoute>
          }
        />

        <Route
          path="/study"
          element={
            <ProtectedRoute>
              <MainLayout>
                <div className="text-muted-foreground">Study page coming soon</div>
              </MainLayout>
            </ProtectedRoute>
          }
        />

        <Route
          path="/stats"
          element={
            <ProtectedRoute>
              <MainLayout>
                <div className="text-muted-foreground">Statistics page coming soon</div>
              </MainLayout>
            </ProtectedRoute>
          }
        />

        <Route path="*" element={<Navigate to="/" replace />} />
      </Routes>
      <Toaster />
    </>
  )
}

export default App

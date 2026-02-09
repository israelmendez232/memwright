import { AuthLayout } from '@/components/layout'
import { LoginForm } from '@/components/auth'

const isRegistrationEnabled =
  import.meta.env.VITE_REGISTRATION_ENABLED === 'true'

export function LoginPage() {
  return (
    <AuthLayout>
      <LoginForm showRegisterLink={isRegistrationEnabled} />
    </AuthLayout>
  )
}

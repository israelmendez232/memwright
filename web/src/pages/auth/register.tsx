import { Link } from 'react-router-dom'
import { AuthLayout } from '@/components/layout'
import { RegisterForm } from '@/components/auth'
import { Button } from '@/components/ui/button'

const isRegistrationEnabled =
  import.meta.env.VITE_REGISTRATION_ENABLED === 'true'

export function RegisterPage() {
  if (!isRegistrationEnabled) {
    return (
      <AuthLayout>
        <div className="space-y-4 text-center">
          <h1 className="text-2xl font-semibold tracking-tight">
            Registration Disabled
          </h1>
          <p className="text-sm text-muted-foreground">
            New account registration is currently disabled.
          </p>
          <Link to="/login">
            <Button variant="outline" className="mt-4">
              Back to Sign in
            </Button>
          </Link>
        </div>
      </AuthLayout>
    )
  }

  return (
    <AuthLayout>
      <RegisterForm />
    </AuthLayout>
  )
}

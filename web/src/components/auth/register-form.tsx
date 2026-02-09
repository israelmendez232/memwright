import { useState } from 'react'
import { Link } from 'react-router-dom'
import { toast } from 'sonner'
import { useAuth } from '@/context'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Loader2 } from 'lucide-react'

interface FormErrors {
  name?: string
  email?: string
  password?: string
  confirmPassword?: string
}

interface RegisterFormProps {
  showLoginLink?: boolean
}

export function RegisterForm({ showLoginLink = true }: RegisterFormProps) {
  const { register } = useAuth()
  const [name, setName] = useState('')
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [confirmPassword, setConfirmPassword] = useState('')
  const [errors, setErrors] = useState<FormErrors>({})
  const [isSubmitting, setIsSubmitting] = useState(false)

  const validateName = (value: string): string | undefined => {
    if (!value.trim()) {
      return 'Display name is required'
    }
    return undefined
  }

  const validateEmail = (value: string): string | undefined => {
    if (!value.trim()) {
      return 'Email is required'
    }
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
    if (!emailRegex.test(value)) {
      return 'Please enter a valid email address'
    }
    return undefined
  }

  const validatePassword = (value: string): string | undefined => {
    if (!value) {
      return 'Password is required'
    }
    if (value.length < 8) {
      return 'Password must be at least 8 characters'
    }
    return undefined
  }

  const validateConfirmPassword = (
    value: string,
    passwordValue: string
  ): string | undefined => {
    if (!value) {
      return 'Please confirm your password'
    }
    if (value !== passwordValue) {
      return 'Passwords do not match'
    }
    return undefined
  }

  const validateForm = (): boolean => {
    const newErrors: FormErrors = {
      name: validateName(name),
      email: validateEmail(email),
      password: validatePassword(password),
      confirmPassword: validateConfirmPassword(confirmPassword, password),
    }
    setErrors(newErrors)
    return (
      !newErrors.name &&
      !newErrors.email &&
      !newErrors.password &&
      !newErrors.confirmPassword
    )
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()

    if (!validateForm()) {
      return
    }

    setIsSubmitting(true)

    try {
      await register({ name, email, password })
    } catch {
      toast.error('Failed to create account. Please try again.')
    } finally {
      setIsSubmitting(false)
    }
  }

  const handleNameChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value
    setName(value)
    if (errors.name) {
      setErrors((prev) => ({ ...prev, name: validateName(value) }))
    }
  }

  const handleEmailChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value
    setEmail(value)
    if (errors.email) {
      setErrors((prev) => ({ ...prev, email: validateEmail(value) }))
    }
  }

  const handlePasswordChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value
    setPassword(value)
    if (errors.password) {
      setErrors((prev) => ({ ...prev, password: validatePassword(value) }))
    }
    if (errors.confirmPassword && confirmPassword) {
      setErrors((prev) => ({
        ...prev,
        confirmPassword: validateConfirmPassword(confirmPassword, value),
      }))
    }
  }

  const handleConfirmPasswordChange = (
    e: React.ChangeEvent<HTMLInputElement>
  ) => {
    const value = e.target.value
    setConfirmPassword(value)
    if (errors.confirmPassword) {
      setErrors((prev) => ({
        ...prev,
        confirmPassword: validateConfirmPassword(value, password),
      }))
    }
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      <div className="space-y-2 text-center">
        <h1 className="text-2xl font-semibold tracking-tight">
          Create an account
        </h1>
        <p className="text-sm text-muted-foreground">
          Enter your details to get started
        </p>
      </div>

      <div className="space-y-4">
        <div className="space-y-2">
          <label htmlFor="name" className="text-sm font-medium">
            Name
          </label>
          <Input
            id="name"
            type="text"
            placeholder="Your name"
            value={name}
            onChange={handleNameChange}
            autoComplete="name"
            disabled={isSubmitting}
            aria-invalid={!!errors.name}
            aria-describedby={errors.name ? 'name-error' : undefined}
          />
          {errors.name && (
            <p id="name-error" className="text-sm text-destructive">
              {errors.name}
            </p>
          )}
        </div>

        <div className="space-y-2">
          <label htmlFor="email" className="text-sm font-medium">
            Email
          </label>
          <Input
            id="email"
            type="email"
            placeholder="you@example.com"
            value={email}
            onChange={handleEmailChange}
            autoComplete="email"
            disabled={isSubmitting}
            aria-invalid={!!errors.email}
            aria-describedby={errors.email ? 'email-error' : undefined}
          />
          {errors.email && (
            <p id="email-error" className="text-sm text-destructive">
              {errors.email}
            </p>
          )}
        </div>

        <div className="space-y-2">
          <label htmlFor="password" className="text-sm font-medium">
            Password
          </label>
          <Input
            id="password"
            type="password"
            placeholder="At least 8 characters"
            value={password}
            onChange={handlePasswordChange}
            autoComplete="new-password"
            disabled={isSubmitting}
            aria-invalid={!!errors.password}
            aria-describedby={errors.password ? 'password-error' : undefined}
          />
          {errors.password && (
            <p id="password-error" className="text-sm text-destructive">
              {errors.password}
            </p>
          )}
        </div>

        <div className="space-y-2">
          <label htmlFor="confirmPassword" className="text-sm font-medium">
            Confirm Password
          </label>
          <Input
            id="confirmPassword"
            type="password"
            placeholder="Confirm your password"
            value={confirmPassword}
            onChange={handleConfirmPasswordChange}
            autoComplete="new-password"
            disabled={isSubmitting}
            aria-invalid={!!errors.confirmPassword}
            aria-describedby={
              errors.confirmPassword ? 'confirmPassword-error' : undefined
            }
          />
          {errors.confirmPassword && (
            <p id="confirmPassword-error" className="text-sm text-destructive">
              {errors.confirmPassword}
            </p>
          )}
        </div>
      </div>

      <Button type="submit" className="w-full" disabled={isSubmitting}>
        {isSubmitting ? (
          <>
            <Loader2 className="mr-2 h-4 w-4 animate-spin" />
            Creating account...
          </>
        ) : (
          'Create account'
        )}
      </Button>

      {showLoginLink && (
        <p className="text-center text-sm text-muted-foreground">
          Already have an account?{' '}
          <Link
            to="/login"
            className="font-medium text-primary underline-offset-4 hover:underline"
          >
            Sign in
          </Link>
        </p>
      )}
    </form>
  )
}

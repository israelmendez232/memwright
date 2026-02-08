import { render, screen } from '@testing-library/react'
import { describe, it, expect } from 'vitest'
import { MemoryRouter } from 'react-router-dom'
import { AuthProvider } from '../../src/context'
import App from '../../src/App'

function renderApp(initialRoute = '/') {
  return render(
    <MemoryRouter initialEntries={[initialRoute]}>
      <AuthProvider>
        <App />
      </AuthProvider>
    </MemoryRouter>
  )
}

describe('App', () => {
  it('redirects unauthenticated users to login', () => {
    renderApp('/')
    expect(
      screen.getByRole('heading', { name: /welcome back/i })
    ).toBeInTheDocument()
  })

  it('renders login page at /login', () => {
    renderApp('/login')
    expect(
      screen.getByRole('heading', { name: /welcome back/i })
    ).toBeInTheDocument()
    expect(screen.getByLabelText(/email/i)).toBeInTheDocument()
    expect(screen.getByLabelText(/password/i)).toBeInTheDocument()
  })

  it('renders register page at /register', () => {
    renderApp('/register')
    expect(
      screen.getByRole('heading', { name: /create an account/i })
    ).toBeInTheDocument()
    expect(screen.getByLabelText(/name/i)).toBeInTheDocument()
  })

  it('shows sign in button on login page', () => {
    renderApp('/login')
    expect(
      screen.getByRole('button', { name: /sign in/i })
    ).toBeInTheDocument()
  })

  it('shows link to register from login page', () => {
    renderApp('/login')
    expect(screen.getByRole('link', { name: /sign up/i })).toBeInTheDocument()
  })
})

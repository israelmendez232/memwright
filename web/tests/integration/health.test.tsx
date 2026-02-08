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

describe('Application Health', () => {
  it('mounts without crashing', () => {
    const { container } = renderApp()
    expect(container).toBeTruthy()
  })

  it('renders login page for unauthenticated users', () => {
    renderApp('/')
    expect(screen.getByRole('heading', { name: /welcome back/i })).toBeInTheDocument()
  })

  it('renders auth layout on login page', () => {
    renderApp('/login')
    expect(screen.getByText('Memwright')).toBeInTheDocument()
    expect(screen.getByText('Spaced repetition learning platform')).toBeInTheDocument()
  })
})

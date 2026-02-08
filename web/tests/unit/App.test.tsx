import { render, screen } from '@testing-library/react'
import { describe, it, expect } from 'vitest'
import App from '../../src/App'

describe('App', () => {
  it('renders the app title', () => {
    render(<App />)
    expect(screen.getByRole('heading', { name: /memwright/i })).toBeInTheDocument()
  })

  it('renders the app description', () => {
    render(<App />)
    expect(screen.getByText(/spaced repetition learning platform/i)).toBeInTheDocument()
  })
})

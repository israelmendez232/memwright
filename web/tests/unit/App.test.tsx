import { render, screen } from '@testing-library/react'
import { describe, it, expect } from 'vitest'
import App from '../../src/App'

describe('App', () => {
  it('renders the welcome heading', () => {
    render(<App />)
    expect(
      screen.getByRole('heading', { name: /welcome to memwright/i })
    ).toBeInTheDocument()
  })

  it('renders the app description', () => {
    render(<App />)
    expect(
      screen.getByText(/your spaced repetition learning platform/i)
    ).toBeInTheDocument()
  })

  it('renders the header with app name', () => {
    render(<App />)
    expect(screen.getByText('Memwright')).toBeInTheDocument()
  })

  it('renders navigation elements', () => {
    render(<App />)
    const navElements = screen.getAllByRole('navigation')
    expect(navElements.length).toBeGreaterThan(0)
    expect(screen.getByRole('link', { name: /study/i })).toBeInTheDocument()
    expect(screen.getByRole('link', { name: /statistics/i })).toBeInTheDocument()
  })
})

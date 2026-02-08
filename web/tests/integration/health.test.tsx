import { render } from '@testing-library/react'
import { describe, it, expect } from 'vitest'
import App from '../../src/App'

describe('Application Health', () => {
  it('mounts without crashing', () => {
    const { container } = render(<App />)
    expect(container).toBeTruthy()
  })

  it('renders with expected structure', () => {
    const { container } = render(<App />)
    const appDiv = container.querySelector('.app')
    expect(appDiv).toBeInTheDocument()
  })
})

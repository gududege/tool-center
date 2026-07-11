import { describe, it, expect } from 'vitest'
import { render, screen } from '@testing-library/react'
import { Button } from './Button'

describe('Button', () => {
  it('renders with default variant and size', () => {
    render(<Button>Click me</Button>)

    const btn = screen.getByRole('button', { name: 'Click me' })
    expect(btn).toBeDefined()
    expect(btn.className).toContain('bg-neutral-900')
    expect(btn.className).toContain('h-9')
  })

  it('renders with custom variant and size', () => {
    render(
      <Button variant="destructive" size="sm">
        Delete
      </Button>,
    )

    const btn = screen.getByRole('button', { name: 'Delete' })
    expect(btn.className).toContain('bg-red-600')
    expect(btn.className).toContain('h-7')
  })

  it('forwards native button attributes', () => {
    render(
      <Button disabled type="submit">
        Submit
      </Button>,
    )

    const btn = screen.getByRole('button', { name: 'Submit' })
    expect((btn as HTMLButtonElement).disabled).toBe(true)
    expect((btn as HTMLButtonElement).type).toBe('submit')
  })

  it('renders secondary variant with correct styles', () => {
    render(<Button variant="secondary">Secondary</Button>)

    const btn = screen.getByRole('button', { name: 'Secondary' })
    expect(btn.className).toContain('bg-neutral-100')
    expect(btn.className).toContain('text-neutral-900')
  })
})

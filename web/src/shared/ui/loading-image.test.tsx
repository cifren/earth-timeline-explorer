import { fireEvent, render, screen } from '@testing-library/react'

import { LoadingImage } from './loading-image'

describe('LoadingImage', () => {
  it('shows a skeleton until the file is loaded', () => {
    render(<LoadingImage src="https://example.org/a.jpg" alt="Harbour" />)
    const image = screen.getByRole('img', { name: 'Harbour' })

    expect(screen.getByTestId('image-skeleton')).toBeInTheDocument()
    expect(image).toHaveAttribute('aria-busy', 'true')

    fireEvent.load(image)

    expect(screen.queryByTestId('image-skeleton')).not.toBeInTheDocument()
    expect(image).toHaveAttribute('aria-busy', 'false')
  })

  it('stops loading when the file fails', () => {
    render(<LoadingImage src="https://example.org/missing.jpg" alt="Missing" />)

    fireEvent.error(screen.getByRole('img', { name: 'Missing' }))

    expect(screen.queryByTestId('image-skeleton')).not.toBeInTheDocument()
  })

  it('shows the skeleton again for a new source', () => {
    const { rerender } = render(<LoadingImage src="https://example.org/a.jpg" alt="Photo" />)
    fireEvent.load(screen.getByRole('img', { name: 'Photo' }))

    rerender(<LoadingImage src="https://example.org/b.jpg" alt="Photo" />)

    expect(screen.getByTestId('image-skeleton')).toBeInTheDocument()
  })
})

import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'

import { wireImage } from '../../../test/api-fixtures'
import { toImage } from '../infrastructure/image-mapper'
import { PhotoDetails } from './photo-details'

describe('PhotoDetails', () => {
  const image = toImage(wireImage('Market day in Paris', '3 April 1953'))

  it('shows the photo with its date, zone, credits and links', () => {
    render(<PhotoDetails image={image} zoneName="Europe" onEnlarge={vi.fn()} />)

    expect(screen.getByRole('img', { name: 'Market day in Paris' })).toBeInTheDocument()
    expect(screen.getByText('3 April 1953')).toBeInTheDocument()
    expect(screen.getByText('Europe')).toBeInTheDocument()
    expect(screen.getByText('CC BY-SA 4.0')).toBeInTheDocument()
    expect(screen.getByText('J. Dupont')).toBeInTheDocument()
    expect(screen.getByText(/Source: Wikimedia Commons · original: Nationaal Archief/)).toBeInTheDocument()
    expect(screen.getByRole('link', { name: 'Paris' })).toHaveAttribute('href', 'https://en.wikipedia.org/wiki/Paris')
  })

  it('asks to enlarge the photo when clicked', async () => {
    const onEnlarge = vi.fn()
    render(<PhotoDetails image={image} zoneName="Europe" onEnlarge={onEnlarge} />)

    await userEvent.click(screen.getByTitle('Click to enlarge'))

    expect(onEnlarge).toHaveBeenCalledOnce()
  })
})

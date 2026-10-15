import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'

import { Timeline } from './timeline'

describe('Timeline', () => {
  it('shows the selected decade', () => {
    render(<Timeline decade={1950} onSelectDecade={vi.fn()} />)

    expect(screen.getByRole('slider', { name: 'Decade' })).toHaveAttribute('aria-valuetext', '1950s')
    expect(screen.getByText('1950–1959')).toBeInTheDocument()
  })

  it('moves by one decade with the arrow keys and jumps to the bounds', async () => {
    const onSelectDecade = vi.fn()
    render(<Timeline decade={1950} onSelectDecade={onSelectDecade} />)
    const slider = screen.getByRole('slider', { name: 'Decade' })

    slider.focus()
    await userEvent.keyboard('{ArrowRight}')
    await userEvent.keyboard('{ArrowLeft}')
    await userEvent.keyboard('{Home}')

    expect(onSelectDecade.mock.calls).toEqual([[1960], [1940], [1840]])
  })

  it('does not go before the first decade', async () => {
    const onSelectDecade = vi.fn()
    render(<Timeline decade={1840} onSelectDecade={onSelectDecade} />)

    screen.getByRole('slider', { name: 'Decade' }).focus()
    await userEvent.keyboard('{ArrowLeft}')

    expect(onSelectDecade).not.toHaveBeenCalled()
  })
})

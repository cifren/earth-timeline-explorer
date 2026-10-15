import { render, screen } from '@testing-library/react'
import { RouterProvider, createMemoryRouter } from 'react-router'

import { AppLayout } from './app-layout'

describe('AppLayout', () => {
  it('shows the logo, linking home, around the page', () => {
    const router = createMemoryRouter([{ Component: AppLayout, children: [{ path: '/', element: <p>Page</p> }] }])

    render(<RouterProvider router={router} />)

    expect(screen.getByText('Page')).toBeInTheDocument()
    expect(screen.getByRole('link', { name: 'Frise Chrono' })).toHaveAttribute('href', '/')
  })
})

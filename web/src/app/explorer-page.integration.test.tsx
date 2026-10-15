import { QueryClientProvider } from '@tanstack/react-query'
import { render, screen, within } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { HttpResponse, http } from 'msw'
import type { SetupServer } from 'msw/node'

import { INITIAL_DECADE, useSelection } from '@core/selection/selection-store'

import { wireImage } from '../test/api-fixtures'
import { API, startApi } from '../test/api-server'
import { ExplorerPage } from './explorer-page'
import { createQueryClient } from './query-client'

let server: SetupServer | undefined

const renderPage = (draws: Parameters<typeof startApi>[0]) => {
  const api = startApi(draws)
  server = api.server
  render(
    <QueryClientProvider client={createQueryClient()}>
      <ExplorerPage />
    </QueryClientProvider>,
  )
  return api
}

const selectZone = async (name: string) => {
  await userEvent.click(await screen.findByRole('button', { name }))
  return screen.findByRole('dialog')
}

beforeEach(() => {
  useSelection.setState({ zoneId: null, decade: INITIAL_DECADE })
})

afterEach(() => {
  server?.close()
})

describe('Explorer page', () => {
  it('draws a photo of the clicked zone for the selected decade', async () => {
    const { requests } = renderPage([[wireImage('Market day in Paris', '3 April 1953')]])

    const panel = await selectZone('Europe')

    expect(await within(panel).findByText('Market day in Paris')).toBeInTheDocument()
    expect(within(panel).getByText('3 April 1953')).toBeInTheDocument()
    expect(within(panel).getByRole('heading', { name: 'Europe · 1950s' })).toBeInTheDocument()
    expect(within(panel).getByRole('link', { name: 'OpenStreetMap' })).toHaveAttribute(
      'href',
      'https://www.openstreetmap.org/?mlat=48.8566&mlon=2.3522',
    )
    expect(requests[0]?.searchParams.toString()).toBe('zone=europe&decade=1950&limit=1')
  })

  it('draws another photo with Renew', async () => {
    renderPage([[wireImage('Market day in Paris', '3 April 1953')], [wireImage('Rotterdam harbour', 'around 1958')]])
    const panel = await selectZone('Europe')
    await within(panel).findByText('Market day in Paris')

    await userEvent.click(within(panel).getByRole('button', { name: 'Renew' }))

    expect(await within(panel).findByText('Rotterdam harbour')).toBeInTheDocument()
  })

  it('draws again for the same zone when the decade changes', async () => {
    const { requests } = renderPage([
      [wireImage('Market day in Paris', '3 April 1953')],
      [wireImage('Berlin wall', '1961')],
    ])
    const panel = await selectZone('Europe')
    await within(panel).findByText('Market day in Paris')

    screen.getByRole('slider', { name: 'Decade' }).focus()
    await userEvent.keyboard('{ArrowRight}')

    expect(await within(panel).findByText('Berlin wall')).toBeInTheDocument()
    expect(requests[1]?.searchParams.get('decade')).toBe('1960')
  })

  it('says when a zone has no photo for the decade', async () => {
    renderPage([[]])

    const panel = await selectZone('Africa')

    expect(await within(panel).findByText('No photo for Africa in the 1950s.')).toBeInTheDocument()
    expect(within(panel).getByRole('button', { name: 'Renew' })).toBeDisabled()
  })

  it('keeps the map usable when photos are unavailable', async () => {
    renderPage([])
    server?.use(
      http.get(`${API}/images`, () =>
        HttpResponse.json({ type: 'about:blank', title: 'Service Unavailable', status: 503 }, { status: 503 }),
      ),
    )

    const panel = await selectZone('Europe')

    expect(await within(panel).findByText('Photos are temporarily unavailable.')).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'Africa' })).toBeEnabled()
  })

  it('closes the panel and deselects the zone', async () => {
    renderPage([[wireImage('Market day in Paris', '3 April 1953')]])
    const panel = await selectZone('Europe')

    await userEvent.click(within(panel).getByRole('button', { name: 'Close' }))

    expect(screen.queryByRole('dialog')).not.toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'Europe' })).toHaveAttribute('aria-pressed', 'false')
  })
})

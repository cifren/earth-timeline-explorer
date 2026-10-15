import { HttpResponse, http } from 'msw'
import { setupServer } from 'msw/node'

import type { Image } from '@api-client/model'

import { zoneGeometry } from './api-fixtures'

export const API = 'http://api.test/api/v1'

/**
 * Fake back answering the zones and availability, drawing the given images in turn, and
 * fetching any of them again by identifier.
 */
export const startApi = (draws: Image[][]) => {
  const requests: URL[] = []
  const server = setupServer(
    http.get(`${API}/zones/geometry`, () => HttpResponse.json(zoneGeometry)),
    http.get(`${API}/decades/:decade/zones`, ({ params }) =>
      HttpResponse.json({ decade: Number(params.decade), zones: ['europe'] }),
    ),
    http.get(`${API}/images`, ({ request }) => {
      requests.push(new URL(request.url))
      return HttpResponse.json({ items: draws[requests.length - 1] ?? [] })
    }),
    http.get(`${API}/images/:imageId`, ({ params }) => {
      const image = draws.flat().find(({ id }) => id === decodeURIComponent(String(params.imageId)))
      return image
        ? HttpResponse.json(image)
        : HttpResponse.json({ type: 'about:blank', title: 'Not Found', status: 404 }, { status: 404 })
    }),
  )
  server.listen({ onUnhandledFrame: 'error' })
  return { server, requests }
}

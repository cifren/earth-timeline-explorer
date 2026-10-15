import type { Image, ZoneFeatureCollection } from '@api-client/model'

const square = (lon: number, lat: number) => [
  [
    [
      [lon, lat],
      [lon + 10, lat],
      [lon + 10, lat + 10],
      [lon, lat + 10],
      [lon, lat],
    ],
  ],
]

export const zoneGeometry: ZoneFeatureCollection = {
  type: 'FeatureCollection',
  features: [
    {
      type: 'Feature',
      properties: { id: 'europe', name: 'Europe' },
      geometry: { type: 'MultiPolygon', coordinates: square(0, 45) },
    },
    {
      type: 'Feature',
      properties: { id: 'africa', name: 'Africa' },
      geometry: { type: 'MultiPolygon', coordinates: square(10, 0) },
    },
  ],
}

export const wireImage = (title: string, dateLabel: string): Image => ({
  id: `wikimedia_commons:File:${title}.jpg`,
  title,
  kind: 'photograph',
  file: {
    thumbnailUrl: `https://upload.wikimedia.org/thumb/${title}.jpg`,
    displayUrl: `https://upload.wikimedia.org/large/${title}.jpg`,
    originalUrl: `https://upload.wikimedia.org/${title}.jpg`,
    width: 1600,
    height: 1200,
  },
  time: { edtf: '1953-04-03', label: dateLabel },
  location: {
    zone: 'europe',
    granularity: 'point',
    point: { type: 'Point', coordinates: [2.3522, 48.8566] },
    mapLinks: [{ label: 'OpenStreetMap', url: 'https://www.openstreetmap.org/?mlat=48.8566&mlon=2.3522' }],
  },
  certification: { level: 'source', authority: 'Wikimedia Commons', originalHolder: 'Nationaal Archief' },
  provenance: {
    source: 'wikimedia_commons',
    sourceId: `File:${title}.jpg`,
    sourceUrl: `https://commons.wikimedia.org/wiki/File:${title}.jpg`,
    license: 'CC BY-SA 4.0',
    author: 'J. Dupont',
    retrievedAt: '2026-10-15T10:00:00Z',
  },
  learnMore: [{ label: 'Paris', url: 'https://en.wikipedia.org/wiki/Paris' }],
})

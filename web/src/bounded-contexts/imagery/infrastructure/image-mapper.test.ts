import { wireImage } from '../../../test/api-fixtures'
import { toImage } from './image-mapper'

describe('toImage', () => {
  it('flattens what the panel needs to show and credit the photo', () => {
    const image = toImage(wireImage('Market day in Paris', '3 April 1953'))

    expect(image).toMatchObject({
      title: 'Market day in Paris',
      thumbnailUrl: 'https://upload.wikimedia.org/thumb/Market day in Paris.jpg',
      dateLabel: '3 April 1953',
      zoneId: 'europe',
      position: { longitude: 2.3522, latitude: 48.8566, precise: true },
      mapLinks: [{ label: 'OpenStreetMap', url: 'https://www.openstreetmap.org/?mlat=48.8566&mlon=2.3522' }],
      authority: 'Wikimedia Commons',
      originalHolder: 'Nationaal Archief',
      author: 'J. Dupont',
      license: 'CC BY-SA 4.0',
      learnMore: [{ label: 'Paris', url: 'https://en.wikipedia.org/wiki/Paris' }],
    })
  })

  it('keeps an imprecise position inside the country', () => {
    const wire = wireImage('Harbour', '1950s')
    const point = { type: 'Point' as const, coordinates: [9.02, 55.97] }
    const image = toImage({ ...wire, location: { ...wire.location, granularity: 'country', point } })

    expect(image.position).toEqual({ longitude: 9.02, latitude: 55.97, precise: false })
  })

  it('has no position when the country has no point', () => {
    const wire = wireImage('Harbour', '1950s')
    const image = toImage({ ...wire, location: { ...wire.location, granularity: 'country', point: undefined } })

    expect(image.position).toBeUndefined()
  })
})

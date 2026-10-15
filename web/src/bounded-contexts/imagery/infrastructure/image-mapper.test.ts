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
      authority: 'Wikimedia Commons',
      originalHolder: 'Nationaal Archief',
      author: 'J. Dupont',
      license: 'CC BY-SA 4.0',
      learnMore: [{ label: 'Paris', url: 'https://en.wikipedia.org/wiki/Paris' }],
    })
  })
})

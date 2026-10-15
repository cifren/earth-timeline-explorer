import { zoneGeometry } from '../../../test/api-fixtures'
import { toZoneOutlines } from './zone-outline-mapper'

describe('toZoneOutlines', () => {
  it('keeps the id, name and polygons of each zone', () => {
    const [europe] = toZoneOutlines(zoneGeometry)

    expect(europe?.id).toBe('europe')
    expect(europe?.name).toBe('Europe')
    expect(europe?.polygons[0]?.[0]?.[0]).toEqual([0, 45])
  })
})

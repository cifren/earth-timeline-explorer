import { isDimmed, zoneStyle } from './zone-style'

const state = { availableZoneIds: new Set(['europe']), selectedZoneId: null, hoveredZoneId: null }

describe('isDimmed', () => {
  it('dims zones without photos, and nothing while availability is unknown', () => {
    expect(isDimmed('europe', state.availableZoneIds)).toBe(false)
    expect(isDimmed('africa', state.availableZoneIds)).toBe(true)
    expect(isDimmed('africa', null)).toBe(false)
  })
})

describe('zoneStyle', () => {
  it('paints a dimmed zone grey and an available one with its colour', () => {
    expect(zoneStyle('africa', state).getFill()?.getColor()).toBe('#ddd2bd')
    expect(zoneStyle('europe', state).getFill()?.getColor()).toBe('#d4b07a')
  })

  it('outlines the selected zone with a thicker accent border', () => {
    const stroke = zoneStyle('europe', { ...state, selectedZoneId: 'europe' }).getStroke()

    expect(stroke?.getColor()).toBe('#a0461c')
    expect(stroke?.getWidth()).toBe(3)
  })
})

import type Circle from 'ol/style/Circle'

import { isDimmed, markerStyle, zoneStyle } from './zone-style'

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

describe('markerStyle', () => {
  it('draws a larger dot and halo when only the country is known', () => {
    const radii = (precise: boolean) => markerStyle(precise).map((style) => (style.getImage() as Circle).getRadius())

    expect(radii(false)[0]).toBeGreaterThan(radii(true)[0] ?? 0)
    expect(radii(false)[1]).toBeGreaterThan(radii(true)[1] ?? 0)
  })
})

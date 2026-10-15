import { Circle, Fill, Stroke, Style } from 'ol/style'

const ZONE_COLORS: Record<string, string> = {
  'northern-america': '#b3a27e',
  'latin-america-caribbean': '#a8a372',
  europe: '#d4b07a',
  africa: '#cf9a6b',
  'middle-east': '#b99282',
  'south-central-asia': '#dcc07a',
  'east-southeast-asia': '#a3a88a',
  oceania: '#c7957f',
  antarctica: '#efe6d2',
}
const DIMMED_COLOR = '#ddd2bd'
const BORDER_COLOR = '#4a3b2a'
const SELECTED_COLOR = '#a0461c'

const roundMark = (dotRadius: number, haloRadius: number) => [
  new Style({ image: new Circle({ radius: haloRadius, fill: new Fill({ color: 'rgba(184, 98, 42, 0.22)' }) }) }),
  new Style({
    image: new Circle({
      radius: dotRadius,
      fill: new Fill({ color: '#b8622a' }),
      stroke: new Stroke({ color: '#faf4e6', width: 2.5 }),
    }),
  }),
]

const PRECISE_MARK = roundMark(6, 13)
const COUNTRY_MARK = roundMark(9, 34)

/** Round mark of a position, a dot in a halo; both grow when only the country is known. */
export const markerStyle = (precise: boolean) => (precise ? PRECISE_MARK : COUNTRY_MARK)

export type ZoneStyleState = {
  /** Zones having photos for the decade; null while unknown, which dims nothing. */
  availableZoneIds: ReadonlySet<string> | null
  selectedZoneId: string | null
  hoveredZoneId: string | null
}

export const isDimmed = (zoneId: string, availableZoneIds: ReadonlySet<string> | null) =>
  availableZoneIds !== null && !availableZoneIds.has(zoneId)

export const zoneStyle = (zoneId: string, state: ZoneStyleState): Style => {
  const selected = zoneId === state.selectedZoneId
  const color = isDimmed(zoneId, state.availableZoneIds) ? DIMMED_COLOR : (ZONE_COLORS[zoneId] ?? DIMMED_COLOR)

  return new Style({
    fill: new Fill({ color: zoneId === state.hoveredZoneId ? `${color}bf` : color }),
    stroke: new Stroke({
      color: selected ? SELECTED_COLOR : BORDER_COLOR,
      width: selected ? 3 : 1.2,
      lineJoin: 'round',
    }),
  })
}

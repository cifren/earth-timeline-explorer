import { Fill, Stroke, Style } from 'ol/style'

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

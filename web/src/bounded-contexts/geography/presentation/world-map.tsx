import 'ol/ol.css'

import Feature from 'ol/Feature'
import Map from 'ol/Map'
import View from 'ol/View'
import { defaults as defaultControls } from 'ol/control/defaults'
import MultiPolygon from 'ol/geom/MultiPolygon'
import Point from 'ol/geom/Point'
import { containsCoordinate } from 'ol/extent'
import { defaults as defaultInteractions } from 'ol/interaction/defaults'
import Graticule from 'ol/layer/Graticule'
import VectorLayer from 'ol/layer/Vector'
import VectorSource from 'ol/source/Vector'
import { Stroke } from 'ol/style'
import { useEffect, useRef, useState } from 'react'

import type { ZoneOutline } from '../domain/zone-outline'
import { isDimmed, markerStyle, zoneStyle, type ZoneStyleState } from './zone-style'

// Far beyond one world horizontally, so that only latitude is bounded.
const VIEW_EXTENT = [-1e7, -90, 1e7, 90]

type Hover = { zoneId: string; x: number; y: number }

type Marker = { longitude: number; latitude: number; precise: boolean }

type WorldMapProps = {
  outlines: ZoneOutline[]
  availableZoneIds: ReadonlySet<string> | null
  selectedZoneId: string | null
  decade: number
  /** Position to point at, shown as a round mark, larger when imprecise. */
  marker?: Marker
  onSelectZone: (zoneId: string | null) => void
}

/** Flat world map (EPSG:4326) repeated horizontally, showing the zones to pick from. */
export const WorldMap = ({
  outlines,
  availableZoneIds,
  selectedZoneId,
  decade,
  marker,
  onSelectZone,
}: WorldMapProps) => {
  const targetRef = useRef<HTMLDivElement>(null)
  const [source] = useState(() => new VectorSource<Feature<MultiPolygon>>())
  const [markerSource] = useState(() => new VectorSource<Feature<Point>>())
  const layerRef = useRef<VectorLayer | null>(null)
  const mapRef = useRef<Map | null>(null)
  const styleStateRef = useRef<ZoneStyleState>({
    availableZoneIds,
    selectedZoneId,
    hoveredZoneId: null,
  })
  const onSelectZoneRef = useRef(onSelectZone)
  const [hover, setHover] = useState<Hover | null>(null)

  useEffect(() => {
    const target = targetRef.current
    if (target === null) {
      return
    }
    const layer = new VectorLayer({
      source,
      style: (feature) => zoneStyle(String(feature.getId()), styleStateRef.current),
    })
    layerRef.current = layer
    const map = new Map({
      target,
      layers: [
        new Graticule({
          strokeStyle: new Stroke({ color: '#cfba93', width: 1 }),
          wrapX: true,
        }),
        layer,
        new VectorLayer({ source: markerSource, style: (feature) => markerStyle(Boolean(feature.get('precise'))) }),
      ],
      view: new View({ projection: 'EPSG:4326', center: [15, 15], zoom: 2, extent: VIEW_EXTENT }),
      controls: defaultControls({ attribution: false, rotate: false }),
      interactions: defaultInteractions({ altShiftDragRotate: false, pinchRotate: false }),
    })
    mapRef.current = map
    const zoneAt = (pixel: number[]) =>
      map.forEachFeatureAtPixel(pixel, (feature) => String(feature.getId()), {
        layerFilter: (candidate) => candidate === layer,
      }) ?? null

    map.on('pointermove', (event) => {
      if (event.dragging) {
        return
      }
      const zoneId = zoneAt(event.pixel)
      const [x = 0, y = 0] = event.pixel
      target.style.cursor = zoneId === null ? '' : 'pointer'
      setHover(zoneId === null ? null : { zoneId, x, y })
      if (zoneId === styleStateRef.current.hoveredZoneId) {
        return
      }
      styleStateRef.current = { ...styleStateRef.current, hoveredZoneId: zoneId }
      layer.changed()
    })
    map.on('click', (event) => onSelectZoneRef.current(zoneAt(event.pixel)))
    const leave = () => {
      setHover(null)
      styleStateRef.current = { ...styleStateRef.current, hoveredZoneId: null }
      layer.changed()
    }
    target.addEventListener('pointerleave', leave)

    return () => {
      target.removeEventListener('pointerleave', leave)
      map.setTarget(undefined)
    }
  }, [source, markerSource])

  useEffect(() => {
    markerSource.clear()
    const view = mapRef.current?.getView()
    if (marker === undefined || view === undefined) {
      return
    }
    const coordinate = [marker.longitude, marker.latitude]
    markerSource.addFeature(new Feature({ geometry: new Point(coordinate), precise: marker.precise }))
    if (containsCoordinate(view.calculateExtent(), coordinate)) {
      return
    }
    view.animate({ center: coordinate, duration: 400 })
  }, [markerSource, marker])

  useEffect(() => {
    onSelectZoneRef.current = onSelectZone
  }, [onSelectZone])

  useEffect(() => {
    styleStateRef.current = { ...styleStateRef.current, availableZoneIds, selectedZoneId }
    layerRef.current?.changed()
  }, [availableZoneIds, selectedZoneId])

  useEffect(() => {
    const features = outlines.map((outline) => {
      const feature = new Feature({ geometry: new MultiPolygon(outline.polygons), name: outline.name })
      feature.setId(outline.id)
      return feature
    })
    source.clear()
    source.addFeatures(features)
  }, [source, outlines])

  const hoveredName = outlines.find((outline) => outline.id === hover?.zoneId)?.name

  return (
    <div className="world-map relative h-full w-full bg-sea">
      <div ref={targetRef} className="h-full w-full" />
      {hover && hoveredName && (
        <div
          className="pointer-events-none absolute rounded bg-night/90 px-2 py-1 text-[13px] whitespace-nowrap text-night-foreground"
          style={{ left: hover.x + 14, top: hover.y + 14 }}
        >
          {isDimmed(hover.zoneId, availableZoneIds) ? `${hoveredName} · no photo in the ${decade}s` : hoveredName}
        </div>
      )}
      <ul className="sr-only" aria-label="Zones">
        {outlines.map((outline) => (
          <li key={outline.id}>
            <button type="button" aria-pressed={outline.id === selectedZoneId} onClick={() => onSelectZone(outline.id)}>
              {outline.name}
            </button>
          </li>
        ))}
      </ul>
    </div>
  )
}

import type { ZoneFeatureCollection } from '@api-client/model'

import type { ZoneOutline } from '../domain/zone-outline'

export const toZoneOutlines = (collection: ZoneFeatureCollection): ZoneOutline[] =>
  collection.features.map((feature) => ({
    id: feature.properties.id,
    name: feature.properties.name,
    polygons: feature.geometry.coordinates,
  }))

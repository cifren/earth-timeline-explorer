import { drawImages, listZonesWithImages } from '@api-client/imagery/imagery'
import { ZoneId } from '@api-client/model'
import { ApiError } from '@shared/http/api-error'

import type { Image } from '../domain/image'
import { toImage } from './image-mapper'

const zoneIds: readonly string[] = Object.values(ZoneId)

const isZoneId = (value: string): value is ZoneId => zoneIds.includes(value)

/** Draws one random image of the zone and decade; null when none exists. */
export const drawImage = async (zoneId: string, decade: number, signal: AbortSignal): Promise<Image | null> => {
  if (!isZoneId(zoneId)) {
    throw new Error(`Unknown zone ${zoneId}`)
  }
  const response = await drawImages({ zone: zoneId, decade, limit: 1 }, { signal })
  if (response.status !== 200) {
    throw new ApiError(response.status, response.data)
  }
  const [image] = response.data.items
  return image === undefined ? null : toImage(image)
}

export const fetchZonesWithImages = async (decade: number, signal: AbortSignal): Promise<string[]> => {
  const response = await listZonesWithImages(decade, { signal })
  if (response.status !== 200) {
    throw new ApiError(response.status, response.data)
  }
  return response.data.zones
}

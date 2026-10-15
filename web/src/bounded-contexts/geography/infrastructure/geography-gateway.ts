import { getZoneGeometry } from '@api-client/geography/geography'
import { ApiError } from '@shared/http/api-error'

import type { ZoneOutline } from '../domain/zone-outline'
import { toZoneOutlines } from './zone-outline-mapper'

export const fetchZoneOutlines = async (signal: AbortSignal): Promise<ZoneOutline[]> => {
  const response = await getZoneGeometry({ signal })
  if (response.status !== 200) {
    throw new ApiError(response.status, response.data)
  }
  return toZoneOutlines(response.data)
}

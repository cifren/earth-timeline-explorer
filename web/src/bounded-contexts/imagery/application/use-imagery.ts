import { skipToken, useQuery } from '@tanstack/react-query'
import { useMemo } from 'react'

import { useSelection } from '@core/selection/selection-store'

import { drawImage, fetchZonesWithImages } from '../infrastructure/imagery-gateway'

const AVAILABILITY_STALE_MS = 60 * 60 * 1000

/**
 * Photos of the selected zone and decade: the drawn image, the zones having photos for
 * the decade, and the intents to draw again or to close.
 */
export const useImagery = () => {
  const zoneId = useSelection((state) => state.zoneId)
  const decade = useSelection((state) => state.decade)
  const selectZone = useSelection((state) => state.selectZone)

  const image = useQuery({
    queryKey: ['images', zoneId, decade],
    queryFn: zoneId === null ? skipToken : ({ signal }) => drawImage(zoneId, decade, signal),
    staleTime: Infinity,
  })
  const availability = useQuery({
    queryKey: ['zones-with-images', decade],
    queryFn: ({ signal }) => fetchZonesWithImages(decade, signal),
    staleTime: AVAILABILITY_STALE_MS,
  })
  const availableZoneIds = useMemo(
    () => (availability.data === undefined ? null : new Set(availability.data)),
    [availability.data],
  )

  return {
    zoneId,
    decade,
    image: image.data,
    isDrawing: image.isFetching,
    isUnavailable: image.isError,
    availableZoneIds,
    renew: () => void image.refetch(),
    close: () => selectZone(null),
  }
}

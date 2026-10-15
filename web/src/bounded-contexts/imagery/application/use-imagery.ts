import { skipToken, useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { useEffect, useEffectEvent, useMemo } from 'react'

import { useSelection } from '@core/selection/use-selection'

import type { Image } from '../domain/image'
import { drawImage, fetchImage, fetchZonesWithImages } from '../infrastructure/imagery-gateway'

const AVAILABILITY_STALE_MS = 60 * 60 * 1000

const imageKey = (imageId: string, decade: number) => ['image', imageId, decade]

/**
 * Photos of the selected zone and decade: the shown image, the zones having photos for
 * the decade, and the intents to draw again or to close. The shown image is the one of
 * the URL; without one, a photo is drawn and its identifier written into the URL, so that
 * every photo has an address and the browser history goes back through the photos seen.
 */
export const useImagery = () => {
  const { zoneId, decade, imageId, selectZone, showImage } = useSelection()
  const queryClient = useQueryClient()

  const show = (image: Image, options?: { replace: boolean }) => {
    queryClient.setQueryData(imageKey(image.id, decade), image)
    showImage(image.id, options)
  }

  const shown = useQuery({
    queryKey: imageKey(imageId ?? '', decade),
    queryFn: imageId === null ? skipToken : ({ signal }) => fetchImage(imageId, decade, signal),
    staleTime: Infinity,
  })
  const draw = useQuery({
    queryKey: ['images', zoneId, decade],
    queryFn: zoneId === null || imageId !== null ? skipToken : ({ signal }) => drawImage(zoneId, decade, signal),
    staleTime: Infinity,
  })
  const renew = useMutation({
    mutationFn: (zone: string) => drawImage(zone, decade),
    onSuccess: (image) => image && show(image),
  })

  const showDrawn = useEffectEvent((image: Image) => show(image, { replace: true }))
  useEffect(() => {
    if (imageId !== null || !draw.data) {
      return
    }
    showDrawn(draw.data)
  }, [imageId, draw.data])

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
    image: imageId === null ? draw.data : shown.data,
    isDrawing: shown.isFetching || draw.isFetching || renew.isPending,
    isUnavailable: shown.isError || draw.isError || renew.isError,
    availableZoneIds,
    renew: () => {
      if (zoneId === null) {
        return
      }
      renew.mutate(zoneId)
    },
    close: () => selectZone(null),
  }
}

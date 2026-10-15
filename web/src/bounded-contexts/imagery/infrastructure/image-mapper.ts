import type { Location, Image as WireImage } from '@api-client/model'

import type { Image, Position } from '../domain/image'

const toPosition = ({ point, granularity }: Location): Position | undefined => {
  if (point === undefined) {
    return undefined
  }
  const [longitude = 0, latitude = 0] = point.coordinates
  return { longitude, latitude, precise: granularity === 'point' }
}

export const toImage = (image: WireImage): Image => ({
  id: image.id,
  title: image.title,
  thumbnailUrl: image.file.thumbnailUrl,
  originalUrl: image.file.originalUrl,
  width: image.file.width,
  height: image.file.height,
  dateLabel: image.time.label,
  zoneId: image.location.zone,
  position: toPosition(image.location),
  mapLinks: image.location.mapLinks,
  authority: image.certification.authority,
  originalHolder: image.certification.originalHolder,
  author: image.provenance.author,
  license: image.provenance.license,
  sourceUrl: image.provenance.sourceUrl,
  learnMore: image.learnMore,
})

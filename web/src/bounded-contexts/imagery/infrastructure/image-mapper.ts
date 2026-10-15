import type { Image as WireImage } from '@api-client/model'

import type { Image } from '../domain/image'

export const toImage = (image: WireImage): Image => ({
  id: image.id,
  title: image.title,
  thumbnailUrl: image.file.thumbnailUrl,
  originalUrl: image.file.originalUrl,
  width: image.file.width,
  height: image.file.height,
  dateLabel: image.time.label,
  zoneId: image.location.zone,
  authority: image.certification.authority,
  originalHolder: image.certification.originalHolder,
  author: image.provenance.author,
  license: image.provenance.license,
  sourceUrl: image.provenance.sourceUrl,
  learnMore: image.learnMore,
})

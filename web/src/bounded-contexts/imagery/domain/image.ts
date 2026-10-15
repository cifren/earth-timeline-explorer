export type Link = { label: string; url: string }

/** WGS 84 position, in decimal degrees: where the photo was taken when precise, otherwise inside its country. */
export type Position = { longitude: number; latitude: number; precise: boolean }

/** A photograph drawn for a zone and a decade, with what is needed to credit it. */
export type Image = {
  id: string
  title: string
  thumbnailUrl: string
  originalUrl: string
  width: number
  height: number
  dateLabel: string
  zoneId: string
  position?: Position
  /** The place on map services. */
  mapLinks: Link[]
  /** Who vouches for the image, e.g. "Wikimedia Commons". */
  authority: string
  originalHolder?: string
  author?: string
  license: string
  sourceUrl: string
  learnMore: Link[]
}

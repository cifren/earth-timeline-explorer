export type Link = { label: string; url: string }

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
  /** Who vouches for the image, e.g. "Wikimedia Commons". */
  authority: string
  originalHolder?: string
  author?: string
  license: string
  sourceUrl: string
  learnMore: Link[]
}

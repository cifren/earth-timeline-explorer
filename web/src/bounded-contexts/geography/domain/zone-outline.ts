/** A [longitude, latitude] position in WGS 84. */
export type Position = number[]

/** The outline of a zone: polygons, each made of rings of positions. */
export type ZoneOutline = {
  id: string
  name: string
  polygons: Position[][][]
}

export type Size = { width: number; height: number }
export type Point = { x: number; y: number }

/** Scale and top-left position of an image inside a frame. */
export type ZoomView = Point & { scale: number }

export const MAX_ZOOM = 8
const FIT_WIDTH_RATIO = 0.9
const FIT_HEIGHT_RATIO = 0.82

/** Scale at which the whole image fits the frame, leaving room for the caption and tools. */
export const fitScale = (image: Size, frame: Size) =>
  Math.min((frame.width * FIT_WIDTH_RATIO) / image.width, (frame.height * FIT_HEIGHT_RATIO) / image.height)

const clampAxis = (position: number, imageLength: number, frameLength: number) => {
  if (imageLength <= frameLength) {
    return (frameLength - imageLength) / 2
  }
  return Math.min(0, Math.max(frameLength - imageLength, position))
}

/** Centres an image smaller than the frame, and keeps a larger one covering it. */
export const clampView = (view: ZoomView, image: Size, frame: Size): ZoomView => ({
  scale: view.scale,
  x: clampAxis(view.x, image.width * view.scale, frame.width),
  y: clampAxis(view.y, image.height * view.scale, frame.height),
})

export const fitView = (image: Size, frame: Size): ZoomView =>
  clampView({ scale: fitScale(image, frame), x: 0, y: 0 }, image, frame)

/** Zooms by a factor around a point of the frame, between the fitted size and MAX_ZOOM. */
export const zoomAt = (view: ZoomView, factor: number, point: Point, image: Size, frame: Size): ZoomView => {
  const fit = fitScale(image, frame)
  const scale = Math.min(Math.max(view.scale * factor, fit), fit * MAX_ZOOM)
  const ratio = scale / view.scale
  return clampView(
    { scale, x: point.x - (point.x - view.x) * ratio, y: point.y - (point.y - view.y) * ratio },
    image,
    frame,
  )
}

export const panBy = (view: ZoomView, dx: number, dy: number, image: Size, frame: Size): ZoomView =>
  clampView({ scale: view.scale, x: view.x + dx, y: view.y + dy }, image, frame)

export const isZoomed = (view: ZoomView, image: Size, frame: Size) => view.scale > fitScale(image, frame) * 1.001

export const zoomLevel = (view: ZoomView, image: Size, frame: Size) =>
  Math.round((view.scale / fitScale(image, frame)) * 100)

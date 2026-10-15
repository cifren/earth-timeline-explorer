import { MAX_ZOOM, fitView, isZoomed, panBy, zoomAt, zoomLevel } from './zoom-pan'

const image = { width: 2000, height: 1000 }
const frame = { width: 1000, height: 1000 }

describe('zoom and pan', () => {
  it('fits the image centred in the frame', () => {
    const view = fitView(image, frame)

    expect(view.scale).toBeCloseTo(0.45)
    expect(view.x).toBeCloseTo(50)
    expect(view.y).toBeCloseTo(275)
    expect(zoomLevel(view, image, frame)).toBe(100)
  })

  it('zooms around a point, between the fitted size and the maximum', () => {
    const fitted = fitView(image, frame)

    expect(zoomAt(fitted, 0.5, { x: 500, y: 500 }, image, frame).scale).toBeCloseTo(fitted.scale)
    expect(zoomLevel(zoomAt(fitted, 100, { x: 500, y: 500 }, image, frame), image, frame)).toBe(MAX_ZOOM * 100)

    const zoomed = zoomAt(fitted, 3, { x: 500, y: 500 }, image, frame)
    expect(isZoomed(zoomed, image, frame)).toBe(true)
    expect((500 - zoomed.x) / zoomed.scale).toBeCloseTo((500 - fitted.x) / fitted.scale)
  })

  it('pans a zoomed image without uncovering the frame', () => {
    const zoomed = zoomAt(fitView(image, frame), 4, { x: 500, y: 500 }, image, frame)

    const panned = panBy(zoomed, 10_000, 10_000, image, frame)

    expect(panned.x).toBe(0)
    expect(panned.y).toBe(0)
  })
})

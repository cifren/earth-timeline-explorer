import { type MouseEvent, type PointerEvent, type WheelEvent, useEffect, useRef, useState } from 'react'

import { cn } from 'cn'
import { Button } from '@shared/ui/button'
import { Dialog, DialogContent, DialogTitle } from '@shared/ui/dialog'

import type { Image } from '../domain/image'
import { type Point, type Size, fitView, isZoomed, panBy, zoomAt, zoomLevel } from './zoom-pan'

const WHEEL_FACTOR = 1.2
const BUTTON_FACTOR = 1.5
const DOUBLE_CLICK_FACTOR = 3

const windowSize = (): Size => ({ width: window.innerWidth, height: window.innerHeight })

type PhotoViewerProps = {
  image: Image
  onClose: () => void
}

/** Full-screen photo with zoom (wheel, buttons, double click) and pan when zoomed. */
export const PhotoViewer = ({ image, onClose }: PhotoViewerProps) => {
  const size: Size = { width: image.width, height: image.height }
  const [frame, setFrame] = useState(windowSize)
  const [view, setView] = useState(() => fitView(size, frame))
  const [animate, setAnimate] = useState(false)
  const dragRef = useRef<Point | null>(null)
  const zoomed = isZoomed(view, size, frame)

  useEffect(() => {
    const resize = () => {
      const next = windowSize()
      setFrame(next)
      setView(fitView({ width: image.width, height: image.height }, next))
    }
    window.addEventListener('resize', resize)
    return () => window.removeEventListener('resize', resize)
  }, [image.width, image.height])

  const zoom = (factor: number, point: Point, animated: boolean) => {
    setAnimate(animated)
    setView((current) => zoomAt(current, factor, point, size, frame))
  }
  const center = { x: frame.width / 2, y: frame.height / 2 }
  const fit = () => {
    setAnimate(true)
    setView(fitView(size, frame))
  }

  const handleWheel = (event: WheelEvent) =>
    zoom(event.deltaY < 0 ? WHEEL_FACTOR : 1 / WHEEL_FACTOR, { x: event.clientX, y: event.clientY }, false)

  const handleDoubleClick = (event: MouseEvent) => {
    if (zoomed) {
      fit()
      return
    }
    zoom(DOUBLE_CLICK_FACTOR, { x: event.clientX, y: event.clientY }, true)
  }

  const handlePointerDown = (event: PointerEvent<HTMLImageElement>) => {
    event.preventDefault()
    event.currentTarget.setPointerCapture(event.pointerId)
    dragRef.current = { x: event.clientX, y: event.clientY }
  }

  const handlePointerMove = (event: PointerEvent<HTMLImageElement>) => {
    const drag = dragRef.current
    if (drag === null || !zoomed) {
      return
    }
    dragRef.current = { x: event.clientX, y: event.clientY }
    setAnimate(false)
    setView((current) => panBy(current, event.clientX - drag.x, event.clientY - drag.y, size, frame))
  }

  const handleBackdropClick = (event: MouseEvent) => {
    if (event.target !== event.currentTarget) {
      return
    }
    onClose()
  }

  return (
    <Dialog open onOpenChange={onClose}>
      <DialogContent
        aria-describedby={undefined}
        className="top-0 left-0 block h-screen w-screen max-w-none translate-x-0 translate-y-0 overflow-hidden rounded-none border-0 bg-night/95 p-0 text-night-foreground shadow-none sm:max-w-none"
        onWheel={handleWheel}
        onClick={handleBackdropClick}
      >
        <DialogTitle className="absolute top-4 left-5 text-sm font-normal">{image.title}</DialogTitle>
        <img
          src={image.originalUrl}
          alt={image.title}
          draggable={false}
          className={cn(
            'absolute top-0 left-0 max-w-none origin-top-left shadow-2xl select-none',
            zoomed ? 'cursor-grab active:cursor-grabbing' : 'cursor-zoom-in',
            animate && 'transition-transform duration-200 ease-out',
          )}
          style={{
            width: image.width,
            height: image.height,
            transform: `translate(${view.x}px, ${view.y}px) scale(${view.scale})`,
          }}
          onDoubleClick={handleDoubleClick}
          onPointerDown={handlePointerDown}
          onPointerMove={handlePointerMove}
          onPointerUp={() => {
            dragRef.current = null
          }}
        />
        <div className="absolute bottom-5 left-1/2 flex -translate-x-1/2 items-center gap-1.5 rounded-md bg-night/90 p-1.5">
          <Button
            variant="secondary"
            size="sm"
            aria-label="Zoom out"
            onClick={() => zoom(1 / BUTTON_FACTOR, center, true)}
          >
            −
          </Button>
          <span className="min-w-14 text-center text-[13px]" aria-live="polite">
            {zoomLevel(view, size, frame)}%
          </span>
          <Button variant="secondary" size="sm" aria-label="Zoom in" onClick={() => zoom(BUTTON_FACTOR, center, true)}>
            +
          </Button>
          <Button variant="secondary" size="sm" onClick={fit}>
            Fit
          </Button>
        </div>
      </DialogContent>
    </Dialog>
  )
}

import { XIcon } from 'lucide-react'
import { type MouseEvent, type PointerEvent, type WheelEvent, useEffect, useRef, useState } from 'react'

import { cn } from 'cn'
import { Button } from '@shared/ui/button'
import { Dialog, DialogClose, DialogContent, DialogTitle } from '@shared/ui/dialog'
import { Skeleton } from '@shared/ui/skeleton'
import { useImageLoading } from '@shared/ui/use-image-loading'

import type { Image } from '../domain/image'
import { PhotoInfo } from './photo-details'
import { PhotoSide } from './photo-side'
import { type Point, type Size, fitView, isZoomed, panBy, zoomAt, zoomLevel } from './zoom-pan'

const WHEEL_FACTOR = 1.2
const BUTTON_FACTOR = 1.5
const DOUBLE_CLICK_FACTOR = 3
const DETAILS_WIDTH = 360

const stageSize = (): Size => ({ width: window.innerWidth - DETAILS_WIDTH, height: window.innerHeight })

type PhotoViewerProps = {
  image: Image
  zoneName: string | undefined
  decade: number
  isDrawing: boolean
  onRenew: () => void
  onClose: () => void
}

/**
 * Full-screen photo with zoom (wheel, buttons, double click) and pan when zoomed, next to
 * its details on the right. The photo fits the stage left of the details and, once zoomed,
 * overflows above them. Renew draws the next photo in place, shown fitted.
 */
export const PhotoViewer = ({ image, zoneName, decade, isDrawing, onRenew, onClose }: PhotoViewerProps) => {
  const size: Size = { width: image.width, height: image.height }
  const [frame, setFrame] = useState(stageSize)
  const [view, setView] = useState(() => fitView(size, frame))
  const [animate, setAnimate] = useState(false)
  const [shownImageId, setShownImageId] = useState(image.id)
  const dragRef = useRef<Point | null>(null)
  const zoomed = isZoomed(view, size, frame)
  const { loading, settle } = useImageLoading(image.displayUrl)

  if (shownImageId !== image.id) {
    setShownImageId(image.id)
    setAnimate(false)
    setView(fitView(size, frame))
  }

  useEffect(() => {
    const resize = () => {
      const next = stageSize()
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
        showCloseButton={false}
        className="top-0 left-0 block h-screen w-screen max-w-none translate-x-0 translate-y-0 overflow-hidden rounded-none border-0 bg-night/95 p-0 text-night-foreground shadow-none sm:max-w-none"
      >
        <div
          role="presentation"
          className="absolute inset-y-0 left-0"
          style={{ right: DETAILS_WIDTH }}
          onWheel={handleWheel}
          onClick={handleBackdropClick}
        >
          <DialogTitle className="absolute top-4 left-4 z-30 max-w-[calc(100%-2rem)] truncate rounded-lg bg-night/90 px-3.5 py-2 text-base font-semibold shadow-lg ring-1 ring-night-foreground/30">
            {image.title}
          </DialogTitle>
          {loading && (
            <Skeleton
              data-testid="viewer-skeleton"
              className="absolute z-20 rounded-none bg-night-foreground/15"
              style={{ left: view.x, top: view.y, width: size.width * view.scale, height: size.height * view.scale }}
            />
          )}
          <img
            src={image.displayUrl}
            alt={image.title}
            draggable={false}
            className={cn(
              'absolute top-0 left-0 z-20 max-w-none origin-top-left shadow-2xl select-none',
              zoomed ? 'cursor-grab active:cursor-grabbing' : 'cursor-zoom-in',
              animate && 'transition-transform duration-200 ease-out',
              loading && 'opacity-0',
            )}
            aria-busy={loading}
            onLoad={settle}
            onError={settle}
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
          <div className="absolute bottom-5 left-1/2 z-30 flex -translate-x-1/2 items-center gap-1.5 rounded-md bg-night/90 p-1.5">
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
            <Button
              variant="secondary"
              size="sm"
              aria-label="Zoom in"
              onClick={() => zoom(BUTTON_FACTOR, center, true)}
            >
              +
            </Button>
            <Button variant="secondary" size="sm" onClick={fit}>
              Fit
            </Button>
          </div>
        </div>
        <aside
          aria-label="Photo details"
          className="absolute inset-y-0 right-0 flex flex-col border-l bg-background text-foreground"
          style={{ width: DETAILS_WIDTH }}
        >
          <PhotoSide zoneName={zoneName} decade={decade} renewDisabled={isDrawing} onRenew={onRenew}>
            <PhotoInfo image={image} zoneName={zoneName} />
          </PhotoSide>
        </aside>
        <DialogClose asChild>
          <Button
            variant="secondary"
            size="icon-lg"
            aria-label="Close"
            className="absolute top-4 right-4 z-30 rounded-full shadow-lg ring-2 ring-night-foreground/70 hover:scale-105 [&_svg:not([class*='size-'])]:size-6"
          >
            <XIcon />
          </Button>
        </DialogClose>
      </DialogContent>
    </Dialog>
  )
}

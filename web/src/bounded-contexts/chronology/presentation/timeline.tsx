import { type KeyboardEvent, type PointerEvent, type WheelEvent, useRef, useState } from 'react'

import { cn } from 'cn'

import { FIRST_DECADE, clampDecade, clampYear, currentDecade, decadeOf } from '../domain/decade'

const PX_PER_YEAR = 14
const DRAG_THRESHOLD_PX = 4
const WHEEL_SNAP_DELAY_MS = 180

const offsetOf = (year: number) => (year - FIRST_DECADE) * PX_PER_YEAR

type Drag = { x: number; startX: number; moved: boolean }

type TimelineProps = {
  decade: number
  onSelectDecade: (decade: number) => void
}

/**
 * Strip graduated in years, dragged under a fixed central cursor. The decade under the
 * cursor is highlighted while moving and selected when the strip is released.
 */
export const Timeline = ({ decade, onSelectDecade }: TimelineProps) => {
  const lastDecade = currentDecade()
  const [centerYear, setCenterYear] = useState(decade + 5)
  const [shownDecade, setShownDecade] = useState(decade)
  const [snapping, setSnapping] = useState(false)
  const [dragging, setDragging] = useState(false)
  const dragRef = useRef<Drag | null>(null)
  const wheelTimerRef = useRef<ReturnType<typeof setTimeout>>(undefined)

  if (decade !== shownDecade) {
    setShownDecade(decade)
    setCenterYear(decade + 5)
  }

  const liveDecade = decadeOf(centerYear)
  const years = Array.from({ length: lastDecade + 10 - FIRST_DECADE + 1 }, (_, i) => FIRST_DECADE + i)

  const snapTo = (target: number) => {
    setSnapping(true)
    setCenterYear(target + 5)
    if (target === decade) {
      return
    }
    onSelectDecade(target)
  }

  const moveTo = (year: number) => {
    setSnapping(false)
    setCenterYear(clampYear(year, lastDecade))
  }

  const handlePointerDown = (event: PointerEvent<HTMLDivElement>) => {
    event.currentTarget.setPointerCapture(event.pointerId)
    dragRef.current = { x: event.clientX, startX: event.clientX, moved: false }
  }

  const handlePointerMove = (event: PointerEvent<HTMLDivElement>) => {
    const drag = dragRef.current
    if (drag === null) {
      return
    }
    const dx = event.clientX - drag.x
    drag.x = event.clientX
    drag.moved ||= Math.abs(event.clientX - drag.startX) > DRAG_THRESHOLD_PX
    if (!drag.moved) {
      return
    }
    setDragging(true)
    moveTo(centerYear - dx / PX_PER_YEAR)
  }

  const handlePointerUp = (event: PointerEvent<HTMLDivElement>) => {
    const drag = dragRef.current
    dragRef.current = null
    setDragging(false)
    if (drag === null || drag.moved) {
      snapTo(liveDecade)
      return
    }
    const rect = event.currentTarget.getBoundingClientRect()
    const clicked = clampYear(centerYear + (event.clientX - rect.left - rect.width / 2) / PX_PER_YEAR, lastDecade)
    snapTo(decadeOf(clicked))
  }

  const handleWheel = (event: WheelEvent<HTMLDivElement>) => {
    const year = clampYear(centerYear + (event.deltaX || event.deltaY) / PX_PER_YEAR, lastDecade)
    moveTo(year)
    clearTimeout(wheelTimerRef.current)
    wheelTimerRef.current = setTimeout(() => snapTo(decadeOf(year)), WHEEL_SNAP_DELAY_MS)
  }

  const handleKeyDown = (event: KeyboardEvent<HTMLDivElement>) => {
    const targets: Record<string, number> = {
      ArrowLeft: decade - 10,
      ArrowRight: decade + 10,
      Home: FIRST_DECADE,
      End: lastDecade,
    }
    const target = targets[event.key]
    if (target === undefined) {
      return
    }
    event.preventDefault()
    snapTo(clampDecade(target, lastDecade))
  }

  return (
    <div
      role="slider"
      tabIndex={0}
      aria-label="Decade"
      aria-valuemin={FIRST_DECADE}
      aria-valuemax={lastDecade}
      aria-valuenow={decade}
      aria-valuetext={`${decade}s`}
      className={cn(
        'relative h-26 shrink-0 touch-none overflow-hidden bg-night text-night-foreground select-none',
        dragging ? 'cursor-grabbing' : 'cursor-grab',
      )}
      onPointerDown={handlePointerDown}
      onPointerMove={handlePointerMove}
      onPointerUp={handlePointerUp}
      onPointerCancel={handlePointerUp}
      onWheel={handleWheel}
      onKeyDown={handleKeyDown}
    >
      <div
        className={cn('absolute top-0 left-1/2 h-full', snapping && 'transition-transform duration-200 ease-out')}
        style={{ width: offsetOf(lastDecade + 10), transform: `translateX(${-offsetOf(centerYear)}px)` }}
      >
        <div
          className="absolute bottom-7.5 h-8 rounded-sm border border-highlight bg-highlight/30"
          style={{ left: offsetOf(liveDecade), width: 10 * PX_PER_YEAR }}
        />
        {years.map((year) => (
          <div key={year}>
            <div
              className={cn('absolute bottom-8.5 w-px', year % 10 === 0 ? 'h-6 bg-tick' : 'h-2.5 bg-tick-minor')}
              style={{ left: offsetOf(year) }}
            />
            {year % 10 === 0 && (
              <div
                className="absolute bottom-15.5 -translate-x-1/2 text-[13px] text-tick"
                style={{ left: offsetOf(year) }}
              >
                {year}
              </div>
            )}
          </div>
        ))}
      </div>
      <div className="pointer-events-none absolute top-2 bottom-6.5 left-1/2 -ml-px w-0.5 bg-highlight" />
      <div className="pointer-events-none absolute bottom-1 left-1/2 -translate-x-1/2 text-sm font-semibold">
        {liveDecade}–{liveDecade + 9}
      </div>
    </div>
  )
}

import { Suspense, lazy, useState } from 'react'

import { Sheet, SheetContent, SheetDescription, SheetFooter, SheetHeader, SheetTitle } from '@shared/ui/sheet'
import { Skeleton } from '@shared/ui/skeleton'

import { useImagery } from '../application/use-imagery'
import type { Image } from '../domain/image'
import { PhotoDetails } from './photo-details'
import { RenewButton } from './renew-button'

const PhotoViewer = lazy(() => import('./photo-viewer').then((module) => ({ default: module.PhotoViewer })))

type PanelBodyProps = {
  image: Image | null | undefined
  isDrawing: boolean
  isUnavailable: boolean
  zoneName: string | undefined
  decade: number
  onEnlarge: () => void
}

const PanelBody = ({ image, isDrawing, isUnavailable, zoneName, decade, onEnlarge }: PanelBodyProps) => {
  if (isDrawing) {
    return (
      <div className="space-y-3" aria-label="Loading photo">
        <Skeleton className="aspect-4/3 w-full" />
        <Skeleton className="h-5 w-3/4" />
        <Skeleton className="h-4 w-1/3" />
      </div>
    )
  }
  if (isUnavailable) {
    return <p className="mt-16 text-center text-muted-foreground">Photos are temporarily unavailable.</p>
  }
  if (!image) {
    return (
      <p className="mt-16 text-center text-muted-foreground">
        No photo for {zoneName} in the {decade}s.
      </p>
    )
  }
  return <PhotoDetails image={image} zoneName={zoneName} onEnlarge={onEnlarge} />
}

type PhotoPanelProps = {
  zoneName: string | undefined
}

/** Side panel showing a random photo of the selected zone and decade. */
export const PhotoPanel = ({ zoneName }: PhotoPanelProps) => {
  const { zoneId, decade, image, isDrawing, isUnavailable, renew, close } = useImagery()
  const [enlarged, setEnlarged] = useState(false)

  const handleOpenChange = (open: boolean) => {
    if (open) {
      return
    }
    close()
  }

  return (
    <Sheet open={zoneId !== null} modal={false} onOpenChange={handleOpenChange}>
      <SheetContent
        side="right"
        className="bottom-26 h-auto w-90 gap-0 sm:max-w-90"
        onInteractOutside={(event) => event.preventDefault()}
        onOpenAutoFocus={(event) => event.preventDefault()}
      >
        <SheetHeader className="border-b">
          <SheetTitle>
            {zoneName} <span className="font-normal text-muted-foreground">· {decade}s</span>
          </SheetTitle>
          <SheetDescription className="sr-only">A random photo of this zone and decade.</SheetDescription>
        </SheetHeader>
        <div className="flex-1 overflow-y-auto p-4">
          <PanelBody
            image={image}
            isDrawing={isDrawing}
            isUnavailable={isUnavailable}
            zoneName={zoneName}
            decade={decade}
            onEnlarge={() => setEnlarged(true)}
          />
        </div>
        <SheetFooter className="border-t">
          <RenewButton disabled={isDrawing || image === null} onClick={renew} />
        </SheetFooter>
        {enlarged && image && (
          <Suspense>
            <PhotoViewer image={image} isDrawing={isDrawing} onRenew={renew} onClose={() => setEnlarged(false)} />
          </Suspense>
        )}
      </SheetContent>
    </Sheet>
  )
}

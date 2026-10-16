import { type ElementType, type ReactNode } from 'react'

import { SheetFooter, SheetHeader } from '@shared/ui/sheet'

import { RenewButton } from './renew-button'

type PhotoSideProps = {
  zoneName: string | undefined
  decade: number
  titleAs?: ElementType
  renewDisabled: boolean
  onRenew: () => void
  children: ReactNode
}

/**
 * Zone and decade heading, scrollable photo content and Renew button. The heading and the
 * footer stay above anything overflowing the content, such as a zoomed photo.
 */
export const PhotoSide = ({
  zoneName,
  decade,
  titleAs: Title = 'h2',
  renewDisabled,
  onRenew,
  children,
}: PhotoSideProps) => (
  <>
    <SheetHeader className="relative z-30 border-b bg-background">
      <Title className="font-semibold text-foreground">
        {zoneName} <span className="font-normal text-muted-foreground">· {decade}s</span>
      </Title>
    </SheetHeader>
    <div className="flex-1 overflow-y-auto p-4">{children}</div>
    <SheetFooter className="relative z-30 border-t bg-background">
      <RenewButton disabled={renewDisabled} onClick={onRenew} />
    </SheetFooter>
  </>
)

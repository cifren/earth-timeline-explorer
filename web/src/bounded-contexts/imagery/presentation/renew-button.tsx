import { type ComponentProps } from 'react'

import { cn } from 'cn'
import { Button } from '@shared/ui/button'

/** Draws another photo of the same zone and decade. */
export const RenewButton = ({ className, ...props }: ComponentProps<typeof Button>) => (
  <Button className={cn('bg-highlight text-primary-foreground hover:bg-highlight-strong', className)} {...props}>
    Renew
  </Button>
)

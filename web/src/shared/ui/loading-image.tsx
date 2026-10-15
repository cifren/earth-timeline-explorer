import { type ComponentProps, useState } from 'react'

import { cn } from 'cn'
import { Skeleton } from '@shared/ui/skeleton'

type LoadingImageProps = ComponentProps<'img'> & { src: string; alt: string }

/** Image covered by a skeleton until its file is loaded (or has failed), for each new source. */
export const LoadingImage = ({ src, alt, className, onLoad, onError, ...props }: LoadingImageProps) => {
  const [settledSrc, setSettledSrc] = useState<string | null>(null)
  const loading = settledSrc !== src

  return (
    <span className="relative block">
      {loading && <Skeleton data-testid="image-skeleton" className={cn('absolute inset-0', className)} />}
      <img
        {...props}
        src={src}
        alt={alt}
        aria-busy={loading}
        className={cn(className, loading && 'opacity-0')}
        onLoad={(event) => {
          setSettledSrc(src)
          onLoad?.(event)
        }}
        onError={(event) => {
          setSettledSrc(src)
          onError?.(event)
        }}
      />
    </span>
  )
}

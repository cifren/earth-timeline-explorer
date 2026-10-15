import { type ComponentProps } from 'react'

import { cn } from 'cn'
import { Skeleton } from '@shared/ui/skeleton'
import { useImageLoading } from '@shared/ui/use-image-loading'

type LoadingImageProps = ComponentProps<'img'> & { src: string; alt: string }

/** Image covered by a skeleton until its file is loaded (or has failed), for each new source. */
export const LoadingImage = ({ src, alt, className, onLoad, onError, ...props }: LoadingImageProps) => {
  const { loading, settle } = useImageLoading(src)

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
          settle()
          onLoad?.(event)
        }}
        onError={(event) => {
          settle()
          onError?.(event)
        }}
      />
    </span>
  )
}

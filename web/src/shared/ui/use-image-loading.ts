import { useState } from 'react'

/** Whether the file of an image source is still loading; each new source loads again. */
export const useImageLoading = (src: string) => {
  const [settledSrc, setSettledSrc] = useState<string | null>(null)
  return { loading: settledSrc !== src, settle: () => setSettledSrc(src) }
}

import { useQuery } from '@tanstack/react-query'

import { fetchZoneOutlines } from '../infrastructure/geography-gateway'

/** Zone outlines, fetched once per session: they only change with a new release. */
export const useGeography = () => {
  const query = useQuery({
    queryKey: ['zones', 'geometry'],
    queryFn: ({ signal }) => fetchZoneOutlines(signal),
    staleTime: Infinity,
  })
  const outlines = query.data ?? []

  return {
    outlines,
    isError: query.isError,
    zoneName: (zoneId: string) => outlines.find((outline) => outline.id === zoneId)?.name,
  }
}

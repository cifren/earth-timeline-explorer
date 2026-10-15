import { useSearchParams } from 'react-router'

export const INITIAL_DECADE = 1950

type Selection = { zone: string | null; decade: number; image: string | null }

const toSearch = ({ zone, decade, image }: Selection) => {
  const search = new URLSearchParams({ decade: String(decade) })
  if (zone !== null) {
    search.set('zone', zone)
  }
  if (image !== null) {
    search.set('image', image)
  }
  return search
}

const decadeOf = (value: string | null) => {
  const decade = Number(value)
  return value === null || !Number.isInteger(decade) ? INITIAL_DECADE : decade
}

/**
 * The user's choice of zone, decade and photo, shared by the map, the timeline and the
 * photos. It lives in the page URL (`?decade=1950&zone=europe&image=…`), so that each
 * photo has its own address, to share or to come back to through the browser history.
 */
export const useSelection = () => {
  const [search, setSearch] = useSearchParams()
  const selection: Selection = {
    zone: search.get('zone'),
    decade: decadeOf(search.get('decade')),
    image: search.get('image'),
  }

  return {
    zoneId: selection.zone,
    decade: selection.decade,
    imageId: selection.image,
    /** Picks a zone, or none, for the current decade; its photo is drawn afresh. */
    selectZone: (zone: string | null) => setSearch(toSearch({ ...selection, zone, image: null })),
    /** Moves to a decade, in the current zone; its photo is drawn afresh. */
    selectDecade: (decade: number) => setSearch(toSearch({ ...selection, decade, image: null })),
    /** Shows a photo of the current zone and decade; replacing keeps a single history entry. */
    showImage: (image: string, { replace = false } = {}) => setSearch(toSearch({ ...selection, image }), { replace }),
  }
}

import { create } from 'zustand'

export const INITIAL_DECADE = 1950

type SelectionState = {
  zoneId: string | null
  decade: number
  selectZone: (zoneId: string | null) => void
  selectDecade: (decade: number) => void
}

/** The user's choice of zone and decade, shared by the map, the timeline and the photos. */
export const useSelection = create<SelectionState>()((set) => ({
  zoneId: null,
  decade: INITIAL_DECADE,
  selectZone: (zoneId) => set({ zoneId }),
  selectDecade: (decade) => set({ decade }),
}))

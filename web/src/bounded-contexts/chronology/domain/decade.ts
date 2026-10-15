/** Oldest decade of the timeline, the first one of photography. */
export const FIRST_DECADE = 1840

export const decadeOf = (year: number) => Math.floor(year / 10) * 10

export const currentDecade = (today: Date = new Date()) => decadeOf(today.getFullYear())

/** Keeps a (fractional) year between the start of the first decade and the end of the last. */
export const clampYear = (year: number, lastDecade: number) =>
  Math.min(Math.max(year, FIRST_DECADE), lastDecade + 9.999)

export const clampDecade = (decade: number, lastDecade: number) => Math.min(Math.max(decade, FIRST_DECADE), lastDecade)

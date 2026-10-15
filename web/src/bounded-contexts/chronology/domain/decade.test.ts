import { FIRST_DECADE, clampDecade, clampYear, currentDecade, decadeOf } from './decade'

describe('decades', () => {
  it('finds the decade of a year', () => {
    expect(decadeOf(1953)).toBe(1950)
    expect(decadeOf(1950)).toBe(1950)
    expect(decadeOf(1959.99)).toBe(1950)
  })

  it('takes the current decade from the date', () => {
    expect(currentDecade(new Date('2026-10-15'))).toBe(2020)
  })

  it('keeps years and decades on the timeline', () => {
    expect(clampYear(1700, 2020)).toBe(FIRST_DECADE)
    expect(clampYear(2100, 2020)).toBeCloseTo(2029.999)
    expect(clampDecade(1830, 2020)).toBe(FIRST_DECADE)
    expect(clampDecade(2030, 2020)).toBe(2020)
    expect(clampDecade(1950, 2020)).toBe(1950)
  })
})

import { ApiError } from './api-error'

describe('ApiError', () => {
  it('uses the detail of a Problem Details body', () => {
    const error = new ApiError(400, { title: 'Bad Request', detail: 'The decade must be a multiple of 10.' })

    expect(error.status).toBe(400)
    expect(error.message).toBe('The decade must be a multiple of 10.')
  })

  it('falls back to the title, then to the status', () => {
    expect(new ApiError(503, { title: 'Service Unavailable' }).message).toBe('Service Unavailable')
    expect(new ApiError(502, undefined).message).toBe('HTTP 502')
  })
})

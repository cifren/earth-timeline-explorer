type Problem = { title: string; detail?: string }

const isProblem = (body: unknown): body is Problem => typeof body === 'object' && body !== null && 'title' in body

const messageOf = (status: number, body: unknown): string => {
  if (!isProblem(body)) {
    return `HTTP ${status}`
  }
  return body.detail ?? body.title
}

/** Error answered by the API, described as Problem Details (RFC 9457) when available. */
export class ApiError extends Error {
  readonly status: number

  constructor(status: number, body: unknown) {
    super(messageOf(status, body))
    this.name = 'ApiError'
    this.status = status
  }
}

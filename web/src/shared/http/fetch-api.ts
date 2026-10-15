const baseUrl = import.meta.env.VITE_API_URL ?? ''

/**
 * Fetch function of the generated API client: prefixes the API origin and returns the
 * status, headers and decoded JSON body, whatever the status.
 */
export const fetchApi = async <T>(path: string, init: RequestInit): Promise<T> => {
  const response = await fetch(`${baseUrl}${path}`, init)
  const text = await response.text()
  const data: unknown = text === '' ? undefined : JSON.parse(text)
  return { data, status: response.status, headers: response.headers } as T
}

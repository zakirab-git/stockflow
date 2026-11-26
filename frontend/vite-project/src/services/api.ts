const API_URL = import.meta.env.VITE_API_URL ?? "/api"

export type ApiErrorBody = {
  code: string
  message?: string
  available?: number
  sku?: string
}

export class ApiError extends Error {
  status: number
  code: string
  available?: number
  sku?: string

  constructor(status: number, body: ApiErrorBody) {
    super(body.message ?? body.code)
    this.name = "ApiError"
    this.status = status
    this.code = body.code
    this.available = body.available
    this.sku = body.sku
  }
}

export async function apiRequest<T>(
  endpoint: string,
  options: RequestInit = {}
): Promise<T> {
  const response = await fetch(`${API_URL}${endpoint}`, {
    ...options,
    headers: {
      "Content-Type": "application/json",
      ...options.headers,
    },
  })

  const responseText = await response.text()

  if (!response.ok) {
    let body: ApiErrorBody = { code: "REQUEST_FAILED", message: response.statusText }

    if (responseText) {
      try {
        const parsed = JSON.parse(responseText) as {
          error?: ApiErrorBody | string
          message?: string
        }
        if (parsed.error && typeof parsed.error === "object") {
          body = parsed.error
        } else if (typeof parsed.error === "string") {
          body = { code: "BAD_REQUEST", message: parsed.error }
        } else if (parsed.message) {
          body = { code: "BAD_REQUEST", message: parsed.message }
        }
      } catch {
        body = { code: "REQUEST_FAILED", message: responseText }
      }
    }

    throw new ApiError(response.status, body)
  }

  if (!responseText) {
    return undefined as T
  }

  return JSON.parse(responseText) as T
}
// Works with Vite proxy in dev and nginx proxy in production
const API_BASE = '/api'

export interface GenerateOrderParams {
  listMenu: string
  currentOrders: string
}

export interface GenerateOrderResponse {
  generatedMessage: string
  error?: string
}

export async function generateOrder(
  params: GenerateOrderParams
): Promise<GenerateOrderResponse> {
  const response = await fetch(`${API_BASE}/generate-order`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(params),
  })

  const data: GenerateOrderResponse = await response.json()

  if (!response.ok) {
    throw new Error(data.error || 'Failed to generate order')
  }

  return data
}

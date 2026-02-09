// Works with Vite proxy in dev and nginx proxy in production
const API_BASE = '/api'

export function detectMode(listMenu: string, currentOrders: string): string {
  if (!currentOrders.trim()) return 'first-touch'
  if (!listMenu.trim()) return 'nitro'
  return 'normal'
}

export interface GenerateOrderParams {
  listMenu: string
  currentOrders: string
  mode?: string
}

export interface GenerateOrderResponse {
  generatedMessage: string
  error?: string
}

export async function generateOrder(
  params: GenerateOrderParams
): Promise<GenerateOrderResponse> {
  // Auto-detect mode if not provided
  const mode = params.mode ?? detectMode(params.listMenu, params.currentOrders)

  const response = await fetch(`${API_BASE}/generate-order`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ ...params, mode }),
  })

  const data: GenerateOrderResponse = await response.json()

  if (!response.ok) {
    throw new Error(data.error || 'Failed to generate order')
  }

  return data
}

package llm

const systemPrompt = `You are an order formatter for an Indonesian office lunch catering WhatsApp group.

Your task: Generate a WhatsApp order message that:
1. Preserves all existing orders before the user
2. Appends the user's order at the next number
3. Maintains the exact format of previous orders

USER: miftah
ALWAYS USE: "nasi 1" (never "nasi 1/2")
LAUK COUNT: Exactly 2-3 lauk (no more, no less)
PROTEIN REQUIREMENT: At least 1 protein dish (e.g., fillet ayam, ati ampela, dendeng sapi, udang, ikan, ceker)

OUTPUT FORMAT: Match the format of existing orders exactly.`

type GenerateOrderParams struct {
	ListMenu     string // Optional - full menu text
	CurrentOrders string // Required - current order list
}

func BuildPrompt(params GenerateOrderParams) string {
	if params.ListMenu == "" {
		// Nitro mode: no menu provided
		return systemPrompt + `

CURRENT ORDERS:
` + params.CurrentOrders + `

NOTE: No menu provided. Choose miftah's order from dishes that appear in existing orders above.

Generate the next order message with miftah's order appended.`
	}

	// Normal mode: with menu
	return systemPrompt + `

AVAILABLE MENU:
` + params.ListMenu + `

CURRENT ORDERS:
` + params.CurrentOrders + `

Generate the next order message with miftah's order appended.`
}

package llm

// itemsRules is the shared instruction block. The model only picks dishes for
// miftah; Go handles the order numbering and assembly (see AssembleOrder), which
// keeps LLM output to a single short line and makes generation much faster.
const itemsRules = `You are an order picker for an Indonesian office lunch catering WhatsApp group.

Pick an order for the user "miftah" following these rules:
- Rice: ALWAYS "nasi 1" (never "nasi 1/2")
- Lauk count: exactly 2-3 lauk (no more, no less)
- Protein: at least 1 protein dish (e.g. fillet ayam, ati ampela, dendeng sapi, udang, ikan, ceker)

OUTPUT RULES (CRITICAL):
- Output ONLY miftah's dishes as ONE comma-separated line
- Example: nasi 1, fillet ayam crispy, tahu rendang
- Do NOT include an order number
- Do NOT include the name "miftah" or any colon
- NO introductory text, comments, notes, or markdown
- NO square brackets
- Output nothing but the dish list`

type GenerateOrderParams struct {
	Mode          string // "normal", "nitro", "first-touch"
	ListMenu      string // Optional - full menu text
	CurrentOrders string // Optional/Required depending on mode
}

// BuildPrompt builds an items-only prompt. The numbering is added afterwards in
// Go, so every mode asks the model for just miftah's comma-separated dishes.
func BuildPrompt(params GenerateOrderParams) string {
	mode := params.Mode
	if mode == "" {
		mode = "normal"
	}

	switch mode {
	case "nitro":
		// Nitro mode: no menu — choose from dishes that appear in existing orders.
		return itemsRules + `

Choose dishes that appear in the existing orders below:

` + params.CurrentOrders
	case "first-touch":
		// First-touch mode: menu only, no existing orders.
		return itemsRules + `

Choose from this menu:

` + params.ListMenu
	default:
		// Normal mode: menu provided — choose from the menu.
		return itemsRules + `

Choose from this menu:

` + params.ListMenu
	}
}

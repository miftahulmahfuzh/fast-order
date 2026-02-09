# First-Touch Mode Design

**Date:** 2025-02-09

**Feature:** Add first-touch mode for users who order first (before anyone else)

---

## Overview

The app currently supports two modes:
1. **Normal Mode**: List menu + current orders → generates order appended to existing list
2. **Nitro Mode**: Current orders only → generates order from existing dishes

This design adds a third mode:
3. **First-Touch Mode**: List menu only → generates order #1

---

## User Workflow

| Mode | User Input | Shortcut | Output |
|------|------------|----------|--------|
| Normal | Menu → TAB → Current Orders | ENTER | Order appended to list |
| Nitro | Current Orders only | ENTER | Order from existing dishes |
| **First-Touch** | Menu only | **Shift+Enter** | Order #1 |

---

## Architecture Changes

### Mode Detection

Frontend determines the mode before sending the API request:
- If `listMenu` is empty → `mode: "nitro"`
- If `currentOrders` is empty → `mode: "first-touch"`
- If both have content → `mode: "normal"`

### Keyboard Shortcuts

- **List Menu textarea**: NEW - Shift+Enter triggers first-touch mode
- **Current Orders textarea**: Enter triggers generate (existing)
- **Global**: Ctrl+Shift+C still generates, respects current mode

---

## Backend Changes

### API Request Structure

```go
type GenerateOrderRequest struct {
    Mode         string `json:"mode"`         // "normal", "nitro", "first-touch"
    ListMenu     string `json:"listMenu"`
    CurrentOrders string `json:"currentOrders"`
}
```

### Validation by Mode

| Mode | listMenu | currentOrders |
|------|----------|---------------|
| first-touch | Required | Optional (ignored if present) |
| nitro | Optional (ignored) | Required |
| normal | Optional | Required |

### First-Touch Mode LLM Prompt

Simplified system prompt (no need to preserve existing orders):

```
You are an order formatter for an Indonesian office lunch catering WhatsApp group.

Your task: Generate the FIRST order for a new lunch order.

USER: miftah
ALWAYS USE: "nasi 1" (never "nasi 1/2")
LAUK COUNT: Exactly 2-3 lauk (no more, no less)
PROTEIN REQUIREMENT: At least 1 protein dish

AVAILABLE MENU:
{menu}

OUTPUT FORMAT: 1. miftah - [nasi 1], [lauk 1], [lauk 2]

CRITICAL OUTPUT RULES:
- Output ONLY the numbered order list
- NO introductory text, comments, or markdown
- Ready to paste directly into WhatsApp
```

---

## Frontend Changes

### App.tsx

1. **New keyboard handler for List Menu** (Shift+Enter)
2. **New generate function** with mode detection and fallback
3. **Updated success notification** to show active mode
4. **Updated placeholder** to mention Shift+Enter shortcut

### API Client (lib/api.ts)

```typescript
export async function generateOrder(params: {
  listMenu: string
  currentOrders: string
  mode?: string
})

function detectMode(listMenu: string, currentOrders: string): string {
  if (!currentOrders.trim()) return 'first-touch'
  if (!listMenu.trim()) return 'nitro'
  return 'normal'
}
```

---

## Error Handling & Edge Cases

| Scenario | Behavior |
|----------|----------|
| Shift+Enter on empty list menu | Error: "List menu required for first-touch mode" |
| Shift+Enter with current orders filled | Falls back to normal mode silently |
| Invalid mode value | Backend defaults to "normal" |

---

## Success Notification

Format: `Order copied to clipboard! ({modeName}) Press Ctrl+V to paste in WhatsApp`

Mode labels: "Normal Mode", "Nitro Mode", "First-Touch Mode"

---

## Implementation Tasks

| ID | Task | Dependencies |
|----|------|--------------|
| 1 | Add mode field to backend API and handler | - |
| 2 | Create first-touch mode LLM prompt | 1 |
| 3 | Update API client to send mode parameter | 1 |
| 4 | Add Shift+Enter keyboard shortcut | 3 |
| 5 | Add mode indication in success notification | 4 |
| 6 | Update global shortcuts hint text | 4 |

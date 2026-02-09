# Modern Card UI Design - Fast Order

**Status:** Design Approved
**Created:** 2025-02-09
**Theme:** Clean & Professional SaaS-style

## Overview

Replacing the brutalist theme with a modern card-based design inspired by Linear, Vercel, and contemporary SaaS applications. The new design prioritizes professionalism, spaciousness, and mobile responsiveness.

---

## Design System

### Color Palette

| Purpose | Color | Hex |
|---------|-------|-----|
| Primary (Actions) | Indigo 500 | `#6366F1` |
| Primary Hover | Indigo 600 | `#4F46E5` |
| Background | Slate 50 | `#F8FAFC` |
| Surface/Card | White | `#FFFFFF` |
| Text Primary | Slate 900 | `#0F172A` |
| Text Secondary | Slate 500 | `#64748B` |
| Border | Slate 200 | `#E2E8F0` |
| Success | Emerald 500 | `#10B981` |
| Error | Red 500 | `#EF4444` |

### Typography

| Element | Weight | Size | Notes |
|---------|--------|------|-------|
| Title | 600 | 1.75rem | No uppercase, normal letter-spacing |
| Labels | 500 | 0.875rem | Slate 700 |
| Body | 400 | 1rem | Line-height: 1.6 |

Font family remains Inter.

### Icons (Lucide)

| Component | Icon |
|-----------|------|
| Header logo | `Clipboard` or `Zap` |
| Clear button | `X` |
| Generate button | `Copy` or `Send` |
| Success | `CheckCircle` |
| Error | `AlertCircle` |
| Loading | `Loader2` (animated) |
| Keyboard hint | `Keyboard` |

---

## Layout Structure

```
body (Slate 50 background)
  └── main-container (full-width)
      ├── header (sticky, white, shadow)
      └── content-wrapper
          ├── horizontal padding: 24px (desktop), 16px (mobile)
          ├── vertical padding: 32px
          └── form section
```

**Key changes from brutalist:**
- No `max-width: 800px` constraint
- No `border: 2px solid` on container
- Uses horizontal padding instead of fixed-width
- Header has subtle shadow instead of border

### Responsive Breakpoints

| Screen Size | Padding | Notes |
|-------------|---------|-------|
| Mobile (< 640px) | 16px | Stacked layout |
| Tablet (640-768px) | 20px | - |
| Desktop (> 768px) | 24px | Max-width: 1200px optional |

---

## Component Specifications

### TextAreaField

```css
Container: No border, transparent background
Label: 500 weight, 0.875rem, Slate 700
Required: Subtle asterisk
Textarea:
  - Border: 1px solid Slate 200
  - Border-radius: 8px
  - Padding: 12px 16px
  - Min-height: 120px
  - Focus: Indigo 500 border + ring shadow
Clear Button:
  - Icon only (X), 16px
  - Hover: Slate 100 background
  - Border-radius: 6px
```

### Generate Button

```css
Background: Indigo 500
Text: White, 600 weight
Border-radius: 8px
Padding: 14px 24px
Height: 48px (touch-friendly)
Icon: Copy/Send on left
Hover: Indigo 600 + translateY(-1px)
Loading: Loader2 icon (spinning)
```

### Status Message

```css
Border-radius: 8px
Padding: 12px 16px
Icon left of text

Success: Emerald 50 bg, Emerald 200 border
Error: Red 50 bg, Red 200 border
Loading: Slate 100 bg
```

---

## Mobile Adaptations

- Minimum touch target: 44px height
- Hide keyboard shortcuts on mobile (< 640px)
- Reduced vertical gaps (16px instead of 24px)
- Compact header (12px padding)
- Viewport meta tag required

---

## Animations

| Interaction | Duration | Easing |
|-------------|----------|--------|
| Button hover | 150ms | ease-out |
| Button active | 100ms | - |
| Input focus | 200ms | - |
| Status appear | 300ms | ease-out |

---

## Implementation Tasks

1. `#1` - Update color system and typography variables
2. `#2` - Redesign layout container for full-width spacious feel [blocked by #1]
3. `#3` - Redesign header with sticky positioning and icon [blocked by #1]
4. `#4` - Redesign form fields with soft borders and rounded corners [blocked by #1]
5. `#5` - Redesign generate button with icon and modern styling [blocked by #1]
6. `#6` - Redesign status messages with icons and soft colors [blocked by #1]
7. `#7` - Add mobile responsive styles and touch targets [blocked by #2, #3, #4, #5, #6]
8. `#8` - Add micro-interactions and animations [blocked by #2, #3, #4, #5, #6]

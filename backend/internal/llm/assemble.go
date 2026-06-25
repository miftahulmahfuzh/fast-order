package llm

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// leadingNumberRe captures the number at the very start of a line, tolerating
// leading whitespace and invisible Unicode format characters (\p{Cf}, the
// word-joiners WhatsApp/phones insert). It does NOT require any separator after
// the number, so all of these match: "1mi", "2. alam", "3-rini", "4  Ervina.",
// "6..farid", "7)Clive", "8Refki", "9, jennifer", "10 Audrey".
var leadingNumberRe = regexp.MustCompile(`^[\s\p{Cf}]*([0-9]+)`)

// itemNumberPrefixRe strips a leading order number AND its trailing separator
// punctuation from raw model output, e.g. "11. ", "11) ", "11 ".
var itemNumberPrefixRe = regexp.MustCompile(`^[\s\p{Cf}]*[0-9]+[.):,\-]*[\s\p{Cf}]*`)

// leadingNameRe matches a leading "miftah :" / "miftah:" prefix so we can strip
// it if the model includes the name despite being told not to.
var leadingNameRe = regexp.MustCompile(`(?i)^[\s\p{Cf}]*miftah[\s\p{Cf}]*:`)

// monthRe detects a date header by month name (Indonesian or English). A line
// like "25 juni" must not be counted as an order number.
var monthRe = regexp.MustCompile(`(?i)\b(jan(uari|uary)?|feb(ruari|ruary)?|mar(et|ch)?|apr(il)?|mei|jun(i|e)?|jul(i|y)?|agu(stus)?|aug(ust)?|sep(t|tember)?|okt(ober)?|oct(ober)?|nov(ember)?|des(ember)?|dec(ember)?)\b`)

// dateNumberRe detects a numeric date header like "25/06" or "25-06-2026":
// a number followed by a separator and ANOTHER number. This intentionally does
// NOT match order lines like "3-rini" (separator followed by a letter).
var dateNumberRe = regexp.MustCompile(`^[\s\p{Cf}]*[0-9]{1,2}\s*[/.\-]\s*[0-9]`)

// isDateLine reports whether a line is a date header rather than an order.
func isDateLine(line string) bool {
	return monthRe.MatchString(line) || dateNumberRe.MatchString(line)
}

// NextOrderNumber returns the next order number for the given list: the highest
// leading number found across all order lines, plus one. Date headers (e.g.
// "25 juni") are skipped, and the leading number is matched regardless of which
// separator — if any — follows it. Returns 1 when no numbered order line is
// present (e.g. an empty list in first-touch mode).
func NextOrderNumber(orders string) int {
	max := 0
	for line := range strings.SplitSeq(orders, "\n") {
		if strings.TrimSpace(line) == "" || isDateLine(line) {
			continue
		}
		m := leadingNumberRe.FindStringSubmatch(line)
		if m == nil {
			continue
		}
		n, err := strconv.Atoi(m[1])
		if err == nil && n > max {
			max = n
		}
	}
	return max + 1
}

// SanitizeOrderItems reduces raw LLM output to a single clean comma-separated
// dish list for miftah. It defensively strips a leading order number, a leading
// "miftah :" prefix, square brackets, and surrounding punctuation/whitespace —
// so the result is safe to assemble into a numbered line ourselves.
func SanitizeOrderItems(raw string) string {
	// Use the first non-empty line; the model should emit exactly one.
	line := ""
	for l := range strings.SplitSeq(raw, "\n") {
		if strings.TrimSpace(l) != "" {
			line = l
			break
		}
	}

	line = itemNumberPrefixRe.ReplaceAllString(line, "")
	line = leadingNameRe.ReplaceAllString(line, "")
	line = strings.ReplaceAll(line, "[", "")
	line = strings.ReplaceAll(line, "]", "")
	line = strings.TrimSpace(line)
	// Drop any stray leading/trailing commas or periods left after stripping.
	line = strings.Trim(line, ",. ")
	return strings.TrimSpace(line)
}

// AssembleOrder appends miftah's order to the existing list with the correct
// numbering, preserving the original list text verbatim. When currentOrders is
// empty (first-touch mode) it returns just the single numbered line.
func AssembleOrder(currentOrders string, number int, items string) string {
	line := fmt.Sprintf("%d. miftah : %s", number, items)

	trimmed := strings.TrimRight(currentOrders, " \t\r\n")
	if trimmed == "" {
		return line
	}
	return trimmed + "\n" + line
}

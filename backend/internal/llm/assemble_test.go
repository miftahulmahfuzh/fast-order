package llm

import "testing"

// realOrders mirrors a messy WhatsApp list: a date header, inconsistent
// spacing ("8.Refki" with no space), and "+" item separators.
const realOrders = `25 juni
1. mi: nasi 1, Telur bulet rendangss, Tahu cabe garam
2. alam: nasi 1/2, sup bayam jagung, jamur crispy, tahu rendang
3. rini : nasi, sup bayam jagung, oseng bakso cabe garam
4. Ervina : Nasi 1/2, Capcai, Perkedel, sambal
5. nabila : nasi 1/2 + fillet ayam crispy mini sambal matah + sambal
6. farid : nasi, jamur crispy, tahu rendangs, acar kuning, capcay
7. Clive: nasi 1/2, oncom leuncah, ayam suwir daun jeruk, sosis bakar bbq
8.Refki: Nasi 1/2, Ikan Cue Sarden, Jamur Crispy, Bihun Goreng, tahu isi 1
9. jennifer: nasi 1/2 + fillet ayam crispy mini sambal matah + supbayam jagung
10. Audrey: nasi merah 1/2 + sup Bayam jagung + tempe tahu bacem`

// messyOrders has a different leading separator on every single line (and a
// date header), reflecting people typing over each other in a WhatsApp group.
const messyOrders = `25 juni
1mi: nasi 1, Telur bulet rendangss, Tahu cabe garam
2. alam: nasi 1/2, sup bayam jagung, jamur crispy, tahu rendang
3-rini : nasi, sup bayam jagung, oseng bakso cabe garam
4   Ervina. Nasi 1/2, Capcai, Perkedel, sambal
5. nabila : nasi 1/2 + fillet ayam crispy mini sambal matah + sambal
6.. farid : nasi, jamur crispy, tahu rendangs, acar kuning, capcay
7)Clive: nasi 1/2, oncom leuncah, ayam suwir daun jeruk, sosis bakar bbq
8Refki: Nasi 1/2, Ikan Cue Sarden, Jamur Crispy, Bihun Goreng, tahu isi 1
9, jennifer: nasi 1/2 + fillet ayam crispy mini sambal matah + supbayam jagung
10 Audrey: nasi merah 1/2 + sup Bayam jagung + tempe tahu bacem`

func TestNextOrderNumber(t *testing.T) {
	tests := []struct {
		name   string
		orders string
		want   int
	}{
		{"empty list starts at 1", "", 1},
		{"whitespace only starts at 1", "   \n\n", 1},
		{"date header is ignored", "25 juni\n1. a: nasi", 2},
		{"no-space number 8.Refki counts", "7. a: x\n8.Refki: y", 9},
		{"two-digit numbers", realOrders, 11},
		{"leading invisible chars tolerated", "⁠1.⁠ a: nasi", 2},
		{"out-of-order uses max not last", "1. a\n5. b\n3. c", 6},
		{"messy mixed separators", messyOrders, 11},
		{"no date header present", "1mi: nasi\n2-alam: tahu", 3},
		{"numeric date header dd/mm ignored", "25/06\n1. a: nasi\n2. b: tahu", 3},
		{"number with no separator at all", "8Refki: nasi", 9},
		{"paren separator", "7)Clive: nasi", 8},
		{"comma separator", "9, jennifer: nasi", 10},
		{"dash separator not confused with numeric date", "3-rini: nasi", 4},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NextOrderNumber(tt.orders); got != tt.want {
				t.Errorf("NextOrderNumber() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestSanitizeOrderItems(t *testing.T) {
	tests := []struct {
		name string
		raw  string
		want string
	}{
		{"clean line untouched", "nasi 1, fillet ayam, tahu rendang", "nasi 1, fillet ayam, tahu rendang"},
		{"strips leading number", "11. nasi 1, fillet ayam", "nasi 1, fillet ayam"},
		{"strips leading name and number", "11. miftah : nasi 1, tahu", "nasi 1, tahu"},
		{"strips leading name only", "miftah: nasi 1, tahu", "nasi 1, tahu"},
		{"removes brackets", "[nasi 1], [fillet ayam]", "nasi 1, fillet ayam"},
		{"takes first non-empty line", "\n\nnasi 1, tahu\nignored second line", "nasi 1, tahu"},
		{"trims trailing period", "nasi 1, fillet ayam.", "nasi 1, fillet ayam"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SanitizeOrderItems(tt.raw); got != tt.want {
				t.Errorf("SanitizeOrderItems(%q) = %q, want %q", tt.raw, got, tt.want)
			}
		})
	}
}

func TestAssembleOrder(t *testing.T) {
	t.Run("appends with correct number", func(t *testing.T) {
		got := AssembleOrder(realOrders, 11, "nasi 1, fillet ayam crispy mini sambal matah, tahu rendang")
		wantLast := "11. miftah : nasi 1, fillet ayam crispy mini sambal matah, tahu rendang"
		lines := splitLines(got)
		if lines[len(lines)-1] != wantLast {
			t.Errorf("last line = %q, want %q", lines[len(lines)-1], wantLast)
		}
		// Original list must be preserved verbatim as the prefix.
		if got[:len(realOrders)] != realOrders {
			t.Errorf("original list was not preserved verbatim")
		}
	})

	t.Run("first-touch returns single line", func(t *testing.T) {
		got := AssembleOrder("", 1, "nasi 1, fillet ayam, tahu")
		want := "1. miftah : nasi 1, fillet ayam, tahu"
		if got != want {
			t.Errorf("AssembleOrder() = %q, want %q", got, want)
		}
	})

	t.Run("trims trailing newlines before appending", func(t *testing.T) {
		got := AssembleOrder("1. a: nasi\n\n", 2, "nasi 1, tahu")
		want := "1. a: nasi\n2. miftah : nasi 1, tahu"
		if got != want {
			t.Errorf("AssembleOrder() = %q, want %q", got, want)
		}
	})
}

func splitLines(s string) []string {
	var out []string
	cur := ""
	for _, r := range s {
		if r == '\n' {
			out = append(out, cur)
			cur = ""
			continue
		}
		cur += string(r)
	}
	out = append(out, cur)
	return out
}

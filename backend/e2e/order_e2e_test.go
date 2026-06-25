//go:build e2e

// Package e2e contains end-to-end tests that exercise a running backend
// (e.g. the dockerized service on :8089). They are gated behind the `e2e`
// build tag so they never run during normal `go test ./...`.
//
// Run against the running container:
//
//	BASE_URL=http://localhost:8089 go test -tags e2e -v ./e2e/...
package e2e

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"testing"
	"time"
)

const nitroCurrentOrders = `25 juni
1. miftah: nasi 1, Telur bulet rendangss, Tahu cabe garam
2. alam: nasi 1/2, sup bayam jagung, jamur crispy, tahu rendang
3. rini : nasi, sup bayam jagung, oseng bakso cabe garam
4. Ervina : Nasi 1/2, Capcai, Perkedel, sambal
5. nabila : nasi 1/2 + fillet ayam crispy mini sambal matah + sambal
6. farid : nasi, jamur crispy, tahu rendangs, acar kuning, capcay
7. Clive: nasi 1/2, oncom leuncah, ayam suwir daun jeruk, sosis bakar bbq
8. Refki: Nasi 1/2, Ikan Cue Sarden, Jamur Crispy, Bihun Goreng, tahu isi 1
9. jennifer: nasi 1/2 + fillet ayam crispy mini sambal matah + supbayam jagung
10. Audrey: nasi merah 1/2 + sup Bayam jagung + tempe tahu bacem`

type generateOrderRequest struct {
	Mode          string `json:"mode"`
	ListMenu      string `json:"listMenu"`
	CurrentOrders string `json:"currentOrders"`
}

type generateOrderResponse struct {
	GeneratedMessage string `json:"generatedMessage"`
	DurationMs       int64  `json:"durationMs"`
	Error            string `json:"error"`
}

func baseURL() string {
	if v := os.Getenv("BASE_URL"); v != "" {
		return v
	}
	return "http://localhost:8089"
}

// TestNitroModeGeneration calls the live backend in nitro mode and reports
// how long the LLM generation took, both as measured by the client (full
// round-trip) and as reported by the server (durationMs).
func TestNitroModeGeneration(t *testing.T) {
	reqBody, err := json.Marshal(generateOrderRequest{
		Mode:          "nitro",
		CurrentOrders: nitroCurrentOrders,
	})
	if err != nil {
		t.Fatalf("failed to marshal request: %v", err)
	}

	url := baseURL() + "/api/generate-order"
	client := &http.Client{Timeout: 60 * time.Second}

	start := time.Now()
	resp, err := client.Post(url, "application/json", bytes.NewReader(reqBody))
	roundTrip := time.Since(start)
	if err != nil {
		t.Fatalf("request to %s failed: %v", url, err)
	}
	defer resp.Body.Close()

	var out generateOrderResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d (error=%q)", resp.StatusCode, out.Error)
	}
	if out.GeneratedMessage == "" {
		t.Fatalf("expected non-empty generatedMessage, got empty")
	}

	t.Logf("nitro mode generation:")
	t.Logf("  server LLM duration : %d ms", out.DurationMs)
	t.Logf("  client round-trip   : %d ms", roundTrip.Milliseconds())
	t.Logf("  generated message:\n%s", out.GeneratedMessage)
}

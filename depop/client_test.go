package depop

import (
	"encoding/json"
	"os"
	"testing"
)

func TestParseSearchResponse(t *testing.T) {
	data, err := os.ReadFile("testdata/search_mock.json")
	if err != nil {
		t.Fatalf("failed to read testdata/search_mock.json: %v", err)
	}

	var resp SearchResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		t.Fatalf("failed to unmarshal JSON: %v", err)
	}

	if len(resp.Products) != 2 {
		t.Fatalf("expected 2 products, got %d", len(resp.Products))
	}

	p1 := resp.Products[0]
	if p1.ID != 4819201 {
		t.Errorf("expected product ID 4819201, got %d", p1.ID)
	}
	if p1.Seller.Username != "avantgardearchive" {
		t.Errorf("expected seller username avantgardearchive, got %s", p1.Seller.Username)
	}
	if p1.Price.Amount != "450.00" || p1.Price.CurrencyCode != "EUR" {
		t.Errorf("expected price 450.00 EUR, got %s %s", p1.Price.Amount, p1.Price.CurrencyCode)
	}
}

func TestNewClient_Options(t *testing.T) {
	// Verify client initialization without proxy
	client, err := NewClient("")
	if err != nil {
		t.Fatalf("expected client to initialize without error, got: %v", err)
	}
	if client == nil || client.httpClient == nil {
		t.Fatalf("expected client and httpClient to be non-nil")
	}
}

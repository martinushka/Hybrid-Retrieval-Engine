package text

import "testing"

func TestNormalize(t *testing.T) {
	got := Normalize("  Lenovo, THINKPAD X1!  ")

	want := "lenovo thinkpad x1"

	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestNormalizeRussian(t *testing.T) {
	got := Normalize("Ноутбук Lenovo!")

	want := "ноутбук lenovo"

	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestTokens(t *testing.T) {
	got := Tokens("Lenovo ThinkPad X1")

	want := []string{"lenovo", "thinkpad", "x1"}

	if len(got) != len(want) {
		t.Fatalf("expected %d tokens, got %d", len(want), len(got))
	}

	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("expected token %q at index %d, got %q", want[i], i, got[i])
		}
	}
}

func TestTokensEmpty(t *testing.T) {
	got := Tokens("!!!")

	if got != nil {
		t.Fatalf("expected nil, got %v", got)
	}
}

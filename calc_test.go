package LuhnCalculator

import (
	"errors"
	"strings"
	"testing"
)

func TestValidate(t *testing.T) {
	for _, tc := range []struct {
		number string
		want   bool
	}{
		{"79927398713", true}, {"79927398710", false},
		{"4242424242424242", true}, {"4242424242424241", false},
		{"00079927398713", true}, {"0", true}, {"1", false},
		{"", false}, {"-18", false}, {" 18", false},
		{"1 8", false}, {"18\n", false}, {"１２", false}, {"1a8", false},
	} {
		t.Run(tc.number, func(t *testing.T) {
			if got := Validate(tc.number); got != tc.want {
				t.Errorf("Validate(%q) = %v, want %v", tc.number, got, tc.want)
			}
		})
	}
}

func TestCheckDigit(t *testing.T) {
	for _, tc := range []struct {
		payload string
		digit   int
	}{
		{"7992739871", 3}, {"424242424242424", 2}, {"1", 8},
		{"12", 5}, {"0", 0}, {"0001", 8}, {strings.Repeat("0", 10000), 0},
	} {
		digit, err := CheckDigit(tc.payload)
		if err != nil || digit != tc.digit {
			t.Errorf("CheckDigit(%q) = %d, %v; want %d", tc.payload, digit, err, tc.digit)
		}
		complete, err := AppendCheckDigit(tc.payload)
		if err != nil || !strings.HasPrefix(complete, tc.payload) || !Validate(complete) {
			t.Errorf("AppendCheckDigit(%q) = %q, %v", tc.payload, complete, err)
		}
	}
}

func TestInvalidPayload(t *testing.T) {
	for _, payload := range []string{"", "-1", "1.2", "1 2", "１２", "\xff", "12\n"} {
		if _, err := CheckDigit(payload); !errors.Is(err, ErrInvalidDigits) {
			t.Errorf("CheckDigit(%q) error = %v", payload, err)
		}
		if got, err := AppendCheckDigit(payload); got != "" || !errors.Is(err, ErrInvalidDigits) {
			t.Errorf("AppendCheckDigit(%q) = %q, %v", payload, got, err)
		}
	}
}

func TestLegacyAPI(t *testing.T) {
	if got := GenerateLastDigit(12); got != 125 {
		t.Errorf("GenerateLastDigit(12) = %d, want complete number 125", got)
	}
	if !IsValidLuhn(125) || IsValidLuhn(126) || IsValidLuhn(-18) {
		t.Fatal("incorrect legacy validation")
	}
	maxInt := int(^uint(0) >> 1)
	for _, input := range []int{-1, maxInt} {
		if got := GenerateLastDigit(input); got != 0 {
			t.Errorf("GenerateLastDigit(%d) = %d, want 0", input, got)
		}
	}
	if !IsValidLuhn(0) || GenerateLastDigit(0) != 0 {
		t.Fatal("zero's historical checksum behavior changed")
	}
}

func FuzzAppendCheckDigit(f *testing.F) {
	for _, seed := range []string{"7992739871", "0001", "", "-18", "１２"} {
		f.Add(seed)
	}
	f.Fuzz(func(t *testing.T, payload string) {
		complete, err := AppendCheckDigit(payload)
		if err != nil {
			if !errors.Is(err, ErrInvalidDigits) {
				t.Fatal(err)
			}
			return
		}
		if len(complete) != len(payload)+1 || !strings.HasPrefix(complete, payload) || !Validate(complete) {
			t.Fatalf("invalid generated number %q for %q", complete, payload)
		}
		for digit := byte('0'); digit <= '9'; digit++ {
			if digit != complete[len(complete)-1] && Validate(payload+string(digit)) {
				t.Fatalf("more than one valid check digit for %q", payload)
			}
		}
	})
}

//nolint:testpackage
package ipcounter

import (
	"errors"
	"testing"
)

func TestExtractAeroKey(t *testing.T) {
	t.Parallel()

	t.Run("extract key success", func(t *testing.T) {
		t.Parallel()

		input := "example::ab cd ef"
		expected := "ab cd ef"
		got, err := extractAeroKey(input)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if got != expected {
			t.Fatalf("expected %q, got %q", expected, got)
		}
	})

	t.Run("extract key fails", func(t *testing.T) {
		t.Parallel()

		input := "example without proper hex"
		_, err := extractAeroKey(input)
		if err == nil {
			t.Fatal("expected an error, got nil")
		}
		if !errors.Is(err, ErrNoHexIPFound) {
			t.Fatalf("expected error %v, got %v", ErrNoHexIPFound, err)
		}
	})
}

func TestConvertHexToIPV4String(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		input := "00 00 00 00 00 00 00 00 00 00 00 00 c0 a8 01 01"
		expected := "192.168.1.1"
		got, err := convertHexToIPV4String(input)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if got != expected {
			t.Fatalf("expected %q, got %q", expected, got)
		}
	})

	t.Run("decode hex error", func(t *testing.T) {
		t.Parallel()

		input := "zy yz"
		_, err := convertHexToIPV4String(input)
		if err == nil {
			t.Fatal("expected an error, got nil")
		}
		if !errors.Is(err, ErrDecodingHexFailed) {
			t.Fatalf("expected error %v, got %v", ErrDecodingHexFailed, err)
		}
	})

	t.Run("invalid IPV4 length", func(t *testing.T) {
		t.Parallel()

		input := "00 00 00 00 00 00 00"
		_, err := convertHexToIPV4String(input)
		if err == nil {
			t.Fatal("expected an error, got nil")
		}
		if !errors.Is(err, ErrInvalidIPV4Length) {
			t.Fatalf("expected error %v, got %v", ErrInvalidIPV4Length, err)
		}
	})
}

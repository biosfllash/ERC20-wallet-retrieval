package gui

import (
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
)

func TestDerivePrivateKey(t *testing.T) {
	// Table-driven tests for derivePrivateKey.
	tests := []struct {
		name           string
		mnemonic       string
		derivationPath string
		wantErr        bool
	}{
		{
			name:           "valid mnemonic and derivation path",
			mnemonic:       "test test test test test test test test test test test junk",
			derivationPath: "m/44'/60'/0'/0/0",
			wantErr:        false,
		},
		{
			name:           "invalid mnemonic",
			mnemonic:       "invalid mnemonic",
			derivationPath: "m/44'/60'/0'/0/0",
			wantErr:        true,
		},
		{
			name:           "invalid derivation path",
			mnemonic:       "test test test test test test test test test test test junk",
			derivationPath: "m/invalid/path",
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			privateKey, err := derivePrivateKey(tt.mnemonic, tt.derivationPath)
			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error but got nil")
				}
				if privateKey != nil {
					t.Errorf("Expected nil private key on error, got: %v", privateKey)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if privateKey == nil {
					t.Errorf("Expected a valid private key, got nil")
				}
			}
		})
	}
}

func TestDeriveAddress(t *testing.T) {
	// Generate a new random key.
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("Failed to generate private key: %v", err)
	}

	address, err := deriveAddress(privateKey)
	if err != nil {
		t.Fatalf("Failed to derive address: %v", err)
	}

	if address.Hex() == "" {
		t.Error("Derived address is empty")
	}

	// Verify that the returned address matches the address from crypto.PubkeyToAddress.
	expected := crypto.PubkeyToAddress(privateKey.PublicKey)
	if address.Hex() != expected.Hex() {
		t.Errorf("Derived address %s does not match expected %s", address.Hex(), expected.Hex())
	}
}

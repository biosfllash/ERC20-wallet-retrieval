package gui

import (
	"crypto/ecdsa"
	"fmt"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
)

const basePath = "m/44'/60'/0'/0/"

// Start launches the GUI for wallet retrieval.
func Start() {
	a := app.New()
	w := a.NewWindow("ERC20 Wallet Retrieval")

	// Build the UI elements.
	mnemonicEntry := widget.NewEntry()
	mnemonicEntry.SetPlaceHolder("Enter your mnemonic")

	addressIndexEntry := widget.NewEntry()
	addressIndexEntry.SetPlaceHolder("Enter address index (number)")

	privateKeyLabel := widget.NewLabel("Private Key: *****")
	addressLabel := widget.NewLabel("Address: ")
	errorLabel := widget.NewLabel("")       // Error area for displaying errors
	errorLabel.Wrapping = fyne.TextWrapWord // Enable word wrapping

	// Variables to store values for copying
	var privateKeyHex string
	var addressHex string

	// Copy buttons (initially disabled)
	copyPrivateKeyButton := widget.NewButton("Copy Private Key", nil)
	copyPrivateKeyButton.Disable()

	copyAddressButton := widget.NewButton("Copy Address", nil)
	copyAddressButton.Disable()

	showButton := widget.NewButton("Show", nil) // Button to toggle private key visibility
	privateKeyHidden := true                    // Track whether the private key is hidden

	retrieveButton := widget.NewButton("Retrieve Wallet", func() {
		mnemonic := strings.TrimSpace(mnemonicEntry.Text)
		indexStr := strings.TrimSpace(addressIndexEntry.Text)
		if mnemonic == "" || indexStr == "" {
			errorLabel.SetText("Please fill in both fields.")
			privateKeyLabel.SetText("Private Key: *****")
			addressLabel.SetText("Address: ")
			copyPrivateKeyButton.Disable()
			copyAddressButton.Disable()
			return
		}

		// Let domain packages validate mnemonic and index formats.

		addrIndex, err := strconv.Atoi(indexStr)
		if err != nil || addrIndex < 0 {
			errorLabel.SetText("Invalid address index. It must be a non-negative number.")
			privateKeyLabel.SetText("Private Key: *****")
			addressLabel.SetText("Address: ")
			copyPrivateKeyButton.Disable()
			copyAddressButton.Disable()
			return
		}

		// Construct the derivation path.
		derivationPath := fmt.Sprintf("%s%d", basePath, addrIndex)
		privateKey, err := derivePrivateKey(mnemonic, derivationPath)
		if err != nil {
			errorLabel.SetText(fmt.Sprintf("Error: %v", err))
			privateKeyLabel.SetText("Private Key: *****")
			addressLabel.SetText("Address: ")
			copyPrivateKeyButton.Disable()
			copyAddressButton.Disable()
			return
		}

		// Store and mask the private key by default.
		privateKeyHex = fmt.Sprintf("%x", crypto.FromECDSA(privateKey))
		privateKeyLabel.SetText("Private Key: *****")
		privateKeyHidden = true
		showButton.SetText("Show") // Reset the button text

		// Set the address.
		addr, err := deriveAddress(privateKey)
		if err != nil {
			errorLabel.SetText(fmt.Sprintf("Error: %v", err))
			addressLabel.SetText("Address: ")
			copyPrivateKeyButton.Disable()
			copyAddressButton.Disable()
			return
		}

		addressHex = addr.Hex()
		addressLabel.SetText(fmt.Sprintf("Address: %s", addressHex))

		// Enable copy buttons
		copyPrivateKeyButton.Enable()
		copyAddressButton.Enable()

		// Set up copy button actions using the window's clipboard.
		copyPrivateKeyButton.OnTapped = func() {
			w.Clipboard().SetContent(privateKeyHex)
			errorLabel.SetText("Private key copied to clipboard!")
		}

		copyAddressButton.OnTapped = func() {
			w.Clipboard().SetContent(addressHex)
			errorLabel.SetText("Address copied to clipboard!")
		}

		// Clear any previous errors.
		errorLabel.SetText("")

		// Update the "Show" button functionality.
		showButton.OnTapped = func() {
			if privateKeyHidden {
				privateKeyLabel.SetText(fmt.Sprintf("Private Key: %s", privateKeyHex))
				showButton.SetText("Hide")
			} else {
				privateKeyLabel.SetText("Private Key: *****")
				showButton.SetText("Show")
			}
			privateKeyHidden = !privateKeyHidden
		}
	})

	// Arrange UI components.
	privateKeyContainer := container.NewHBox(privateKeyLabel, showButton, copyPrivateKeyButton)
	addressContainer := container.NewHBox(addressLabel, copyAddressButton)

	content := container.NewVBox(
		widget.NewLabel("Enter wallet details:"),
		mnemonicEntry,
		addressIndexEntry,
		retrieveButton,
		privateKeyContainer,
		addressContainer,
		errorLabel, // Error area
	)

	w.Resize(fyne.NewSize(500, 350))
	w.SetContent(content)
	w.ShowAndRun()
}

// derivePrivateKey returns the private key derived from the mnemonic and derivation path.
func derivePrivateKey(mnemonic string, derivationPath string) (*ecdsa.PrivateKey, error) {
	wallet, err := hdwallet.NewFromMnemonic(mnemonic)
	if err != nil {
		return nil, fmt.Errorf("failed to create wallet from mnemonic: %w", err)
	}

	path, err := hdwallet.ParseDerivationPath(derivationPath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse derivation path: %w", err)
	}

	account, err := wallet.Derive(path, false)
	if err != nil {
		return nil, fmt.Errorf("failed to derive account: %w", err)
	}

	privateKey, err := wallet.PrivateKey(account)
	if err != nil {
		return nil, fmt.Errorf("failed to get private key: %w", err)
	}

	return privateKey, nil
}

// deriveAddress returns the public address associated with the given private key.
func deriveAddress(privateKey *ecdsa.PrivateKey) (common.Address, error) {
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return common.Address{}, fmt.Errorf("error casting public key to ECDSA")
	}

	address := crypto.PubkeyToAddress(*publicKeyECDSA)
	return address, nil
}

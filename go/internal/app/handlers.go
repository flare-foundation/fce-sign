package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sign-extension/internal/base"

	secp256k1 "github.com/decred/dcrd/dcrec/secp256k1/v4"
)

// state holds the mutable state for the extension.
// The framework serializes all handler calls, so no additional locking is needed.
var (
	privateKey *secp256k1.PrivateKey
	signPort   string
	httpClient = http.DefaultClient
)

// SetSignPort sets the sign port for communicating with the TEE node.
func SetSignPort(port string) {
	signPort = port
}

// Register registers the handlers and initial state with the framework.
func Register(f *base.Framework) {
	f.Handle(OpTypeKey, OpCommandUpdate, handleKeyUpdate)
	f.Handle(OpTypeKey, OpCommandSign, handleKeySign)
}

// ReportState returns a JSON snapshot of the current state.
func ReportState() json.RawMessage {
	hasKey := privateKey != nil
	data, _ := json.Marshal(map[string]interface{}{
		"hasKey":  hasKey,
		"version": Version,
	})
	return data
}

// handleKeyUpdate decrypts the original message using the TEE node's key, then
// stores the decrypted value as an ECDSA private key.
func handleKeyUpdate(msg string) (data *string, status int, err error) {
	if msg == "" {
		return nil, 0, fmt.Errorf("originalMessage is empty")
	}

	// originalMessage is a hex string (hexutil.Bytes JSON serialization).
	// Hex-decode to get the raw ECIES ciphertext bytes.
	ciphertext, hexErr := base.HexToBytes(msg)
	if hexErr != nil {
		return nil, 0, fmt.Errorf("invalid hex in originalMessage: %v", hexErr)
	}

	// Decrypt via TEE node — sends ciphertext bytes (JSON-serialized as base64).
	keyBytes, decryptErr := decryptViaNode(ciphertext)
	if decryptErr != nil {
		return nil, 0, fmt.Errorf("decryption failed: %v", decryptErr)
	}

	privKey, parseErr := parseSecp256k1PrivateKey(keyBytes)
	if parseErr != nil {
		return nil, 0, fmt.Errorf("invalid private key: %v", parseErr)
	}

	privateKey = privKey
	log.Printf("private key updated")
	return nil, 1, nil
}

// handleKeySign signs the original message with the stored private key.
// Returns the message and signature in data as ABI-encoded (bytes, bytes).
func handleKeySign(msg string) (data *string, status int, err error) {
	if privateKey == nil {
		return nil, 0, fmt.Errorf("no private key stored")
	}

	if msg == "" {
		return nil, 0, fmt.Errorf("originalMessage is empty")
	}

	msgBytes, hexErr := base.HexToBytes(msg)
	if hexErr != nil {
		return nil, 0, fmt.Errorf("invalid hex in originalMessage: %v", hexErr)
	}

	sig, signErr := signECDSA(privateKey, msgBytes)
	if signErr != nil {
		return nil, 0, fmt.Errorf("signing failed: %v", signErr)
	}

	encoded, abiErr := abiEncodeTwo(msgBytes, sig)
	if abiErr != nil {
		return nil, 0, fmt.Errorf("ABI encoding failed: %v", abiErr)
	}

	dataHex := base.BytesToHex(encoded)
	return &dataHex, 1, nil
}

// decryptViaNode calls the TEE node's /decrypt endpoint.
// ciphertext is the raw ECIES ciphertext bytes; it is JSON-serialized as base64
// in the request, matching the tee-node's DecryptRequest.EncryptedMessage []byte field.
// Returns the decrypted plaintext bytes (also base64-serialized by tee-node).
func decryptViaNode(ciphertext []byte) ([]byte, error) {
	url := fmt.Sprintf("http://localhost:%s/decrypt", signPort)
	reqBody, _ := json.Marshal(DecryptRequest{EncryptedMessage: ciphertext})

	resp, err := httpClient.Post(url, "application/json", bytes.NewReader(reqBody))
	if err != nil {
		return nil, fmt.Errorf("request error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("node returned %d: %s", resp.StatusCode, string(b))
	}

	var dr DecryptResponse
	if err := json.NewDecoder(resp.Body).Decode(&dr); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return dr.DecryptedMessage, nil
}

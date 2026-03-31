package app

import (
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"sign-extension/internal/base"
)

func setupTestServer(mockNodeURL string) *base.Server {
	// Override the signPort and httpClient for testing.
	parts := strings.Split(mockNodeURL, ":")
	port := parts[len(parts)-1]

	// Save and restore globals.
	origSignPort := signPort
	origClient := httpClient
	signPort = port
	httpClient = http.DefaultClient

	srv := base.New("0", port, Version, Register, ReportState)

	// Restore the original signPort so Register doesn't re-read env.
	signPort = origSignPort
	httpClient = origClient

	// Re-set for this test.
	signPort = port

	return srv
}

func makeActionBody(opType, opCommand, originalMessage string) string {
	df := map[string]interface{}{
		"instructionId":   "0x0000000000000000000000000000000000000000000000000000000000000001",
		"teeId":           "0x0000000000000000000000000000000001",
		"timestamp":       1234567890,
		"opType":          opType,
		"opCommand":       opCommand,
		"originalMessage": originalMessage,
	}
	dfJSON, _ := json.Marshal(df)

	action := map[string]interface{}{
		"data": map[string]interface{}{
			"id":            "0x0000000000000000000000000000000000000000000000000000000000000001",
			"type":          "instruction",
			"submissionTag": "submit",
			"message":       base.BytesToHex(dfJSON),
		},
	}
	body, _ := json.Marshal(action)
	return string(body)
}

func opTypeHex(s string) string {
	return base.VersionToHex(s) // reuse the same stringToBytes32Hex
}

func TestActionKeyUpdateAndSign(t *testing.T) {
	// Reset state.
	privateKey = nil

	privKeyBytes := big.NewInt(12345).FillBytes(make([]byte, 32))

	mockNode := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/decrypt" {
			var req DecryptRequest
			json.NewDecoder(r.Body).Decode(&req)
			json.NewEncoder(w).Encode(DecryptResponse{DecryptedMessage: privKeyBytes})
			return
		}
		http.Error(w, "not found", 404)
	}))
	defer mockNode.Close()

	// Point signPort and httpClient at mock.
	parts := strings.Split(mockNode.URL, ":")
	testPort := parts[len(parts)-1]
	signPort = testPort
	httpClient = http.DefaultClient

	srv := base.New("0", testPort, Version, Register, ReportState)

	// Step 1: Update key.
	updateBody := makeActionBody(opTypeHex("KEY"), opTypeHex("UPDATE"), base.BytesToHex([]byte("encrypteddata")))
	req := httptest.NewRequest(http.MethodPost, "/action", strings.NewReader(updateBody))
	w := httptest.NewRecorder()
	srv.Handler().ServeHTTP(w, req)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		t.Fatalf("update: status %d, body: %s", resp.StatusCode, body)
	}

	var updateResult base.ActionResult
	if err := json.Unmarshal(body, &updateResult); err != nil {
		t.Fatalf("unmarshal update result: %v", err)
	}
	if updateResult.Status != 1 {
		t.Fatalf("update failed: status=%d log=%v", updateResult.Status, updateResult.Log)
	}

	// Step 2: Sign a message.
	messageHex := base.BytesToHex([]byte("hello"))
	signBody := makeActionBody(opTypeHex("KEY"), opTypeHex("SIGN"), messageHex)
	req = httptest.NewRequest(http.MethodPost, "/action", strings.NewReader(signBody))
	w = httptest.NewRecorder()
	srv.Handler().ServeHTTP(w, req)

	resp = w.Result()
	body, _ = io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		t.Fatalf("sign: status %d, body: %s", resp.StatusCode, body)
	}

	var signResult base.ActionResult
	if err := json.Unmarshal(body, &signResult); err != nil {
		t.Fatalf("unmarshal sign result: %v", err)
	}
	if signResult.Status != 1 {
		t.Fatalf("sign failed: status=%d log=%v", signResult.Status, signResult.Log)
	}
	if signResult.Data == nil {
		t.Fatal("sign result data is nil")
	}

	// Decode the ABI-encoded (message, signature).
	dataBytes, err := base.HexToBytes(*signResult.Data)
	if err != nil {
		t.Fatalf("hex decode result data: %v", err)
	}

	msg, sig, err := abiDecodeTwo(dataBytes)
	if err != nil {
		t.Fatalf("abi decode: %v", err)
	}

	if string(msg) != "hello" {
		t.Errorf("message mismatch: got %q, want %q", string(msg), "hello")
	}
	if len(sig) != 65 {
		t.Errorf("expected 65-byte signature, got %d", len(sig))
	}
}

func TestActionSignWithoutKey(t *testing.T) {
	privateKey = nil
	signPort = "9999"

	srv := base.New("0", signPort, Version, Register, ReportState)

	messageHex := base.BytesToHex([]byte("hello"))
	body := makeActionBody(opTypeHex("KEY"), opTypeHex("SIGN"), messageHex)
	req := httptest.NewRequest(http.MethodPost, "/action", strings.NewReader(body))
	w := httptest.NewRecorder()
	srv.Handler().ServeHTTP(w, req)

	var result base.ActionResult
	json.NewDecoder(w.Result().Body).Decode(&result)
	if result.Status != 0 {
		t.Errorf("expected status 0 (error), got %d", result.Status)
	}
	if result.Log == nil || !strings.Contains(*result.Log, "no private key") {
		t.Errorf("expected 'no private key' error, got %v", result.Log)
	}
}

func TestActionUnknownOperation(t *testing.T) {
	privateKey = nil
	signPort = "9999"

	srv := base.New("0", signPort, Version, Register, ReportState)

	body := makeActionBody(opTypeHex("UNKNOWN"), opTypeHex("OP"), "0xdeadbeef")
	req := httptest.NewRequest(http.MethodPost, "/action", strings.NewReader(body))
	w := httptest.NewRecorder()
	srv.Handler().ServeHTTP(w, req)

	if w.Code != http.StatusNotImplemented {
		t.Errorf("expected 501, got %d", w.Code)
	}
}

func TestActionUpdateEmptyMessage(t *testing.T) {
	privateKey = nil
	signPort = "9999"

	srv := base.New("0", signPort, Version, Register, ReportState)

	body := makeActionBody(opTypeHex("KEY"), opTypeHex("UPDATE"), "")
	req := httptest.NewRequest(http.MethodPost, "/action", strings.NewReader(body))
	w := httptest.NewRecorder()
	srv.Handler().ServeHTTP(w, req)

	var result base.ActionResult
	json.NewDecoder(w.Result().Body).Decode(&result)
	if result.Status != 0 {
		t.Errorf("expected status 0, got %d", result.Status)
	}
	if result.Log == nil || !strings.Contains(*result.Log, "originalMessage is empty") {
		t.Errorf("expected 'originalMessage is empty' error, got %v", result.Log)
	}
}

func TestActionMethodNotAllowed(t *testing.T) {
	srv := base.New("0", "9999", Version, Register, ReportState)

	req := httptest.NewRequest(http.MethodGet, "/action", nil)
	w := httptest.NewRecorder()
	srv.Handler().ServeHTTP(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", w.Code)
	}
}

func TestStateEndpoint(t *testing.T) {
	privateKey = nil

	srv := base.New("0", "9999", Version, Register, ReportState)

	req := httptest.NewRequest(http.MethodGet, "/state", nil)
	w := httptest.NewRecorder()
	srv.Handler().ServeHTTP(w, req)

	var resp base.StateResponse
	json.NewDecoder(w.Result().Body).Decode(&resp)
	if resp.StateVersion == "" {
		t.Error("stateVersion is empty")
	}

	var state map[string]interface{}
	json.Unmarshal(resp.State, &state)
	if state["hasKey"] != false {
		t.Error("expected hasKey=false")
	}
}

func TestStateMethodNotAllowed(t *testing.T) {
	srv := base.New("0", "9999", Version, Register, ReportState)

	req := httptest.NewRequest(http.MethodPost, "/state", nil)
	w := httptest.NewRecorder()
	srv.Handler().ServeHTTP(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", w.Code)
	}
}

func TestActionDecryptionFailure(t *testing.T) {
	privateKey = nil

	mockNode := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, `{"message":"decryption error"}`)
	}))
	defer mockNode.Close()

	parts := strings.Split(mockNode.URL, ":")
	signPort = parts[len(parts)-1]
	httpClient = http.DefaultClient

	srv := base.New("0", signPort, Version, Register, ReportState)

	body := makeActionBody(opTypeHex("KEY"), opTypeHex("UPDATE"), base.BytesToHex([]byte("baddata")))
	req := httptest.NewRequest(http.MethodPost, "/action", strings.NewReader(body))
	w := httptest.NewRecorder()
	srv.Handler().ServeHTTP(w, req)

	var result base.ActionResult
	json.NewDecoder(w.Result().Body).Decode(&result)
	if result.Status != 0 {
		t.Errorf("expected status 0, got %d", result.Status)
	}
	if result.Log == nil || !strings.Contains(*result.Log, "decryption failed") {
		t.Errorf("expected decryption failure, got %v", result.Log)
	}
}

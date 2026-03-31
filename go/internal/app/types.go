package app

// DecryptRequest is sent to the TEE node's /decrypt endpoint.
// EncryptedMessage is []byte so it JSON-marshals as base64, matching the tee-node's
// DecryptRequest which also uses []byte.
type DecryptRequest struct {
	EncryptedMessage []byte `json:"encryptedMessage"`
}

// DecryptResponse is returned from the TEE node's /decrypt endpoint.
// DecryptedMessage is []byte so it JSON-unmarshals from base64, matching the tee-node's
// DecryptResponse which also uses []byte.
type DecryptResponse struct {
	DecryptedMessage []byte `json:"decryptedMessage"`
}

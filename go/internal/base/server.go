package base

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
)

// Server is the TEE extension HTTP server.
type Server struct {
	extPort     string
	signPort    string
	version     string
	versionHex  string
	framework   *Framework
	reportState ReportStateFunc
	mu          sync.Mutex
	mux         *http.ServeMux
}

// New creates and configures the extension server.
func New(extPort, signPort, version string, register RegisterFunc, reportState ReportStateFunc) *Server {
	f := &Framework{}
	register(f)

	s := &Server{
		extPort:     extPort,
		signPort:    signPort,
		version:     version,
		versionHex:  VersionToHex(version),
		framework:   f,
		reportState: reportState,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/action", s.handleAction)
	mux.HandleFunc("/state", s.handleState)
	s.mux = mux

	return s
}

// ListenAndServe starts the HTTP server.
func (s *Server) ListenAndServe() error {
	return http.ListenAndServe(":"+s.extPort, s.mux)
}

// Handler returns the http.Handler (for testing).
func (s *Server) Handler() http.Handler {
	return s.mux
}

func (s *Server) handleAction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "failed to read body", http.StatusBadRequest)
		return
	}

	var action Action
	if err := json.Unmarshal(body, &action); err != nil {
		http.Error(w, "invalid action JSON", http.StatusBadRequest)
		return
	}

	// action.Data.Message is the hex-encoded JSON bytes of the DataFixed struct.
	msgBytes, err := HexToBytes(action.Data.Message)
	if err != nil {
		http.Error(w, "invalid hex in message", http.StatusBadRequest)
		return
	}
	var df DataFixed
	if err := json.Unmarshal(msgBytes, &df); err != nil {
		http.Error(w, "invalid DataFixed JSON in message", http.StatusBadRequest)
		return
	}

	handler := s.framework.Lookup(df.OpType, df.OpCommand)
	if handler == nil {
		http.Error(w, "unsupported op type", http.StatusNotImplemented)
		return
	}

	// Serialize handler calls with exclusive lock.
	s.mu.Lock()
	data, status, handlerErr := handler(df.OriginalMessage)
	s.mu.Unlock()

	result := ActionResult{
		ID:            action.Data.ID,
		SubmissionTag: action.Data.SubmissionTag,
		OpType:        df.OpType,
		OpCommand:     df.OpCommand,
		Version:       s.versionHex,
		Status:        status,
		Data:          data,
	}

	switch {
	case status == 0:
		errMsg := "error: unknown"
		if handlerErr != nil {
			errMsg = fmt.Sprintf("error: %v", handlerErr)
		}
		result.Log = &errMsg
	case status == 1:
		okMsg := "ok"
		result.Log = &okMsg
	default:
		pendingMsg := "pending"
		result.Log = &pendingMsg
	}

	log.Printf("action %s: opType=%s opCommand=%s status=%d",
		action.Data.ID, OpTypeToString(df.OpType), OpTypeToString(df.OpCommand), status)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (s *Server) handleState(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	s.mu.Lock()
	stateData := s.reportState()
	s.mu.Unlock()

	resp := StateResponse{
		StateVersion: s.versionHex,
		State:        stateData,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

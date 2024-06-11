// Package protocol provides the protocol implementation.
// File: response.go
package protocol

const (
	SUCCESS = "success" // SUCCESS is the success status.
	FAILED  = "failed"  // FAILED is the failed status.
	ERROR   = "error"   // ERROR is the error status.
)

// Response is the response structure.
type Response struct {
	Message string                 `json:"message"`
	Status  string                 `json:"status"`
	Data    map[string]interface{} `json:"data"`
}

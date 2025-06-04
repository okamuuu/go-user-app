// internal/handler/response.go
package handler

// ErrorResponse はエラーレスポンスの共通構造です。
// swagger:response ErrorResponse
type ErrorResponse struct {
	Error string `json:"error" example:"invalid request"`
}

// LoginResponse はログイン成功時のレスポンスです。
type LoginResponse struct {
	Token string `json:"token" example:"your-jwt-token"`
}

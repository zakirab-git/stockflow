package dto

type APIErrorBody struct {
	Code      string `json:"code"`
	Message   string `json:"message,omitempty"`
	Available *int   `json:"available,omitempty"`
	SKU       string `json:"sku,omitempty"`
}

type APIErrorEnvelope struct {
	Error APIErrorBody `json:"error"`
}
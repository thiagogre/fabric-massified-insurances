package dto

type DocsResponse[T any] struct {
	Docs []T `json:"docs"`
}

type SuccessResponse struct {
	Success bool `json:"success"`
}

type ErrorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

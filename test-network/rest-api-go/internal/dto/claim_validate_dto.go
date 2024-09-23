package dto

type ClaimValidateRequest struct {
	Username   string `json:"username"`
	IsApproved bool   `json:"isApproved"`
}

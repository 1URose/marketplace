package dto

type StillValidResponse struct {
	StillValid bool   `json:"still_valid"`
	Detail     string `json:"detail"`
}

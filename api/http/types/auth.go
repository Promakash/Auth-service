package types

import (
	"auth_service/domain"
	"encoding/json"
	"github.com/google/uuid"
	"net/http"
)

type GenerateRequest struct {
	UserID uuid.UUID
}

func CreateGenerateRequest(r *http.Request) (*GenerateRequest, error) {
	var req GenerateRequest
	query := r.URL.Query()

	param := query.Get("uuid")

	if param == "" {
		return nil, domain.ErrMissingParameter
	}

	userID, err := uuid.Parse(param)
	if err != nil {
		return nil, err
	}
	req.UserID = userID
	return &req, nil
}

type GenerateResponse struct {
	AccessToken  domain.AccessToken  `json:"access_token"`
	RefreshToken domain.RefreshToken `json:"refresh_token"`
}

type RefreshRequest struct {
	RefreshToken domain.RefreshToken `json:"refresh_token"`
}

func CreateRefreshRequest(r *http.Request) (*RefreshRequest, error) {
	var req RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return &req, nil
}

type RefreshResponse struct {
	AccessToken  domain.AccessToken  `json:"access_token"`
	RefreshToken domain.RefreshToken `json:"refresh_token"`
}

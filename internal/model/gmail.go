package model

import "time"

const (
	oauth2Endpoint = ""
)

type RefreshResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

type OAuth2Get struct {
	TokenName    string    `json:"token_name"`
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiredAt    time.Time `json:"expired_at"`
}

type OAuth2Update struct {
	TokenName    string    `json:"token_name"`
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiredIn    int64     `json:"expired_in"`
	RefreshURL   string    `json:"refresh_url"`
	RequestAt    time.Time `json:"requested_at"`
}

func NewOAuth2Update(prevToken OAuth2Get, refreshResp RefreshResponse, RequestAt time.Time) (oau OAuth2Update) {
	oau.TokenName = prevToken.TokenName
	oau.AccessToken = refreshResp.AccessToken
	oau.RefreshToken = prevToken.RefreshToken
	oau.ExpiredIn = int64(refreshResp.ExpiresIn)
	oau.RefreshURL = oauth2Endpoint
	oau.RequestAt = RequestAt
	return oau
}

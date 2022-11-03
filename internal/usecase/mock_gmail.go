package usecase

import (
	"fmt"
	"gmail-sender/internal/model"

	"golang.org/x/oauth2"
)

type MockGmailClient struct {
	ErrFetchNewAccessToken bool
	ErrSend                bool
}

func (m *MockGmailClient) FetchNewAccessToken(refreshToken string) (resp model.RefreshResponse, err error) {
	if m.ErrFetchNewAccessToken {
		return model.RefreshResponse{}, fmt.Errorf("error")
	}

	resp = model.RefreshResponse{
		AccessToken:  "accesstoken",
		TokenType:    "Bearer",
		RefreshToken: "refreshtoken",
		ExpiresIn:    3599,
	}
	return resp, nil
}

func (m *MockGmailClient) Send(b []byte) (err error) {
	if m.ErrSend {
		return fmt.Errorf("error")
	}
	return nil
}

func (m *MockGmailClient) SetToken(tk oauth2.Token) {
	return
}

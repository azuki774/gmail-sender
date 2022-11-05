package usecase

import (
	"fmt"
	"gmail-sender/internal/model"
	"time"
)

type MockTokenRepo struct {
	ErrNotify bool
	ErrGet    bool
}

func (m *MockTokenRepo) Notify(oa model.OAuth2Update) error {
	if m.ErrNotify {
		return fmt.Errorf("error")
	}
	return nil
}

func (m *MockTokenRepo) Get() (model.OAuth2Get, error) {
	if m.ErrGet {
		return model.OAuth2Get{}, fmt.Errorf("error")
	}

	oa := model.OAuth2Get{
		TokenName:    "gmail",
		AccessToken:  "accesstoken",
		RefreshToken: "refreshtoken",
		ExpiredAt:    time.Date(2000, 1, 2, 1, 0, 0, 0, time.Local),
	}
	return oa, nil
}

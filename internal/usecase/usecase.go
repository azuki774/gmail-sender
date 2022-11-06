package usecase

import (
	"context"
	"fmt"
	"gmail-sender/internal/model"
	"time"

	"go.uber.org/zap"
	"golang.org/x/oauth2"
)

type Usecase struct {
	Logger      *zap.Logger
	TokenRepo   TokenRepo
	GmailClient GmailClient
}

type TokenRepo interface {
	Notify(model.OAuth2Update) error
	Get() (model.OAuth2Get, error)
}

type GmailClient interface {
	FetchNewAccessToken(refreshToken string) (resp model.RefreshResponse, err error)
	Send(context.Context, []byte) (err error)
	SetToken(tk oauth2.Token)
}

func (u *Usecase) Send(ctx context.Context, mc model.MailContent) (err error) {
	oa, err := u.TokenRepo.Get()
	if err != nil {
		u.Logger.Error("failed to get token from DB", zap.String("err", err.Error()), zap.Error(err))
		return err
	}
	u.Logger.Info("get access token from DB")

	if oa.ExpiredAt.Before(time.Now()) || oa.AccessToken == "" {
		u.Logger.Info("access token expired")
		err = u.RefineNewToken(ctx)
		if err != nil {
			u.Logger.Error("failed to fetch new token", zap.String("err", err.Error()), zap.Error(err))
			return err
		}
	}

	b := []byte(fmt.Sprintf("From: 'me'\r\n" +
		"To: azuki774s@gmail.com\r\n" +
		"Subject: TestSubject\r\n" +
		"\r\n" + "TestBody"))

	err = u.GmailClient.Send(ctx, b)
	if err != nil {
		return err
	}

	return nil
}

func (u *Usecase) RefineNewToken(ctx context.Context) (err error) {
	oa, err := u.TokenRepo.Get()
	if err != nil {
		return err
	}
	u.Logger.Info("get access expired token from DB")

	// access token refresh
	t := time.Now()
	resp, err := u.GmailClient.FetchNewAccessToken(oa.RefreshToken)
	if err != nil {
		u.Logger.Error("failed to fetch new token", zap.String("err", err.Error()), zap.Error(err))
		return err
	}

	oau := model.NewOAuth2Update(oa, resp, t)
	err = u.TokenRepo.Notify(oau)
	if err != nil {
		u.Logger.Error("failed to notify new token to DB", zap.String("err", err.Error()), zap.Error(err))
		return err
	}

	u.Logger.Info("update access token")
	u.GmailClient.SetToken(oauth2.Token{
		AccessToken: resp.AccessToken,
		TokenType:   "Bearer",
		Expiry:      t.Add(time.Duration(resp.ExpiresIn) * time.Second),
	})

	return nil
}

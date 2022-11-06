package usecase

import (
	"context"
	"gmail-sender/internal/model"
	"testing"

	"go.uber.org/zap"
)

func NewTestLogger() *zap.Logger {
	config := zap.NewProductionConfig()
	// config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	config.EncoderConfig = zap.NewProductionEncoderConfig()
	l, _ := config.Build()

	l.WithOptions(zap.AddStacktrace(zap.ErrorLevel))
	return l
}

func TestUsecase_Send(t *testing.T) {
	type fields struct {
		Logger         *zap.Logger
		TokenRepo      TokenRepo
		GmailClient    GmailClient
		DefaultContent func() model.MailContent
	}
	type args struct {
		ctx context.Context
		mc  model.MailContent
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "ok",
			fields: fields{
				Logger:      NewTestLogger(),
				TokenRepo:   &MockTokenRepo{},
				GmailClient: &MockGmailClient{},
				DefaultContent: func() model.MailContent {
					return model.MailContent{
						From:  "'me'",
						To:    "TO",
						Title: "TITLE",
					}
				},
			},
			args: args{
				ctx: context.Background(),
			},
			wantErr: false,
		},
		{
			name: "failed to get token",
			fields: fields{
				Logger:      NewTestLogger(),
				TokenRepo:   &MockTokenRepo{ErrGet: true},
				GmailClient: &MockGmailClient{},
				DefaultContent: func() model.MailContent {
					return model.MailContent{
						From:  "'me'",
						To:    "TO",
						Title: "TITLE",
					}
				},
			},
			args: args{
				ctx: context.Background(),
			},
			wantErr: true,
		},
		{
			name: "failed to fetch new token",
			fields: fields{
				Logger:      NewTestLogger(),
				TokenRepo:   &MockTokenRepo{},
				GmailClient: &MockGmailClient{ErrFetchNewAccessToken: true},
				DefaultContent: func() model.MailContent {
					return model.MailContent{
						From:  "'me'",
						To:    "TO",
						Title: "TITLE",
					}
				},
			},
			args: args{
				ctx: context.Background(),
			},
			wantErr: true,
		},
		{
			name: "failed to send email",
			fields: fields{
				Logger:      NewTestLogger(),
				TokenRepo:   &MockTokenRepo{},
				GmailClient: &MockGmailClient{ErrSend: true},
				DefaultContent: func() model.MailContent {
					return model.MailContent{
						From:  "'me'",
						To:    "TO",
						Title: "TITLE",
					}
				},
			},
			args: args{
				ctx: context.Background(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &Usecase{
				Logger:         tt.fields.Logger,
				TokenRepo:      tt.fields.TokenRepo,
				GmailClient:    tt.fields.GmailClient,
				DefaultContent: tt.fields.DefaultContent,
			}
			if err := u.Send(tt.args.ctx, tt.args.mc); (err != nil) != tt.wantErr {
				t.Errorf("Usecase.Send() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUsecase_RefineNewToken(t *testing.T) {
	type fields struct {
		Logger      *zap.Logger
		TokenRepo   TokenRepo
		GmailClient GmailClient
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "ok",
			fields: fields{
				Logger:      NewTestLogger(),
				TokenRepo:   &MockTokenRepo{},
				GmailClient: &MockGmailClient{},
			},
			args:    args{ctx: context.Background()},
			wantErr: false,
		},
		{
			name: "failed to new token",
			fields: fields{
				Logger:      NewTestLogger(),
				TokenRepo:   &MockTokenRepo{},
				GmailClient: &MockGmailClient{ErrFetchNewAccessToken: true},
			},
			args:    args{ctx: context.Background()},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &Usecase{
				Logger:      tt.fields.Logger,
				TokenRepo:   tt.fields.TokenRepo,
				GmailClient: tt.fields.GmailClient,
			}
			if err := u.RefineNewToken(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("Usecase.RefineNewToken() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

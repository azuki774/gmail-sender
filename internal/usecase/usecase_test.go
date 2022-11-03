package usecase

import (
	"context"
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
				Logger:      tt.fields.Logger,
				TokenRepo:   tt.fields.TokenRepo,
				GmailClient: tt.fields.GmailClient,
			}
			if err := u.Send(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("Usecase.Send() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

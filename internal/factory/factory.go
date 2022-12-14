package factory

import (
	"fmt"
	"gmail-sender/internal/client"
	"gmail-sender/internal/model"
	"gmail-sender/internal/server"
	"gmail-sender/internal/usecase"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func JSTTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	const layout = "2006-01-02T15:04:05+09:00"
	jst := time.FixedZone("Asia/Tokyo", 9*60*60)
	enc.AppendString(t.In(jst).Format(layout))
}

func NewLogger() (*zap.Logger, error) {
	config := zap.NewProductionConfig()
	// config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	config.EncoderConfig = zap.NewProductionEncoderConfig()
	config.EncoderConfig.EncodeTime = JSTTimeEncoder
	l, err := config.Build()

	l.WithOptions(zap.AddStacktrace(zap.ErrorLevel))
	if err != nil {
		fmt.Printf("failed to create logger: %v\n", err)
	}
	return l, err
}

func NewGmailConf() *oauth2.Config {
	conf := &oauth2.Config{
		ClientID:     os.Getenv("OAUTH_CLIENT_ID"),
		ClientSecret: os.Getenv("OAUTH_CLIENT_SECRET"),
		RedirectURL:  "http://localhost:8080/",
		Scopes:       []string{"https://mail.google.com/"},
		Endpoint:     google.Endpoint,
	}

	return conf
}

func NewGmailClient() *client.GmailClient {
	gc := client.GmailClient{
		Conf: NewGmailConf(),
	}
	return &gc
}

func NewDefaultContent() func() model.MailContent {
	return func() model.MailContent {
		return model.MailContent{
			From:  "'me'",
			To:    os.Getenv("SEND_EMAIL_TO"),
			Title: os.Getenv("SEND_EMAIL_TITLE"),
		}
	}
}

func NewTokenRepo(host string, port string) *client.TokenRepo {
	return &client.TokenRepo{Host: host, Port: port}
}

func NewUsecase(l *zap.Logger, t *client.TokenRepo, g *client.GmailClient, df func() model.MailContent) *usecase.Usecase {
	return &usecase.Usecase{Logger: l, TokenRepo: t, GmailClient: g, DefaultContent: df}
}

func NewServer(l *zap.Logger, port string, uc *usecase.Usecase) *server.Server {
	return &server.Server{Logger: l, Port: port, Usecase: uc}
}

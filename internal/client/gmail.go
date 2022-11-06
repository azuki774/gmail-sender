package client

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"gmail-sender/internal/model"
	"io"
	"net/http"
	"strings"

	"golang.org/x/oauth2"
	"google.golang.org/api/gmail/v1"
)

type GmailClient struct {
	Conf  *oauth2.Config
	Token oauth2.Token
}

func (g *GmailClient) FetchNewAccessToken(refreshToken string) (resp model.RefreshResponse, err error) {
	endpoint := "https://www.googleapis.com/oauth2/v4/token"

	reqData := model.RefreshRequest{
		ClientId:     g.Conf.ClientID,
		ClientSecret: g.Conf.ClientSecret,
		RefreshToken: refreshToken,
		RedirectUri:  g.Conf.RedirectURL,
		GrantType:    "refresh_token",
	}
	jsonData, err := json.Marshal(reqData)
	if err != nil {
		return model.RefreshResponse{}, err
	}

	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return model.RefreshResponse{}, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		resBody, _ := io.ReadAll(res.Body)
		fmt.Println(string(resBody))
		return model.RefreshResponse{}, fmt.Errorf("unexpected status code: %v", res.StatusCode)
	}
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return model.RefreshResponse{}, err
	}

	err = json.Unmarshal(resBody, &resp)
	if err != nil {
		return model.RefreshResponse{}, err
	}

	return resp, nil
}

func (g *GmailClient) Send(ctx context.Context, b []byte) (err error) {
	client := g.Conf.Client(ctx, &g.Token)
	srv, err := gmail.New(client)
	if err != nil {
		return fmt.Errorf("unable to retrieve gmail Client: %w", err)
	}

	var message gmail.Message
	message.Raw = base64.StdEncoding.EncodeToString(b)
	message.Raw = strings.Replace(message.Raw, "/", "_", -1)
	message.Raw = strings.Replace(message.Raw, "+", "-", -1)
	message.Raw = strings.Replace(message.Raw, "=", "", -1)

	_, err = srv.Users.Messages.Send("me", &message).Do()
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

func (g *GmailClient) SetToken(tk oauth2.Token) {
	g.Token = tk
}

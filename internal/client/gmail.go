package client

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"gmail-sender/internal/model"
	"io"
	"log"
	"net/http"
	"strings"

	"golang.org/x/oauth2"
	"google.golang.org/api/gmail/v1"
)

type GmailClient struct {
	AppKey    string
	AppSecret string
	Conf      *oauth2.Config
	Token     oauth2.Token
}

func (g *GmailClient) FetchNewAccessToken(refreshToken string) (resp model.RefreshResponse, err error) {
	endpoint := ""
	reqbody := fmt.Sprintf("refresh_token=%s&grant_type=refresh_token", refreshToken)
	reader := strings.NewReader(reqbody)

	req, err := http.NewRequest("POST", endpoint, reader)
	req.SetBasicAuth(g.AppKey, g.AppSecret)
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return model.RefreshResponse{}, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
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

func (g *GmailClient) Send(b []byte) (err error) {
	client := g.Conf.Client(oauth2.NoContext, &g.Token)
	srv, err := gmail.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve gmail Client %v", err)
	}

	var message gmail.Message
	message.Raw = base64.StdEncoding.EncodeToString(b)
	message.Raw = strings.Replace(message.Raw, "/", "_", -1)
	message.Raw = strings.Replace(message.Raw, "+", "-", -1)
	message.Raw = strings.Replace(message.Raw, "=", "", -1)

	_, err = srv.Users.Messages.Send("me", &message).Do()
	if err != nil {
		fmt.Printf("%v", err)
	}

	return nil
}

func (g *GmailClient) SetToken(tk oauth2.Token) {
	g.Token = tk
}

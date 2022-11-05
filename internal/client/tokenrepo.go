package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gmail-sender/internal/model"
	"io"
	"net"
	"net/http"
)

type TokenRepo struct {
	Host string
	Port string
}

const (
	tokenName = "gmail"
)

func (t *TokenRepo) Notify(oa model.OAuth2Update) error {
	endPoint := fmt.Sprintf("http://%s/", net.JoinHostPort(t.Host, t.Port)) + "oauth2/"
	client := &http.Client{}
	jsonData, err := json.Marshal(oa)
	if err != nil {
		return err
	}
	r, err := http.NewRequest("PUT", endPoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	r.Header.Set("Content-Type", "application/json")
	res, err := client.Do(r)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return fmt.Errorf("status code is not 200, but %v", res.StatusCode)
	}

	return nil
}

func (t *TokenRepo) Get() (model.OAuth2Get, error) {
	endPoint := fmt.Sprintf("http://%s/", net.JoinHostPort(t.Host, t.Port)) + "oauth2/" + tokenName
	res, err := http.Get(endPoint)
	if err != nil {
		return model.OAuth2Get{}, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return model.OAuth2Get{}, fmt.Errorf("status code is not 200, but %v", res.StatusCode)
	}

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return model.OAuth2Get{}, err
	}

	var oa model.OAuth2Get
	err = json.Unmarshal(b, &oa)
	if err != nil {
		return model.OAuth2Get{}, err
	}

	return oa, nil
}

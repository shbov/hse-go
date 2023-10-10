package Structs

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

const maxTimeoutRequest = 15 * time.Second

type Client struct {
	config *Config
}

func NewClient(config *Config) *Client {
	return &Client{config: config}
}

func (client *Client) MakeURL(path string) string {
	protocol := "http"

	if client.config.Https {
		protocol = "https"
	}

	return fmt.Sprintf("%s://%s:%d/%s", protocol, client.config.Host, client.config.Port, path)
}

func (client *Client) SendLongRequest(path string) (string, error) {
	ctx, stop := context.WithTimeout(context.Background(), maxTimeoutRequest)
	defer stop()

	req, err := http.NewRequestWithContext(ctx, "GET", client.MakeURL(path), nil)
	if err != nil {
		return "", err
	}

	httpClient := http.DefaultClient
	resp, err := httpClient.Do(req)
	if err != nil {
		if !errors.Is(err, context.DeadlineExceeded) {
			return "", err
		}

		return "false, timeout", nil
	}

	return fmt.Sprintf("true, %d", resp.StatusCode), nil
}

func (client *Client) SendGetRequest(path string) (string, error) {
	resp, err := http.Get(client.MakeURL(path))
	if err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		log.Warning(string(body))
	}

	return string(body), nil
}

func (client *Client) SendPostRequest(path string, params []byte) ([]byte, error) {
	req, err := http.NewRequest("POST", client.MakeURL(path), bytes.NewBuffer(params))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Error(err)
		}
	}(resp.Body)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

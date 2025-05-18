package transport

import (
	"aletheia-go-sdk/sdk/config"
	"aletheia-go-sdk/sdk/models"
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type HTTPTransport struct {
	client *http.Client
	config *config.AletheiaConfig
	logger config.Logger
}

func NewHTTPTransport(cfg config.AletheiaConfig) (*HTTPTransport, error) {
	tr, err := buildTransport(cfg)
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Transport: tr,
		Timeout:   cfg.RPCTimeout,
	}

	return &HTTPTransport{
		client: client,
		config: &cfg,
		logger: cfg.Logger,
	}, nil
}

func buildTransport(cfg config.AletheiaConfig) (*http.Transport, error) {
	if !cfg.UseTLS {
		return http.DefaultTransport.(*http.Transport), nil
	}

	caCert, err := ioutil.ReadFile(cfg.CaCertPath)
	if err != nil {
		return nil, fmt.Errorf("cannot read CA certificate: %w", err)
	}
	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(caCert) {
		return nil, fmt.Errorf("failed to append CA cert")
	}

	return &http.Transport{
		TLSClientConfig: &tls.Config{RootCAs: certPool},
	}, nil
}

func (t *HTTPTransport) SendEvent(event interface{}) error {
	var url string
	var bodyBytes []byte
	var err error

	switch req := event.(type) {
	case models.EventRequest:
		url = fmt.Sprintf("%s/api/v1/event", t.config.CollectorAddress)
		bodyBytes, err = json.Marshal(req)
	case models.ResourceRequest:
		url = fmt.Sprintf("%s/api/v1/resource", t.config.CollectorAddress)
		bodyBytes, err = json.Marshal(req)
	case models.ErrorEvent:
		url = fmt.Sprintf("%s/api/v1/error", t.config.CollectorAddress)
		bodyBytes, err = json.Marshal(req)

	default:
		return fmt.Errorf("unsupported event type")
	}

	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	httpReq, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(bodyBytes))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	if t.config.AuthToken != "" {
		httpReq.Header.Set("Authorization", "Bearer "+t.config.AuthToken)
	}

	resp, err := t.client.Do(httpReq)
	if err != nil {
		return fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var response struct {
		Success bool `json:"success"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	if !response.Success {
		return fmt.Errorf("server returned failure")
	}

	return nil
}

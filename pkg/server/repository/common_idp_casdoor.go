package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/ryo-arima/cmn-core/pkg/config"
	"github.com/ryo-arima/cmn-core/pkg/entity/model"
)

// ---- casdoorManager --------------------------------------------------------

type casdoorManager struct {
	cfg    config.CasdoorConfig
	client *http.Client
}

func newCasdoorManager(cfg config.CasdoorConfig) IdPManager {
	return &casdoorManager{
		cfg:    cfg,
		client: &http.Client{Timeout: 15 * time.Second},
	}
}

func (rcvr *casdoorManager) apiURL(path string) string {
	return fmt.Sprintf("%s%s", rcvr.cfg.BaseURL, path)
}

// authParams appends clientId and clientSecret to query params.
// Casdoor admin API authenticates via these query parameters.
func (rcvr *casdoorManager) authParams(q url.Values) url.Values {
	if q == nil {
		q = url.Values{}
	}
	q.Set("clientId", rcvr.cfg.ClientID)
	q.Set("clientSecret", rcvr.cfg.ClientSecret)
	return q
}

// do performs an authenticated request to the Casdoor admin API.
// Authentication is via clientId/clientSecret query parameters.
func (rcvr *casdoorManager) do(ctx context.Context, method, rawURL string, payload interface{}) (int, []byte, error) {
	// Append auth params to URL.
	u, err := url.Parse(rawURL)
	if err != nil {
		return 0, nil, fmt.Errorf("casdoor: parse URL: %w", err)
	}
	u.RawQuery = rcvr.authParams(u.Query()).Encode()

	var bodyReader io.Reader
	if payload != nil {
		b, err := json.Marshal(payload)
		if err != nil {
			return 0, nil, fmt.Errorf("casdoor: marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(b)
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), bodyReader)
	if err != nil {
		return 0, nil, fmt.Errorf("casdoor: build request: %w", err)
	}
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := rcvr.client.Do(req)
	if err != nil {
		return 0, nil, fmt.Errorf("casdoor: %s %s: %w", method, u.Path, err)
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)
	return resp.StatusCode, b, nil
}

// checkCdResponse unmarshals the Casdoor generic response and returns an error
// if the status is not "ok".
func checkCdResponse(body []byte) (*model.CdResponse, error) {
	var r model.CdResponse
	if err := json.Unmarshal(body, &r); err != nil {
		return nil, fmt.Errorf("casdoor: parse response: %w", err)
	}
	if r.Status != "ok" {
		return nil, fmt.Errorf("casdoor: API error: %s", r.Msg)
	}
	return &r, nil
}

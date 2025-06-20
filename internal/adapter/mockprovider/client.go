package mockprovider

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/e-wrobel/cash_ai/internal/domain"
)

type Client struct {
	baseURL string
	hc      *http.Client
}

func New(baseURL string) *Client {
	return &Client{baseURL: baseURL, hc: &http.Client{Timeout: 10 * time.Second}}
}

func (c *Client) Fetch() ([]domain.Transaction, error) {
	u := fmt.Sprintf("%s/transactions", c.baseURL)
	resp, err := c.hc.Get(u)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var txs []domain.Transaction
	err = json.NewDecoder(resp.Body).Decode(&txs)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response from %s: %w", u, err)
	}
	return txs, err
}

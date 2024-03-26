package exchangeRate

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/pkg/errors"
)

type Client struct {
	httpClient http.Client
	baseUrl    string
	apiKey     string
}

func (c *Config) NewClient() *Client {
	client := http.Client{Timeout: c.Timeout}

	return &Client{
		httpClient: client,
		baseUrl:    c.BaseUrl,
		apiKey:     c.ApiKey,
	}
}

func (c *Client) GetCurrencyRate(ctx context.Context, base string, target string) (float32, error) {
	url := c.createRequestUrl(base, target)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return 0, errors.Wrap(err, "create get currency rate http request")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return 0, errors.Wrap(err, "get currency rate http request failed")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, errors.Errorf("get exchange rate for pair %s/%s", base, target)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, errors.Wrap(err, "read response body from get exchange rate http response")
	}

	var currencyRate struct {
		Rate float32 `json:"conversion_rate"`
	}

	err = json.Unmarshal(body, &currencyRate)
	if err != nil {
		return 0, errors.Wrap(err, "unmarshal response body")
	}

	return currencyRate.Rate, nil
}

func (c *Client) createRequestUrl(base string, target string) string {
	return fmt.Sprintf("%s/%s/pair/%s/%s", c.baseUrl, c.apiKey, base, target)
}

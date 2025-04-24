package stockapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/RobertCastro/stock-insights-api/internal/domain/models"
)

const (
	defaultBaseURL   = "https://api.stockapi.com/v1/stocks"
	defaultAuthToken = ""
)

type APIResponse struct {
	Items    []models.Stock `json:"items"`
	NextPage string         `json:"next_page"`
}

// Representa un error de la API
type APIError struct {
	StatusCode int
	Body       string
	URL        string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("API returned status %d for URL %s: %s", e.StatusCode, e.URL, e.Body)
}

type Client struct {
	httpClient *http.Client
	baseURL    string
	authToken  string
}

func NewClient() *Client {
	baseURL := os.Getenv("STOCK_API_BASE_URL")
	if baseURL == "" {
		baseURL = defaultBaseURL
	}

	authToken := os.Getenv("STOCK_API_AUTH_TOKEN")
	if authToken == "" {
		authToken = defaultAuthToken
	}

	return &Client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL:   baseURL,
		authToken: authToken,
	}
}

func (c *Client) FetchStocks(nextPage string) ([]models.Stock, string, error) {
	if c.authToken == "" {
		return nil, "", fmt.Errorf("no se ha configurado el token de autenticación (STOCK_API_AUTH_TOKEN)")
	}

	// Parámetros de paginación
	reqURL := c.baseURL
	if nextPage != "" {
		params := url.Values{}
		params.Add("next_page", nextPage)
		reqURL = fmt.Sprintf("%s?%s", c.baseURL, params.Encode())
	}

	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, "", fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Add("Authorization", "Bearer "+c.authToken)
	req.Header.Add("Content-Type", "application/json")

	// Realizar la solicitud
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, "", fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", fmt.Errorf("error reading response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusGone {
			return nil, "", fmt.Errorf("API resource is no longer available (410 Gone). The API endpoint might have been deprecated or moved")
		}

		return nil, "", &APIError{
			StatusCode: resp.StatusCode,
			Body:       string(bodyBytes),
			URL:        reqURL,
		}
	}

	var apiResp APIResponse
	if err := json.Unmarshal(bodyBytes, &apiResp); err != nil {
		return nil, "", fmt.Errorf("error decoding response: %w", err)
	}

	return apiResp.Items, apiResp.NextPage, nil
}

// Recuperamos todos los stocks paginando
func (c *Client) FetchAllStocks() ([]models.Stock, error) {
	var allStocks []models.Stock
	nextPage := ""
	maxRetries := 3
	retryCount := 0

	for {
		stocks, newNextPage, err := c.FetchStocks(nextPage)
		if err != nil {
			if err.Error() == "API resource is no longer available (410 Gone). The API endpoint might have been deprecated or moved" {
				return nil, err
			}

			retryCount++
			if retryCount <= maxRetries {
				time.Sleep(2 * time.Second)
				continue
			}
			return nil, err
		}

		retryCount = 0

		if len(stocks) > 0 {
			allStocks = append(allStocks, stocks...)
		}

		if newNextPage == "" {
			break
		}

		nextPage = newNextPage
	}

	return allStocks, nil
}

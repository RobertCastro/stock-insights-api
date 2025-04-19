package stockapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/RobertCastro/stock-insights-api/internal/domain/models"
)

const (
	baseURL   = "https://8j5baasof2.execute-api.us-west-2.amazonaws.com/production/swechallenge/list"
	authToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdHRlbXB0cyI6MTcsImVtYWlsIjoib25lcm9iZXJ0aEBnbWFpbC5jb20iLCJleHAiOjE3NDQyNDMyMjUsImlkIjoiMCIsInBhc3N3b3JkIjoiJyBPUiAnMSc9JzEifQ.rLkT-QxbUQYL0fGjIDfz0EHhD_wqKS3xH0cyfFR6ZCM"
)

type APIResponse struct {
	Items    []models.Stock `json:"items"`
	NextPage string         `json:"next_page"`
}

type Client struct {
	httpClient *http.Client
}

func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{},
	}
}

func (c *Client) FetchStocks(nextPage string) ([]models.Stock, string, error) {
	// Parámetros de paginación
	reqURL := baseURL
	if nextPage != "" {
		params := url.Values{}
		params.Add("next_page", nextPage)
		reqURL = fmt.Sprintf("%s?%s", baseURL, params.Encode())
	}

	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, "", fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Add("Authorization", "Bearer "+authToken)
	req.Header.Add("Content-Type", "application/json")

	// Realizar la solicitud
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, "", fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, "", fmt.Errorf("API returned non-200 status: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	var apiResp APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, "", fmt.Errorf("error decoding response: %w", err)
	}

	return apiResp.Items, apiResp.NextPage, nil
}

// Recuperamos todos los stocks paginando
func (c *Client) FetchAllStocks() ([]models.Stock, error) {
	var allStocks []models.Stock
	nextPage := ""

	for {
		stocks, newNextPage, err := c.FetchStocks(nextPage)
		if err != nil {
			return nil, err
		}

		allStocks = append(allStocks, stocks...)

		if newNextPage == "" {
			break
		}

		nextPage = newNextPage
	}

	return allStocks, nil
}

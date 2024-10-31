package util

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

var apiKey = os.Getenv("API_KEY")

const baseURL = "https://api.exchangerate-api.com/v4/latest/"

type ExchangeRateResponse struct {
	Rates map[string]float64 `json:"rates"`
}

func getExchangeRates(baseCurrency string) (map[string]float64, error) {
	url := fmt.Sprintf("%s%s?apikey=%s", baseURL, baseCurrency, apiKey)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error fetching exchange rates: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error fetching data: %s", resp.Status)
	}

	var result ExchangeRateResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return result.Rates, nil
}

func ConvertINRtoUSD(amountINR float64) (float64, error) {
	rates, err := getExchangeRates("INR")
	if err != nil {
		return 0, fmt.Errorf("failed to get exchange rates: %w", err)
	}

	rate, ok := rates["USD"]
	if !ok {
		return 0, fmt.Errorf("USD rate not found")
	}

	return amountINR * rate, nil
}

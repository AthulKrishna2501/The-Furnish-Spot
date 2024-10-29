package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const apiKey = "c469eaee3c47fc40a9cf9154"                     // Replace with your actual API key
const baseURL = "https://api.exchangerate-api.com/v4/latest/" // Example base URL

type ExchangeRateResponse struct {
	Rates map[string]float64 `json:"rates"`
}

func getExchangeRates(baseCurrency string) (map[string]float64, error) {
	url := fmt.Sprintf("%s%s?apikey=%s", baseURL, baseCurrency)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error fetching data: %s", resp.Status)
	}

	var result ExchangeRateResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Rates, nil
}

func convertINRtoUSD(amountINR float64) (float64, error) {
	rates, err := getExchangeRates("INR")
	if err != nil {
		return 0, err
	}

	rate, ok := rates["USD"]
	if !ok {
		return 0, fmt.Errorf("USD rate not found")
	}

	return amountINR * rate, nil
}

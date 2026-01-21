package external

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"ewallet-ums/helpers"
)

type Wallet struct {
	ID      int     `json:"id"`
	UserID  int     `json:"user_id"`
	Balance float64 `json:"balance"`
}

type ExtWallet struct{}

func (e *ExtWallet) CreateWallet(ctx context.Context, userID int) (*Wallet, error) {
	req := Wallet{UserID: userID}

	payload, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal json: %v", err)
	}

	url := helpers.GetEnv("WALLET_HOST", "") + helpers.GetEnv("WALLET_ENDPOINT_CREATE", "")

	httpReq, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, fmt.Errorf("failed to create wallet http request: %v", err)
	}

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to connect wallet service: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("got error response from wallet service %d", resp.StatusCode)
	}

	result := &Wallet{}
	err = json.NewDecoder(resp.Body).Decode(result)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}
	defer resp.Body.Close()

	return result, nil
}

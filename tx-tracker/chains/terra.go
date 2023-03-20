package chains

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/wormhole-foundation/wormhole-explorer/txtracker/config"
)

func fetchTerraTx(
	ctx context.Context,
	cfg *config.RpcProviderSettings,
	txHash string,
) (*TxDetail, error) {

	// build the HTTP request
	url := fmt.Sprintf("%s/txs/%s", cfg.TerraBaseUrl, txHash)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	// read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// check the response status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Unexpected HTTP status code: %d (body %s)", resp.StatusCode, string(body))
	}

	// deserialize the response body
	var terraResponse terraResponse
	err = json.Unmarshal(body, &terraResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal terra response from API: %w", err)
	}

	// get the tx timestamp
	txDetail := TxDetail{
		NativeTxHash: terraResponse.Tx.TxHash,
	}
	txDetail.Timestamp, err = time.Parse("2006-01-02T15:04:05Z", terraResponse.Timestamp)
	if err != nil {
		return nil, fmt.Errorf("failed to parse tx timestamp: %w", err)
	}

	// get the tx sender
	if len(terraResponse.Tx.Value.Msg) > 0 {
		txDetail.Signer = terraResponse.Tx.Value.Msg[0].Value.Sender
	}
	if txDetail.Signer == "" {
		return nil, errors.New("can't find tx sender")
	}

	return &txDetail, nil
}

type terraResponse struct {
	Tx        terraTx `json:"tx"`
	Timestamp string  `json:"timestamp"`
}

type terraTx struct {
	Type_  string       `json:"type"`
	Value  terraTxValue `json:"value"`
	TxHash string       `json:"txhash"`
}

type terraTxValue struct {
	Memo string       `json:"memo"`
	Msg  []terraTxMsg `json:"msg"`
}

type terraTxMsg struct {
	Type_ string              `json:"type"`
	Value terraTxMessageValue `json:"value"`
}

type terraTxMessageValue struct {
	Contract string `json:"contract"`
	Sender   string `json:"sender"`
}

package chaininfo

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Chain struct {
	Name    string   `json:"name"`
	ChainID int64    `json:"chainId"`
	RpcURL  []string `json:"rpc"`
}

func (c *Chain) checkRpcURL(url string) (bool, error) {
	reqBody := []byte(`{"jsonrpc": "2.0", "method": "eth_blockNumber", "params": [], "id": 1}`)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return false, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	var result struct {
		JsonRPC string `json:"jsonrpc"`
		Result  string `json:"result"`
		Id      int    `json:"id"`
	}

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return false, err
	}

	return result.JsonRPC == "2.0" && result.Result != "", nil
}

func fetchChainsData() ([]Chain, error) {
	resp, err := http.Get("https://chainid.network/chains.json")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	chainsData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var chains []Chain
	err = json.Unmarshal(chainsData, &chains)
	if err != nil {
		return nil, err
	}

	return chains, nil
}

func GetChainInfo(chainID int64) (*Chain, error) {
	chains, err := fetchChainsData()
	if err != nil {
		return nil, fmt.Errorf("无法获取 JSON 数据：%v", err)
	}

	for _, chain := range chains {
		if chain.ChainID != chainID {
			continue
		}

		for _, url := range chain.RpcURL {
			if strings.Contains(url, "INFURA_API_KEY") {
				continue
			}

			isValid, err := chain.checkRpcURL(url)
			if err != nil {
				continue
			}

			if isValid {
				return &Chain{
					Name:    chain.Name,
					ChainID: chain.ChainID,
					RpcURL:  []string{url},
				}, nil
			}
		}

		break
	}

	return nil, errors.New("未找到有效的链或 RPC URL")
}

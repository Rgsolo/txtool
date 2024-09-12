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

type ChainInfo struct {
	Name           string         `json:"name"`
	Chain          string         `json:"chain"`
	Icon           string         `json:"icon"`
	RPC            []string       `json:"rpc"`
	Features       []Feature      `json:"features"`
	Faucets        []string       `json:"faucets"`
	NativeCurrency NativeCurrency `json:"nativeCurrency"`
	InfoURL        string         `json:"infoURL"`
	ShortName      string         `json:"shortName"`
	ChainID        int            `json:"chainId"`
	NetworkID      int            `json:"networkId"`
	Slip44         int            `json:"slip44"`
	ENS            ENS            `json:"ens"`
	Explorers      []Explorer     `json:"explorers"`
}

type Feature struct {
	Name string `json:"name"`
}

type NativeCurrency struct {
	Name     string `json:"name"`
	Symbol   string `json:"symbol"`
	Decimals int    `json:"decimals"`
}

type ENS struct {
	Registry string `json:"registry"`
}

type Explorer struct {
	Name     string `json:"name"`
	URL      string `json:"url"`
	Icon     string `json:"icon,omitempty"`
	Standard string `json:"standard"`
}

//type Chain struct {
//	Name    string   `json:"name"`
//	ChainID int64    `json:"chainId"`
//	RpcURL  []string `json:"rpc"`
//}

func (c *ChainInfo) checkRpcURL(url string) (bool, error) {
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

func fetchChainsData() ([]ChainInfo, error) {
	resp, err := http.Get("https://chainid.network/chains.json")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	chainsData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var chains []ChainInfo
	err = json.Unmarshal(chainsData, &chains)
	if err != nil {
		return nil, err
	}

	return chains, nil
}

func GetChainInfo(chainID int) (*ChainInfo, error) {
	chain := GetDefaultChainsData(chainID)
	if chain != nil {
		return chain, nil
	}

	chains, err := fetchChainsData()
	if err != nil {
		return nil, fmt.Errorf("无法获取 JSON 数据：%v", err)
	}

	for _, chain := range chains {
		if chain.ChainID != chainID {
			continue
		}

		for _, url := range chain.RPC {
			if strings.Contains(url, "INFURA_API_KEY") {
				continue
			}

			isValid, err := chain.checkRpcURL(url)
			if err != nil {
				continue
			}

			if isValid {
				return &ChainInfo{
					Name:    chain.Name,
					ChainID: chain.ChainID,
					RPC:     []string{url},
				}, nil
			}
		}

		break
	}

	return nil, errors.New("未找到有效的链或 RPC URL")
}

func GetChainInfoByName(chainName string) (*ChainInfo, error) {

	// 忽略大小写查找本地 CMap
	chain := GetChainByName(chainName)
	if chain != nil {
		return chain, nil
	}

	// 如果本地找不到，则从 fetchChainsData 查找
	chains, err := fetchChainsData()
	if err != nil {
		return nil, fmt.Errorf("无法获取 JSON 数据：%v", err)
	}

	for _, chain := range chains {
		// 忽略大小写匹配名称
		if strings.EqualFold(chain.Chain, chainName) {
			return &chain, nil
		}
	}

	return nil, errors.New("未找到有效的链或 RPC URL")
}

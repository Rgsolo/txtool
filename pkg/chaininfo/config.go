package chaininfo

import "strings"

var CMap = map[int]*ChainInfo{
	//1:          {Name: "ETH", ChainID: 1, RPC: []string{"https://rpc.ankr.com/eth"}},
	//56:         {Name: "BSC", ChainID: 56, RPC: []string{"https://rpc.ankr.com/bsc"}},
	//137:        {Name: "POLYGON", ChainID: 137, RPC: []string{"https://rpc.ankr.com/polygon"}},
	//42161:      {Name: "Arbitrum", ChainID: 42161, RPC: []string{"https://rpc.ankr.com/arbitrum"}},
	//43114:      {Name: "Avalanche", ChainID: 43114, RPC: []string{"https://rpc.ankr.com/avalanche"}},
	//10:         {Name: "Optimism", ChainID: 10, RPC: []string{"https://rpc.ankr.com/optimism"}},
	323:        {Name: "MCOIN", ChainID: 323, RPC: []string{"https://rpc1.m20chain.com"}},
	250:        {Name: "Fantom", ChainID: 250, RPC: []string{"https://rpc.ankr.com/fantom"}},
	42220:      {Name: "Celo", ChainID: 42220, RPC: []string{"https://rpc.ankr.com/celo"}},
	1284:       {Name: "Moonbeam", ChainID: 1284, RPC: []string{"https://rpc.ankr.com/moonbeam"}},
	592:        {Name: "Astar", ChainID: 592, RPC: []string{"https://evm.astar.network"}},
	1370:       {Name: "Rama", ChainID: 1370, RPC: []string{"https://blockchain.ramestta.com"}},
	14:         {Name: "FLR", ChainID: 14, RPC: []string{"https://flare-api.flare.network/ext/C/rpc"}},
	1501795822: {Name: "APEX", ChainID: 1501795822, RPC: []string{"https://rpc.theapexchain.org"}},
	5500:       {Name: "GODE", ChainID: 5500, RPC: []string{"https://rpc.godechain.net"}},
	256256:     {Name: "CMP", ChainID: 256256, RPC: []string{"https://mainnet.block.caduceus.foundation"}},
	1131:       {Name: "OZONE", ChainID: 1131, RPC: []string{"https://chain.ozonescan.com/"}},
	21:         {Name: "CTXC", ChainID: 21, RPC: []string{"http://10.95.6.218:8545"}},
}

func GetDefaultChainsData(chainID int) *ChainInfo {
	return CMap[chainID]
}

func GetChainByName(chainName string) *ChainInfo {
	for _, chain := range CMap {
		if strings.EqualFold(chain.Name, chainName) {
			return chain
		}
	}
	return nil
}

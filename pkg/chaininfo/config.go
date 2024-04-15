package chaininfo

var CMap = map[int64]*Chain{
	1:          {Name: "ETH", ChainID: 1, RpcURL: []string{"https://rpc.ankr.com/eth"}},
	56:         {Name: "BSC", ChainID: 56, RpcURL: []string{"https://rpc.ankr.com/bsc"}},
	137:        {Name: "POLYGON", ChainID: 137, RpcURL: []string{"https://rpc.ankr.com/polygon"}},
	42161:      {Name: "Arbitrum", ChainID: 42161, RpcURL: []string{"https://rpc.ankr.com/arbitrum"}},
	43114:      {Name: "Avalanche", ChainID: 43114, RpcURL: []string{"https://rpc.ankr.com/avalanche"}},
	10:         {Name: "Optimism", ChainID: 10, RpcURL: []string{"https://rpc.ankr.com/optimism"}},
	250:        {Name: "Fantom", ChainID: 250, RpcURL: []string{"https://rpc.ankr.com/fantom"}},
	42220:      {Name: "Celo", ChainID: 42220, RpcURL: []string{"https://rpc.ankr.com/celo"}},
	1284:       {Name: "Moonbeam", ChainID: 1284, RpcURL: []string{"https://rpc.ankr.com/moonbeam"}},
	592:        {Name: "Astar", ChainID: 592, RpcURL: []string{"https://evm.astar.network"}},
	1370:       {Name: "Rama", ChainID: 1370, RpcURL: []string{"https://blockchain.ramestta.com"}},
	14:         {Name: "FLR", ChainID: 14, RpcURL: []string{"https://flare-api.flare.network/ext/C/rpc"}},
	1501795822: {Name: "APEX", ChainID: 1501795822, RpcURL: []string{"https://rpc.theapexchain.org"}},
	5500:       {Name: "GODE", ChainID: 5500, RpcURL: []string{"https://rpc.godechain.net"}},
	256256:     {Name: "CMP", ChainID: 256256, RpcURL: []string{"https://mainnet.block.caduceus.foundation"}},
	1131:       {Name: "OZONE", ChainID: 1131, RpcURL: []string{"https://chain.ozonescan.com/"}},
	21:         {Name: "CTXC", ChainID: 21, RpcURL: []string{"http://10.95.6.218:8545"}},
}

func GetDefaultChainsData(chainID int64) *Chain {
	return CMap[chainID]
}

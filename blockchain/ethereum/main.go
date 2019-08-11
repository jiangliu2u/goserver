package ethereum

import (
	"fmt"
	"github.com/ethereum/go-ethereum/ethclient"
)

var EthClient ethclient.Client

func Init() {
	EthClient, err := ethclient.Dial("https://testnet.infura.io")
	if err != nil {
		panic("Eth init error")
	}
	fmt.Println(EthClient)
}

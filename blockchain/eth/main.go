package eth

import (
	"fmt"
	"github.com/ethereum/go-ethereum/ethclient"
	"os"
)

var Client *ethclient.Client

func Init() {
	var err error
	Client, err = ethclient.Dial(os.Getenv("ETH_API_TESTNET"))
	if err != nil {
		panic("Eth init error")
	}
	fmt.Println("init ethclient")
}

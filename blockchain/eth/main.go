package eth

import (
	"fmt"
	"github.com/ethereum/go-ethereum/ethclient"
)

var Client *ethclient.Client

func Init() {
	Client, _ = ethclient.Dial("https://kovan.infura.io")
	//if err != nil {
	//	panic("Eth init error")
	//}
	fmt.Println("init ethclient")
	//fmt.Println(EthClient)
}

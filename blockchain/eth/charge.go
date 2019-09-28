package eth

import (
	"c-server/cache"
	"c-server/model"
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/core/types"
	"math/big"
	"strconv"
	"sync"
	"time"
)

const CONFIRM_NUMBER = 1

/*
获取最新的区块数
*/
func GetBlock() (*big.Int, error) {
	header, err := Client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		fmt.Println("获取区块错误")
		return big.NewInt(0), err
	}
	return header.Number, nil

}

/*
获取区块里的交易
@param blockNumber 区块高度
*/
func GetTransaction(blockNumber *big.Int) types.Transactions {
	block, err := Client.BlockByNumber(context.Background(), blockNumber)
	if err != nil {
		return nil

	}
	return block.Transactions()
}

var AlreadySynced sync.Map

func StartSync(ws *sync.WaitGroup) {
	defer ws.Done()
	ticker := time.NewTicker(time.Second)
	blockNumber := big.NewInt(0)
	blockNumberChannel := make(chan *big.Int)
	go SyncTx(blockNumberChannel)
	for range ticker.C {
		KEY := "eth:kovan:recent_synced"
		ex, err := cache.RedisClient.Exists(KEY).Result()
		if ex == 0 && err == nil {
			bn, err := GetBlock()
			if err != nil {
				fmt.Println("redis 设置错误")
				fmt.Println(err)
				continue
			}
			cache.RedisClient.Set(KEY, bn.String(), 0)
		}
		lastSynced, err := cache.RedisClient.Get(KEY).Result()
		if err != nil {
			fmt.Println("redis 获取错误")
			fmt.Println(err)
		}
		syncedBlock, err := strconv.ParseInt(lastSynced, 10, 64)
		if err != nil {
			fmt.Println(err)
			continue
		}
		blockNumber = big.NewInt(syncedBlock)
		break
	}
	ticker2 := time.NewTicker(time.Second)

	for range ticker2.C {
		bln, err := GetBlock()
		if err != nil {
			fmt.Println("获取新区块错误")
			fmt.Println(err)
			continue
		}
		if bln.Cmp(blockNumber) == 0 || bln.Cmp(blockNumber) == -1 {
			fmt.Println("无新区块..")
			continue
		}
		var add int64 = 1
		blockNumber = blockNumber.Add(blockNumber, big.NewInt(add))
		blockNumberChannel <- blockNumber
	}

}
func SyncTx(blockNumberChannel chan *big.Int) {
	for {
		select {
		case i := <-blockNumberChannel:
			wg := new(sync.WaitGroup)
			wg.Add(1)
			go func() {
				GetTx(i)
				wg.Done()
			}()
			wg.Wait()
		}
	}
}
func GetTx(i *big.Int) {
	if _, ok := AlreadySynced.Load("eth:SyncedBlock"); ok {
		fmt.Println("同步过了", i.String())
		return
	}
	wg := new(sync.WaitGroup)
	wg.Add(1)

	go func() {
		txes := GetTransaction(i)
		fmt.Println(txes)
		for _, tx := range txes {
			to := tx.To()
			if to == nil {
				continue
			}
			toAddr := to.Hex()
			value := tx.Value().String()
			chainID, err := Client.NetworkID(context.Background())
			if err != nil {
				fmt.Println(err)
				continue
			}
			var from string
			if msg, err := tx.AsMessage(types.NewEIP155Signer(chainID)); err == nil {
				from = msg.From().Hex() //
			}

			receipt, err := Client.TransactionReceipt(context.Background(), tx.Hash())
			if err != nil {
				fmt.Println(err)
				continue
			}
			if receipt.Status == 1 {
				var player model.Player
				if err := model.DB.Where("ethChargeAddress = ?", toAddr).First(&player).Error; err != nil {
					continue
				}
				fmt.Println(player.EthChargeAddress)
				fmt.Printf("玩家%d有充值,充值来源%s,大小:%s\n", player.Id, from, value)
			}
		}
		//Save <- i
		AlreadySynced.Store(i.String(), 1)
		KEY := "eth:kovan:recent_synced"
		cache.RedisClient.Set(KEY, i.String(), 0)
		wg.Done()
	}()
	wg.Wait()

}

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
func GetBlock() *big.Int {
	header, err := Client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return GetBlock()
	}
	return header.Number

}

/*
获取区块里的交易
@param blockNumber 区块高度
*/
func GetTransaction(blockNumber *big.Int) types.Transactions {
	block, err := Client.BlockByNumber(context.Background(), blockNumber)
	if err != nil {
		//time.Sleep(4 * time.Second)
		fmt.Println("同步出错区块:",blockNumber.String())
		return GetTransaction(blockNumber)

	}
	return block.Transactions()
}

var Synced sync.Map



var Blo chan *big.Int  //区块数
var Save chan *big.Int //已经获取玩交易的区块数

func StartSync() {
	time.Sleep(2 * time.Second)
	bln := GetBlock()
	KEY := "eth:kovan:recent_synced"
	ex, err := cache.RedisClient.Exists(KEY).Result()
	if ex == 0 && err == nil {
		cache.RedisClient.Set(KEY, bln.String(), 0)
	}
	lastSynced, err := cache.RedisClient.Get(KEY).Result()
	if err != nil {
		fmt.Println(err)
		//StartSync()
		return

	}
	n, _ := strconv.Atoi(lastSynced)
	ls := big.NewInt(int64(n + 5))
	if bln.Cmp(ls) == 0 {
		time.Sleep(2 * time.Second)
		fmt.Println("无新区块..")
		StartSync()
		return
	}
	end := bln.Int64()
	for i :=n; i < int(end)-1; i++ {
		b:= big.NewInt(int64(i))
		Blo <- b
		cache.RedisClient.Set(KEY,b.String(),0)
	}

	time.Sleep(25 * time.Second)
	StartSync()
}
func SyncTx() {
	for {
		select {
		case i := <-Blo:
			go GetTx(i)
		}
	}
}
func GetTx(i *big.Int) {
	if ok,_:= cache.RedisClient.HExists("eth:SyncedBlock",i.String()).Result();ok{
		//fmt.Println("同步过了",i.String())
		return
	}
	txes := GetTransaction(i)
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
	cache.RedisClient.HSet("eth:SyncedBlock",i.String(),1)
}

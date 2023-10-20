package bifcore

import (
	"fmt"
	"testing"
)

func TestNewClient(t *testing.T) {
	url := "http://192.168.4.111:27002"
	client, err := NewClient(url)
	if err != nil {
		t.Error(err.Error())
		return
	}
	blockNumber, err := client.GetBlockNumber()
	if err != nil {
		t.Error(err.Error())
		return
	}
	t.Log("blockNumber=", blockNumber)
	//blockNumber = 3171122
	header, err := client.GetBlockHeaderByNumber(blockNumber)
	if err != nil {
		t.Error(err.Error())
		return
	}
	t.Log(fmt.Sprintf("result=%+v\n", header))
	transactionInBlock, err := client.GetTransactionsByNumber(blockNumber)
	if err != nil {
		t.Error(err.Error())
		return
	}
	t.Log(fmt.Sprintf("transactionInBlock=%+v\n", transactionInBlock))

	t.Log(fmt.Sprintf("tx_count=%+v", transactionInBlock.Transactions[0].Transaction.Operations[0].SendGas.Input))

}

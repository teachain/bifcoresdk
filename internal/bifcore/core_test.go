package bifcore

import (
	"fmt"
	"testing"
)

func TestNewClient(t *testing.T) {
	//url := "http://test-bif-core.xinghuo.space"
	url := "http://192.168.3.89:27002"
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
	blockNumber = 100000
	header, err := client.GetBlockHeaderByNumber(blockNumber)
	if err != nil {
		t.Error("number", blockNumber, err.Error())
		return
	}
	t.Log(fmt.Sprintf("result=%+v\n", header))
	transactionInBlock, err := client.GetTransactionsByNumber(blockNumber)
	if err != nil {
		t.Error(err.Error())
		return
	}
	t.Log(fmt.Sprintf("transactionInBlock=%+v\n", transactionInBlock))

	t.Log(fmt.Sprintf("tx_count=%+v", transactionInBlock.TotalCount))

}

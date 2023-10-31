package bifcoresdk

import (
	"encoding/hex"
	"fmt"
	"github.com/teachain/bifcoresdk/proto"
	protoBuff "google.golang.org/protobuf/proto"
	"testing"
)

func TestNewSubscriptionClient(t *testing.T) {
	url := "ws://test-bif-core.xinghuo.space:7053"
	client, err := NewObserver(url, 200)
	if err != nil {
		t.Error(err.Error())
		return
	}

	//OnChainTxEnvStore := func(msg []byte) {
	//	fmt.Println("OnChainTxEnvStore")
	//	transactionEnvStore := new(proto.TransactionEnvStore)
	//	err := protoBuff.Unmarshal(msg, transactionEnvStore)
	//	if err != nil {
	//		//proto:cannot parse invalid wire-format data
	//		fmt.Println(err.Error())
	//		fmt.Println("data=", hex.EncodeToString(msg))
	//		return
	//	}
	//	fmt.Println("OnChainTxEnvStore number:", transactionEnvStore.LedgerSeq)
	//}
	OnChainLedgerTxs := func(msg []byte) {
		ledgerTxs := new(proto.LedgerTxs)
		err := protoBuff.Unmarshal(msg, ledgerTxs)
		if err != nil {
			//proto:cannot parse invalid wire-format data
			fmt.Println(err.Error())
			fmt.Println("data=", hex.EncodeToString(msg))
			return
		}
		fmt.Println("OnChainLedgerTxs number:", ledgerTxs.Header.Seq)
	}
	//onLedgerHeader := func(msg []byte) {
	//	ledgerHeader := new(proto.LedgerHeader)
	//	err := protoBuff.Unmarshal(msg, ledgerHeader)
	//	if err != nil {
	//		//proto:cannot parse invalid wire-format data
	//		fmt.Println(err.Error())
	//		return
	//	}
	//
	//	fmt.Println("onLedgerHeader number:", ledgerHeader.Seq)
	//}

	//client.AddResponseHandler(int64(proto.ChainMessageType_CHAIN_HELLO), OnChainHello)
	//client.AddRequestHandler(int64(proto.ChainMessageType_CHAIN_TX_ENV_STORE), OnChainTxEnvStore)
	//client.AddRequestHandler(int64(proto.ChainMessageType_CHAIN_LEDGER_TXS), OnChainLedgerTxs)
	//client.AddRequestHandler(int64(proto.ChainMessageType_CHAIN_LEDGER_HEADER), onLedgerHeader)
	client.RegisterChainLedgerTxs(OnChainLedgerTxs)
	err = client.SayHello()
	if err != nil {
		t.Error(err.Error())
	}
	client.Wait()
}

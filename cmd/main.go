package main

import (
	"encoding/hex"
	"fmt"
	"github.com/teachain/bifcoresdk"
	"github.com/teachain/bifcoresdk/proto"
	protoBuff "google.golang.org/protobuf/proto"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	url := "ws://test-bif-core.xinghuo.space:7053"
	client, err := bifcoresdk.NewObserver(url, 200)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	OnChainLedgerTxs := func(msg []byte) {
		ledgerTxs := new(proto.LedgerTxs)
		err := protoBuff.Unmarshal(msg, ledgerTxs)
		if err != nil {
			fmt.Println("Unmarshal", err.Error())
			fmt.Println("data=", hex.EncodeToString(msg))
			return
		}
		fmt.Println("OnChainLedgerTxs number:", ledgerTxs.Header.Seq, "time", time.Now().Unix())
	}
	client.RegisterChainLedgerTxs(OnChainLedgerTxs)
	err = client.SayHello()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-ch
	client.Stop()
	fmt.Println("application exited")

}

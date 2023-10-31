package bifcoresdk

import (
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/teachain/bifcoresdk/proto"
	protoBuff "google.golang.org/protobuf/proto"
	"sync"
	"time"
)

type Observer interface {
	AddRequestHandler(msgType int64, f func(msg []byte))
	AddResponseHandler(msgType int64, f func(msg []byte))
	SayHello() error
	RegisterChainLedgerTxs(f func(msg []byte))
	EnqueueMessage(msgType int64, request bool, sequence int64, data []byte) error
	GetSequence() int64
	Wait()
	Stop()
}
type observer struct {
	url              string
	mutex            sync.Mutex
	requestHandlers  map[int64]func(msg []byte)
	responseHandlers map[int64]func(msg []byte)
	sequence         int64
	conn             *websocket.Conn
	sendQueue        chan []byte
	wg               sync.WaitGroup
	cancelFunc       context.CancelFunc
}

func NewObserver(url string, cacheSize int) (Observer, error) {
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithCancel(context.Background())
	client := &observer{
		url:              url,
		requestHandlers:  make(map[int64]func(msg []byte)),
		sequence:         0,
		conn:             conn,
		sendQueue:        make(chan []byte, cacheSize),
		responseHandlers: make(map[int64]func(msg []byte)),
		cancelFunc:       cancel,
	}
	client.wg.Add(3)

	client.AddRequestHandler(int64(proto.OVERLAY_MESSAGE_TYPE_OVERLAY_MSGTYPE_PING), client.onReceivePingMessage)
	client.AddResponseHandler(int64(proto.OVERLAY_MESSAGE_TYPE_OVERLAY_MSGTYPE_PING), client.onReceivePongMessage)
	client.AddResponseHandler(int64(proto.ChainMessageType_CHAIN_HELLO), client.onChainHello)
	go client.readMessage(ctx, &(client.wg))
	go client.sendHeartbeat(ctx, &(client.wg))
	go client.sendMessage(ctx, &(client.wg))
	return client, nil
}

// AddRequestHandler 添加作为请求信息的处理函数
func (o *observer) AddRequestHandler(msgType int64, f func(msg []byte)) {
	o.mutex.Lock()
	defer o.mutex.Unlock()
	if o.requestHandlers == nil {
		o.requestHandlers = make(map[int64]func(msg []byte))
	}
	o.requestHandlers[msgType] = f
}

// AddResponseHandler 添加作为响应信息的处理函数
func (o *observer) AddResponseHandler(msgType int64, f func(msg []byte)) {
	o.mutex.Lock()
	defer o.mutex.Unlock()
	if o.responseHandlers == nil {
		o.responseHandlers = make(map[int64]func(msg []byte))
	}
	o.responseHandlers[msgType] = f
}

func (o *observer) buildPingMsg() ([]byte, error) {
	ping := &proto.Ping{
		Nonce: time.Now().UnixNano(),
	}
	return protoBuff.Marshal(ping)
}
func (o *observer) EnqueueMessage(msgType int64, request bool, sequence int64, data []byte) error {
	data, err := o.buildWsMessage(msgType, request, sequence, data)
	if err != nil {
		return err
	}
	o.sendQueue <- data
	return nil
}

func (o *observer) buildWsMessage(msgType int64, request bool, sequence int64, data []byte) ([]byte, error) {
	msg := &proto.WsMessage{
		Type:     msgType,
		Request:  request, //这条消息是否是请求，是请求为true,是响应则为false
		Sequence: sequence,
		Data:     make([]byte, 0, len(data)),
	}
	copy(msg.Data, data)
	return protoBuff.Marshal(msg)
}

func (o *observer) sendHeartbeat(parent context.Context, wg *sync.WaitGroup) {
	ctx, cancel := context.WithCancel(parent)
	defer func() {
		cancel()
		wg.Done()
	}()
	for {
		select {
		case <-ctx.Done():
			fmt.Println("sendHeartbeat exited")
			return
		default:
		}
		time.Sleep(time.Second * 5)
		msg, err := o.buildPingMsg()
		if err != nil {
			fmt.Println("buildPingMsg", err.Error())
			continue
		}
		err = o.EnqueueMessage(int64(proto.OVERLAY_MESSAGE_TYPE_OVERLAY_MSGTYPE_PING), true, o.sequence, msg)
		if err == nil {
			o.sequence++
		}
	}
}

// 把数据发送出去
func (o *observer) sendMessage(parent context.Context, wg *sync.WaitGroup) {
	ctx, cancel := context.WithCancel(parent)
	defer func() {
		cancel()
		wg.Done()
	}()
	for {
		select {
		case <-ctx.Done():
			fmt.Println("sendMessage exited")
			return
		case msg := <-o.sendQueue:
			err := o.conn.WriteMessage(websocket.BinaryMessage, msg)
			if err != nil {
				fmt.Println("WriteMessage:", err.Error())
				o.sendQueue <- msg
			}
		}
	}
}

// 读取信息
func (o *observer) readMessage(parent context.Context, wg *sync.WaitGroup) {
	ctx, cancel := context.WithCancel(parent)
	defer func() {
		cancel()
		wg.Done()
	}()
	for {
		select {
		case <-ctx.Done():
			fmt.Println("readMessage exited")
			return
		default:
		}
		_, msg, err := o.conn.ReadMessage()
		if err != nil {
			fmt.Println("readMessage:", err.Error())
			time.Sleep(time.Second * 1)
			continue
		}
		o.onReceiveMessage(msg)
	}
}
func (o *observer) onReceiveMessage(msg []byte) {
	o.mutex.Lock()
	defer o.mutex.Unlock()
	wsMessage := new(proto.WsMessage)
	err := protoBuff.Unmarshal(msg, wsMessage)
	if err != nil {
		fmt.Println("onReceiveMessage Unmarshal:", err.Error())
		return
	}
	if f, ok := o.requestHandlers[wsMessage.Type]; ok && wsMessage.Request {
		f(wsMessage.Data)
	} else if f, ok := o.responseHandlers[wsMessage.Type]; ok && !wsMessage.Request {
		f(wsMessage.Data)
	}
}

func (o *observer) onReceivePingMessage(msg []byte) {
	ping := new(proto.Ping)
	err := protoBuff.Unmarshal(msg, ping)
	if err != nil {
		fmt.Println("onReceivePingMessage Unmarshal:", err.Error())
		return
	}
	pong := &proto.Pong{
		Nonce: ping.Nonce,
	}
	data, err := protoBuff.Marshal(pong)
	if err != nil {
		fmt.Println("onReceivePingMessage Marshal:", err.Error())
		return
	}
	err = o.EnqueueMessage(int64(proto.OVERLAY_MESSAGE_TYPE_OVERLAY_MSGTYPE_PING), false, o.sequence, data)
	if err != nil {
		fmt.Println("onReceivePingMessage EnqueueMessage:", err.Error())
	}
}
func (o *observer) onReceivePongMessage(msg []byte) {
	//fmt.Println("onReceivePongMessage", time.Now().Unix())
}
func (o *observer) Wait() {
	o.wg.Wait()
}
func (o *observer) GetSequence() int64 {
	return o.sequence
}
func (o *observer) Stop() {
	if o.cancelFunc != nil {
		o.cancelFunc()
	}
	o.wg.Wait()
}
func (o *observer) RegisterChainLedgerTxs(f func(msg []byte)) {
	o.AddRequestHandler(int64(proto.ChainMessageType_CHAIN_LEDGER_TXS), f)
}
func (o *observer) SayHello() error {
	chainHello := &proto.ChainHello{
		Timestamp: time.Now().Unix(),
	}
	msg, err := protoBuff.Marshal(chainHello)
	if err != nil {
		return err
	}
	//必须say hello
	err = o.EnqueueMessage(int64(proto.ChainMessageType_CHAIN_HELLO), true, o.GetSequence(), msg)
	if err != nil {
		return err
	}
	return nil
}
func (o *observer) onChainHello(msg []byte) {
	chainStatus := new(proto.ChainStatus)
	err := protoBuff.Unmarshal(msg, chainStatus)
	if err != nil {
		fmt.Println("onChainHello", err.Error())
		return
	}
	fmt.Println("连接的节点地址:", chainStatus.SelfAddr)
	fmt.Println("区块版本号:", chainStatus.LedgerVersion)
	fmt.Println("监控程序版本号:", chainStatus.MonitorVersion)
	fmt.Println("星火链程序版本号:", chainStatus.ChainVersion)
	fmt.Println("时间戳:", chainStatus.Timestamp)
	fmt.Println("账号前缀:", chainStatus.AddressPrefix)
}

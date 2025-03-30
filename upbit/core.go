package upbit

import (
	"bnt/commons"
	"bnt/config"
	"bnt/tgmanager"
	"bnt/websocketmanager"

	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

var (
	exchange string
	pingMsg  string = "PING"
)

func pongWs(done <-chan struct{}) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			websocketmanager.SendMsg(exchange, pingMsg)
		case <-done:
			return
		}
	}
}

func subscribeWs(pairs []string, wg *sync.WaitGroup) {
	defer wg.Done()
	time.Sleep(time.Second * 1)

	var streamSlice []string
	for _, pair := range pairs {
		pairInfo := strings.Split(pair, ":")
		market, symbol := strings.ToUpper(pairInfo[0]), strings.ToUpper(pairInfo[1])
		streamSlice = append(streamSlice, fmt.Sprintf("'%s-%s'", market, symbol))
	}

	streams := strings.Join(streamSlice, ",")
	uuid := uuid.NewString()
	msg := fmt.Sprintf("[{'ticket':'%s'}, {'type': 'orderbook', 'codes': [%s]}, {'type': 'trade', 'is_only_realtime': 'true', 'codes': [%s]}]", uuid, streams, streams)

	websocketmanager.SendMsg(exchange, msg)
	fmt.Printf(websocketmanager.SubscribeMsg, exchange)
}

func receiveWs(done <-chan struct{}, msgQueue chan<- []byte) {
	for {
		select {
		case <-done:
			return
		default:
			_, msgBytes, err := websocketmanager.Conn(exchange).ReadMessage()
			if err != nil {
				tgmanager.HandleErr(exchange, err)
			}
			msgQueue <- msgBytes
		}
	}
}

func processWsMessages(done <-chan struct{}, msgQueue <-chan []byte) {
	for {
		select {
		case <-done:
			return
		case msgBytes := <-msgQueue:
			if strings.Contains(string(msgBytes), "status") {
				fmt.Println("UPB PONG") // {"status":"UP"}
			} else {
				var rJson interface{}
				commons.Bytes2Json(msgBytes, &rJson)
				rType := rJson.(map[string]interface{})["type"].(string)
				switch rType {
				case "orderbook":
					go SetOrderbook(exchange, rJson.(map[string]interface{}))
				case "trade":
					go SetTrade(exchange, rJson.(map[string]interface{}))
				}
			}
		}
	}
}

func Run(e string) {
	exchange = e
	pairs := config.GetPairs(exchange)
	var wg sync.WaitGroup
	done := make(chan struct{})
	wsQueue := make(chan []byte, 1) // WebSocket 메시지 큐

	// ping
	wg.Add(1)
	go func() {
		defer wg.Done()
		pongWs(done)
	}()

	// subscribe websocket stream
	wg.Add(1)
	go subscribeWs(pairs, &wg)

	// receive websocket msg
	wg.Add(1)
	go func() {
		defer wg.Done()
		receiveWs(done, wsQueue)
	}()

	// process websocket messages
	wg.Add(1)
	go func() {
		defer wg.Done()
		processWsMessages(done, wsQueue)
	}()

	wg.Wait()
	close(done)
}
